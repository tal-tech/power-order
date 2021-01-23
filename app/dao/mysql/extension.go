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

type ExtensionDao struct {
	Dao
	strExtName string
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func NewExtensionDao2(uAppId uint, uUserId uint64, strExtName string, session *torm.Session, ctx context.Context) *ExtensionDao {
	object := new(ExtensionDao)
	object.InitExtName(strExtName)
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
func NewExtensionDao(ctx context.Context, uAppId uint, uUserId uint64, strExtName string, writer string) *ExtensionDao {
	object := new(ExtensionDao)
	object.InitExtName(strExtName)

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
func (this *ExtensionDao) InitDatabaseName(uAppId uint, uUserId uint64) {
	this.uAppId = uAppId
	this.uUserId = uUserId
	this.strDatabaseName = utils.GenDatabaseName(uUserId)
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *ExtensionDao) InitExtName(strExtName string) {
	this.strExtName = strExtName
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *ExtensionDao) InitTableName() {
	this.strTableName, _ = utils.GenOrderExtensionTableName(this.uUserId, this.uAppId, this.strExtName)
}

/**
 * @description:
 * @params {type}
 * @return:
 */

func (this *ExtensionDao) Get(arrOrderIds []string) (ret []map[string]interface{}, err error) {
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
