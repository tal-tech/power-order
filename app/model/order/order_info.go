/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-17 16:58:44
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-02-23 01:54:01
 * @Description:
 */
package order

import (
	"context"
	logger "github.com/tal-tech/loggerX"
	"powerorder/app/constant"
	"powerorder/app/dao/mysql"
	"powerorder/app/utils"
)

type OrderInfo struct {
	uAppId  uint
	uUserId uint64
	ctx     context.Context
}

/**
 * @description:初始化赋值
 * @params {ctx}
 * @params {uAppId}
 * @params {uUserId}
 * @return:
 */
func NewOrderInfo(ctx context.Context, uAppId uint, uUserId uint64) *OrderInfo {
	object := new(OrderInfo)
	object.uAppId = uAppId
	object.uUserId = uUserId
	object.ctx = ctx
	return object
}

/**
 * @description: 组装返回数据，根据redis缓存的订单信息【info】
 * @params {OrderInfo} 可能已拿到的OrderInfo，概率有点小
 * @params {OrderIds} 想要获得OrderInfo的OrderId
 * @params {OrderCached} 可能已从cache中读取到的cache
 * @return: 返回OrderInfo集合
 * @return: 返回需要缓存Info数据
 */
func (this *OrderInfo) Get(OrderInfo []mysql.OrderInfo, OrderIds []string, OrderCached map[string]map[string]interface{}) (ret map[string]mysql.OrderInfo, PipeData map[string]map[string]interface{}, err error) {

	if this == nil {
		err = errors.New("system recover error")
		return
	}

	ret = make(map[string]mysql.OrderInfo)
	PipeData = make(map[string]map[string]interface{})
	var OrderIdsWithoutOrderInfo []string
	//遍历要获得的OrderId

	for i := 0; i < len(OrderIds); i++ {
		PipeData[OrderIds[i]] = make(map[string]interface{})
		//如果某个OrderId在Cache中存在

		if Order, ok := OrderCached[OrderIds[i]]; ok {

			var Value interface{}
			if Value, ok = Order[constant.Order_HashSubKey_Info]; ok {
				var Info mysql.OrderInfo
				if Info, ok = Value.(mysql.OrderInfo); ok {
					ret[OrderIds[i]] = Info
				}
				continue
			}
			//且在Cache中的类型为mysql.OrderInfo

		}
		//缓存不存在，则需要种缓存
		exist := false
		for j := 0; j < len(OrderInfo); j++ {

			if OrderInfo[j].OrderId == OrderIds[i] {
				//设置在返回值中
				ret[OrderIds[i]] = OrderInfo[j]
				//种缓存

				PipeData[OrderIds[i]][constant.Info] = OrderInfo[j]
				exist = true
				break
			}
		}

		if exist == false {
			OrderIdsWithoutOrderInfo = append(OrderIdsWithoutOrderInfo, OrderIds[i])
			PipeData[OrderIds[i]][constant.Order_HashSubKey_Info] = constant.Order_HashFieldDefault
		}
	}

	if len(OrderIdsWithoutOrderInfo) == 0 {
		return
	}
	// 缓存没命中，查询只读库mysql
	var OrderInfoFromMysql []mysql.OrderInfo
	if len(OrderIdsWithoutOrderInfo) > utils.SQL_WHERE_IN_MAX {
		OrderInfoFromMysql = this.getBatchByOrderIds(OrderIdsWithoutOrderInfo)
	} else {
		OrderInfoDao := mysql.NewOrderInfoDao(this.ctx, this.uAppId, this.uUserId, constant.DBReader)
		OrderInfoFromMysql, err = OrderInfoDao.Get(OrderIdsWithoutOrderInfo, "", []uint{constant.TxStatusCommitted}, "", "")
	}

	if err != nil {
		logger.Ex(this.ctx, "ModelOrderInfo", "get from db error =%v", err)
		return
	}

	for i := 0; i < len(OrderInfoFromMysql); i++ {
		ret[OrderInfoFromMysql[i].OrderId] = OrderInfoFromMysql[i]
		//种缓存
		PipeData[OrderInfoFromMysql[i].OrderId][constant.Order_HashSubKey_Info] = OrderInfoFromMysql[i]

	}

	return

}

func (this *OrderInfo) getBatchByOrderIds(OrderIds []string) (datas []mysql.OrderInfo) {
	defer utils.Catch()
	datas = make([]mysql.OrderInfo, 0)
	tmpOrderIdArr := make([]string, 0)
	for i := 0; i < len(OrderIds); i++ {
		tmpOrderIdArr = append(tmpOrderIdArr, OrderIds[i])
		if len(tmpOrderIdArr) == utils.SQL_WHERE_IN_MAX {
			OrderInfoDao := mysql.NewOrderInfoDao(this.ctx, this.uAppId, this.uUserId, constant.DBReader)
			OrderInfoFromMysql, err := OrderInfoDao.Get(tmpOrderIdArr, "", []uint{constant.TxStatusCommitted}, "", "")
			if err != nil {
				logger.Ex(this.ctx, "model.toc.orderinfo.getBatchByOrderIds", "getBatchByOrderIds error:%v", err)
				continue
			}
			if len(OrderInfoFromMysql) > 0 {
				datas = append(datas, OrderInfoFromMysql...)
			}
			tmpOrderIdArr = make([]string, 0)
		}
	}
	if len(tmpOrderIdArr) > 0 {
		OrderInfoDao := mysql.NewOrderInfoDao(this.ctx, this.uAppId, this.uUserId, constant.DBReader)
		OrderInfoFromMysql, err := OrderInfoDao.Get(tmpOrderIdArr, "", []uint{constant.TxStatusCommitted}, "", "")
		if err != nil {
			logger.Ex(this.ctx, "model.toc.orderinfo.getBatchByOrderIds", "getBatchByOrderIds error:%v", err)
		}
		if len(OrderInfoFromMysql) > 0 {
			datas = append(datas, OrderInfoFromMysql...)
		}
	}
	return
}
