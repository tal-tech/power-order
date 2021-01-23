/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-04-09 16:11:58
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-04-13 02:05:08
 * @Description:
 */

package finishing

import (
	"context"
	logger "github.com/tal-tech/loggerX"
	"github.com/tal-tech/torm"
	"powerorder/app/constant"
	"powerorder/app/dao/mysql"
	"powerorder/app/dao/redis"
	imysql "powerorder/app/i/dao/mysql"
	"powerorder/app/utils"
)

type Insertion struct {
	action  uint8
	ctx     context.Context
	uAppId  uint
	uUserId uint64
	txId    string
}

func NewInsertion(uAppId uint, uUserId uint64, ctx context.Context, txId string) *Insertion {

	object := new(Insertion)
	object.uAppId = uAppId
	object.uUserId = uUserId
	object.ctx = ctx
	object.txId = txId
	return object
}

/**
 * @description:  根据orderids 将数据写入mysql表
 * @params {params}
 * @return:
 */
func (this *Insertion) InsertOrder(txInfo redis.TxInfo) (OrderIds []string, err error) {

	defer utils.Catch()

	session := mysql.NewSessionIdx(this.ctx, this.uAppId, this.uUserId, constant.DBWriter)
	if session == nil {
		logger.Ex(this.ctx, "model.finishing.insertion.InsertOrder error", "insertOrderInfo session err", "%+v", session)

		return
	}

	defer session.Close()

	data := make(map[string]map[string]interface{})
	Extensions := make([]string, 0)
	//logger.Ix(this.ctx,"model.finishing.insertion.InsertOrder","InsertOrderIdx start","txinfo %+v, OrderIds %+v", txInfo, OrderIds)

	if err1 := session.Begin(); err1 != nil {
		err = err1
		return
	}
	// 将info数据写入中间表
	OrderIds, err1 := this.insertOrderInfo(&data, &Extensions, txInfo, session)

	if err1 != nil {
		err = err1
		if err := session.Rollback(); err != nil {
			logger.Ex(this.ctx, "model.finishing.insertion.InsertOrder error", "InsertOrderIdx err", "err %+v , txinfo %+v", err, txInfo)
		}
		return
	}
	// 将detail数据写入中间表
	err = this.insertOrderDetail(&data, OrderIds, txInfo, session)
	if err != nil {
		if err := session.Rollback(); err != nil {
			logger.Ex(this.ctx, "model.finishing.insertion.InsertOrder error", "insertOrderDetailIdx error", "err %+v , OrderIds %+v", err, OrderIds)
		}
		return
	}

	// 将extension数据写入中间表
	err = this.insertOrderExtension(&data, OrderIds, txInfo, Extensions, session)
	if err != nil {
		if err := session.Rollback(); err != nil {
			logger.Ex(this.ctx, "model.finishing.insertion.InsertOrder error", "insertOrderExtensionIdx error", "err %+v , OrderIds %+v , extension %+v", err, OrderIds, Extensions)
		}
		return
	}
	if err := session.Commit(); err != nil {
		if err := session.Rollback(); err != nil {
			logger.Ex(this.ctx, "model.finishing.insertion.InsertOrder error", "insertOrderExtensionIdx error", "err %+v , OrderIds %+v , extension %+v", err, OrderIds, Extensions)
		}
		logger.Ex(this.ctx, "model.finishing.insertion.InsertOrder error", "InsertOrderIdx failed", "OrderIds %+v", OrderIds)
	}
	return
}

func (this *Insertion) getBatchInfos(OrderIds []string, txInfo redis.TxInfo) (ret []mysql.OrderInfo) {
	noCacheOrderIds := make([]string, 0)
	ret = make([]mysql.OrderInfo, 0)
	for i := 0; i < len(OrderIds); i++ {
		if _, ok := txInfo.OrderInfos[OrderIds[i]]; ok {
			ret = append(ret, txInfo.OrderInfos[OrderIds[i]].Info)
		} else {
			noCacheOrderIds = append(noCacheOrderIds, OrderIds[i])
		}
	}
	if len(noCacheOrderIds) >= 1 {
		// 根据orderIds 查询所有的detail信息
		infoDao := mysql.NewOrderInfoDao(this.ctx, this.uAppId, this.uUserId, constant.DBWriter)
		details, err := infoDao.Get(noCacheOrderIds, "", []uint{constant.TxStatusCommitted}, "", "")
		if err != nil {
			logger.Ex(this.ctx, "model.finishing.insertion.getBatchInfos", "getOrderInfo get err", "orderids %+v , err %+v", OrderIds, err)
			return
		}
		for i := 0; i < len(details); i++ {
			ret = append(ret, details[i])
		}
	}
	return
}

