/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-31 09:54:50
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-04-03 11:24:51
 * @Description:
 */

package es

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/olivere/elastic"
	"net/http"
	"powerorder/app/constant"
	"powerorder/app/dao/es/query"
	"powerorder/app/output"
	"powerorder/app/params"
	"strings"
)

type OrderInfoSearch struct {
	Search
}

// 查询ES时，使用keyword类型。
var keywordFieldMap = map[interface{}]int{
	"parent_order_id": 1,
	"order_id":        1,
}

/**
 * @description:
 * @params {uAppId}
 * @params {header}
 * @return:
 */
func NewOrderInfoSearch(uAppId uint, header http.Header) *OrderInfoSearch {
	object := new(OrderInfoSearch)
	object.Search = *NewSearch(uAppId, header)
	object.InitIndexName(constant.Info)
	return object
}

/**
 * @description:
 * @params {params}
 * @return:
 */
func (this *OrderInfoSearch) Get(param params.TobReqSearch) (o []output.Order, total uint, err error) {
	bq := elastic.NewBoolQuery()

	if err := GenOrderInfoQuery(param, bq); err != nil {
		return
	}
	if err := GenOrderDetailQuery(param, bq); err != nil {
		return
	}
	if err := GenExtensionQuery(param, bq); err != nil {
		return
	}

	sortArr := Sort(param.Sorter, extensionKeywordFieldMap)

	service := EsClient().Search(this.strIndexName)

	for i := 0; i < len(sortArr); i++ {
		for k, v := range sortArr[i] {
			service = service.Sort(k, v)
		}
	}

	ret, err := service.Query(bq).
		RequestCache(true).
		From(int(param.Start)).
		Size(int(param.End - param.Start)).
		Headers(this.header).
		Type(constant.OrderInfoType).
		Do(context.Background())

	if err != nil {
		return
	}

	ExtCount := 0
	for i := 0; i < len(param.Fields); i++ {
		if param.Fields[i] == constant.OrderInfo || param.Fields[i] == constant.OrderDetail {
		} else {
			ExtCount++
		}
	}

	total = uint(ret.TotalHits())
	datas := ret.Hits.Hits
	num := len(ret.Hits.Hits)

	o = make([]output.Order, num)
	for i := 0; i < num; i++ {
		source := datas[i].Source
		data := string(*source)

		data = strings.Replace(data, "_time\": \"0000-00-00 00:00:00\"", "_time\": \"0001-01-01 00:00:00\"", -1)

		for _, field := range param.Fields {
			if field == constant.Detail {
				err = json.Unmarshal([]byte(data), &(o[i]))
				//解析details
				if err != nil {
					return
				}
				break
			}
		}

		for _, field := range param.Fields {
			if field == constant.Info {
				err = json.Unmarshal([]byte(data), &(o[i].Info))
				//解析info
				if err != nil {
					return
				}
				break
			}
		}

		if ExtCount > 0 {
			var res *simplejson.Json
			res, err = simplejson.NewJson([]byte(*source))

			if err != nil {
				return
			}
			o[i].Extensions = make(map[string][]map[string]interface{}, 0)

			for _, field := range param.Fields {

				if field == constant.OrderDetail || field == constant.OrderInfo {
					continue
				}
				o[i].Extensions[field] = make([]map[string]interface{}, 0)
				rows, _ := res.Get(field).Array()

				for _, row := range rows {
					//对每个row获取其类型，每个row相当于 C++/Golang 中的map、Python中的dict
					//每个row对应一个map，该map类型为map[string]interface{}，也即key为string类型，value是interface{}类型
					if slice, ok := row.(map[string]interface{}); ok {
						o[i].Extensions[field] = append(o[i].Extensions[field], slice)
					}
				}

			}

		}
	}

	return
}

func GenOrderInfoQuery(param params.TobReqSearch, bq *elastic.BoolQuery) (err error) {
	if len(param.InfoFilter) < 1 {
		return
	}
	_bq := elastic.NewBoolQuery()
	for i := 0; i < len(param.InfoFilter); i++ {
		q := query.NewQuery(nil)
		for j := 0; j < len(param.InfoFilter[i]); j++ {
			field := param.InfoFilter[i][j][0]
			if _, ok := keywordFieldMap[field]; ok {
				param.InfoFilter[i][j][0] = fmt.Sprintf("%s.keyword", param.InfoFilter[i][j][0])
				err = q.Add(param.InfoFilter[i][j])
			} else {
				err = q.Add(param.InfoFilter[i][j])
			}
			if err != nil {
				return
			}
		}
		_bq.Should(q.Get())
	}
	_bq.MinimumNumberShouldMatch(1)
	bq.Filter(_bq)
	return
}

func GenOrderDetailQuery(param params.TobReqSearch, bq *elastic.BoolQuery) (err error) {
	if len(param.DetailFilter) < 1 {
		return
	}
	_bq := elastic.NewBoolQuery()
	for i := 0; i < len(param.DetailFilter); i++ {
		nq := query.NewNestedQuery(constant.Order_HashSubKey_Detail, nil)
		for j := 0; j < len(param.DetailFilter[i]); j++ {
			err = nq.Add(param.DetailFilter[i][j])
			if err != nil {
				return
			}
		}
		_bq.Should(nq.Get())
	}
	_bq.MinimumNumberShouldMatch(1)
	bq.Filter(_bq)
	return
}

func GenExtensionQuery(param params.TobReqSearch, bq *elastic.BoolQuery) (err error) {
	if len(param.ExtFilter) < 1 {
		return
	}
	for extension, Filter := range param.ExtFilter {
		_bq := elastic.NewBoolQuery()
		for i := 0; i < len(Filter); i++ {
			nq := query.NewNestedQuery(extension, nil)
			for j := 0; j < len(Filter[i]); j++ {
				err = nq.Add(Filter[i][j])
				if err != nil {
					return
				}
			}
			_bq.Should(nq.Get())
		}
		_bq.MinimumNumberShouldMatch(1)
		bq.Filter(_bq)
	}
	return
}
