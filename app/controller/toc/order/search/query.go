/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-21 02:46:03
 * @LastEditors  : lichanglin@tal.com
 * @LastEditTime : 2020-02-05 23:19:32
 * @Description:
 */

package search

import (
	"fmt"
	"github.com/gin-gonic/gin"
	logger "github.com/tal-tech/loggerX"
	"net/http"
	"powerorder/app/constant"
	"powerorder/app/model/toc"
	"powerorder/app/params"
	"powerorder/app/utils"
)

type Query struct {
	//controller.Base
}

/**
 * @description:查找订单接口(2c)-筛选订单
 * @params {ctx}
 * @return: 获得订单详情列表
 */
func (q Query) Index(ctx *gin.Context) {
	var param params.ReqQuery
	if err := ctx.ShouldBindJSON(&param); err != nil {
		resp := utils.Error(err)
		ctx.JSON(http.StatusOK, resp)
		return
	}
	if param.End == 0 {
		param.End = constant.MaxOrderCountAtSearchTime
	}
	q.validate(ctx, param)

	search := toc.NewSearchQuery(utils.TransferToContext(ctx), param.AppId, param.UserId)
	result, total, err := search.Get(param)
	//result, total, err = search.GetV2(params)
	if err != nil {
		resp := utils.Error(err)
		ctx.JSON(http.StatusOK, resp)
		return
	}
	ctx.JSON(http.StatusOK, utils.Success(map[string]interface{}{"total": total, "result": result}))
	return
}

func (q Query) validate(ctx *gin.Context, param params.ReqQuery) {
	if param.AppId == 0 || param.UserId == 0 {
		resp := utils.Error(logger.NewError("error app_id or user_id"))
		ctx.JSON(http.StatusOK, resp)
		return
	}
	if len(param.Fields) == 0 {
		resp := utils.Error(logger.NewError("empty fields"))
		ctx.JSON(http.StatusOK, resp)
		return
	}
	if param.StartDate == 0 || param.EndDate == 0 {
		resp := utils.Error(logger.NewError("empty start_date or end_date"))
		ctx.JSON(http.StatusOK, resp)
		return
	}
	difference := param.EndDate - param.StartDate
	if difference > utils.TOC_QUERY_MAX_DAY || difference < utils.TOC_QUERY_MIN_DAY {
		resp := utils.Error(logger.NewError("out of range start_date or end_date"))
		ctx.JSON(http.StatusOK, resp)
		return
	}
	if param.End-param.Start > utils.APP_SEA_MAXSIZE {
		resp := utils.Error(logger.NewError(fmt.Sprintf("max order count at search time is %d", utils.APP_SEA_MAXSIZE)))
		ctx.JSON(http.StatusOK, resp)
		return
	}

	i := 0
	if len(param.DetailFilter) > 0 {
		for i = 0; i < len(param.Fields); i++ {
			if constant.OrderDetail == param.Fields[i] {
				break
			}
		}
		if i > len(param.Fields) {
			resp := utils.Error(logger.NewError("detail filter exist but detail field not"))
			ctx.JSON(http.StatusOK, resp)
			return
		}
	}

	if len(param.InfoFilter) > 0 {
		for i = 0; i < len(param.Fields); i++ {
			if constant.OrderInfo == param.Fields[i] {
				break
			}
		}

		if i > len(param.Fields) {
			resp := utils.Error(logger.NewError("info filter exist but info field not"))
			ctx.JSON(http.StatusOK, resp)
			return
		}
	}
}
