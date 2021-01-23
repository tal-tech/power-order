/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-31 09:54:59
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-04-03 11:25:20
 * @Description:
 */
package es

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic"
	"net/http"
	"powerorder/app/constant"
	"powerorder/app/dao/mysql"
	"powerorder/app/output"
	"powerorder/app/params"
	"strings"
)

type OrderDetailSearch struct {
	Search
}

// 查询ES时，使用keyword类型。
var detailKeywordFieldMap = map[interface{}]int{
	"order_id": 1,
}

/**
 * @description:
 * @params {uAppId}
 * @params {header}
 * @return:
 */
func NewOrderDetailSearch(uAppId uint, header http.Header) *OrderDetailSearch {
	object := new(OrderDetailSearch)
	object.Search = *NewSearch(uAppId, header)
	object.InitIndexName(constant.Detail)
	return object
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *OrderDetailSearch) Get(param params.TobReqSearch) (o []output.Order, total uint, err error) {
	bq := elastic.NewBoolQuery()

	if err := GenOrderDetailQuery(param, bq); err != nil {
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
		From(int(param.Start)).
		Size(int(param.End - param.Start)).
		Headers(this.header).
		Type(constant.OrderDetailType).
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

		o[i].Detail = make([]mysql.OrderDetail, 1)
		err = json.Unmarshal([]byte(data), &(o[i].Detail[0]))

		//解析details
		if err != nil {
			return
		}

	}

	return
}
