/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-04-08 12:49:49
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-04-23 10:25:56
 * @Description:
 */

package orderid

import (
	"context"
	logger "github.com/tal-tech/loggerX"
	"powerorder/app/constant"
	"powerorder/app/dao/mysql"
	"powerorder/app/output"
	"powerorder/app/params"
	"powerorder/app/utils"
)

type Search struct {
	uAppId   uint
	strTable string
	ctx      context.Context
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func NewSearch(ctx context.Context, uAppId uint, strTable string) *Search {
	object := new(Search)
	object.uAppId = uAppId
	object.strTable = strTable
	object.ctx = ctx
	return object
}

/**
 * @description: 根据条件查询订单id
 * @params {params}
 * @return: 订单id列表
 */
func (this *Search) Get(param params.ReqOrderIdSearch) (o output.OrderId, total int64, err error) {

	if ok := this.validate(param); !ok {
		return
	}
	if ok := this.verifyFilterParam(&param); !ok {
		return
	}

	if param.Table == constant.Info {
		ret, err := this.getOrderIdByInfo(param)
		if err != nil {
			return
		}
		o.OrderId = ret
	} else if param.Table == constant.Detail {
		ret, err := this.getOrderIdByDetail(param)
		if err != nil {
			return
		}
		o.OrderId = ret
	} else {
		ret, err := this.getOrderIdByExtension(param)
		if err != nil {
			return
		}
		o.OrderId = ret
	}
	//logger.Ix(this.ctx,"model.toc.search.Get","model_toc_orderid_search", "end, table %s ", param.Table )
	return
}

func (this *Search) getOrderIdByInfo(param params.ReqOrderIdSearch) (ret []mysql.OrderInfoIdx, err error) {
	whereStr, whereValue := utils.And2Where(param.Filter[0])
	ret = make([]mysql.OrderInfoIdx, 0)
	obj := mysql.NewOrderInfoIdxDao(this.ctx, this.uAppId, constant.DBReader)
	err = obj.GetColsWhere(&ret, nil, whereStr, whereValue, 0, int(param.Num))
	return
}

func (this *Search) getOrderIdByDetail(param params.ReqOrderIdSearch) (ret []mysql.OrderDetailIdx, err error) {
	whereStr, whereValue := utils.And2Where(param.Filter[0])
	ret = make([]mysql.OrderDetailIdx, 0)
	obj := mysql.NewOrderDetailIdxDao(this.ctx, this.uAppId, constant.DBReader)
	err = obj.GetColsWhere(&ret, nil, whereStr, whereValue, 0, int(param.Num))
	return
}

func (this *Search) getOrderIdByExtension(param params.ReqOrderIdSearch) (ret []map[string]interface{}, err error) {
	whereStr, whereValue := utils.And2Where(param.Filter[0])
	ret = make([]map[string]interface{}, 0)
	obj := mysql.NewExtensionIdxDao(this.ctx, this.uAppId, param.Table, constant.DBReader)
	err = obj.GetColsWhere(&ret, nil, whereStr, whereValue, 0, int(param.Num))
	for i := 0; i < len(ret); i++ {
		for k, v := range ret[i] {
			switch v.(type) {
			case []byte:
				ret[i][k] = string(v.([]byte))
			default:
			}
		}
	}
	return
}

/**
 * @description:
 * @params {params}
 * @return: 订单id列表
 */
func (this *Search) validate(param params.ReqOrderIdSearch) (ok bool) {

	if _, ok = utils.ExtensionTableShardingInfo[this.uAppId]; !ok {
		logger.Ex(this.ctx, "model.toc.search.Get error", "search error app_id", "")
		return
	}

	if _, ok = utils.ExtensionTableShardingInfo[this.uAppId][this.strTable]; !ok {
		if this.strTable != constant.Info && this.strTable != constant.Detail {
			logger.Ex(this.ctx, "model.toc.search.Get error", "search error app_id", "")
			return
		}
	}

	if _, ok = utils.IdxdbCollectionIndexInfo[this.uAppId]; !ok {
		logger.Ex(this.ctx, "model.toc.search.Get error", "search error app_id", "")
		return
	}

	if _, ok = utils.IdxdbCollectionIndexInfo[this.uAppId][this.strTable]; !ok {
		logger.Ex(this.ctx, "model.toc.search.Get error", "search error field", "")
		return
	}
	ok = true
	return
}

/**
 * @description:
 * @params {params}
 * @return: 订单id列表
 */
func (this *Search) verifyFilterParam(param *params.ReqOrderIdSearch) (ok bool) {
	ok = false
	indexs, _ := utils.IdxdbCollectionIndexInfo[this.uAppId][this.strTable]

	var useridExisted = false //如果出现id的判断，则userid的判断必须有
	var useridEqual = false   // 如果出现id的判断，则userid的判断必须是==
	var index string
	var comparer uint
	var valueType int
	var j = 0
	for i := 0; i < len(param.Filter[0]); i++ {
		//logger.Dx(this.ctx,"model.toc.search.Get","i  ", "%+v, %+v", index, i)

		if len(param.Filter[0][i]) != 3 {
			logger.Ex(this.ctx, "model.toc.search.Get error", "search error filter", "")
			return
		}

		switch param.Filter[0][i][0].(type) {

		case string:
			index = param.Filter[0][i][0].(string)
			break
		default:
			logger.Ex(this.ctx, "model.toc.search.Get error", "search error index type", "")
			return
		}
		logger.Dx(this.ctx, "model.toc.search.Get", "model_toc_orderid_search", "i2, %+v, %+v", index, i)

		switch param.Filter[0][i][1].(type) {
		case float64:
			comparer, _ = utils.JsonInterface2UInt(param.Filter[0][i][1])
			break
		default:
			logger.Ex(this.ctx, "model.toc.search.Get error", "search error index type", "")
			return
		}

		if index == "user_id" {
			useridExisted = true
			if comparer == constant.IntEqual {
				useridEqual = true
			}
		} else if index == "id" {
			if useridExisted == false || useridEqual == false {
				logger.Ex(this.ctx, "model.toc.search.Get error", "search error order", "")
				return
			}
			//如果是id，比较符只能是大于、大于或等于
			if comparer != constant.IntGreater && comparer != constant.IntGreaterOrEqual {
				logger.Ex(this.ctx, "model.toc.search.Get error", "search error index comparer", "")

				return
			}
			//如果是id，则必须是最后一个
			if i != len(param.Filter[0])-1 {
				logger.Ex(this.ctx, "model.toc.search.Get error", "search error index comparer", "")

				return
			}
		}
		logger.Dx(this.ctx, "model.toc.search.Get", "i", "%+v, %+v", index, i)

		if i > 0 && (index == "user_id" || index == "id") {
			if comparer != constant.IntGreater && comparer != constant.IntGreaterOrEqual && comparer != constant.IntEqual {
				logger.Ex(this.ctx, "model.toc.search.Get error", "search error index comparer", "")

				return
			}
			valueType = constant.TypeInt
		} else {

			if comparer != constant.IntLess && comparer != constant.IntLessOrEqual && comparer != constant.IntEqual && comparer != constant.IntGreater && comparer != constant.IntGreaterOrEqual {
				logger.Ex(this.ctx, "model.toc.search.Get error", "search error index comparer", "")

				return
			}
			logger.Dx(this.ctx, "model.toc.search.Get", "index  ", "%+v, %+v", index, utils.IdxdbCollectionFieldInfo[this.uAppId][this.strTable][index])

			if valueType, ok = utils.IdxdbCollectionFieldInfo[this.uAppId][this.strTable][index]; !ok {
				logger.Ex(this.ctx, "model.toc.search.Get error", "search error index", "")

				return
			}
		}

		//var ret bool
		switch param.Filter[0][i][2].(type) {

		case string:
			if valueType == constant.TypeTimestamp2 {
				param.Filter[0][i][2], _ = utils.JsonInterface2Timestamp(param.Filter[0][i][2])
				break
			}
			if valueType != constant.TypeString {
				logger.Ex(this.ctx, "model.toc.search.Get error", "search error index", "")
				return
			}
			break
		case float64:
			if valueType != constant.TypeInt {
				logger.Ex(this.ctx, "model.toc.search.Get error", "search error index", "")
				return
			}
			break
		default:
			return
		}

		if i == 0 {
			for j = 0; j < len(indexs); j++ {
				if index == indexs[j][0] {
					break
				}
			}
			if j == len(indexs) {
				logger.Ex(this.ctx, "model.toc.search.Get error", "search error index", "")
				return
			}

		} else if i == 1 {

			if index != indexs[j][1] {
				logger.Ex(this.ctx, "model.toc.search.Get error", "search error index", "")
				return
			}
		}
	}
	ok = true
	//logger.Ix(this.ctx,"model.toc.search.Get","model_toc_orderid_search", "end, table %s ", param.Table )
	return

}