/**
 * @description:  根据orderid 查询orderinfo信息，然后将固定的字段组装后 写入mysql表
 * @params {params}
 * @return:
 */
func (this *Insertion) insertOrderInfo(data *map[string]map[string]interface{}, Extensions *[]string, txInfo redis.TxInfo, session *torm.Session) (OrderIds []string, err error) {

	//logger.Ix(this.ctx, "model.finishing.insertion.insertOrderInfo","insertOrderInfo start","txinfo %+v", txInfo)

	MysqlInfoDao := mysql.NewOrderInfoDao(this.ctx, this.uAppId, this.uUserId, constant.DBWriter)

	OrderIds = make([]string, 0)

	var MysqlOrderInfo []mysql.OrderInfo
	var IdxOrderInfo []mysql.OrderInfoIdx
	// 若orderIds为空，则根据 txId 查询本次提交的orderId
	if (len(txInfo.OrderIds)) == 0 {
		MysqlOrderInfo, err = MysqlInfoDao.Get([]string{}, this.txId, []uint{constant.TxStatusCommitted}, "", "")
		for i := 0; i < len(MysqlOrderInfo); i++ {
			OrderIds = append(OrderIds, MysqlOrderInfo[i].OrderId)
		}

		if ExtensionInfo, ok := utils.ExtensionTableShardingInfo[this.uAppId]; ok {
			for extension, _ := range ExtensionInfo {
				(*Extensions) = append((*Extensions), extension)
			}
		}
	} else {
		OrderIds = txInfo.OrderIds
		MysqlOrderInfo = this.getBatchInfos(OrderIds, txInfo)
		for extension, _ := range txInfo.Extensions {
			(*Extensions) = append((*Extensions), extension)
		}
	}

	// 若查询mysql出现错误 直接返回
	if err != nil {
		logger.Ex(this.ctx, "model.finishing.insertion.insertOrderInfo error", "insertOrderInfo MysqlOrderInfo get err", "%+v", err)
		return
	}

	IdxOrderInfo = make([]mysql.OrderInfoIdx, len(MysqlOrderInfo))
	for i := 0; i < len(OrderIds); i++ {
		(*data)[OrderIds[i]] = make(map[string]interface{})
	}

	for i := 0; i < len(MysqlOrderInfo); i++ {
		var info mysql.OrderInfoIdx
		(*data)[MysqlOrderInfo[i].OrderId][constant.Order_HashSubKey_Info] = MysqlOrderInfo[i]
		(*data)[MysqlOrderInfo[i].OrderId][constant.Order_HashSubKey_Detail] = make([]mysql.OrderDetail, 0)

		info.Bid = MysqlOrderInfo[i].Id
		info.UserId = MysqlOrderInfo[i].UserId
		info.OrderId = MysqlOrderInfo[i].OrderId
		info.CreatedTime = MysqlOrderInfo[i].CreatedTime
		info.UpdatedTime = MysqlOrderInfo[i].UpdatedTime
		info.ExpiredTime = MysqlOrderInfo[i].ExpiredTime
		info.PaidTime = MysqlOrderInfo[i].PaidTime
		info.PayCreatedTime = MysqlOrderInfo[i].PayCreatedTime
		info.CancelledTime = MysqlOrderInfo[i].CancelledTime
		IdxOrderInfo[i] = info
	}

	orderDao := mysql.NewOrderInfoIdxDao2(this.uAppId, session, this.ctx)
	ret, err := imysql.Create(orderDao, IdxOrderInfo)
	if err != nil {
		panic(err)
		return
	}
	logger.Ix(this.ctx, "model.finishing.insertion.insertOrderInfo", "insertOrderInfoIdx succ", "this %+v, session +%v, ret: %d", this, session, ret)

	return
}

func (this *Insertion) getBatchDetails(OrderIds []string, txInfo redis.TxInfo) (ret []mysql.OrderDetail) {
	noCacheOrderIds := make([]string, 0)
	ret = make([]mysql.OrderDetail, 0)
	for i := 0; i < len(OrderIds); i++ {
		if _, ok := txInfo.OrderInfos[OrderIds[i]]; ok {
			o := txInfo.OrderInfos[OrderIds[i]].Detail
			for j := 0; j < len(o); j++ {
				ret = append(ret, o[j])
			}
		} else {
			noCacheOrderIds = append(noCacheOrderIds, OrderIds[i])
		}
	}
	if len(noCacheOrderIds) >= 1 {
		// 根据orderIds 查询所有的detail信息
		detailDao := mysql.NewOrderDetailDao(this.ctx, this.uAppId, this.uUserId, constant.DBWriter)
		details, err := detailDao.Get(noCacheOrderIds)
		if err != nil {
			logger.Ex(this.ctx, "model.finishing.insertion.getBatchDetails", "insertOrderDetail OrderDetail get err", "orderids %+v , err %+v", OrderIds, err)
			return
		}
		for i := 0; i < len(details); i++ {
			ret = append(ret, details[i])
		}
	}
	return
}

