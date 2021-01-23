/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-21 02:51:14
 * @LastEditors  : lichanglin@tal.com
 * @LastEditTime : 2020-02-05 04:49:23
 * @Description:
 */

package toc

import (
	"context"
	logger "github.com/tal-tech/loggerX"
	"powerorder/app/dao/mysql"
	"powerorder/app/output"
	"powerorder/app/params"
	"powerorder/app/utils"
	"sort"
	"sync"
	"time"
)

type Query struct {
	uAppId  uint
	uUserId uint64
	ctx     context.Context
	lock    sync.RWMutex
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func NewSearchQuery(ctx context.Context, uAppId uint, uUserId uint64) *Query {
	object := new(Query)
	object.uAppId = uAppId
	object.uUserId = uUserId
	object.ctx = ctx
	return object
}

/**
 * @description: 根据条件查询订单详情
 * @params {params}
 * @return: 订单详情列表
 */
func (this *Query) Get(param params.ReqQuery) (o []output.Order, total uint, err error) {

	startDate := time.Unix(param.StartDate, 0).Format("2006-01-02")
	endDate := time.Unix(param.EndDate, 0).Format("2006-01-02")

	var filterInfo *utils.Filter = nil   //info筛选工具
	var filterDetail *utils.Filter = nil //detail筛选工具
	var sorter *utils.Sorter = nil       //排序工具

	var t map[string]output.Order
	t = make(map[string]output.Order)

	filterInfo = nil

	ret := utils.CheckRules(param.Sorter)
	if !ret {
		err = logger.NewError("checkrules failed")
		return
	}

	// 校验向前提
	if len(param.InfoFilter) != 0 {
		filterInfo = utils.NewFilter(param.InfoFilter)
		if filterInfo == nil {
			logger.Ex(this.ctx, "model.toc.query.Get error", "init filterInfo error", "")
			return
		}
	}
	if len(param.DetailFilter) != 0 {
		filterDetail = utils.NewFilter(param.DetailFilter)
		if filterDetail == nil {
			logger.Ex(this.ctx, "model.toc.query.Get error", "init filterDetail error", "")
			return
		}
	}

	getModel := NewGet(this.ctx, this.uAppId, this.uUserId)
	OrderInfo, OrderIds, err := getModel.GetV2([]string{}, param.Fields, startDate, endDate)
	if err != nil {
		logger.Ex(this.ctx, "model.toc.query.Get error", "toc search in db error = %v", err)
		return
	}

	o = make([]output.Order, 0)
	OrderInfoMap := make(map[string]interface{})
	OrderDetailMap := make(map[string]interface{})

	var real bool
	var oi uint = 0
	EligibleOrderInfo := make([]interface{}, 0)

	for j := 0; j < len(OrderIds); j++ {
		OrderId := OrderIds[j]
		if _, ok := OrderInfo[OrderId]; !ok {
			continue
		}
		Value := OrderInfo[OrderId]
		if filterInfo != nil || len(param.Sorter) > 0 {
			OrderInfoMap, err = mysql.OrderInfo2Map(Value.Info)

			if err != nil {
				return
			}
		}
		real = false
		if filterInfo != nil {
			real, err = filterInfo.Execute(OrderInfoMap)

			if err != nil {
				return
			}

			if real == false {
				continue
			}
		}

		if filterDetail != nil {
			var i int
			for i = 0; i < len(Value.Detail); i++ {
				OrderDetailMap, err = mysql.OrderDetail2Map(Value.Detail[i])

				real, err = filterDetail.Execute(OrderDetailMap)

				if err != nil {
					return
				}
				if real == true {
					break
				}
			}

			if i >= len(Value.Detail) {
				continue
			}
		}

		total++
		if len(param.Sorter) > 0 {

			EligibleOrderInfo = append(EligibleOrderInfo, OrderInfoMap)
			t[OrderId] = Value
		} else {
			if oi >= param.Start && oi < param.End {
				o = append(o, Value)
			}
			oi++
		}

	}

	if len(EligibleOrderInfo) == 0 {
		return
	}
	if len(param.Sorter) > 0 {
		sorter, err = utils.NewSorter(&EligibleOrderInfo, param.Sorter)
		if err != nil {
			return
		}
		if sorter == nil {
			logger.Ex(this.ctx, "model.toc.query error", "init sorter error", "")
			return
		}
	}

	if sorter != nil {
		sort.Sort(sorter)
		for i := 0; i < len(*sorter.Slice); i++ {
			if uint(i) >= param.Start && uint(i) < param.End {
				OrderId := (*sorter.Slice)[i].(map[string]interface{})["order_id"].(string)
				o = append(o, t[OrderId])
			}
		}
	}

	return

}
