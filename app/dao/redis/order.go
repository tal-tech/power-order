/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-12 15:16:47
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-04-09 23:50:05
 * @Description:
 */
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	logger "github.com/tal-tech/loggerX"
	"github.com/tal-tech/xredis"
	"powerorder/app/constant"
	"powerorder/app/dao/mysql"
	"powerorder/app/utils"
	"time"
)

type Order struct {
	Dao
	uUserId uint64
	ctx     context.Context
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func NewOrderDao(ctx context.Context, uAppId uint, uUserId uint64) *Order {
	object := new(Order)
	object.uAppId = uAppId
	object.uUserId = uUserId
	object.ctx = ctx
	return object
}

/**
 * @description:组装订单信息存储，hash的Key
 * @params {OrderId} 订单ID
 * @return: hash的Key[xes_platform_order_{$app_id}_{$order_id}]
 */
func (this *Order) GenHashKey(OrderId string) string {
	pts := utils.GetPts(this.ctx)
	if pts {
		return fmt.Sprintf("pts_%s%03d_%s_%d", constant.Order_HashKeyPrx, this.uAppId, OrderId, this.uUserId)
	}
	return fmt.Sprintf("%s%03d_%s_%d", constant.Order_HashKeyPrx, this.uAppId, OrderId, this.uUserId)
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *Order) HMDel(OrderId string, keys []string) (back string, err error) {
	hashKey := this.GenHashKey(OrderId)

	var fields []interface{}
	for i := 0; i < len(keys); i++ {
		fields = append(fields, keys[i])
	}

	back, err = xredis.NewSimpleXesRedis(this.ctx, constant.Order_Redis_Cluster).HDel(hashKey, []interface{}{}, fields)
	return
}

/**
 * @description:根据OrderId获取哈希表
 * @params {OrderId} 订单ID
 * @return:
 */
func (this *Order) HMGet(OrderId string) (result map[string]interface{}, err error) {
	hashKey := this.GenHashKey(OrderId)
	back, err := xredis.NewSimpleXesRedis(this.ctx, constant.Order_Redis_Cluster).HGetAll(hashKey, []interface{}{})
	if err == nil {
		result, err = this.Decode(back)
	}
	return
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *Order) HMSet(OrderId string, value map[string]interface{}) error {

	if len(value) == 0 {
		logger.Ix(this.ctx, "dao.redis.order.HMSet", "hmset empty value", "")
		return nil
	}
	expire := this.GenExpireTime(value)

	if expire == constant.Order_ExpiredTimeNoCache {
		//不需要设置缓存
		return nil
	}
	var HashSubKeys []string
	var ok bool

	if HashSubKeys, ok = constant.Order_HashSubKey[this.uAppId]; !ok {
		logger.Ex(this.ctx, "dao.redis.order.HMSet error", "hmset error appid", "")
		return nil
	}

	var i int
	for k, v := range value {
		for i = 0; i < len(HashSubKeys); i++ {
			if k == HashSubKeys[i] {
				break
			}
		}

		if i == len(HashSubKeys) {
			logger.Wx(this.ctx, "dao.redis.order.HMSet error", "hmset err key", "")
			return nil
		}

		if k == constant.Order_HashSubKey_Info {
			var info mysql.OrderInfo
			if info, ok = v.(mysql.OrderInfo); ok {
				value[constant.Order_HashSubKey_UserId] = info.UserId
			}
		}

	}

	for k, v := range value {

		if v == constant.Order_HashFieldDefault {
			value[k] = v
			continue
		}
		if k == constant.Order_HashSubKey_UserId {
			value[k] = v
			continue
		}
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return err
		}
		value[k] = string(jsonBytes)

	}

	hashKey := this.GenHashKey(OrderId)

	_, err := xredis.NewSimpleXesRedis(this.ctx, constant.Order_Redis_Cluster).HMSet(hashKey, []interface{}{}, value)
	if err != nil {
		logger.Ex(this.ctx, "dao.redis.order.HMSet error", "hmset error", "")
		return err
	}
	if expire == constant.Order_ExpiredTimeNoOperation {
		//不需要操作缓存时间
		//logger.Dx(this.ctx, "dao.redis.order.HMSet","no operateion time","")
		return nil
	}
	_, err = xredis.NewSimpleXesRedis(this.ctx, constant.Order_Redis_Cluster).Expire(hashKey, []interface{}{}, expire)
	return err
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *Order) Decode(value map[string]string) (result map[string]interface{}, err error) {
	result = make(map[string]interface{})

	for k, v := range value {

		if k == constant.Order_HashSubKey_Info {
			if v == constant.Order_HashFieldDefault {
				result[k] = v
				continue
			}
			var info mysql.OrderInfo
			err = json.Unmarshal([]byte(v), &info)
			result[k] = info
		} else if k == constant.Order_HashSubKey_Detail {
			if v == constant.Order_HashFieldDefault {
				result[k] = v
				continue
			}
			var details []mysql.OrderDetail
			err = json.Unmarshal([]byte(v), &details)
			result[k] = details
		} else if k == constant.Order_HashSubKey_UserId {
			result[k] = v
		} else {
			var vtemp []map[string]interface{}

			if v != constant.Order_HashFieldDefault {
				err = json.Unmarshal([]byte(v), &vtemp)

			} else {
				vtemp = make([]map[string]interface{}, 0)
			}
			result[k] = vtemp
		}

	}

	return
}

