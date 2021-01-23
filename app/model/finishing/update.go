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
	"fmt"
	logger "github.com/tal-tech/loggerX"
	"github.com/tal-tech/torm"
	"powerorder/app/constant"
	"powerorder/app/dao/mysql"
	"powerorder/app/params"
	"powerorder/app/utils"
	"strconv"
)

type Update struct {
	action  uint8
	ctx     context.Context
	uAppId  uint
	uUserId uint64
	txId    string
}

func NewUpdate(uAppId uint, uUserId uint64, ctx context.Context, txId string) *Update {

	object := new(Update)
	object.uAppId = uAppId
	object.uUserId = uUserId
	object.ctx = ctx
	object.txId = txId
	return object
}

/**
 * @description:  根据orderids 将数据写入mongo中。
 * @params {params}
 * @return:
 */
func (this *Update) UpdateOrder(param params.ReqUpdate) (err error) {

	defer utils.Catch()

	// todo这里需要使用事务
	session := mysql.NewSessionIdx(this.ctx, this.uAppId, this.uUserId, constant.DBWriter)
	if session == nil {
		logger.Ex(this.ctx, "model.finishing.update.UpdateOrder", "insertOrderInfo session err", "%+v", session)

		return
	}

	defer session.Close()

	data := make(map[string]map[string]interface{})
	Extensions := make([]string, 0)
	//logger.Ix(this.ctx,"model.finishing.update.UpdateOrder","UpdateOrder start","params %+v, OrderIds %+v", params, OrderIds)

	// 启动事务
	session.Begin()

	// 将info数据写入索引库
	OrderIds, err1 := this.updateOrderInfo(&data, &Extensions, param, session)

	if err1 != nil {
		err = err1
		logger.Ex(this.ctx, "model.finishing.update.UpdateOrder error", "UpdateOrder updateOrderInfo err", "err %+v", err)
		session.Rollback()
		return
	}
	// 将detail数据写入索引库
	err = this.updateOrderDetail(&data, OrderIds, param, session)
	if err != nil {
		logger.Ex(this.ctx, "model.finishing.update.UpdateOrder error", "UpdateOrder updateOrderDetail err", "err %+v", err)
		session.Rollback()
		return
	}

	if len(Extensions) > 0 {
		// 将extension数据写入索引库
		err = this.updateOrderExtension(&data, OrderIds, Extensions, param, session)
		if err != nil {
			logger.Ex(this.ctx, "model.finishing.update.UpdateOrder error", "UpdateOrder updateOrderExtension err", "err %+v", err)
			session.Rollback()
			return
		}
	}

	session.Commit()

	// 将数据写入缓存
	//OrderRedisDao := redis.NewOrderDao(this.ctx, this.uAppId, this.uUserId, constant.Order_Redis_Cluster)
	//for OrderId, OrderData := range data {
	//	OrderRedisDao.HMSet(OrderId, OrderData)
	//}
	logger.Ix(this.ctx, "model.finishing.update.UpdateOrder", "UpdateOrder success", "OrderIds %+v, params %+v", OrderIds, param)

	return
}

/**
 * @description:  根据orderid 查询orderinfo信息，然后将固定的字段组装后 写入索引库
 * @params {params}t
 * @return:
 */
