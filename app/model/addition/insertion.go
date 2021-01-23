/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-07 18:04:10
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-04-10 16:33:26
 * @Description:
 */
package addition

import (
	"context"
	"fmt"
	logger "github.com/tal-tech/loggerX"
	"powerorder/app/constant"
	"powerorder/app/dao/mysql"
	"powerorder/app/dao/redis"
	imysql "powerorder/app/i/dao/mysql"
	"powerorder/app/model/finishing"
	"powerorder/app/output"
	"powerorder/app/params"
	"powerorder/app/utils"
	"time"
)

type Insertion struct {
	action uint8
	ctx    context.Context
}

func NewInsertion(action uint8, ctx context.Context) *Insertion {
	object := new(Insertion)
	object.action = action
	object.ctx = ctx
	return object
}

/**
 * @description:
 * @params {params}
 * @return:
 */
func (this *Insertion) Insert(param params.ReqBegin, TxId string) (OrderIds []string, err error) {
	var txInfo redis.TxInfo
	txInfo.OrderIds = make([]string, 0)
	txInfo.Extensions = make(map[string]bool)
	txInfo.OrderInfos = make(map[string]output.Order, 0)

	CreatedTime := utils.Time(time.Now())
	//data := make(map[string]map[string]interface{})

	infos := make([]mysql.OrderInfo, 0)
	details := make([]mysql.OrderDetail, 0)
	extensions := make(map[string][]map[string]interface{})

	for j := 0; j < len(param.Additions); j++ {
		var orderData output.Order

		addition := param.Additions[j]

		var OrderId string
		if len(addition.Info.OrderId) == 0 {
			OrderId = utils.GenOrderId(param.AppId, time.Time(CreatedTime), param.UserId)
		} else {
			OrderId = addition.Info.OrderId
		}

		OrderIds = append(OrderIds, OrderId)
		//if data[OrderId] == nil {
		//	data[OrderId] = make(map[string]interface{})
		//}

		orderData.Detail = this.assembleOrderDetail(addition.Detail, OrderId, param.UserId, param.AppId)
		details = append(details, orderData.Detail...)

		orderData.Extensions = this.assembleExtension(addition.Extensions, &txInfo, OrderId, param.UserId)
		for k, val := range orderData.Extensions {
			if _, ok := extensions[k]; ok == false {
				extensions[k] = []map[string]interface{}{}
			}
			extensions[k] = append(extensions[k], val...)
		}

		info := this.assembleOrderInfo(addition.Info)
		info.UserId = param.UserId
		info.AppId = param.AppId
		info.TxId = TxId

		txInfo.OrderIds = append(txInfo.OrderIds, OrderId)
		infos = append(infos, info)
		orderData.Info = info

		txInfo.OrderInfos[OrderId] = orderData
	}

	session := mysql.NewSession(this.ctx, param.AppId, param.UserId, constant.DBWriter)
	if session == nil {
		err = logger.NewError("init session error")
		return
	}
	defer session.Close()

	if err1 := session.Begin(); err1 != nil {
		err = err1
		return
	}

	var dao imysql.Dao

	for key, value := range extensions {
		dao = mysql.NewExtensionDao2(param.AppId, param.UserId, key, session, this.ctx)
		if _, err1 := imysql.Create(dao, value); err1 != nil {
			if err1 := session.Rollback(); err1 != nil {
				err = err1
			}
			return
		}
	}
	dao = mysql.NewOrderDetailDao2(param.AppId, param.UserId, session, this.ctx)
	if _, err1 := imysql.Create(dao, details); err1 != nil {
		if err1 = session.Rollback(); err1 != nil {
			err = err1
		}
		return
	}

	dao = mysql.NewOrderInfoDao2(param.AppId, param.UserId, session, this.ctx)
	if _, err1 := imysql.Create(dao, infos); err1 != nil {
		if err1 = session.Rollback(); err1 != nil {
			err = err1
		}
		return
	}
	if err1 := session.Commit(); err1 != nil {
		err = err1
		return
	}

	//	将数据临时写入到redis，减少对mysql的访问
	if this.action == constant.Addition {
		TxInfoDao := redis.NewTxOrderDao(this.ctx, param.AppId, TxId, constant.Order_Redis_Cluster)
		err = TxInfoDao.Set(txInfo)
		return
	}
	// 异步处理相关数据
	go this.finishing(param, TxId, txInfo)
	return
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *Insertion) finishing(param params.ReqBegin, txId string, txInfo redis.TxInfo) {
	// 数据同步到索引库中
	finishingModel := finishing.NewInsertion(param.AppId, param.UserId, this.ctx, txId)
	OrderIds, _ := finishingModel.InsertOrder(txInfo)

	// 清理redis数据
	orderIdDao := redis.NewOrderIdDao(this.ctx, param.AppId, param.UserId, constant.Order_Redis_Cluster)
	if _, err := orderIdDao.SDel([]string{}); err != nil {

	}
	logger.Ix(this.ctx, "model.addition.insertion.finishing", "insert finishing", "txId %s , params %+v , txInfo %+v", txId, param, OrderIds)

}