/**
 * @description: 获取数据的缓存时长
 * @params {value}
 * @return: 返回数据的缓存时长
 */
func (this *Order) GenExpireTime(value map[string]interface{}) (expire int) {
	expire = constant.Order_ExpiredTimeNoOperation
	var secondsPassed int64
	if temp, ok := value[constant.Order_HashSubKey_Info]; ok {

		switch temp.(type) {
		case mysql.OrderInfo:
			OrderInfo := temp.(mysql.OrderInfo)

			secondsPassed = time.Now().Unix() - OrderInfo.CreatedTime.Unix()
			if secondsPassed > constant.Order_SecondsPassedFromCreatedOneMonth {
				//订单创建时间超过1个月 不设置缓存
				expire = constant.Order_ExpiredTimeOneDay
				return
			} else if secondsPassed > constant.Order_SecondsPassedFromCreatedOneDay {
				//订单创建时间超过1天 设置缓存时间为1天
				expire = constant.Order_ExpiredTimeOneDay
				return
			} else {
				expire = constant.Order_ExpiredTimeTwoMinutes
				return
			}
		default: // ?? 正常数据会出现这种情况吗
			//info值为-2,缓存时间设置为2分钟
			expire = constant.Order_ExpiredTimeTwoMinutes
			return
		}

	}
	return expire
}

/**
 * @description:获取redis缓存
 * @params {OrderIds} 订单ID
 * @return: redis批量订单信息
 */
func (this *Order) BatchGet(OrderIds []string) (ret map[string]map[string]interface{}, err error) {
	ret = make(map[string]map[string]interface{}, 0)
	for i := 0; i < len(OrderIds); i++ {
		orderData, err1 := this.HMGet(OrderIds[i])
		if err1 != nil {
			err = err1
			return
		}
		CacheValue, err := utils.ToStringStringMap(orderData)
		CacheValueDecoded, err := this.Decode(CacheValue)
		if err == nil && len(CacheValueDecoded) > 0 {
			ret[OrderIds[i]] = CacheValueDecoded
		}
	}
	return
}

/**
 * @description:  异步设置过期时间，解决tw无序问题。
 * @return:
 */
func (this *Order) SetExpire(value map[string]map[string]interface{}) {

	defer utils.Catch()
	if len(value) == 0 {
		logger.Ix(this.ctx, "dao.redis.order.SetExpire", "SetExpire empty value", "")
		return
	}
	for orderId, orderData := range value {

		if len(orderData) == 0 {
			continue
		}
		expire := this.GenExpireTime(orderData)
		if expire == constant.Order_ExpiredTimeNoCache {
			//不需要设置缓存
			continue
		}
		if _, ok := constant.Order_HashSubKey[this.uAppId]; !ok {
			logger.Ex(this.ctx, "dao.redis.order.SetExpire error", "SetExpire error appid", "")
			continue
		}
		hashKey := this.GenHashKey(orderId)
		if expire == constant.Order_ExpiredTimeNoOperation {
			//不需要操作缓存时间
			logger.Dx(this.ctx, "dao.redis.order.SetExpire", "no operateion time", "")
			continue
		}
		// 最多尝试设置3次
		for i := 0; i < 3; i++ {
			ret, err := xredis.NewSimpleXesRedis(this.ctx, constant.Order_Redis_Cluster).Expire(hashKey, []interface{}{}, expire)
			if ret && err == nil {
				break
			} else {
				time.Sleep(500 * time.Microsecond)
			}
		}
	}

	return
}
