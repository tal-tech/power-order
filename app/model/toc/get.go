/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-12 22:32:33
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime : 2020-02-14 19:35:27
 * @Description:
 */

package toc

import (
	"context"
	logger "github.com/tal-tech/loggerX"
	"powerorder/app/constant"
	"powerorder/app/dao/mysql"
	"powerorder/app/dao/redis"
	"powerorder/app/model/order"
	"powerorder/app/output"
)

type Get struct {
	uAppId  uint
	uUserId uint64
	ctx     context.Context
}

/**
 * @description:初始化赋值
 * @params {ctx} 上下文
 * @params {uAppId} 业务线Id
 * @params {uUserId} 用户Id
 * @return:
 */
func NewGet(ctx context.Context, uAppId uint, uUserId uint64) *Get {
	object := new(Get)
	object.uAppId = uAppId
	object.uUserId = uUserId
	object.ctx = ctx
	return object
}

/**
 * @description: 批量查询订单 【toc/search/(get|search) 公用】
 * @params {OrderIds} 订单ID
 * @params {Fields} 查询的字段
 * @return: 订单详情列表【info detail extensions】
 */
func (this *Get) Get(OrderIds []string, Fields []string) (ret map[string]output.Order, Oids []string, err error) {
	var InfoToBeCached = make(map[string]map[string]interface{}, 0)
	var InfoCached map[string]map[string]interface{}

	OrderInfo := make([]mysql.OrderInfo, 0)
	if len(OrderIds) == 0 {
		_tmpOrderIds, err := this.getOrderIds()
		if err != nil {
			return
		}
		OrderIds = _tmpOrderIds
	}
	Oids = OrderIds

	if err != nil {
		logger.Ex(this.ctx, "model.toc.get.Get error", "model toc/get batchget orderids error = %+v", err)
	}
	ret = make(map[string]output.Order, 0)

	OrderRedisDao := redis.NewOrderDao(this.ctx, this.uAppId, this.uUserId)

	OrderInfoMapped := make(map[string]mysql.OrderInfo, 0)
	OrderDetailMapped := make(map[string][]mysql.OrderDetail, 0)
	var PipeData map[string]map[string]interface{}
	for i := 0; i < len(Fields); i++ {
		if Fields[i] == constant.OrderInfo {

			ModelOrderInfo := order.NewOrderInfo(this.ctx, this.uAppId, this.uUserId)
			OrderInfoMapped, PipeData, _ = ModelOrderInfo.Get(OrderInfo, OrderIds, InfoCached)

			for OrderId, Pipe := range PipeData {

				for field, value := range Pipe {
					if _, ok := InfoToBeCached[OrderId]; !ok {
						InfoToBeCached[OrderId] = make(map[string]interface{}, 0)
					}
					InfoToBeCached[OrderId][field] = value
				}
			}
			OrderIds = make([]string, 0)
			for OrderId, OrderInfo := range OrderInfoMapped {
				OrderIds = append(OrderIds, OrderId)
				ret[OrderId] = output.Order{OrderInfo, make([]mysql.OrderDetail, 0), make(map[string][]map[string]interface{})}
			}
			break
		}
	}

	for i := 0; i < len(Fields); i++ {
		if Fields[i] == constant.OrderDetail {
			ModelOrderDetail := order.NewOrderDetail(this.ctx, this.uAppId, this.uUserId)
			OrderDetailMapped, PipeData, _ = ModelOrderDetail.Get(OrderIds, InfoCached)

			for OrderId, Pipe := range PipeData {
				for field, value := range Pipe {
					if _, ok := InfoToBeCached[OrderId]; !ok {
						InfoToBeCached[OrderId] = make(map[string]interface{}, 0)
					}
					InfoToBeCached[OrderId][field] = value
				}
			}
			OrderIds = make([]string, 0)

			for OrderId, OrderDetails := range OrderDetailMapped {
				OrderIds = append(OrderIds, OrderId)
				if _, ok := ret[OrderId]; !ok {
					ret[OrderId] = output.Order{mysql.OrderInfo{}, OrderDetails, make(map[string][]map[string]interface{})}
				} else {
					ret[OrderId] = output.Order{ret[OrderId].Info, OrderDetails, make(map[string][]map[string]interface{})}
				}
			}
			break
		}
	}

	for i := 0; i < len(Fields); i++ {
		if Fields[i] != constant.Detail && Fields[i] != constant.Info {

			ExtensionModel := order.NewExtensionModel(this.ctx, this.uAppId, this.uUserId, Fields[i])
			ExtDataMapped, PipeData, _ := ExtensionModel.Get(OrderIds, InfoCached)

			for OrderId, Pipe := range PipeData {
				for field, value := range Pipe {
					if _, ok := InfoToBeCached[OrderId]; !ok {
						InfoToBeCached[OrderId] = make(map[string]interface{}, 0)
					}
					InfoToBeCached[OrderId][field] = value
				}
			}
			for OrderId, ExtData := range ExtDataMapped {
				if _, ok := ret[OrderId]; !ok {
					ret[OrderId] = output.Order{mysql.OrderInfo{}, make([]mysql.OrderDetail, 0), make(map[string][]map[string]interface{})}
				}
				ret[OrderId].Extensions[Fields[i]] = ExtData
			}
		}

	}
	for OrderId, Value := range InfoToBeCached {
		if err := OrderRedisDao.HMSet(OrderId, Value); err != nil {

		}
	}
	return

}

