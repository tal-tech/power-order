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
	"github.com/spf13/cast"
	logger "github.com/tal-tech/loggerX"
	"github.com/tal-tech/torm"
	"github.com/tal-tech/xtools/traceutil"
	"powerorder/app/utils"
	"time"
)

type OrderInfo struct {
	Id             int64      `json:"id" xorm:"not null pk autoincr INT(11)"`
	BId            int64      `json:"bid" xorm:"not null default 0 bigint(20)"`
	UserId         uint64     `json:"user_id" xorm:"not null default 0 comment('user_id，如果业务线user_id为字符串，则可以使用36进制或59进制法转成10进制') index(user_app) BIGINT(20)"`
	AppId          uint       `json:"app_id" xorm:"not null default 0 comment('接入的商户ID  appid') unique(order_app) index(user_app) INT(11)"`
	OrderId        string     `json:"order_id" xorm:"not null default '' comment('订单号') unique(order_app) CHAR(32)"`
	OrderType      uint       `json:"order_type" xorm:"not null default 0 comment('订单类型') TINYINT(4)"`
	Status         uint       `json:"status" xorm:"not null default 1 comment('状态(如：1:未付款，2:支付中，3:支付成功，4:用户手动取消，5:已过期脚本取消)') TINYINT(2)"`
	Source         uint       `json:"source" xorm:"not null default 1 comment('订单来源（如： 1: 商城, 2: 购物车, 3:续报列表）') TINYINT(2)"`
	OrderDevice    uint       `json:"order_device" xorm:"not null default 1 comment('订单设备（如1：pc，2：iPad，3：Touch，4：APP，7：IOS，8：Android）') TINYINT(2)"`
	CreatedTime    utils.Time `json:"created_time" xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	CancelledTime  utils.Time `json:"cancelled_time" xorm:"not null default '1970-00-00 00:00:00' comment('取消时间') TIMESTAMP"`
	UpdatedTime    utils.Time `json:"updated_time" xorm:"not null default '1970-00-00 00:00:00' comment('修改时间') TIMESTAMP"`
	Extras         string     `json:"extras" xorm:"not null comment('订单产生中附属信息存储 比如交易快照之类的，不会来查询，不会用来检索，格式任意') TEXT"`
	Price          uint       `json:"price" xorm:"not null default 0 comment('订单金额，单位：分') INT(11)"`
	PromotionPrice uint       `json:"promotion_price" xorm:"not null default 0 comment('促销总金额，单位：分') INT(11)"`
	RealPrice      uint       `json:"real_price" xorm:"not null default 0 comment('实际缴费金额，单位：分') INT(11)"`
	PayDevice      uint       `json:"pay_device" xorm:"not null default 1 comment('支付订单时的设备（如1：pc，2：iPad，3：Touch，4：APP，7：IOS，8：Android）') TINYINT(2)"`
	PayCreatedTime utils.Time `json:"pay_created_time" xorm:"not null default '1970-00-00 00:00:00' comment('支付开始时间') datetime"`
	PaidTime       utils.Time `json:"paid_time" xorm:"not null default 'CURRENT_TIMESTAMP' comment('支付时间') TIMESTAMP"`
	ExpiredTime    utils.Time `json:"expired_time" xorm:"default null comment('过期时间') datetime"`
	TxId           string     `json:"-" xorm:"not null default '' comment('事务id，格式为traceId的后32位') VARCHAR(32)"`
	TxStatus       uint       `json:"-" xorm:"not null default 0 comment('事务状态：0，刚创建 1、已提交 2~50、被回滚（可自定义回滚原因）') TINYINT(2)"`
	Version        uint       `json:"version"  xorm:"not null default 0 comment('版本号') TINYINT(2)"`
}

type OrderInfoDao struct {
	Dao
}

/**
 * @description: 类型转换，OrderDetail => map
 * @params {detail} OrderDetail
 * @return:
 */
func OrderInfo2Map(info OrderInfo) (res map[string]interface{}, err error) {
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
func NewOrderInfoDao2(uAppId uint, uUserId uint64, session *torm.Session, ctx context.Context) *OrderInfoDao {
	object := new(OrderInfoDao)
	object.InitDatabaseName(uAppId, uUserId)
	object.UpdateEngine(session)
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
func NewOrderInfoDao(ctx context.Context, uAppId uint, uUserId uint64, writer string) *OrderInfoDao {
	object := new(OrderInfoDao)
	object.InitDatabaseName(uAppId, uUserId)
	writer = utils.DbShadowHandler(ctx, writer)
	if ins := torm.GetDbInstance(object.strDatabaseName, writer); ins != nil {
		object.UpdateEngine(ins.Engine)
	} else {
		return nil
	}
	object.UpdateEngine(uUserId, writer)
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
func (this *OrderInfoDao) InitTableName() {
	this.strTableName = utils.GenOrderInfoTableName(this.uUserId)

}

/**
 * @description: mysql查询订单信息 [已提交状态]
 * @params {arrOrderIds} 订单ID集合
 * @params {strTxId} 事务ID
 * @params {txStatus} 事务状态
 * @return: 订单信息集合
 */

func (this *OrderInfoDao) Get(arrOrderIds []string, strTxId string, txStatus []uint, startDate, endDate string) (ret []OrderInfo, err error) {

	if this == nil {
		err = logger.NewError("system error dao instance is nil")
		return
	}
	span, _ := traceutil.Trace(this.ctx, "dao_mysql_order_info")
	if span != nil {
		defer span.Finish()
		span.Tag("t1", cast.ToString(time.Now().UnixNano()))
	}
	ret = make([]OrderInfo, 0)
	this.InitSession()
	this.InitTableName()
	this.SetTable(this.strTableName)
	this.BuildQuery(this.uUserId, "user_id")
	this.BuildQuery(this.uAppId, "app_id")

	if len(arrOrderIds) > 0 {
		this.BuildQuery(torm.CastToParamIn(arrOrderIds), "order_id")
	}

	if len(strTxId) > 0 {
		this.BuildQuery(strTxId, "tx_id")
	}

	if len(txStatus) > 0 {
		this.BuildQuery(torm.CastToParamIn(txStatus), "tx_status")
	}
	if len(startDate) > 0 && len(endDate) > 0 {
		this.BuildQuery(torm.ParamRange{startDate, endDate}, "created_time")
	}

	if span != nil {
		span.Tag("t2", cast.ToString(time.Now().UnixNano()))
	}
	err = this.Session.Find(&ret)

	if err != nil {
		logger.Ex(this.ctx, "DaoMysqlOrderInfoDao", "get from db error = %v", err)
	}
	if span != nil {
		span.Tag("t3", cast.ToString(time.Now().UnixNano()))
	}
	return
}
