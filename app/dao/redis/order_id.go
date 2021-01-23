/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-12 15:16:47
 * @LastEditors  : lichanglin@tal.com
 * @LastEditTime : 2020-02-06 19:45:51
 * @Description:
 */
package redis

import (
	"context"
	"fmt"
	logger "github.com/tal-tech/loggerX"
	"github.com/tal-tech/xredis"
	"powerorder/app/constant"
	"powerorder/app/utils"
	"time"
)

type OrderId struct {
	Dao
	uUserId uint64
	ctx     context.Context
}

func NewOrderIdDao(ctx context.Context, uAppId uint, uUserId uint64, instance string) *OrderId {
	object := new(OrderId)
	// object.uAppId = uAppId
	object.uUserId = uUserId
	object.ctx = ctx
	object.Dao.uAppId = uAppId
	object.Dao.instance = instance
	return object
}

/**
 * @description: 获取order_id列表的redis无序集合的key
 * @params {}
 * @return: xes_platform_order_id_001_1
 */
func (this *OrderId) GenSetKey() string {
	pts := utils.GetPts(this.ctx)
	if pts {
		return fmt.Sprintf("pts_%s%03d_%v", constant.OrderId_SetKeyPrx, this.uAppId, this.uUserId)
	}
	return fmt.Sprintf("%s%03d_%v", constant.OrderId_SetKeyPrx, this.uAppId, this.uUserId)
}

/**
 * @description: 获取order_id列表的redis无序集合的key
 * @params {}
 * @return: xes_platform_order_id_001_1_20200922
 */
func (this *OrderId) GenSetKeyV2() string {
	pts := utils.GetPts(this.ctx)
	today := time.Now().Format("20060102") // 拼接当天的日期
	if pts {
		return fmt.Sprintf("pts_%s%03d_%v_%s", constant.OrderId_SetKeyPrx, this.uAppId, this.uUserId, today)
	}

	return fmt.Sprintf("%s%03d_%v_%s", constant.OrderId_SetKeyPrx, this.uAppId, this.uUserId, today)
}

/**
 * @description:
 * @params {OrderIds}
 * @return:
 */
func (this *OrderId) SDel(OrderIds []string) (back int64, err error) {
	setKey := this.GenSetKey()

	if len(OrderIds) == 0 {
		//删除所有的数据
		_, err1 := xredis.NewSimpleXesRedis(this.ctx, constant.Order_Redis_Cluster).Del(setKey, nil)
		if err1 != nil {
			err = err1
			logger.Ex(this.ctx, "dao.redis.orderid.sdel error", "sdel failed!", err)
		}
		return
	}
	params := make([]interface{}, 0)
	for i := 0; i < len(OrderIds); i++ {
		params = append(params, OrderIds[i])
	}
	back, err = xredis.NewSimpleXesRedis(this.ctx, constant.Order_Redis_Cluster).SRem(setKey, []interface{}{}, params)
	return
}

/**
 * @description: 获取order_id列表，Redis 无序集合key列表
 * @params {}
 * @return: 订单ID集合
 */
func (this *OrderId) SGet() (back []string, err error) {
	setKey := this.GenSetKey()
	back, err = xredis.NewSimpleXesRedis(this.ctx, constant.Order_Redis_Cluster).SMembers(setKey, []interface{}{})
	return
}

/**
 * @description: 新增order_id列表
 * @params {OrderIds}
 * @return:
 */
func (this *OrderId) SAdd(OrderIds []string) (back int64, err error) {

	var bDefaultMemberExisted = false
	if len(OrderIds) == 0 {
		//logger.Ix(this.ctx,"dao.redis.orderid.SAdd","sadd empty OrderIds","")
		return
	}

	members := make([]interface{}, 0)
	for i := 0; i < len(OrderIds); i++ {
		members = append(members, OrderIds[i])
		if OrderIds[i] == constant.OrderId_MemberDefault {
			bDefaultMemberExisted = true
		}
	}
	setKey := this.GenSetKey()
	back, err = xredis.NewSimpleXesRedis(this.ctx, constant.Order_Redis_Cluster).SAdd(setKey, []interface{}{}, members)
	if bDefaultMemberExisted {
		//如果是设置一个空的数据，则设置过期时间
		_, err = xredis.NewSimpleXesRedis(this.ctx, constant.Order_Redis_Cluster).Expire(setKey, []interface{}{}, constant.OrderId_ExpiredTimeDefault)
	} else {
		_, err = xredis.NewSimpleXesRedis(this.ctx, constant.Order_Redis_Cluster).Expire(setKey, []interface{}{}, constant.OrderId_ExpiredTime30Days)
	}
	return

}

/**
 * @description: 获取order_id列表，Redis 无序集合key列表
 * @params {}
 * @return: 订单ID集合
 */
func (this *OrderId) SGetV2() (back []string, err error) {
	setKey := this.GenSetKeyV2()
	back, err = xredis.NewSimpleXesRedis(this.ctx, constant.Order_Redis_Cluster).SMembers(setKey, []interface{}{})
	return
}

/**
 * @description: 新增order_id列表
 * @params {OrderIds}
 * @return:
 */
func (this *OrderId) SAddV2(OrderIds []string) (back int64, err error) {

	if len(OrderIds) == 0 {
		logger.Ix(this.ctx, "dao.redis.orderid.SAddV2", "sadd empty OrderIds", "")
		return
	}
	var members []interface{}
	for i := 0; i < len(OrderIds); i++ {
		members = append(members, OrderIds[i])
	}
	setKey := this.GenSetKeyV2()

	back, err = xredis.NewSimpleXesRedis(this.ctx, constant.Order_Redis_Cluster).SAdd(setKey, []interface{}{}, members)
	_, err = xredis.NewSimpleXesRedis(this.ctx, constant.Order_Redis_Cluster).Expire(setKey, []interface{}{}, constant.OrderId_ExpiredTimeDefault)
	return
}

/**
 * @description:
 * @params {OrderIds}
 * @return:
 */
func (this *OrderId) SDelV2(OrderIds []string) (back int64, err error) {
	setKey := this.GenSetKeyV2()
	if len(OrderIds) == 0 {
		//删除所有的数据
		return xredis.NewSimpleXesRedis(this.ctx, constant.Order_Redis_Cluster).Del(setKey, nil)
	}
	var params []interface{}
	for i := 0; i < len(OrderIds); i++ {
		params = append(params, OrderIds[i])
	}
	return xredis.NewSimpleXesRedis(this.ctx, constant.Order_Redis_Cluster).SRem(setKey, []interface{}{}, params)
}
