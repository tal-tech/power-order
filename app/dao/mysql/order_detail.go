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
	"encoding/json"
	"github.com/spf13/cast"
	"github.com/tal-tech/torm"
	"github.com/tal-tech/xtools/traceutil"
	"powerorder/app/utils"
	"time"
)

type OrderDetail struct {
	Id                int64      `json:"id" xorm:"not null pk autoincr INT(11)"`
	AppId             uint       `json:"app_id" xorm:"not null default 0 comment('接入的商户ID  app_id') index(order_app) INT(11)"`
	UserId            uint64     `json:"user_id" xorm:"not null default 0 comment('user_id') BIGINT(20)"`
	OrderId           string     `json:"order_id" xorm:"not null default '' comment('订单号') index(order_app) CHAR(32)"`
	ProductId         uint       `json:"product_id" xorm:"not null default 0 comment('商品ID') INT(11)"`
	ProductType       uint       `json:"product_type" xorm:"not null default 1 comment('商品类别') TINYINT(4)"`
	ProductName       string     `json:"product_name" xorm:"not null default '' comment('商品名称') VARCHAR(100)"`
	ProductNum        uint       `json:"product_num" xorm:"not null default 1 comment('商品数量') INT(11)"`
	ProductPrice      uint       `json:"product_price" xorm:"not null default 0 comment('商品销售金额') INT(11)"`
	CouponPrice       uint       `json:"coupon_price" xorm:"not null default 0 comment('优惠券分摊金额') INT(11)"`
	PromotionPrice    uint       `json:"promotion_price" xorm:"not null default 0 comment('促销分摊金额') INT(10)"`
	PromotionId       string     `json:"promotion_id" xorm:"not null default 0 comment('促销id') char(24)"`
	PromotionType     string     `json:"promotion_type" xorm:"not null default '' comment('促销类型') VARCHAR(100)"`
	ParentProductId   uint       `json:"parent_product_id" xorm:"not null default 0 comment('父商品ID') INT(11)"`
	ParentProductType uint       `json:"parent_product_type" xorm:"not null default 1 comment('父商品类别，业务线可自己定义') TINYINT(4)"`
	Extras            string     `json:"extras" xorm:"not null comment('订单商品中附属信息存储 比如促销的关键不变更信息存储之类的，不会来查询，不会用来检索') TEXT"`
	CreatedTime       utils.Time `json:"created_time" xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	UpdatedTime       utils.Time `json:"updated_time" xorm:"not null default 'CURRENT_TIMESTAMP' comment('修改时间') TIMESTAMP"`
	Version           uint       `json:"version"  xorm:"not null default 0 comment('版本号') TINYINT(2)"`
	SourceId          string     `json:"source_id" xorm:"not null default '' comment('热点数据') varchar(255)"`
}

type OrderDetailDao struct {
	Dao
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func NewOrderDetailDao2(uAppId uint, uUserId uint64, session *torm.Session, ctx context.Context) *OrderDetailDao {
	object := new(OrderDetailDao)
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
func NewOrderDetailDao(ctx context.Context, uAppId uint, uUserId uint64, writer string) *OrderDetailDao {
	object := new(OrderDetailDao)
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
 * @description: 格式转换，OrderInfo => Map
 * @params {info} OrderInfo
 * @return: Map
 */
func OrderDetail2Map(detail OrderDetail) (res map[string]interface{}, err error) {
	j, err := json.Marshal(detail)
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
func (this *OrderDetailDao) InitTableName() {
	this.strTableName = utils.GenOrderDetailTableName(this.uUserId)

}

/**
 * @description:根据orderId查询mysql
 * @params {arrOrderIds} 订单ID集合
 * @return: 订单detail集合
 */

func (this *OrderDetailDao) Get(arrOrderIds []string) (ret []OrderDetail, err error) {
	span, _ := traceutil.Trace(this.ctx, "dao_mysql_order_detail")
	if span != nil {
		defer span.Finish()
		span.Tag("t1", cast.ToString(time.Now().UnixNano()))
	}
	ret = make([]OrderDetail, 0)
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
