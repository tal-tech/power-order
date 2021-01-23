/*
 * @Author: lichanglin@tal.com
 * @Date: 2019-12-31 11:07:17
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-04-09 14:34:03
 * @Description:
 */
package addition

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openzipkin/zipkin-go/idgenerator"
	logger "github.com/tal-tech/loggerX"
	"net/http"
	"powerorder/app/constant"
	"powerorder/app/model/addition"
	"powerorder/app/params"
	"powerorder/app/utils"
)

type Begin struct {
	//controller.Base
}

/**
* @description:  创建订单开始
* @params {type}
* @return:
 */
func (b Begin) Index(ctx *gin.Context) {
	var param params.ReqBegin

	if err := ctx.ShouldBindJSON(&param); err != nil {
		resp := utils.Error(err)
		ctx.JSON(http.StatusOK, resp)
		return
	}
	txId := GetTxId(ctx)

	insertModel := addition.NewInsertion(constant.Addition, utils.TransferToContext(ctx))
	orderIds, err := insertModel.Insert(param, txId)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, utils.Success(map[string]interface{}{"tx_id": txId, "order_ids": orderIds}))
}

func GetTxId(ctx *gin.Context) (txId string) {
	return idgenerator.NewRandom64().TraceID().String()
}

func (b Begin) validate(ctx *gin.Context, param params.ReqBegin) {
	if param.AppId == 0 || param.UserId == 0 {
		resp := utils.Error(logger.NewError("error app_id or error user_id"))
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