func (this *Get) getOrderIds() (OrderIds []string, err error) {
	OrderIdModel := order.NewOrderIdModel(this.uAppId, this.uUserId, this.ctx)
	OrderIds, _, err = OrderIdModel.Get()
	if err != nil {
		return
	}
	return
}

func (this *Get) getOrdersCache(orderIds []string) (infos map[string]map[string]interface{}, err error) {
	OrderRedisDao := redis.NewOrderDao(this.ctx, this.uAppId, this.uUserId)
	infos, err = OrderRedisDao.BatchGet(orderIds)
	return
}

/**
 * @description: 批量查询订单 【toc/search/(get|search) 公用】
 * @param {OrderIds} 订单ID
 * @param {Fields} 查询的字段
 * @return: 订单详情列表【info detail extensions】
 */
func (this *Get) GetV2(OrderIds []string, Fields []string, startDate, endDate string) (ret map[string]output.Order, Oids []string, err error) {
	InfoToBeCached := make(map[string]map[string]interface{}, 0)
	InfoCached := make(map[string]map[string]interface{}, 0)
	OrderInfo := make([]mysql.OrderInfo, 0)

	if len(OrderIds) == 0 {
		OrderIdModel := order.NewOrderIdModel(this.uAppId, this.uUserId, this.ctx)
		OrderIds, OrderInfo, err = OrderIdModel.GetV2(startDate, endDate)
		if err != nil {
			return
		}
		if len(OrderIds) == 0 {
			return
		}
	}
	Oids = OrderIds

	OrderRedisDao := redis.NewOrderDao(this.ctx, this.uAppId, this.uUserId)
	InfoCached, err = OrderRedisDao.BatchGet(OrderIds)
	ret = make(map[string]output.Order, 0)

	var OrderInfoMapped map[string]mysql.OrderInfo
	var OrderDetailMapped map[string][]mysql.OrderDetail

	var PipeData map[string]map[string]interface{}
	for i := 0; i < len(Fields); i++ {
		if Fields[i] == constant.OrderInfo {

			ModelOrderInfo := order.NewOrderInfo(this.ctx, this.uAppId, this.uUserId)
			OrderInfoMapped, PipeData, _ = ModelOrderInfo.Get(OrderInfo, OrderIds, InfoCached)

			for OrderId, Pipe := range PipeData {

				for field, value := range Pipe {
					if _, ok := InfoToBeCached[OrderId]; !ok {
						InfoToBeCached[OrderId] = make(map[string]interface{}, 0)
					}
					InfoToBeCached[OrderId][field] = value
				}
			}
			OrderIds = make([]string, 0)
			for OrderId, OrderInfo := range OrderInfoMapped {
				OrderIds = append(OrderIds, OrderId)
				ret[OrderId] = output.Order{OrderInfo, make([]mysql.OrderDetail, 0), make(map[string][]map[string]interface{})}
			}
			break
		}
	}

	for i := 0; i < len(Fields); i++ {

		if Fields[i] == constant.OrderDetail {
			ModelOrderDetail := order.NewOrderDetail(this.ctx, this.uAppId, this.uUserId)
			OrderDetailMapped, PipeData, _ = ModelOrderDetail.Get(OrderIds, InfoCached)

			for OrderId, Pipe := range PipeData {
				for field, value := range Pipe {
					if _, ok := InfoToBeCached[OrderId]; !ok {
						InfoToBeCached[OrderId] = make(map[string]interface{}, 0)
					}
					InfoToBeCached[OrderId][field] = value
				}
			}
			OrderIds = make([]string, 0)

			for OrderId, OrderDetails := range OrderDetailMapped {
				OrderIds = append(OrderIds, OrderId)
				if _, ok := ret[OrderId]; !ok {
					ret[OrderId] = output.Order{mysql.OrderInfo{}, OrderDetails, make(map[string][]map[string]interface{})}
				} else {
					ret[OrderId] = output.Order{ret[OrderId].Info, OrderDetails, make(map[string][]map[string]interface{})}
				}
			}
			break
		}
	}

	for i := 0; i < len(Fields); i++ {
		if Fields[i] != constant.Detail && Fields[i] != constant.Info {

			ExtensionModel := order.NewExtensionModel(this.ctx, this.uAppId, this.uUserId, Fields[i])
			ExtDataMapped, PipeData, _ := ExtensionModel.Get(OrderIds, InfoCached)

			for OrderId, Pipe := range PipeData {
				for field, value := range Pipe {
					if _, ok := InfoToBeCached[OrderId]; !ok {
						InfoToBeCached[OrderId] = make(map[string]interface{}, 0)
					}
					InfoToBeCached[OrderId][field] = value
				}
			}
			for OrderId, ExtData := range ExtDataMapped {
				if _, ok := ret[OrderId]; !ok {
					ret[OrderId] = output.Order{mysql.OrderInfo{}, make([]mysql.OrderDetail, 0), make(map[string][]map[string]interface{})}
				}
				ret[OrderId].Extensions[Fields[i]] = ExtData
			}
		}

	}
	for OrderId, Value := range InfoToBeCached {
		if err := OrderRedisDao.HMSet(OrderId, Value); err != nil {

		}

	}
	return

}
