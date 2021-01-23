/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-21 02:51:14
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime : 2020-02-16 09:55:33
 * @Description:
 */

package tob

import (
	"context"
	"github.com/pkg/errors"
	"net/http"
	"powerorder/app/constant"
	"powerorder/app/dao/es"
	"powerorder/app/output"
	"powerorder/app/params"
	"powerorder/app/utils"
)

type Search struct {
	uAppId      uint
	ctx         context.Context
	dissociaton bool
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func NewSearch(ctx context.Context, uAppId uint, dissociaton bool) *Search {
	object := new(Search)
	object.uAppId = uAppId
	object.ctx = ctx
	object.dissociaton = dissociaton
	return object
}

/**
 * @description:
 * @params {params}
 * @return:
 */
func (this *Search) Get(param params.TobReqSearch) (o []output.Order, total uint, err error) {
	ret := utils.CheckRules(param.Sorter)
	if !ret {
		err = errors.New("checkrules failed")
		return
	}

	header := http.Header{}
	if this.dissociaton == false {
		s := es.NewOrderInfoSearch(this.uAppId, header)
		o, total, err = s.Get(param)
	} else if constant.Detail == param.Fields[0] {
		s := es.NewOrderDetailSearch(this.uAppId, header)
		o, total, err = s.Get(param)
	} else {
		s := es.NewExtensionSearch(this.uAppId, header)
		o, total, err = s.Get(param)
	}
	return
}