/**
 * @description:  根据orderid 查询orderdetail信息，然后将固定的字段组装后 写入索引库
 * @params {params}
 * @return:
 */
func (this *Insertion) insertOrderDetail(data *map[string]map[string]interface{}, OrderIds []string, txInfo redis.TxInfo, session *torm.Session) (err error) {

	var details []mysql.OrderDetail

	details = this.getBatchDetails(OrderIds, txInfo)

	// 组装order_detail 并写入索引库中
	IdxOrderDetails := make([]mysql.OrderDetailIdx, 0)
	for i := 0; i < len(details); i++ {
		var info mysql.OrderDetailIdx
		info.Bid = details[i].Id
		info.UserId = details[i].UserId
		info.OrderId = details[i].OrderId
		info.ProductName = details[i].ProductName
		info.ProductId = details[i].ProductId
		info.ParentProductId = details[i].ParentProductId
		info.CreatedTime = details[i].CreatedTime
		info.UpdatedTime = details[i].UpdatedTime
		IdxOrderDetails = append(IdxOrderDetails, info)
		(*data)[info.OrderId][constant.Order_HashSubKey_Detail] = append((*data)[info.OrderId][constant.Order_HashSubKey_Detail].([]mysql.OrderDetail), details[i])
	}

	dao := mysql.NewOrderDetailIdxDao2(this.uAppId, session, this.ctx)
	ret, err := imysql.Create(dao, IdxOrderDetails)

	if err != nil {
		panic(err)
		return
	}
	logger.Ix(this.ctx, "model.finishing.insertion.insertOrderDetail", "insertOrderDetail succ", "ret %+v, orderids:%+v", ret, OrderIds)

	return
}

func (this *Insertion) getBatchExtensions(OrderIds []string, extension string, txInfo redis.TxInfo) (ret []map[string]interface{}) {
	noCacheOrderIds := make([]string, 0)
	ret = make([]map[string]interface{}, 0)
	toUpdateFields := utils.IdxdbCollectionFieldInfo[this.uAppId][extension] // 获取extension对应的字段

	for i := 0; i < len(OrderIds); i++ {
		if _, ok := txInfo.OrderInfos[OrderIds[i]]; ok {
			o := txInfo.OrderInfos[OrderIds[i]].Extensions
			tmp := o[extension]
			for j := 0; j < len(tmp); j++ {
				_ret := true
				for k, _ := range toUpdateFields {
					if tmp[j][k] == nil { //  判断索引字段是否为空
						noCacheOrderIds = append(noCacheOrderIds, OrderIds[i])
						_ret = false
						break
					}
				}
				if _ret == false {
					continue
				}
				ret = append(ret, tmp[j])
			}
		} else {
			noCacheOrderIds = append(noCacheOrderIds, OrderIds[i])
		}
	}
	if len(noCacheOrderIds) >= 1 {
		// 根据orderIds 查询所有的detail信息
		extDao := mysql.NewExtensionDao(this.ctx, this.uAppId, this.uUserId, extension, constant.DBWriter)
		details, err := extDao.Get(noCacheOrderIds)
		if err != nil {
			logger.Ex(this.ctx, "model.finishing.insertion.getBatchExtensions", "getExtetnion err", "noCacheOrderIds: %+v ,extension:%s , err:%+v", noCacheOrderIds, extension, err)
			return
		}
		for i := 0; i < len(details); i++ {
			ret = append(ret, details[i])
		}
	}
	return
}

/**
 * @description: 根据orderid 查询order extension信息，然后将固定的字段组装后 写入索引库
 * @params {params}
 * @return:
 */
