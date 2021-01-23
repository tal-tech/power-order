/*
 * @Author: lichanglin@tal.com
 * @Date: 2019-12-31 11:07:17
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime : 2020-02-04 14:10:49
 * @Description:
 */
package search

import (
	"github.com/gin-gonic/gin"
	logger "github.com/tal-tech/loggerX"
	"net/http"
	"powerorder/app/constant"
	"powerorder/app/model/tob"
	"powerorder/app/params"
	"powerorder/app/utils"
)

type Get struct {
	//controller.Base
}

/**
 * @description: 根据订单号查询，默认订单号可以解析分库分表的信息
 * @params {ctx}
 * @return:
 */
func (g Get) Index(ctx *gin.Context) {
	var param params.ReqGet
	if err := ctx.ShouldBindJSON(&param); err != nil {
		resp := utils.Error(err)
		ctx.JSON(http.StatusOK, resp)
		return
	}

	g.validate(ctx, param)
	getModel := tob.NewGet(utils.TransferToContext(ctx), param.AppId)
	result, err := getModel.Get(param.OrderIds, param.Fields)

	if err != nil {
		resp := utils.Error(err)
		ctx.JSON(http.StatusOK, resp)
	} else {
		ctx.JSON(http.StatusOK, utils.Success(result))
	}
}

func (g Get) validate(ctx *gin.Context, param params.ReqGet) {
	if param.AppId == 0 {
		resp := utils.Error(logger.NewError("error app_id"))
		ctx.JSON(http.StatusOK, resp)
		return
	}

	if len(param.OrderIds) == 0 {
		resp := utils.Error(logger.NewError("empty order_ids"))
		ctx.JSON(http.StatusOK, resp)
		return
	}

	if len(param.OrderIds) > constant.MaxBatchOrderCount {
		resp := utils.Error(logger.NewError("count of order_ids is limited"))
		ctx.JSON(http.StatusOK, resp)
		return
	}

	if len(param.Fields) == 0 {
		resp := utils.Error(logger.NewError("empty fields"))
		ctx.JSON(http.StatusOK, resp)
		return
	}

}
