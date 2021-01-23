/*
 * @Author: lichanglin@tal.com
 * @Date: 2019-12-31 11:07:17
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-04-11 00:54:55
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

type Commit struct {
	//Rollback
}

/**
 * @description: 创建订单-提交
 * @params {type}
 * @return:
 */
func (c Commit) Index(ctx *gin.Context) {
	var param params.ReqRollback

	if err := ctx.ShouldBindJSON(&param); err != nil {
		resp := utils.Error(err)
		ctx.JSON(http.StatusOK, resp)
		return
	}
	param.TxStatus = constant.TxStatusCommitted

	c.validate(ctx, param)

	object := addition.NewUpdate(utils.TransferToContext(ctx), param.AppId)
	if err := object.UpdateStatus(param); err != nil {
		resp := utils.Error(err)
		ctx.JSON(http.StatusOK, resp)
		return
	}
	ctx.JSON(http.StatusOK, utils.Success(nil))
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (c Commit) validate(ctx *gin.Context, param params.ReqRollback) {
	if param.AppId == 0 || param.UserId == 0 || len(param.TxId) == 0 {
		resp := utils.Error(logger.NewError("error app_id or error user_id"))
		ctx.JSON(http.StatusOK, resp)
		return
	}
}
