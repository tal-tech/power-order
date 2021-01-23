/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-02-02 18:27:11
 * @LastEditors  : lichanglin@tal.com
 * @LastEditTime : 2020-02-02 20:37:24
 * @Description:
 */

package query

type IQuery interface {
	key(key string) string
}
type GenKey struct {
	query IQuery
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func NewGenKey(query IQuery) *GenKey {
	object := new(GenKey)
	object.query = query
	return object
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *GenKey) key(key string) string {
	return this.query.key(key)
}
