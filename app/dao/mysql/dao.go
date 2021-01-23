/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-03 18:06:07
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-04-11 09:43:02
 * @Description:
 */
package mysql

import (
	"context"
	"github.com/spf13/cast"
	logger "github.com/tal-tech/loggerX"
	"github.com/tal-tech/torm"
	"github.com/tal-tech/xtools/traceutil"
	"powerorder/app/utils"
	"strconv"
	"time"
)

var TiDefaultTimestamp utils.Time

func init() {
	formatTime, _ := time.Parse("2006-01-02 15:04:05", "2001-01-01 00:00:00")

	TiDefaultTimestamp = utils.Time(formatTime)
}

type Dao struct {
	torm.DbBaseDao
	strDatabaseName string
	strTableName    string
	uUserId         uint64
	uAppId          uint
	ctx             context.Context
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func NewSession(ctx context.Context, uAppId uint, uUserId uint64, writer string) *torm.Session {
	// 如果是压测 走子库连接
	writer = utils.DbShadowHandler(ctx, writer)
	object := new(Dao)
	object.InitDatabaseName(uAppId, uUserId)
	if ins := torm.GetDbInstance(object.strDatabaseName, writer); ins != nil {
		return ins.GetSession()
	} else {
		return nil
	}
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func NewSessionIdx(ctx context.Context, uAppId uint, uUserId uint64, writer string) *torm.Session {
	// 如果是压测 走影子库连接
	writer = utils.DbShadowHandler(ctx, writer)
	object := new(Dao)
	object.uAppId = uAppId
	object.uUserId = uUserId
	object.strDatabaseName = utils.GenDatabaseIdxName()
	if ins := torm.GetDbInstance(object.strDatabaseName, writer); ins != nil {
		return ins.GetSession()
	} else {
		return nil
	}
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *Dao) InitDatabaseName(uAppId uint, uUserId uint64) {
	this.uAppId = uAppId
	this.uUserId = uUserId

	this.strDatabaseName = utils.GenDatabaseName(uUserId)
	//fmt.Printf("InitDatabaseName user_id = %d, databasename = %s", uUserId, this.strDatabaseName)
	logger.Dx(this.ctx, "dao.mysql.dao.InitDatabaseName", "InitDatabaseName", "user_id = %d, databasename = %s", uUserId, this.strDatabaseName)
	//fmt.Printf()
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *Dao) InitCtx(ctx context.Context) {
	this.ctx = ctx
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *Dao) InitTableName() {

}

/**
 * @description:
 * @params {type}
 * @return:
 */

func (this *Dao) GetTable() string {
	return this.strTableName
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *Dao) UpdateColsWhere(bean interface{}, cols []string, query interface{}, args []interface{}) (int64, error) {
	if this.Session == nil {
		return this.Engine.Cols(cols...).Where(query, args...).Update(bean)
	} else {
		return this.Session.Cols(cols...).Where(query, args...).Update(bean)
	}
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *Dao) Where(query interface{}, args ...interface{}) *torm.Session {
	if this.Session == nil {
		return this.Engine.Where(query, args...)
	} else {
		return this.Session.Where(query, args...)
	}
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *Dao) SetTableName() {
	this.SetTable(this.GetTable())
}

/**
 * @description:
 * @params {type}
 * @return:
 */

func (this *Dao) Create(bean interface{}) (int64, error) {
	n := time.Now().UnixNano()
	ns := strconv.FormatInt(n, 10)
	span, _ := traceutil.Trace(this.ctx, "create_"+ns)
	if span != nil {
		//节点参数注入 可以从链路追踪界面查看节点数据
		span.Tag("db", this.strDatabaseName)
		span.Tag("table", this.strTableName)
		span.Tag("timeStamp1", cast.ToString(time.Now().UnixNano()))
		//切记要回收span
		defer span.Finish()
	}
	r1, r2 := this.DbBaseDao.Create(bean)

	//span2, _ := traceutil.Trace(this.ctx, "create2_" + ns)
	if span != nil {
		//节点参数注入 可以从链路追踪界面查看节点数据
		span.Tag("timeStamp2", cast.ToString(time.Now().UnixNano()))
	}

	return r1, r2
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *Dao) GetColsWhere(bean interface{}, cols []string, query interface{}, args []interface{}, start, limit int) error {
	if this.Session == nil {
		return this.Engine.Cols(cols...).Where(query, args...).Limit(limit, start).Find(bean)
	} else {
		return this.Session.Cols(cols...).Where(query, args...).Limit(limit, start).Find(bean)
	}
}
