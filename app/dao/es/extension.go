/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-31 09:55:04
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-04-03 11:25:34
 * @Description:
 */

package es

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic"
	"net/http"
	"powerorder/app/constant"
	query2 "powerorder/app/dao/es/query"
	"powerorder/app/output"
	"powerorder/app/params"
	"strings"
)

type ExtensionSearch struct {
	Search
}

// 查询ES时，使用keyword类型。
var extensionKeywordFieldMap = map[interface{}]int{
	"order_id": 1,
}

/**
 * @description:
 * @params {uAppId}
 * @params {header}
 * @return:
 */
func NewExtensionSearch(uAppId uint, header http.Header) *ExtensionSearch {
	object := new(ExtensionSearch)

	object.Search = *NewSearch(uAppId, header)
	return object
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *ExtensionSearch) Get(param params.TobReqSearch) (o []output.Order, total uint, err error) {
	var ExtName string
	for i := 0; i < len(param.Fields); i++ {
		if param.Fields[i] == constant.Order_HashSubKey_Detail || param.Fields[i] == constant.Order_HashSubKey_Info {
			continue
		}
		ExtName = param.Fields[i]
		break
	}

	this.InitIndexName(ExtName)

	var Filter [][][]interface{}

	Filter, _ = param.ExtFilter[ExtName]

	rq := elastic.NewBoolQuery()
	r := elastic.NewBoolQuery()

	for i := 0; i < len(Filter); i++ {
		nq := query2.NewQuery(nil)
		for j := 0; j < len(Filter[i]); j++ {
			err = nq.Add(Filter[i][j])
			if err != nil {
				return
			}
		}
		r.Should(nq.Get())

	}
	sortArr := Sort(param.Sorter, extensionKeywordFieldMap)
	r.MinimumNumberShouldMatch(1)
	rq.Filter(r)
	service := EsClient().Search(this.strIndexName)

	for i := 0; i < len(sortArr); i++ {
		for k, v := range sortArr[i] {
			service = service.Sort(k, v)
		}
	}

	ret, err := service.Query(rq).
		From(int(param.Start)).
		Size(int(param.End - param.Start)).
		Headers(this.header).
		Type(ExtName).
		Do(context.Background())

	if err != nil {
		return
	}
	datas := ret.Hits.Hits
	total = uint(ret.TotalHits())
	num := len(ret.Hits.Hits)
	o = make([]output.Order, num)
	for i := 0; i < num; i++ {
		source := datas[i].Source

		data := string(*source)

		data = strings.Replace(data, "_time\": \"0000-00-00 00:00:00\"", "_time\": \"0001-01-01 00:00:00\"", -1)

		o[i].Extensions = make(map[string][]map[string]interface{}, 0)

		o[i].Extensions[ExtName] = make([]map[string]interface{}, 1)
		err = json.Unmarshal([]byte(data), &(o[i].Extensions[ExtName][0]))

		if err != nil {
			return
		}
	}

	return
}
