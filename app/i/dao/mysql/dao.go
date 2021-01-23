/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-08 16:44:48
 * @LastEditors  : lichanglin@tal.com
 * @LastEditTime : 2020-01-10 15:01:05
 * @Description:
 */
package mysql

type Dao interface {
	Create(bean interface{}) (int64, error)
	InitDatabaseName(uAppId uint, uUserId uint64)
	InitTableName()
	SetTable(strTableName string)
	GetTable() string
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func Create(instance Dao, bean interface{}) (int64, error) {
	instance.InitTableName()
	instance.SetTable(instance.GetTable())

	return instance.Create(bean)
}
