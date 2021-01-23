/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-17 16:58:44
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-02-23 01:52:49
 * @Description:
 */
package order

import (
	"context"
	"github.com/spf13/cast"
	logger "github.com/tal-tech/loggerX"
	"powerorder/app/dao/mysql"
	"time"
)

type OrderDetail struct {
	uAppId  uint
	uUserId uint64
	ctx     context.Context
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func NewOrderDetail(ctx context.Context, uAppId uint, uUserId uint64) *OrderDetail {
	object := new(OrderDetail)
	object.uAppId = uAppId
	object.uUserId = uUserId
	object.ctx = ctx
	return object
}

/**
 * @description: 组装返回数据，根据redis缓存的订单信息【detail】
 * @params {OrderIds} 想要获得OrderInfo的OrderId
 * @params {OrderCached} 可能已从cache中读取到的cache
 * @return: 返回OrderDetail集合
 * @return: 返回需要种缓存OrderDetail集合
 */
func (this *OrderDetail) Get(OrderIds []string, OrderCached map[string]map[string]interface{}) (ret map[string][]mysql.OrderDetail, PipeData map[string]map[string]interface{}, err error) {

	ret = make(map[string][]mysql.OrderDetail)
	PipeData = make(map[string]map[string]interface{})
	var OrderIdsWithoutOrderDetail []string
	//遍历要获得的OrderId

	logger.Dx(this.ctx, "model.toc.orderdetail.Get", "model order_detail order_ids = %v ", OrderIds)

	for i := 0; i < len(OrderIds); i++ {
		PipeData[OrderIds[i]] = make(map[string]interface{})

		//如果某个OrderId在Cache中存在
		if Order, ok := OrderCached[OrderIds[i]]; ok {
			var Value interface{}

			if Value, ok = Order[sconstant.Detail]; ok {
				var Detail []mysql.OrderDetail
				//且在Cache中的类型为[]mysql.OrderDetail

				if Detail, ok = Value.([]mysql.OrderDetail); ok {
					ret[OrderIds[i]] = Detail
				}
				continue
			}

		}

		OrderIdsWithoutOrderDetail = append(OrderIdsWithoutOrderDetail, OrderIds[i])
		PipeData[OrderIds[i]][sconstant.Detail] = sconstant.Order_HashFieldDefault
	}
	if span != nil {
		span.Tag("t2", cast.ToString(time.Now().UnixNano()))
	}
	if len(OrderIdsWithoutOrderDetail) == 0 {
		return
	}
	// 超过200个订单 则分批次获取
	var OrderDetailFromMysql []mysql.OrderDetail
	if len(OrderIdsWithoutOrderDetail) > util.SQL_WHERE_IN_MAX {
		OrderDetailFromMysql = this.getBatchByOrderIds(OrderIdsWithoutOrderDetail)
	} else {
		OrderDetailDao := mysql.NewOrderDetailDao(this.ctx, this.uAppId, this.uUserId, sconstant.DBReader)

		OrderDetailFromMysql, err = OrderDetailDao.Get(OrderIdsWithoutOrderDetail)
	}

	if err != nil {
		return
	}
	if span != nil {
		span.Tag("t3", cast.ToString(time.Now().UnixNano()))
	}
	for i := 0; i < len(OrderDetailFromMysql); i++ {

		ret[OrderDetailFromMysql[i].OrderId] = append(ret[OrderDetailFromMysql[i].OrderId], OrderDetailFromMysql[i])
		//种缓存
		if Pipe, ok := PipeData[OrderDetailFromMysql[i].OrderId][sconstant.Detail]; ok {
			switch Pipe.(type) {
			case string:
				PipeData[OrderDetailFromMysql[i].OrderId][sconstant.Detail] = make([]mysql.OrderDetail, 0)
			}
		} else {
			PipeData[OrderDetailFromMysql[i].OrderId][sconstant.Detail] = make([]mysql.OrderDetail, 0)
		}

		PipeData[OrderDetailFromMysql[i].OrderId][sconstant.Detail] = append(PipeData[OrderDetailFromMysql[i].OrderId][sconstant.Order_HashSubKey_Detail].([]mysql.OrderDetail), OrderDetailFromMysql[i])

	}
	if span != nil {
		span.Tag("t4", cast.ToString(time.Now().UnixNano()))
	}
	return

}

func (this *OrderDetail) getBatchByOrderIds(OrderIds []string) (datas []mysql.OrderDetail) {
	defer util2.Catch()
	datas = make([]mysql.OrderDetail, 0)
	tmpOrderIdArr := make([]string, 0)
	for i := 0; i < len(OrderIds); i++ {
		tmpOrderIdArr = append(tmpOrderIdArr, OrderIds[i])
		if len(tmpOrderIdArr) == util.SQL_WHERE_IN_MAX {
			OrderDetailDao := mysql.NewOrderDetailDao(this.ctx, this.uAppId, this.uUserId, sconstant.DBReader)
			OrderDetailFromMysql, err := OrderDetailDao.Get(tmpOrderIdArr)
			if err != nil {
				logger.Ex(this.ctx, "model.toc.orderdetail.getBatchByOrderIds", "getBatchByOrderIds error:%v", err)
				continue
			}
			if len(OrderDetailFromMysql) > 0 {
				datas = append(datas, OrderDetailFromMysql...)
			}
			tmpOrderIdArr = make([]string, 0)
		}
	}
	if len(tmpOrderIdArr) > 0 {
		OrderDetailDao := mysql.NewOrderDetailDao(this.ctx, this.uAppId, this.uUserId, sconstant.DBReader)
		OrderDetailFromMysql, err := OrderDetailDao.Get(tmpOrderIdArr)
		if err != nil {
			logger.Ex(this.ctx, "model.toc.orderdetail.getBatchByOrderIds", "getBatchByOrderIds error:%v", err)
		}
		if len(OrderDetailFromMysql) > 0 {
			datas = append(datas, OrderDetailFromMysql...)
		}
	}
	return
}
