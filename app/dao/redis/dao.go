/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-11 22:32:41
 * @LastEditors  : lichanglin@tal.com
 * @LastEditTime : 2020-01-22 11:32:24
 * @Description:
 */
package redis

type Dao struct {
	uAppId   uint
	instance string
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func NewDao(uAppId uint, instance string) *Dao {

	object := new(Dao)
	object.uAppId = uAppId
	object.instance = instance

	return object
}
