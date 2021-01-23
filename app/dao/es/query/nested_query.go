/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-02-01 01:32:06
 * @LastEditors  : lichanglin@tal.com
 * @LastEditTime : 2020-02-04 01:41:30
 * @Description:
 */

package query

import (
	"fmt"
	"github.com/olivere/elastic"
)

type NestedQuery struct {
	Query
	strKeyPrex string
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func NewNestedQuery(strKeyPrex string, query *elastic.BoolQuery) *NestedQuery {
	object := new(NestedQuery)
	object.Query = *NewQuery(query)
	object.genKey = *NewGenKey(object)

	object.strKeyPrex = strKeyPrex
	return object
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *NestedQuery) key(key string) string {

	return fmt.Sprintf("%s.%s", this.strKeyPrex, key)
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *NestedQuery) Get() elastic.Query {
	return elastic.NewNestedQuery(this.strKeyPrex, this.query)
}
