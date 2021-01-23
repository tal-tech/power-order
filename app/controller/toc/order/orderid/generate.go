/*
 * @Author: lichanglin@tal.com
 * @Date: 2019-12-31 11:07:17
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-04-08 01:15:58
 * @Description:
 */
package orderid

import (
	"github.com/gin-gonic/gin"
	logger "github.com/tal-tech/loggerX"
	"net/http"
	"powerorder/app/constant"
	"powerorder/app/model/toc"
	"powerorder/app/params"
	"powerorder/app/utils"
	"time"
)

type Generate struct {
	//controller.Base
}

/**
 * @description: 生成订单号
 * @params {ctx}
 * @return:
 */
func (g Generate) Index(ctx *gin.Context) {
	var param params.ReqGenerate
	if err := ctx.ShouldBindJSON(&param); err != nil {
		resp := utils.Error(err)
		ctx.JSON(http.StatusOK, resp)
		return
	}

	var OrderIds []string
	now := time.Now()
	for i := 0; i < int(param.Num); i++ {
		OrderIds = append(OrderIds, utils.GenOrderId(param.AppId, now, param.UserId))
	}
	ctx.JSON(http.StatusOK, utils.Success(map[string]interface{}{"order_ids": OrderIds}))

}

func (g Generate) validate(ctx *gin.Context, param params.ReqGenerate) {
	if param.AppId == 0 || param.UserId == 0 || param.Num == 0 || param.Num > constant.MaxBatchOrderCount {
		resp := utils.Error(logger.NewError("error app_id or user_id or num", logger.PARAM_ERROR))
		ctx.JSON(http.StatusOK, resp)
		return
	}
}
