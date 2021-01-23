/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-19 23:54:27
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-04-12 09:19:38
 * @Description:
 */

package update

import (
	"context"
	"fmt"
	logger "github.com/tal-tech/loggerX"
	"github.com/tal-tech/torm"
	"powerorder/app/constant"
	"powerorder/app/dao/mysql"
	"powerorder/app/dao/redis"
	"powerorder/app/model/finishing"
	"powerorder/app/params"
	"powerorder/app/utils"
	"strconv"
)

type Update struct {
	uAppId  uint
	uUserId uint64
	ctx     context.Context
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func NewUpdate(ctx context.Context, uAppId uint, uUserId uint64) *Update {
	object := new(Update)
	object.ctx = ctx
	object.uAppId = uAppId
	object.uUserId = uUserId
	return object
}

/**
 * @description:
 * @params {params} 更新的参数
 * @return: err
 */
func (this *Update) Update(param params.ReqUpdate) (err error) {
	session := mysql.NewSession(this.ctx, param.AppId, param.UserId, constant.DBWriter)
	pts := utils.GetPts(this.ctx) // 压测标识
	if session == nil {
		err = logger.NewError("init session error")
		return
	}
	//logger.Dx(this.ctx,"model.update.update.Update","model update params = %+v", param)
	defer session.Close()
	OrderInfoDao := mysql.NewOrderInfoDao2(param.AppId, param.UserId, session, this.ctx)
	OrderDetailDao := mysql.NewOrderDetailDao2(param.AppId, param.UserId, session, this.ctx)

	session.Begin()
	var res int64
	//OrderDetailDao := mysql.NewOrderDetailDao2(params.AppId, params.UserId, session)
	for i := 0; i < len(param.Orders); i++ {
		var where *torm.Session

		for extName, ExtData := range param.Orders[i].Extensions {
			ExtDao := mysql.NewExtensionDao2(param.AppId, param.UserId, extName, session, this.ctx)

			for j := 0; j < len(ExtData.Updates); j++ {

				Version := ExtData.Updates[j].Version
				Update := ExtData.Updates[j].Update
				for field, value := range Update {
					switch value.(type) {
					case string:
						Update[field] = value.(string)
					case float64:
						Update[field], _ = strconv.Atoi(strconv.FormatFloat(value.(float64), 'f', -1, 64))
					}
				}

				if Version <= 0 {
					err = logger.NewError("error version")
					session.Rollback()
					return
				}
				if Version > 0 {
					ExtDao.SetTableName()
					where = ExtDao.Where("order_id = ?", param.Orders[i].OrderId)

					if Id, exist := Update["id"]; exist && len(Id.(string)) > 0 {
						where = where.And("id = ?", Id)
						//delete(Update, "id") // 移除更新字段 id   会影响数据
					}
					if Version < 0 {
						delete(Update, "version")
					} else {
						Update["version"] = strconv.Itoa(Version + 1)
						if pts != true { // 压测去掉这个条件
							where = where.And("version = ?", Version)
						}
					}

					res, err = where.Update(Update)
					if Version > 0 {
						if res == 0 {
							err = logger.NewError(fmt.Sprintf("%s nothing updated", extName))
							session.Rollback()
							return
						}
					}

					if err != nil {
						session.Rollback()
						return
					}
				} else {
					err = logger.NewError(fmt.Sprintf("%s error version", extName))
					session.Rollback()
					return
				}
			}

			for j := 0; j < len(ExtData.Insertions); j++ {
				Insertion := ExtData.Insertions[j]
				for field, value := range Insertion {
					switch value.(type) {
					case string:
						Insertion[field] = value.(string)
					case float64:
						Insertion[field], _ = strconv.Atoi(strconv.FormatFloat(value.(float64), 'f', -1, 64))
					}
				}
				Id := utils.GenId()
				((param.Orders[i].Extensions[extName]).Insertions)[j]["id"] = Id // ID是生成的，因此需要回写
				Insertion["order_id"] = param.Orders[i].OrderId
				Insertion["id"] = Id
				Insertion["version"] = 1
				Insertion["user_id"] = param.UserId
				ExtDao.SetTableName()

				res, err = ExtDao.Create(Insertion)

				if err != nil {
					session.Rollback()
					return
				}
			}
		}

		for j := 0; j < len(param.Orders[i].Detail); j++ {
			Version := param.Orders[i].Detail[j].Version
			if Version <= 0 {
				err = logger.NewError("error version")
				session.Rollback()
				return
			}
			if Version > 0 {
				Detail := param.Orders[i].Detail[j].Detail

				OrderDetailDao.SetTableName()

				where = OrderDetailDao.Where("order_id = ? AND app_id = ?", param.Orders[i].OrderId, param.AppId)

				if Detail.Id > 0 {
					where = where.And("id = ?", Detail.Id)
				}
				if Detail.ProductId > 0 {
					where = where.And("product_id = ?", Detail.ProductId)
				}
				if Detail.ProductType > 0 {
					where = where.And("product_type = ?", Detail.ProductType)
				}
				if Version == -1 {
					Detail.Version = 0
				} else if Version > 0 {
					Detail.Version = uint(Version) + 1
					where = where.And("version = ?", Version)
				}
				res, err = where.Update(Detail)

				if Version > 0 {
					if res == 0 {
						err = logger.NewError("detail nothing updated")
						session.Rollback()
						return
					}
				}

				if err != nil {
					session.Rollback()
					return
				}
			} else {
				err = logger.NewError("detail error version")
				session.Rollback()
				return
			}
		}

		Version := param.Orders[i].Info.Version
		if Version <= 0 {
			err = logger.NewError("error version")
			session.Rollback()
			return
		}
		if Version > 0 {
			Info := param.Orders[i].Info.Info

			OrderInfoDao.SetTableName()
			where = OrderInfoDao.Where("user_id = ? AND app_id = ? AND order_id = ?", param.UserId, param.AppId, param.Orders[i].OrderId)

			if Version == -1 {
				Info.Version = 0
			} else if Version > 0 {
				Info.Version = uint(Version) + 1
				where = where.And("version = ?", Version)
			}

			res, err = where.Update(Info)

			if Version > 0 {
				if res == 0 {
					err = logger.NewError("info nothing updated")
					session.Rollback()
					return
				}
			}

			if err != nil {
				session.Rollback()
				return
			}
		} else {
			err = logger.NewError("info error version")
			logger.Ex(this.ctx, "model.update.update.Update", "info error version", "")
			session.Rollback()
			return
		}
	}

	err = session.Commit()
	this.updateCache(param)
	if err == nil {
		go this.finishing(param)
	}
	return
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *Update) finishing(param params.ReqUpdate) {

	finishing := finishing.NewUpdate(param.AppId, param.UserId, this.ctx, "")
	err := finishing.UpdateOrder(param)
	if err != nil {

	}
}

/**
 * 先设置缓存
 */
func (this *Update) updateCache(param params.ReqUpdate) {
	if len(param.Orders) < 1 {
		return
	}
	finishing := finishing.NewInsertion(param.AppId, param.UserId, this.ctx, "")
	var txInfo redis.TxInfo
	OrderIds := make([]string, len(param.Orders))
	for i := 0; i < len(param.Orders); i++ {
		OrderIds[i] = param.Orders[i].OrderId
	}
	ok := finishing.SetCache(OrderIds, txInfo)
	logger.Ix(this.ctx, "model.update.update.updateCache", "update.updateCache", "ret:%+v", ok)
	return
}