func (this *Insertion) insertOrderExtension(data *map[string]map[string]interface{}, OrderIds []string, txInfo redis.TxInfo, Extensions []string, session *torm.Session) (err error) {

	for i := 0; i < len(Extensions); i++ { // ["group_info","iextension"]
		extension := Extensions[i]

		ExtInfo := this.getBatchExtensions(OrderIds, extension, txInfo)

		if err != nil || len(ExtInfo) < 1 {
			logger.Wx(this.ctx, "model.finishing.insertion.insertOrderExtension", "insertOrderExtension extDao get err", "OrderIds %+v, err %+v", OrderIds, err)
			continue
		}

		dao := mysql.NewExtensionIdxDao2(this.uAppId, extension, session, this.ctx)

		toUpdateFields := utils.IdxdbCollectionFieldInfo[this.uAppId][extension] // 获取extension对应的字段

		for j := 0; j < len(ExtInfo); j++ {
			OrderId := ExtInfo[j]["order_id"].(string)
			if _, ok := (*data)[OrderId][extension]; !ok {
				(*data)[OrderId][extension] = make([]map[string]interface{}, 0)
			}
			(*data)[OrderId][extension] = append((*data)[OrderId][extension].([]map[string]interface{}), ExtInfo[j])

			if toUpdateFields == nil { // 未配置 索引库， 则不需要更新
				continue
			}
			// 根据ordersdk里面配置的fields 进行赋值
			tmpFieldMap := make(map[string]interface{}, len(toUpdateFields)+3)
			tmpFieldMap["bid"] = ExtInfo[j]["id"]
			tmpFieldMap["order_id"] = OrderId
			tmpFieldMap["user_id"] = ExtInfo[j]["user_id"]
			for k, v := range toUpdateFields {
				if v == constant.TypeTimestamp { // 需要将时间转为为unix时间戳
					tmpFieldMap[k], _ = utils.JsonInterface2Timestamp(ExtInfo[j][k])
				} else {
					tmpFieldMap[k] = ExtInfo[j][k]
				}
			}
			_, err = imysql.Create(dao, tmpFieldMap)
			if err != nil {
				logger.Ex(this.ctx, "model.finishing.insertion.insertOrderExtension err", "insertOrderExtension err", "orderid:%s ,extension:%s, err:%+v", tmpFieldMap["bid"], extension, err)
			}
		}
	}
	logger.Ix(this.ctx, "model.finishing.insertion.insertOrderExtension", "insertOrderExtension succ", "extension %+v", Extensions)

	return
}

/**
 * 设置缓存
 */
func (this *Insertion) SetCache(orderIdArr []string, txInfo redis.TxInfo) (ok bool) {
	ok = false
	data := make(map[string]map[string]interface{})
	OrderIds := make([]string, 0)
	if len(orderIdArr) < 1 {
		OrderIds = txInfo.OrderIds
	} else {
		OrderIds = orderIdArr
	}

	if len(OrderIds) < 1 {
		return
	}
	MysqlOrderInfo := this.getBatchInfos(OrderIds, txInfo)

	// 组装info数据
	for i := 0; i < len(MysqlOrderInfo); i++ {
		data[MysqlOrderInfo[i].OrderId] = make(map[string]interface{})
		data[MysqlOrderInfo[i].OrderId][constant.Order_HashSubKey_Info] = MysqlOrderInfo[i]
		data[MysqlOrderInfo[i].OrderId][constant.Order_HashSubKey_Detail] = make([]mysql.OrderDetail, 0)
	}

	// 组装order_detail 并写入索引库中
	details := this.getBatchDetails(OrderIds, txInfo)
	for i := 0; i < len(details); i++ {
		data[details[i].OrderId][constant.Order_HashSubKey_Detail] = append(data[details[i].OrderId][constant.Order_HashSubKey_Detail].([]mysql.OrderDetail), details[i])
	}

	extensions := utils.ExtensionTableShardingInfo[this.uAppId]

	for extension, _ := range extensions {
		ExtInfo := this.getBatchExtensions(OrderIds, extension, txInfo)
		if len(ExtInfo) < 1 {
			logger.Wx(this.ctx, "model.finishing.insertion.SetCache", "insertion.SetCache.getBatchExtensions empty", "OrderIds:%+v extension:%s", OrderIds, extension)
			continue
		}
		for j := 0; j < len(ExtInfo); j++ {
			OrderId := ExtInfo[j]["order_id"].(string)
			if _, ok := data[OrderId][extension]; !ok {
				data[OrderId][extension] = make([]map[string]interface{}, 0)
			}
			data[OrderId][extension] = append(data[OrderId][extension].([]map[string]interface{}), ExtInfo[j])
		}
	}
	// 将数据写入缓存
	OrderRedisDao := redis.NewOrderDao(this.ctx, this.uAppId, this.uUserId)
	for OrderId, OrderData := range data {
		if err := OrderRedisDao.HMSet(OrderId, OrderData); err != nil {
		}
	}
	ok = true
	//logger.Ix(this.ctx,"model.finishing.insertion.SetCache","insertion.SetCache","OrderIds:%+v , idarrlen:%d", OrderIds, len(orderIdArr))
	return
}
