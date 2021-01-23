/*
 * @Author: lichanglin@tal.com
 * @Date: 2019-12-31 11:07:17
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-02-27 23:17:38
 * @Description:
 */
package addition

import (
	"github.com/gin-gonic/gin"
	logger "github.com/tal-tech/loggerX"
	"net/http"
	"powerorder/app/constant"
	"powerorder/app/model/addition"
	"powerorder/app/params"
	"powerorder/app/utils"
)

type Rollback struct {
	//controller.Base
}

/**
 * @description:  回滚订单
 * @params {type}
 * @return:
 */
func (r Rollback) Index(ctx *gin.Context) {
	var param params.ReqRollback
	if err := ctx.ShouldBindJSON(&param); err != nil {
		resp := utils.Error(err)
		ctx.JSON(http.StatusOK, resp)
		return
	}
	if param.TxStatus == constant.TxStatusUnCommitted {
		param.TxStatus = 2 // rollback 接口 默认 tx_status为2
	}

	r.validate(ctx, param)

	updateModel := addition.NewUpdate(ctx, param.AppId)
	if err := updateModel.UpdateStatus(param); err != nil {
		resp := utils.Error(err)
		ctx.JSON(http.StatusOK, resp)
		return
	}
	ctx.JSON(http.StatusOK, utils.Success(nil))
}

func (r Rollback) validate(ctx *gin.Context, param params.ReqRollback) {
	if param.AppId == 0 || param.UserId == 0 || len(param.TxId) == 0 {
		resp := utils.Error(logger.NewError("error app_id or error user_id"))
		ctx.JSON(http.StatusOK, resp)
		return
	}

	if param.TxStatus == constant.TxStatusCommitted || param.TxStatus == constant.TxStatusUnCommitted {
		resp := utils.Error(logger.NewError("error tx_status"))
		ctx.JSON(http.StatusOK, resp)
		return
	}
}
