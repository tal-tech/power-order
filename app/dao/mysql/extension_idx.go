/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-01 16:45:57
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-02-23 01:44:12
 * @Description:
 */
package mysql

import (
	"context"
	"github.com/tal-tech/torm"
	"powerorder/app/utils"
)

type ExtensionIdxDao struct {
	Dao
	strExtName string
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func NewExtensionIdxDao2(uAppId uint, strExtName string, session *torm.Session, ctx context.Context) *ExtensionIdxDao {
	object := new(ExtensionIdxDao)
	object.InitExtName(strExtName)
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
func NewExtensionIdxDao(ctx context.Context, uAppId uint, strExtName string, writer string) *ExtensionIdxDao {
	object := new(ExtensionIdxDao)
	object.InitExtName(strExtName)

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
func (this *ExtensionIdxDao) InitDatabaseName(uAppId uint, id uint64) {
	this.uAppId = uAppId
	this.strDatabaseName = utils.GenDatabaseIdxName()
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *ExtensionIdxDao) InitExtName(strExtName string) {
	this.strExtName = strExtName
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *ExtensionIdxDao) InitTableName() {
	this.strTableName, _ = utils.GenOrderExtensionIdxTableName(this.uAppId, this.strExtName)
}

/**
 * @description:
 * @params {type}
 * @return:
 */

func (this *ExtensionIdxDao) Get(arrOrderIds []string) (ret []map[string]interface{}, err error) {
	ret = make([]map[string]interface{}, 0)
	this.InitSession()
	this.InitTableName()
	this.SetTable(this.strTableName)

	if len(arrOrderIds) > 0 {
		this.BuildQuery(torm.CastToParamIn(arrOrderIds), "order_id")
	}

	err = this.Session.Find(&ret)

	for i := 0; i < len(ret); i++ {
		for k, v := range ret[i] {
			switch v.(type) {
			case []byte:
				ret[i][k] = string(v.([]byte))
			default:
			}
		}
	}
	return
}

/**
 * @description:
 * @params {type}
 * @return:
 */

func (this *ExtensionIdxDao) IsExists(userId uint64, Id int64) (ret bool, err error) {
	this.InitSession()
	this.InitTableName()
	this.SetTable(this.strTableName)
	this.Where("user_id = ? and id = ?", userId, Id)
	ret, err = this.Session.Exist()
	return
}

/**
 * @description:
 * @params {type}
 * @return:
 */

func (this *ExtensionIdxDao) GetList(cols []string, query interface{}, args []interface{}, start, limit int) (ret []map[string]interface{}, err error) {
	ret = make([]map[string]interface{}, 0)
	this.InitSession()
	this.InitTableName()
	this.SetTable(this.strTableName)
	this.Where(query, args)
	err = this.Session.Limit(0, limit).Find(&ret)
	for i := 0; i < len(ret); i++ {
		for k, v := range ret[i] {
			switch v.(type) {
			case []byte:
				ret[i][k] = string(v.([]byte))
			default:
			}
		}
	}
	return
}
