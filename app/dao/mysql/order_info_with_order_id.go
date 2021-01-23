/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-02-11 23:56:28
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-02-23 01:50:36
 * @Description:
 */

package mysql

import (
	"context"
	logger "github.com/tal-tech/loggerX"
	"github.com/tal-tech/torm"
	"powerorder/app/constant"
	"powerorder/app/utils"
)

type OrderInfoWithOrderIdDao struct {
	Dao
	strOrderId string
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func NewOrderInfoWithOrderIdDao(ctx context.Context, uAppId uint, strOrderId string, writer string) *OrderInfoWithOrderIdDao {
	object := new(OrderInfoWithOrderIdDao)
	object.InitDatabaseName(uAppId, strOrderId)
	writer = utils.DbShadowHandler(ctx, writer)
	if ins := torm.GetDbInstance(object.strDatabaseName, writer); ins != nil {
		object.UpdateEngine(ins.Engine)
	} else {
		return nil
	}
	object.UpdateEngine(writer)
	object.Engine.ShowSQL(true)
	object.InitTableName()
	object.SetTable(object.GetTable())
	object.Dao.InitCtx(ctx)

	return object
}

/**
 * @description:批量获得OrderInfo，同库同表的一起查询
 * @params {uAppId 业务线ID}
 * @params {arrOrderIds 订单ID数组}
 * @return:
 */
func BGet(ctx context.Context, uAppId uint, arrOrderIds []string) (info []OrderInfo, err error) {

	info = make([]OrderInfo, 0)
	err = nil
	OrderIdsMapped := make(map[string]map[string][]string)
	for i := 0; i < len(arrOrderIds); i++ {
		strDbName := utils.GenDatabaseNameByOrderId(arrOrderIds[i])
		strTableName := utils.GenOrderInfoTableNameByOrderId(arrOrderIds[i])

		if _, ok := OrderIdsMapped[strDbName]; !ok {
			OrderIdsMapped[strDbName] = make(map[string][]string)
		}

		if _, ok := OrderIdsMapped[strDbName][strTableName]; !ok {
			OrderIdsMapped[strDbName][strTableName] = make([]string, 0)
		}
		OrderIdsMapped[strDbName][strTableName] = append(OrderIdsMapped[strDbName][strTableName], arrOrderIds[i])
	}

	for _, OrderIdsInSameDb := range OrderIdsMapped {
		for _, OrderIdsInSameTable := range OrderIdsInSameDb {
			object := NewOrderInfoWithOrderIdDao(ctx, uAppId, OrderIdsInSameTable[0], constant.DBReader)
			ret, err := object.Get(OrderIdsInSameTable)
			if err != nil {
				return []OrderInfo{}, err
			}
			for _, Info := range ret {
				info = append(info, Info)
			}
		}
	}

	return
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *OrderInfoWithOrderIdDao) InitDatabaseName(uAppId uint, strOrderId string) {
	this.uAppId = uAppId
	this.strOrderId = strOrderId
	this.strDatabaseName = utils.GenDatabaseNameByOrderId(strOrderId)
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *OrderInfoWithOrderIdDao) InitTableName() {
	this.strTableName = utils.GenOrderInfoTableNameByOrderId(this.strOrderId)
}

/**
 * @description:
 * @params {type}
 * @return:
 */

func (this *OrderInfoWithOrderIdDao) Get(arrOrderIds []string) (ret []OrderInfo, err error) {
	ret = make([]OrderInfo, 0)
	this.InitSession()
	this.InitTableName()
	this.SetTable(this.strTableName)

	if len(arrOrderIds) > 0 {
		this.BuildQuery(torm.CastToParamIn(arrOrderIds), "order_id")
	} else {
		logger.Ex(this.ctx, "dao.mysql.orderinfowithorderid.get error", "OrderInfoWithOrderIdDao empty order_ids", "")
		return
	}
	this.BuildQuery(this.uAppId, "app_id")

	this.BuildQuery(constant.TxStatusCommitted, "tx_status")

	err = this.Session.Find(&ret)

	return
}
