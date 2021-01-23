/*
 * @Author: lichanglin@tal.com
 * @Date: 2019-12-31 11:07:17
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-02-27 23:30:52
 * @Description:
 */
package addition

import (
	"fmt"
	"github.com/gin-gonic/gin"
	logger "github.com/tal-tech/loggerX"
	"net/http"
	"powerorder/app/constant"
	"powerorder/app/model/addition"
	"powerorder/app/params"
	"powerorder/app/utils"
)

type Sync struct {
	//controller.Base
}

/**
 * @description: 历史数据同步(2c)
 * @params {ctx}
 * @return:
 */
func (s Sync) Index(ctx *gin.Context) {
	var param params.ReqBegin

	if err := ctx.ShouldBindJSON(&param); err != nil {
		resp := utils.Error(err)
		ctx.JSON(http.StatusOK, resp)
		return
	}
	s.validate(ctx, param)

	txId := GetTxId(ctx)

	insertModel := addition.NewInsertion(constant.Sync, utils.TransferToContext(ctx))
	if insertModel == nil {
		resp := utils.Error(logger.NewError("system error", logger.SYSTEM_DEFAULT))
		ctx.JSON(http.StatusOK, resp)
		return
	}
	orderIds, err := insertModel.Insert(param, txId)

	if err != nil {
		ctx.JSON(http.StatusOK, utils.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, utils.Success(map[string]interface{}{"order_ids": orderIds}))
}

func (s Sync) validate(ctx *gin.Context, param params.ReqBegin) {
	if param.AppId == 0 || param.UserId == 0 {
		resp := utils.Error(logger.NewError("error app_id or error user_id", logger.PARAM_ERROR))
		ctx.JSON(http.StatusOK, resp)
		return
	}

	if len(param.Additions) == 0 {
		resp := utils.Error(logger.NewError("error additions"))
		ctx.JSON(http.StatusOK, resp)
		return
	} else if len(param.Additions) > int(constant.MaxOrderCountAtInsertionTime) {
		resp := utils.Error(logger.NewError(fmt.Sprintf("max order count at insertion time is %d", constant.MaxOrderCountAtInsertionTime)))
		ctx.JSON(http.StatusOK, resp)
		return
	}
}