func (this *Update) updateOrderInfo(data *map[string]map[string]interface{}, Extensions *[]string, param params.ReqUpdate, session *torm.Session) (OrderIds []string, err error) {
	logger.Ix(this.ctx, "model.finishing.update.updateOrderInfo", "updateOrderInfo start", "params %+v", param)

	OrderIds = make([]string, len(param.Orders))
	var MysqlOrderInfo []mysql.OrderInfo

	for i := 0; i < len(param.Orders); i++ {
		OrderIds[i] = param.Orders[i].OrderId
		(*data)[OrderIds[i]] = make(map[string]interface{})
	}
	// 对应appid需要的扩展字段
	if ExtensionInfo, ok := utils.ExtensionTableShardingInfo[this.uAppId]; ok {
		for extension, _ := range ExtensionInfo {
			(*Extensions) = append((*Extensions), extension)
		}
	}

	// 根据orderids从mysql查询对应的订单信息
	MysqlInfoDao := mysql.NewOrderInfoDao(this.ctx, this.uAppId, this.uUserId, constant.DBWriter)
	MysqlOrderInfo, err = MysqlInfoDao.Get(OrderIds, "", []uint{constant.TxStatusCommitted}, "", "")

	if err != nil {
		logger.Ex(this.ctx, "model.finishing.update.updateOrderInfo error", "updateOrderInfo MysqlOrderInfo get err", "%+v", err)
		return
	}

	// 构造数据data hashmap，以便更新到索引库中
	for i := 0; i < len(MysqlOrderInfo); i++ {
		(*data)[MysqlOrderInfo[i].OrderId][constant.Order_HashSubKey_Info] = MysqlOrderInfo[i]
		(*data)[MysqlOrderInfo[i].OrderId][constant.Order_HashSubKey_Detail] = make([]mysql.OrderDetail, 0)
	}

	orderDao := mysql.NewOrderInfoIdxDao2(this.uAppId, session, this.ctx)
	for i := 0; i < len(param.Orders); i++ {
		if param.Orders[i].Info.Version != 0 {
			Info := (*data)[param.Orders[i].OrderId][constant.Order_HashSubKey_Info].(mysql.OrderInfo)
			Update := map[string]interface{}{
				"updated_time":     Info.UpdatedTime.String(),
				"cancelled_time":   Info.CancelledTime.String(),
				"expired_time":     Info.ExpiredTime.String(),
				"pay_created_time": Info.PayCreatedTime.String(),
				"paid_time":        Info.PaidTime.String(),
				"created_time":     Info.CreatedTime.String(),
			}
			ret, err := orderDao.UpdateColsWhere(Update, nil, " user_id = ? and bid = ? ", []interface{}{this.uUserId, Info.Id})
			logger.Ix(this.ctx, "model.finishing.update.updateOrderInfo", "updateOrderInfo orderInfo update", "ret %+v, userid %d , bid %d , err %+v", ret, param.UserId, Info.Id, err)
		}
	}
	return
}

/**
 * @description:  根据orderid 查询orderdetail信息，然后将固定的字段组装后 写入索引库
 * @params {params}
 * @return:
 */
func (this *Update) updateOrderDetail(data *map[string]map[string]interface{}, OrderIds []string, param params.ReqUpdate, session *torm.Session) (err error) {

	//logger.Ix(this.ctx,"model.finishing.update.updateOrderDetail","updateOrderDetail start","OrderIds %+v", OrderIds)

	var details []mysql.OrderDetail
	// 根据orderIds 查询所有的detail信息
	detailDao := mysql.NewOrderDetailDao(this.ctx, this.uAppId, this.uUserId, constant.DBWriter)
	details, err = detailDao.Get(OrderIds)
	if err != nil {
		logger.Ex(this.ctx, "model.finishing.update.updateOrderDetail error", "updateOrderDetail OrderDetail get err", "orderids %+v , err %+v", OrderIds, err)
		return
	}
	// tmp details hash
	detailsMap := make(map[int64]mysql.OrderDetail, len(details))
	for i := 0; i < len(details); i++ {
		detailsMap[details[i].Id] = details[i]
		(*data)[details[i].OrderId][constant.Order_HashSubKey_Detail] = append((*data)[details[i].OrderId][constant.Order_HashSubKey_Detail].([]mysql.OrderDetail), details[i])
	}

	dao := mysql.NewOrderDetailIdxDao2(this.uAppId, session, this.ctx)

	// 组装order_detail 并写入索引库中
	for i := 0; i < len(param.Orders); i++ {

		tmpDetails := param.Orders[i].Detail
		detailLen := len(tmpDetails)
		if detailLen < 1 {
			break
		}
		for j := 0; j < detailLen; j++ {
			if tmpDetails[j].Version == 0 {
				continue
			}
			t1 := tmpDetails[j].Detail
			if t1.Id < 1 {
				continue
			}
			Detail := detailsMap[t1.Id]

			Update := map[string]interface{}{
				"created_time": Detail.CreatedTime.String(),
				"updated_time": Detail.UpdatedTime.String(),
				"product_name": Detail.ProductName,
				"product_id":   Detail.ProductId,
			}
			ret, err := dao.UpdateColsWhere(Update, nil, " user_id = ? and bid = ? ", []interface{}{this.uUserId, t1.Id})
			logger.Ix(this.ctx, "model.finishing.update.updateOrderDetail", "updateOrderDetail UpdateOrderDetailOne", "ret %+v, err %+v", ret, err)
		}
	}
	return
}

