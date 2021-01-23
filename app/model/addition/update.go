/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-07 18:04:10
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-04-11 00:41:47
 * @Description:
 * @Todo : 更新时清除缓存、更新缓存应该改成事务
 */
package addition

import (
	"context"
	logger "github.com/tal-tech/loggerX"
	"powerorder/app/constant"
	"powerorder/app/dao/mysql"
	"powerorder/app/dao/redis"
	"powerorder/app/model/finishing"
	"powerorder/app/params"
)

type Update struct {
	uAppId uint
	ctx    context.Context
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func NewUpdate(ctx context.Context, uAppId uint) *Update {
	object := new(Update)
	object.ctx = ctx
	object.uAppId = uAppId
	return object
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *Update) UpdateStatus(param params.ReqRollback) (err error) {
	instance := mysql.NewOrderInfoDao(this.ctx, param.AppId, param.UserId, constant.DBWriter)
	if instance == nil {
		err = logger.NewError("init dao error!")
		return
	}
	var info mysql.OrderInfo
	info.TxStatus = param.TxStatus
	instance.InitTableName()
	instance.SetTable(instance.GetTable())

	ret, err := instance.UpdateColsWhere(info, []string{"tx_status"}, " user_id = ? and app_id = ? and tx_id = ? and tx_status = ? ", []interface{}{param.UserId, param.AppId, param.TxId, constant.TxStatusUnCommitted})
	if err != nil || ret == 0 {
		logger.Ix(this.ctx, "model.addition.update.update failed", "model addition/update error = %+v ,ret = %d", err, ret)
		return
	}

	if param.TxStatus != constant.TxStatusCommitted {
		return
	}

	this.setCache(param)
	return
}

func (this *Update) setCache(param params.ReqRollback) {
	finishingModel := finishing.NewInsertion(param.AppId, param.UserId, this.ctx, param.TxId)
	TxOrderDao := redis.NewTxOrderDao(this.ctx, param.AppId, param.TxId, constant.Order_Redis_Cluster)
	TxInfo, _ := TxOrderDao.Get()
	ok := finishingModel.SetCache(make([]string, 0), TxInfo)
	if !ok {
		logger.Ex(this.ctx, "model.addition.update.update failed", "update.update.SetCache", "ok:%+v", ok)
	}
	go this.insertOrderIdx(param, TxInfo)
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *Update) insertOrderIdx(param params.ReqRollback, txInfo redis.TxInfo) {

	finishingModel := finishing.NewInsertion(param.AppId, param.UserId, this.ctx, param.TxId)
	OrderIds, _ := finishingModel.InsertOrder(txInfo)

	orderIdDao := redis.NewOrderIdDao(this.ctx, param.AppId, param.UserId, constant.Order_Redis_Cluster)
	if _, err := orderIdDao.SDel([]string{}); err != nil {
		logger.Ex(this.ctx, "model.addition.update.insertOrderIdx", "sdel order_ids  ", "%+v", OrderIds)
	}
	// todo  去掉 | 加注释
	if _, err := orderIdDao.SDelV2([]string{}); err != nil {
		logger.Ex(this.ctx, "model.addition.update.insertOrderIdx", "sdelv2 order_ids  ", "%+v", OrderIds)
	}
}