func (this *Insertion) assembleOrderInfo(reqInfo mysql.OrderInfo) (info mysql.OrderInfo) {
	var CreatedTime = utils.Time(time.Now())
	info = reqInfo
	info.BId = utils.GenId()
	if info.Source < 1 {
		info.Source = 1 // 默认值
	}
	if constant.ZeroTime == reqInfo.PaidTime.Unix() {
		info.PaidTime = mysql.TiDefaultTimestamp
	}
	if constant.ZeroTime == reqInfo.PayCreatedTime.Unix() {
		info.PayCreatedTime = mysql.TiDefaultTimestamp
	}
	if constant.ZeroTime == reqInfo.CancelledTime.Unix() {
		info.CancelledTime = mysql.TiDefaultTimestamp
	}
	if constant.ZeroTime == reqInfo.ExpiredTime.Unix() {
		info.ExpiredTime = mysql.TiDefaultTimestamp
	}
	if this.action == constant.Addition {
		info.CreatedTime = CreatedTime
		info.UpdatedTime = CreatedTime
		info.TxStatus = constant.TxStatusUnCommitted
	} else {
		if constant.ZeroTime == reqInfo.CreatedTime.Unix() {
			info.CreatedTime = mysql.TiDefaultTimestamp
		}
		if constant.ZeroTime == reqInfo.UpdatedTime.Unix() {
			info.UpdatedTime = mysql.TiDefaultTimestamp
		}
		info.TxStatus = constant.TxStatusCommitted
	}
	info.Version = 1

	return
}

func (this *Insertion) assembleOrderDetail(reqDetails []mysql.OrderDetail, orderId string, userId uint64, appId uint) []mysql.OrderDetail {
	CreatedTime := utils.Time(time.Now())
	tmpDetails := make([]mysql.OrderDetail, len(reqDetails))
	for i := 0; i < len(reqDetails); i++ {
		reqDetails[i].AppId = appId
		reqDetails[i].UserId = userId
		reqDetails[i].OrderId = orderId

		if reqDetails[i].ProductType < 1 {
			reqDetails[i].ProductType = 1 // 默认值
		}

		if this.action == constant.Addition {
			reqDetails[i].CreatedTime = CreatedTime //.Format("2006-01-02 15:04:05")
			reqDetails[i].UpdatedTime = CreatedTime
		} else {
			if constant.ZeroTime == reqDetails[i].CreatedTime.Unix() {
				reqDetails[i].CreatedTime = mysql.TiDefaultTimestamp
			}
			if constant.ZeroTime == reqDetails[i].UpdatedTime.Unix() {
				reqDetails[i].UpdatedTime = mysql.TiDefaultTimestamp
			}
		}
		reqDetails[i].Version = 1
	}
	return tmpDetails
}

// todo 注释
//
//
//
func (this *Insertion) assembleExtension(extensions map[string][]map[string]interface{}, txInfo *redis.TxInfo, orderId string, userId uint64) map[string][]map[string]interface{} {
	tmpExtensions := make(map[string][]map[string]interface{}, 0)
	for key, value := range extensions {
		(*txInfo).Extensions[key] = true
		for i := 0; i < len(value); i++ {
			value[i]["order_id"] = orderId
			value[i]["user_id"] = userId
			value[i]["version"] = 1
			value[i]["id"] = fmt.Sprintf("%d", utils.GenId())
			if _, ok := extensions[key]; ok == false {
				extensions[key] = []map[string]interface{}{}
				tmpExtensions[key] = []map[string]interface{}{}
			}
			extensions[key] = append(extensions[key], value[i])
			tmpExtensions[key] = append(tmpExtensions[key], value[i])
		}
	}
	return tmpExtensions
}
