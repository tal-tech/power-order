/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-02-01 02:30:41
 * @LastEditors  : lichanglin@tal.com
 * @LastEditTime : 2020-02-16 09:36:59
 * @Description:
 */

package search

import (
	"fmt"
	"github.com/gin-gonic/gin"
	logger "github.com/tal-tech/loggerX"
	"net/http"
	"powerorder/app/constant"
	"powerorder/app/model/tob"
	"powerorder/app/params"
	"powerorder/app/utils"
)

type Search struct {
	//controller.Base
}

/**
 * @description: 查找订单接口(2b)-高级
 * @params {ctx}
 * @return:
 */
func (s Search) Index(ctx *gin.Context) {
	var param params.TobReqSearch
	if err := ctx.ShouldBindJSON(&param); err != nil {
		resp := utils.Error(err)
		ctx.JSON(http.StatusOK, resp)
		return
	}

	if param.End == 0 {
		param.End = constant.MaxOrderCountAtSearchTime
	}
	s.validate(ctx, param)
	var i int
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

	for ext, _ := range param.ExtFilter {
		for i = 0; i < len(param.Fields); i++ {
			if ext == param.Fields[i] {
				break
			}
		}

		if i > len(param.Fields) {
			resp := utils.Error(logger.NewError(fmt.Sprintf("%s filter exist but %s field not", ext, ext)))
			ctx.JSON(http.StatusOK, resp)
			return
		}
	}

	searchModel := tob.NewSearch(utils.TransferToContext(ctx), param.AppId, false)
	result, total, err := searchModel.Get(param)

	if err != nil {
		resp := utils.Error(err)
		ctx.JSON(http.StatusOK, resp)
	} else {
		ctx.JSON(http.StatusOK, utils.Success(map[string]interface{}{"total": total, "result": result}))
	}
}

func (s Search) validate(ctx *gin.Context, param params.TobReqSearch) {
	if param.AppId == 0 {
		resp := utils.Error(logger.NewError("error app_id"))
		ctx.JSON(http.StatusOK, resp)
		return
	}
	if len(param.Fields) == 0 {
		resp := utils.Error(logger.NewError("empty fields"))
		ctx.JSON(http.StatusOK, resp)
		return
	}
	if param.End-param.Start > utils.APP_SEA_MAXSIZE {
		resp := utils.Error(logger.NewError(fmt.Sprintf("max order count at search time is %d", utils.APP_SEA_MAXSIZE)))
		ctx.JSON(http.StatusOK, resp)
		return
	}
}
