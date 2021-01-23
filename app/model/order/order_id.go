/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-12 22:32:33
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-02-23 01:52:02
 * @Description:
 */

package order

import (
	"context"
	"powerorder/app/constant"
	"powerorder/app/dao/mysql"
	"powerorder/app/dao/redis"
)

type OrderId struct {
	uAppId  uint
	uUserId uint64
	ctx     context.Context
}

/**
 * @description:初始化
 * @params {uAppId}
 * @params {uUserId}
 * @params {ctx}
 * @return:
 */
func NewOrderIdModel(uAppId uint, uUserId uint64, ctx context.Context) *OrderId {
	object := new(OrderId)
	object.uAppId = uAppId
	object.uUserId = uUserId
	object.ctx = ctx
	return object
}

/**
 * @description: 根据userId和appId，查询所有订单
 * @params {}
 * @return: {OrderIds} 订单ID集合
 * @return {OrderInfo} 订单OrderInfo集合
 */
func (this *OrderId) Get() (OrderIds []string, OrderInfo []mysql.OrderInfo, err error) {
	OrderIdRedisDao := redis.NewOrderIdDao(this.ctx, this.uAppId, this.uUserId, constant.Order_Redis_Cluster)
	value, err := OrderIdRedisDao.SGet()
	if err != nil || len(value) == 0 {
		// 缓存无数据，可能发生了错误，需要读取mysql
		OrderInfo, err = this.GetOrderInfo("", "")
		if err != nil {
			return
		}
		for i := 0; i < len(OrderInfo); i++ {
			OrderIds = append(OrderIds, OrderInfo[i].OrderId)
		}
		if len(OrderIds) == 0 {
			OrderIds = []string{constant.OrderId_MemberDefault}
		}
		_, err := OrderIdRedisDao.SAdd(OrderIds)
		if err != nil {
			return
		}
		return
	}
	for i := 0; i < len(value); i++ {
		if value[i] != constant.OrderId_MemberDefault {
			OrderIds = append(OrderIds, value[i])
		}
	}
	return

}

/**
 * @description: 根据userId和appId，查询所有订单
 * @params {}
 * @return: {OrderIds} 订单ID集合
 * @return {OrderInfo} 订单OrderInfo集合
 */
func (this *OrderId) GetV2(startDate, endDate string) (OrderIds []string, OrderInfo []mysql.OrderInfo, err error) {
	var value []string
	OrderIdRedisDao := redis.NewOrderIdDao(this.ctx, this.uAppId, this.uUserId, constant.Order_Redis_Cluster)

	value, err = OrderIdRedisDao.SGetV2()

	if err != nil || len(value) == 0 {
		// 缓存无数据，可能发生了错误，需要读取mysql
		OrderInfo, err = this.GetOrderInfo(startDate, endDate)
		if err != nil {
			return
		}
		for i := 0; i < len(OrderInfo); i++ {
			OrderIds = append(OrderIds, OrderInfo[i].OrderId)
		}
		if len(OrderIds) == 0 {
			OrderIds = []string{constant.OrderId_MemberDefault}
		}
		_, err := OrderIdRedisDao.SAddV2(OrderIds)
		if err != nil {
			return
		}
		return
	}

	for i := 0; i < len(value); i++ {
		if value[i] != constant.OrderId_MemberDefault {
			OrderIds = append(OrderIds, value[i])
		}
	}
	return

}

/**
 * @description:根据userId和appID，读取mysql,
 * @params {}
 * @return: OrderInfo集合
 */
func (this *OrderId) GetOrderInfo(startDate, endDate string) (OrderInfo []mysql.OrderInfo, err error) {
	OrderInfoDao := mysql.NewOrderInfoDao(this.ctx, this.uAppId, this.uUserId, constant.DBReader)
	OrderInfo, err = OrderInfoDao.Get([]string{}, "", []uint{constant.TxStatusCommitted}, startDate, endDate)
	return
}
