/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-30 20:07:48
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-03-31 13:54:01
 * @Description:
 */
package es

import (
	"fmt"
	"net/http"
	"powerorder/app/constant"
	"powerorder/app/output"
	"powerorder/app/params"
)

type ISearch interface {
	Get(param params.TobReqSearch) (o []output.Order, total uint, err error)
}
type Search struct {
	uAppId       uint
	strIndexName string
	header       http.Header
}

/**
 * @description:
 * @params {}
 * @return:
 */
func init() {
	_ = InitEngine()
}

/**
 * @description:
 * @params {from}
 * @return:
 */
func Sort(from []string, keywordFieldMap map[interface{}]int) []map[string]bool {
	var to []map[string]bool
	var field string

	to = make([]map[string]bool, len(from)/2)
	for i := 0; i < len(from); i += 2 {
		if _, ok := keywordFieldMap[from[i]]; !ok {
			field = from[i]
		} else {
			field = fmt.Sprintf("%s.keyword", from[i])
		}

		if from[i+1] == constant.Desc {
			to[i/2] = map[string]bool{field: false}
		} else {
			to[i/2] = map[string]bool{field: true}
		}
	}
	return to
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func NewSearch(uAppId uint, header http.Header) *Search {
	object := new(Search)
	object.header = header
	object.uAppId = uAppId
	return object

}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *Search) InitIndexName(strPostfix string) {
	this.strIndexName = fmt.Sprintf("%s%03d_%s", constant.IndexPrex, this.uAppId, strPostfix)
}
