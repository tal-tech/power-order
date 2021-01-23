/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-01 16:44:42
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-02-23 01:44:05
 * @Description:
 */
package mysql

import (
	"context"
	"github.com/spf13/cast"
	"github.com/tal-tech/torm"
	"github.com/tal-tech/xtools/traceutil"
	"powerorder/app/utils"
	"time"
)

type OrderDetailIdx struct {
	Id              int64      `json:"-" xorm:"-"`
	Bid             int64      `json:"bid" xorm:"not null default '0' comment('对应原库表的id') BIGINT(20)"`
	UserId          uint64     `json:"user_id" xorm:"not null default 0 comment('user_id') BIGINT(20)"`
	OrderId         string     `json:"order_id" xorm:"not null default '' comment('订单号') index(order_app) CHAR(32)"`
	ProductId       uint       `json:"product_id" xorm:"not null default 0 comment('商品ID') INT(11)"`
	ProductName     string     `json:"product_name" xorm:"not null default '' comment('商品名称') VARCHAR(100)"`
	PromotionType   string     `json:"promotion_type" xorm:"not null default '' comment('促销类型') VARCHAR(100)"`
	ParentProductId uint       `json:"parent_product_id" xorm:"not null default 0 comment('父商品ID') INT(11)"`
	CreatedTime     utils.Time `json:"created_time" xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	UpdatedTime     utils.Time `json:"updated_time" xorm:"not null default 'CURRENT_TIMESTAMP' comment('修改时间') TIMESTAMP"`
}

type OrderDetailIdxDao struct {
	Dao
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func NewOrderDetailIdxDao2(uAppId uint, session *torm.Session, ctx context.Context) *OrderDetailIdxDao {
	object := new(OrderDetailIdxDao)
	object.InitDatabaseName(uAppId, 0)
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
func NewOrderDetailIdxDao(ctx context.Context, uAppId uint, writer string) *OrderDetailIdxDao {
	object := new(OrderDetailIdxDao)
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
func (this *OrderDetailIdxDao) InitDatabaseName(uAppId uint, id uint64) {
	this.uAppId = uAppId
	this.strDatabaseName = utils.GenDatabaseIdxName()
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *OrderDetailIdxDao) InitTableName() {
	this.strTableName = utils.GenOrderDetailIdxTableName(this.uAppId)
}

/**
 * @description:根据orderId查询mysql
 * @params {arrOrderIds} 订单ID集合
 * @return: 订单detail集合
 */

func (this *OrderDetailIdxDao) Get(arrOrderIds []string) (ret []OrderDetailIdx, err error) {
	span, _ := traceutil.Trace(this.ctx, "dao_mysql_order_detail")
	if span != nil {
		defer span.Finish()
		span.Tag("t1", cast.ToString(time.Now().UnixNano()))
	}
	ret = make([]OrderDetailIdx, 0)
	this.InitSession()
	this.InitTableName()
	this.SetTable(this.strTableName)
	this.BuildQuery(this.uUserId, "user_id")

	if len(arrOrderIds) > 0 {
		this.BuildQuery(torm.CastToParamIn(arrOrderIds), "order_id")
	}

	this.BuildQuery(this.uAppId, "app_id")
	if span != nil {
		span.Tag("t2", cast.ToString(time.Now().UnixNano()))
	}

	err = this.Session.Find(&ret)

	if span != nil {
		span.Tag("t3", cast.ToString(time.Now().UnixNano()))
	}

	return
}
