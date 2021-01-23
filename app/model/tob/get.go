/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-23 13:13:33
 * @LastEditors  : lichanglin@tal.com
 * @LastEditTime : 2020-02-12 10:29:53
 * @Description:
 */

package tob

import (
	"context"
	logger "github.com/tal-tech/loggerX"
	"powerorder/app/constant"
	"powerorder/app/dao/mysql"
	"powerorder/app/dao/redis"
	"powerorder/app/model/order"
	"powerorder/app/output"
	"strconv"
)

type Get struct {
	uAppId uint
	ctx    context.Context
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func NewGet(ctx context.Context, uAppId uint) *Get {
	object := new(Get)
	object.uAppId = uAppId
	object.ctx = ctx
	return object
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *Get) Get(OrderIds []string, Fields []string) (ret map[string]output.Order, err error) {
	err = nil

	var OrderInfo []mysql.OrderInfo
	var InfoToBeCached = make(map[string]map[string]interface{}, 0)
	var InfoCached map[string]map[string]interface{}

	if len(OrderIds) == 0 {
		err = logger.NewError("empty order_ids")
		logger.Ex(this.ctx, "model.tob.get error", "empty order_ids")
		return
	}
	OrderRedisDao := redis.NewOrderDao(this.ctx, this.uAppId, 0)
	InfoCached, err = OrderRedisDao.BatchGet(OrderIds)

	OrderIdsCached := make(map[uint64][]string, 0)
	OrderIdsWithoutUserId := make([]string, 0)

	//logger.Ix(this.ctx,"model.tob.get.get","orderIds: %v, infoCached:= %d", OrderIds, len(InfoCached))
	for _, orderId := range OrderIds {
		if value, ok := InfoCached[orderId]; ok && len(value) > 0 {
			if UserIdStr, ok := value[constant.Order_HashSubKey_UserId].(string); ok {
				UserIdTmp, _ := strconv.ParseInt(UserIdStr, 10, 64)
				UserId := uint64(UserIdTmp)
				if _, ok := OrderIdsCached[UserId]; !ok {
					OrderIdsCached[UserId] = make([]string, 0)
				}
				OrderIdsCached[UserId] = append(OrderIdsCached[UserId], orderId)
			}
		} else {
			OrderIdsWithoutUserId = append(OrderIdsWithoutUserId, orderId)
		}
	}

	//logger.Dx(this.ctx,"model tob/get order_ids = %v, OrderIdsCached= %v, OrderIdsWithoutUserId= %v", OrderIds, OrderIdsCached, OrderIdsWithoutUserId)

	if len(OrderIdsWithoutUserId) != 0 {
		OrderInfo, err = mysql.BGet(this.ctx, this.uAppId, OrderIdsWithoutUserId)
		if err != nil {
			return
		}
		for _, info := range OrderInfo {
			UserId := info.UserId
			OrderId := info.OrderId
			if _, ok := OrderIdsCached[UserId]; !ok {
				OrderIdsCached[UserId] = make([]string, 0)
			}
			OrderIdsCached[UserId] = append(OrderIdsCached[UserId], OrderId)
		}
	}

	ret = make(map[string]output.Order, 0)

	var OrderInfoMapped map[string]mysql.OrderInfo
	var OrderDetailMapped map[string][]mysql.OrderDetail

	var PipeData map[string]map[string]interface{}
	for i := 0; i < len(Fields); i++ {
		if Fields[i] == constant.OrderInfo {

			for UserId, OrderIds := range OrderIdsCached {
				ModelOrderInfo := order.NewOrderInfo(this.ctx, this.uAppId, UserId)
				OrderInfoMapped, PipeData, _ = ModelOrderInfo.Get(OrderInfo, OrderIds, InfoCached)

				for OrderId, Pipe := range PipeData {

					for field, value := range Pipe {
						if _, ok := InfoToBeCached[OrderId]; !ok {
							InfoToBeCached[OrderId] = make(map[string]interface{}, 0)
						}
						InfoToBeCached[OrderId][field] = value
					}
				}

				for OrderId, OrderInfo := range OrderInfoMapped {
					ret[OrderId] = output.Order{OrderInfo, make([]mysql.OrderDetail, 0), make(map[string][]map[string]interface{})}
				}

			}
			break

		}
	}

	for i := 0; i < len(Fields); i++ {

		if Fields[i] == constant.OrderDetail {
			for UserId, OrderIds := range OrderIdsCached {

				ModelOrderDetail := order.NewOrderDetail(this.ctx, this.uAppId, UserId)
				OrderDetailMapped, PipeData, _ = ModelOrderDetail.Get(OrderIds, InfoCached)

				for OrderId, Pipe := range PipeData {
					for field, value := range Pipe {
						if _, ok := InfoToBeCached[OrderId]; !ok {
							InfoToBeCached[OrderId] = make(map[string]interface{}, 0)
						}
						InfoToBeCached[OrderId][field] = value
					}
				}

				for OrderId, OrderDetails := range OrderDetailMapped {
					if _, ok := ret[OrderId]; !ok {
						ret[OrderId] = output.Order{mysql.OrderInfo{}, OrderDetails, make(map[string][]map[string]interface{})}
					} else {
						ret[OrderId] = output.Order{ret[OrderId].Info, OrderDetails, make(map[string][]map[string]interface{})}
					}
				}
			}
			break

		}
	}

	for i := 0; i < len(Fields); i++ {
		if Fields[i] == constant.OrderDetail || Fields[i] == constant.OrderInfo {

		} else {
			for UserId, OrderIds := range OrderIdsCached {

				ExtensionModel := order.NewExtensionModel(this.ctx, this.uAppId, UserId, Fields[i])
				ExtDataMapped, PipeData, _ := ExtensionModel.Get(OrderIds, InfoCached)

				for OrderId, Pipe := range PipeData {
					for field, value := range Pipe {
						if _, ok := InfoToBeCached[OrderId]; !ok {
							InfoToBeCached[OrderId] = make(map[string]interface{}, 0)
						}
						InfoToBeCached[OrderId][field] = value
					}
					//OrderRedisDao.HMSet(OrderId, Pipe)
				}
				for OrderId, ExtData := range ExtDataMapped {
					if _, ok := ret[OrderId]; !ok {
						ret[OrderId] = output.Order{mysql.OrderInfo{}, make([]mysql.OrderDetail, 0), make(map[string][]map[string]interface{})}
					}
					ret[OrderId].Extensions[Fields[i]] = ExtData
				}
			}
		}

	}

	for OrderId, Value := range InfoToBeCached {
		err := OrderRedisDao.HMSet(OrderId, Value)
		if err != nil {
			logger.Ex(this.ctx, "model.tob.get error", "hmset error", err)
			continue
		}
	}
	return
}
