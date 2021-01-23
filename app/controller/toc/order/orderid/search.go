/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-04-08 01:14:43
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-04-13 00:54:59
 * @Description:
 */

package orderid

import (
	"github.com/gin-gonic/gin"
	logger "github.com/tal-tech/loggerX"
	"net/http"
	"powerorder/app/model/toc/orderid"
	"powerorder/app/params"
	"powerorder/app/utils"
)

type Search struct {
	//controller.Base
}

/**
 * @description: 根据条件查询订单号
 * @params {ctx}
 * @return:
 */
func (s Search) Index(ctx *gin.Context) {
	var param params.ReqOrderIdSearch
	if err := ctx.ShouldBindJSON(&param); err != nil {
		resp := utils.Error(err)
		ctx.JSON(http.StatusOK, resp)
		return
	}
	s.validate(ctx, param)

	searchModel := orderid.NewSearch(utils.TransferToContext(ctx), param.AppId, param.Table)
	ret, _, err := searchModel.Get(param)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, utils.Success(ret))
}

func (s Search) validate(ctx *gin.Context, param params.ReqOrderIdSearch) {
	if param.AppId == 0 || param.Num == 0 || param.Num > 256 {
		resp := utils.Error(logger.NewError("error app_id or num", logger.PARAM_ERROR))
		ctx.JSON(http.StatusOK, resp)
		return
	}
}
