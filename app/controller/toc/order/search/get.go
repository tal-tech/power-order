/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-02-11 23:01:18
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime : 2020-02-12 21:31:19
 * @Description:
 */

package search

import (
	"github.com/gin-gonic/gin"
	logger "github.com/tal-tech/loggerX"
	"net/http"
	"powerorder/app/constant"
	"powerorder/app/model/toc"
	"powerorder/app/output"
	"powerorder/app/params"
	"powerorder/app/utils"
)

type Get struct {
	//controller.Base
}

/**
 * @description:查找订单接口(2c)-获得订单详情
 * @params {ctx}
 * @return: 获得订单详情列表
 */
func (g Get) Index(ctx *gin.Context) {
	var param params.ReqBGet

	if err := ctx.ShouldBindJSON(&param); err != nil {
		resp := utils.Error(err)
		ctx.JSON(http.StatusOK, resp)
		return
	}
	g.validate(ctx, param)

	out, err := g.getOrders(ctx, param)
	if err != nil {
		resp := utils.Error(err)
		ctx.JSON(http.StatusOK, resp)
		return
	}
	ctx.JSON(http.StatusOK, utils.Success(out))
}

func (g Get) getOrders(ctx *gin.Context, param params.ReqBGet) (ret map[string]output.Order, err error) {
	ret = make(map[string]output.Order, 0)
	for i := 0; i < len(param.Get); i++ {
		order := toc.NewGet(utils.TransferToContext(ctx), param.AppId, param.Get[i].UserId)
		result, _, err := order.Get(param.Get[i].OrderIds, param.Fields)
		if err != nil {
			return
		}
		for OrderId, OrderInfo := range result {
			ret[OrderId] = OrderInfo
		}
	}
	return
}

func (g Get) validate(ctx *gin.Context, param params.ReqBGet) {
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

	if len(param.Get) == 0 {
		resp := utils.Error(logger.NewError("empty get"))
		ctx.JSON(http.StatusOK, resp)
		return
	}

	if len(param.Get) > constant.MaxBatchOrderCount {
		resp := utils.Error(logger.NewError("count of get is limited"))
		ctx.JSON(http.StatusOK, resp)
		return
	}

	for i := 0; i < len(param.Get); i++ {
		if len(param.Get[i].OrderIds) > constant.MaxBatchOrderCount {
			resp := utils.Error(logger.NewError("count of order_id is limited"))
			ctx.JSON(http.StatusOK, resp)
			return
		}
	}

}