/**
 * @description: 根据orderid 查询order extension信息，然后将固定的字段组装后 写入索引库
 * 更新的时候，传入的数据会有 update 或者 insert，但是最终会落库到db，因此这里直接从db查询。不需要再区分 insert或者update的数据。
 * @params {params}
 * @return:
 */
func (this *Update) updateOrderExtension(data *map[string]map[string]interface{}, OrderIds []string, Extensions []string, param params.ReqUpdate, session *torm.Session) (err error) {

	for i := 0; i < len(param.Orders); i++ {
		OrderId := param.Orders[i].OrderId
		extensions := param.Orders[i].Extensions
		if len(extensions) < 1 {
			continue
		}
		for extName, extData := range extensions {
			if _, ok := utils.IdxdbCollectionFieldInfo[this.uAppId][extName]; !ok {
				logger.Wx(this.ctx, "model.finishing.update.updateOrderExtension", "no such extension!", "extension %+v, appId %+v", extName, this.uAppId)
				continue
			}
			extDao := mysql.NewExtensionIdxDao2(this.uAppId, extName, session, this.ctx)
			toUpdateFields := utils.IdxdbCollectionFieldInfo[this.uAppId][extName] // 获取extension对应的字段

			for j := 0; j < len(extData.Updates); j++ {

				Version := extData.Updates[j].Version
				Update := extData.Updates[j].Update
				for field, value := range Update {
					switch value.(type) {
					case string:
						Update[field] = value.(string)
					case float64:
						Update[field], _ = strconv.Atoi(strconv.FormatFloat(value.(float64), 'f', -1, 64))
					}
				}
				// 根据ordersdk里面配置的fields 进行赋值
				tmpFieldMap := make(map[string]interface{}, len(toUpdateFields)+3)
				for k, _ := range toUpdateFields {
					if _, ok := Update[k]; !ok {
						//logger.Wx(this.ctx,"model.finishing.update.updateOrderExtension","extension no such fields!","extension %+v, k %+v", extName, k)
						continue
					}
					tmpFieldMap[k] = Update[k]
				}

				Id := Update["id"]
				if Version <= 0 {
					err = logger.NewError("error version")
					return
				}
				extDao.SetTable(extDao.GetTable())
				where := extDao.Where("user_id = ?", this.uUserId)
				where = where.And("bid = ?", Id)

				res, err1 := where.Update(tmpFieldMap)
				if res == 0 {
					err = logger.NewError(fmt.Sprintf("%s nothing updated", extName))
					return
				}

				if err1 != nil {
					//session.Rollback()
					err = err1
					return
				}
			}

			// 需要新增的数据
			for j := 0; j < len(extData.Insertions); j++ {
				Insert := extData.Insertions[j]
				for field, value := range Insert {
					switch value.(type) {
					case string:
						Insert[field] = value.(string)
					case float64:
						Insert[field], _ = strconv.Atoi(strconv.FormatFloat(value.(float64), 'f', -1, 64))
					}
				}
				// 根据ordersdk里面配置的fields 进行赋值
				tmpFieldMap := make(map[string]interface{}, len(toUpdateFields)+3)
				for k, _ := range toUpdateFields {
					if _, ok := Insert[k]; !ok {
						//logger.Wx(this.ctx,"model.finishing.update.updateOrderExtension","extension no such fields!","extension %+v, k %+v", extName, k)
						continue
					}
					tmpFieldMap[k] = Insert[k]
				}
				tmpFieldMap["bid"] = Insert["id"]
				tmpFieldMap["order_id"] = OrderId
				tmpFieldMap["user_id"] = this.uUserId

				extDao.SetTable(extDao.GetTable())
				res, err1 := extDao.Create(tmpFieldMap)

				if res == 0 {
					err = logger.NewError(fmt.Sprintf("%s nothing updated", extName))
					return
				}

				if err1 != nil {
					err = err1
					return
				}
			}
		}

	}

	return
}
