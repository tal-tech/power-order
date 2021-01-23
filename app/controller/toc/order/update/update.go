/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-19 23:07:48
 * @LastEditors  : lichanglin@tal.com
 * @LastEditTime : 2020-01-28 01:35:48
 * @Description:
 */

package update

import (
	"github.com/gin-gonic/gin"
	logger "github.com/tal-tech/loggerX"
	"net/http"
	"powerorder/app/model/update"
	"powerorder/app/params"
	"powerorder/app/utils"
)

type Update struct {
	//controller.Base
}

/**
 * @description: 更新订单接口(2c)
 * @params {ctx}
 * @return:
 */
func (u Update) Index(ctx *gin.Context) {
	var param params.ReqUpdate
	if err := ctx.ShouldBindJSON(&param); err != nil {
		resp := utils.Error(err)
		ctx.JSON(http.StatusOK, resp)
		return
	}
	u.validate(ctx, param)
	updateModel := update.NewUpdate(utils.TransferToContext(ctx), param.AppId, param.UserId)
	if err := updateModel.Update(param); err != nil {
		resp := utils.Error(err)
		ctx.JSON(http.StatusOK, resp)
		return
	}
	ctx.JSON(http.StatusOK, utils.Success([]interface{}{}))
}

func (u Update) validate(ctx *gin.Context, param params.ReqUpdate) {
	if param.AppId == 0 || param.UserId == 0 {
		resp := logger.NewError("error app_id or error user_id")
		ctx.JSON(http.StatusOK, utils.Error(resp))
		return
	}
}
