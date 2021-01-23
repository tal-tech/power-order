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
	"sync/atomic"
)

type Search struct {
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
func NewSearch(ctx context.Context, uAppId uint, uUserId uint64) *Search {
	object := new(Search)
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
func (this *Search) Get(param params.ReqSearch) (o []output.Order, total uint, err error) {
	var filterInfo *utils.Filter = nil   //info筛选工具
	var filterDetail *utils.Filter = nil //detail筛选工具
	var sorter *utils.Sorter = nil       //排序工具
	//logger.Dx(this.ctx,"model.toc.search.Get","model toc/search params = %v", param)

	t := make(map[string]output.Order, 0)

	ret := utils.CheckRules(param.Sorter)
	if !ret {
		err = logger.NewError("checkrules failed")
		return
	}
	// 校验向前提
	if len(param.InfoFilter) != 0 {
		filterInfo = utils.NewFilter(param.InfoFilter)
		if filterInfo == nil {
			logger.Ex(this.ctx, "model.toc.search.Get error", "init filterInfo error", "")
			return
		}
	}
	if len(param.DetailFilter) != 0 {
		filterDetail = utils.NewFilter(param.DetailFilter)
		if filterDetail == nil {
			logger.Ex(this.ctx, "model.toc.search.Get error", "init filterDetail error", "")
			return
		}
	}

	getModel := NewGet(this.ctx, this.uAppId, this.uUserId)
	OrderInfo, OrderIds, err := getModel.Get([]string{}, param.Fields)
	if err != nil {
		//logger.Ex(this.ctx, "model.toc.search.Get error","toc search in db error = %v", err)
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
			logger.Ex(this.ctx, "model.toc.search error", "init sorter error", "")
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

/**
 * @description: 根据条件查询订单详情 协程版本
 * @params {params}
 * @return: 订单详情列表
 */
func (this *Search) GetV2(param params.ReqSearch) (o []output.Order, total uint32, err error) {

	var filterInfo *utils.Filter = nil   //info筛选工具
	var filterDetail *utils.Filter = nil //detail筛选工具
	var sorter *utils.Sorter = nil       //排序工具
	//logger.Dx(this.ctx,"model.toc.search.GetV2","model toc/search params = %v", param)

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
			//logger.Ex(this.ctx,"model.toc.search.GetV2 error","init filterInfo error","")
			return
		}
	}
	if len(param.DetailFilter) != 0 {
		filterDetail = utils.NewFilter(param.DetailFilter)
		if filterDetail == nil {
			//logger.Ex(this.ctx,"model.toc.search.GetV2 error","init filterDetail error","")
			return
		}
	}

	getModel := NewGet(this.ctx, this.uAppId, this.uUserId)
	OrderInfo, OrderIds, err := getModel.Get([]string{}, param.Fields)
	if err != nil {
		//logger.Ex(this.ctx, "model.toc.search.GetV2 error","toc search in db error = %v", err)
		return
	}

	o = make([]output.Order, 0)

	EligibleOrderInfo := make([]interface{}, 0)

	tmpOrderMap := make(map[string]output.Order, 0)
	tmpEligibleOrderInfoMap := make(map[string]interface{}, 0)

	defer utils.Catch()
	// 协程等待
	gp := utils.NewGPool(100)
	for i := 0; i < len(OrderIds); i++ {
		OrderId := OrderIds[i]
		Value := OrderInfo[OrderId]
		gp.Add(1)
		go this.getOrderInfo(gp, OrderId, Value, param, filterInfo, filterDetail, &tmpOrderMap, &total, &tmpEligibleOrderInfoMap)
	}
	gp.Wait()

	if len(tmpEligibleOrderInfoMap) == 0 {
		return
	}
	if len(param.Sorter) > 0 {
		for i := 0; i < len(OrderIds); i++ {
			OrderId := OrderIds[i]
			if _, ok := tmpEligibleOrderInfoMap[OrderId]; !ok {
				continue
			}
			EligibleOrderInfo = append(EligibleOrderInfo, tmpEligibleOrderInfoMap[OrderId])
		}
		sorter, err = utils.NewSorter(&EligibleOrderInfo, param.Sorter)
		if err != nil {
			return
		}
		if sorter == nil {
			logger.Ex(this.ctx, "model.toc.search.GetV2 error", "init sorter error", "")
			return
		}
	} else {
		for i := 0; i < len(OrderIds); i++ {
			OrderId := OrderIds[i]
			if _, ok := tmpOrderMap[OrderId]; !ok {
				continue
			}
			o = append(o, tmpOrderMap[OrderId])
		}
	}

	if sorter != nil {
		sort.Sort(sorter)
		for i := 0; i < len(*sorter.Slice); i++ {
			if uint(i) >= param.Start && uint(i) < param.End {
				OrderId := (*sorter.Slice)[i].(map[string]interface{})["order_id"].(string)
				o = append(o, tmpOrderMap[OrderId])
			}
		}
	}

	return

}

/**
 * @description: 根据条件查询订单详情  对total的操作需要是原子的
 * @params {params}
 * @return: 订单详情列表
 */
func (this *Search) getOrderInfo(gp *utils.Gpool, orderId string, order output.Order, param params.ReqSearch, filterInfo, filterDetail *utils.Filter, o *map[string]output.Order, total *uint32, EligibleOrderInfo *map[string]interface{}) {

	defer gp.Done()

	var err error
	var real bool

	OrderInfoMap := make(map[string]interface{})
	OrderDetailMap := make(map[string]interface{})

	if filterInfo != nil || len(param.Sorter) > 0 {
		OrderInfoMap, err = mysql.OrderInfo2Map(order.Info)

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
			return
		}
	}

	if filterDetail != nil {
		var i int
		for i = 0; i < len(order.Detail); i++ {
			OrderDetailMap, err = mysql.OrderDetail2Map(order.Detail[i])

			real, err = filterDetail.Execute(OrderDetailMap)

			if err != nil {
				return
			}
			if real == true {
				break
			}
		}

		if i >= len(order.Detail) {
			return
		}
	}
	// 多协程修改map数据，增加原子锁
	atomic.AddUint32(total, 1)
	this.lock.Lock()
	defer this.lock.Unlock()
	if len(param.Sorter) > 0 {
		(*EligibleOrderInfo)[orderId] = OrderInfoMap
		(*o)[orderId] = order
	} else {
		if uint(*total) >= param.Start && uint(*total) < param.End {
			(*o)[orderId] = order
		}
	}

}
