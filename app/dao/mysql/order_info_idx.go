/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-01 16:45:33
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-04-11 01:25:50
 * @Description:
 */
package mysql

import (
	"context"
	"encoding/json"
	logger "github.com/tal-tech/loggerX"
	"github.com/tal-tech/torm"
	"powerorder/app/utils"
)

type OrderInfoIdx struct {
	Id             int64      `json:"-" xorm:"-"`
	Bid            int64      `json:"bid" xorm:"not null default '0' comment('对应原库表的id') BIGINT(20)"`
	UserId         uint64     `json:"user_id" xorm:"not null default 0 comment('user_id，如果业务线user_id为字符串，则可以使用36进制或59进制法转成10进制') index(user_app) BIGINT(20)"`
	OrderId        string     `json:"order_id" xorm:"not null default '' comment('订单号') unique(order_app) CHAR(32)"`
	CreatedTime    utils.Time `json:"created_time" xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	CancelledTime  utils.Time `json:"cancelled_time" xorm:"not null default '1970-00-00 00:00:00' comment('取消时间') TIMESTAMP"`
	UpdatedTime    utils.Time `json:"updated_time" xorm:"not null default '1970-00-00 00:00:00' comment('修改时间') TIMESTAMP"`
	PayCreatedTime utils.Time `json:"pay_created_time" xorm:"not null default '1970-00-00 00:00:00' comment('支付开始时间') datetime"`
	PaidTime       utils.Time `json:"paid_time" xorm:"not null default 'CURRENT_TIMESTAMP' comment('支付时间') TIMESTAMP"`
	ExpiredTime    utils.Time `json:"expired_time" xorm:"default null comment('过期时间') datetime"`
}

type OrderInfoIdxDao struct {
	Dao
}

/**
 * @description: 类型转换，OrderDetail => map
 * @params {detail} OrderDetail
 * @return:
 */
func OrderInfoIdx2Map(info OrderInfo) (res map[string]interface{}, err error) {
	j, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}
	res = make(map[string]interface{})
	json.Unmarshal(j, &res)
	return
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func NewOrderInfoIdxDao2(uAppId uint, session *torm.Session, ctx context.Context) *OrderInfoIdxDao {
	object := new(OrderInfoIdxDao)
	object.InitDatabaseName(uAppId, 0)
	logger.Ix(ctx, "session", "session = %v", session)

	object.UpdateEngine(session)
	logger.Ix(ctx, "engine", "engine = %v", object.Engine)

	object.InitTableName()
	object.SetTable(object.GetTable())
	object.Dao.InitCtx(ctx)

	//object.Engine.ShowSQL(true)

	return object
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func NewOrderInfoIdxDao(ctx context.Context, uAppId uint, writer string) *OrderInfoIdxDao {
	object := new(OrderInfoIdxDao)
	object.InitDatabaseName(uAppId, 0)
	writer = utils.DbShadowHandler(ctx, writer)
	if ins := torm.GetDbInstance(object.strDatabaseName, writer); ins != nil {
		object.UpdateEngine(ins.Engine)
	} else {
		return nil
	}
	object.Engine.ShowSQL(true)
	object.InitTableName()
	object.SetTable(object.GetTable())
	object.Dao.InitCtx(ctx)

	return object
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *OrderInfoIdxDao) InitDatabaseName(uAppId uint, id uint64) {
	this.uAppId = uAppId
	this.strDatabaseName = utils.GenDatabaseIdxName()
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *OrderInfoIdxDao) InitTableName() {
	this.strTableName = utils.GenOrderInfoIdxTableName(this.uAppId)

}

/**
 * @description: mysql查询订单信息 [已提交状态]
 * @params {arrOrderIds} 订单ID集合
 * @params {strTxId} 事务ID
 * @params {txStatus} 事务状态
 * @return: 订单信息集合
 */

func (this *OrderInfoIdxDao) Get(arrOrderIds []string) (ret []OrderInfoIdx, err error) {

	if this == nil {
		err = logger.NewError("system error dao instance is nil")
		return
	}
	ret = make([]OrderInfoIdx, 0)
	this.InitSession()
	this.InitTableName()
	this.SetTable(this.strTableName)

	if len(arrOrderIds) > 0 {
		this.BuildQuery(torm.CastToParamIn(arrOrderIds), "order_id")
	}

	err = this.Session.Find(&ret)

	if err != nil {
		logger.Ex(this.ctx, "DaoMysqlOrderInfoDao", "get from db error = %v", err)
	}
	return
}
