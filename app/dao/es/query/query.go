/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-02-01 01:32:06
 * @LastEditors  : lichanglin@tal.com
 * @LastEditTime : 2020-02-03 23:03:43
 * @Description:
 */

package query

import (
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
	"powerorder/app/constant"
	"powerorder/app/utils"
)

type Query struct {
	genKey GenKey
	query  *elastic.BoolQuery
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func NewQuery(query *elastic.BoolQuery) *Query {
	object := new(Query)
	if query == nil {
		query = elastic.NewBoolQuery()
	}
	object.query = query
	object.genKey = *NewGenKey(object)
	return object
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *Query) key(key string) string {
	return key
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *Query) Get() elastic.Query {
	return this.query
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *Query) Add(rule []interface{}) (err error) {

	err = nil
	if len(rule) != 3 {
		return
	}

	var key string
	var comparer uint
	var field interface{}
	var fields []interface{}
	var ok bool
	var ret bool

	key, ret = utils.JsonInterface2String(rule[0])
	if ret == false {
		return
	}
	comparer, ret = utils.JsonInterface2UInt(rule[1])

	if err != nil {
		err = errors.New("error JsonInterface2String")
		return
	}

	field = rule[2]
	key = this.genKey.key(key)

	switch comparer {
	case constant.IntGreater:
		this.query = this.query.Must(elastic.NewRangeQuery(key).Gt(field))
		break
	case constant.IntGreaterOrEqual:
		this.query = this.query.Must(elastic.NewRangeQuery(key).Gte(field))
		break
	case constant.IntLess:
		this.query = this.query.Must(elastic.NewRangeQuery(key).Lt(field))
		break
	case constant.IntLessOrEqual:
		this.query = this.query.Must(elastic.NewRangeQuery(key).Lte(field))
		break
	case constant.IntEqual:
		this.query = this.query.Must(elastic.NewTermQuery(key, field))
		break
	case constant.IntNotEqual:
		this.query = this.query.MustNot(elastic.NewTermQuery(key, field))
		break
	case constant.IntNotWithIn:

		if fields, ok = field.([]interface{}); !ok {
			err = errors.New("error fields")
			return
		}

		this.query = this.query.MustNot(elastic.NewTermsQuery(key, fields...))
		break
	case constant.IntWithIn:
		if fields, ok = field.([]interface{}); !ok {
			err = errors.New("error fields")
			return
		}

		this.query = this.query.Must(elastic.NewTermsQuery(key, fields...))
		break
	default:
		err = errors.New("error comparer")
	}

	return nil
}
