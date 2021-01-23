/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-02-01 02:30:41
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime : 2020-02-16 10:12:29
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

type DSearch struct {
	//controller.Base
}

func (i DSearch) Index(ctx *gin.Context) {

	var param params.ReqDSearch

	if err := ctx.ShouldBindJSON(&param); err != nil {
		resp := utils.Error(err)
		ctx.JSON(http.StatusOK, resp)
		return
	}
	reqParam, _ := i.initData(param)
	i.validate(ctx, reqParam)

	searchModel := tob.NewSearch(utils.TransferToContext(ctx), reqParam.AppId, true)
	result, total, err := searchModel.Get(reqParam)

	if err != nil {
		resp := utils.Error(err)
		ctx.JSON(http.StatusOK, resp)
	} else {
		ctx.JSON(http.StatusOK, utils.Success(map[string]interface{}{"total": total, "result": result}))
	}
}

func (i DSearch) validate(ctx *gin.Context, param params.TobReqSearch) {
	if param.AppId == 0 {
		resp := utils.Error(logger.NewError("error app_id"))
		ctx.JSON(http.StatusOK, resp)
		return
	}
	if len(param.Fields[0]) == 0 {
		resp := utils.Error(logger.NewError("empty field"))
		ctx.JSON(http.StatusOK, resp)
		return
	}
	if param.End-param.Start > constant.MaxOrderCountAtSearchTime {
		resp := utils.Error(logger.NewError(fmt.Sprintf("max order count at search time is %d", constant.MaxOrderCountAtSearchTime)))
		ctx.JSON(http.StatusOK, resp)
		return
	}
}

func (i DSearch) initData(param params.ReqDSearch) (req params.TobReqSearch, err error) {
	if param.End == 0 {
		param.End = constant.MaxOrderCountAtSearchTime
	}
	req.AppId = param.AppId
	req.End = param.End
	req.Start = param.Start
	req.Sorter = param.Sorter
	req.Fields = make([]string, 1)
	req.Fields[0] = param.Field

	if param.Field == constant.Detail {
		req.DetailFilter = param.Filter
	} else {
		req.ExtFilter = make(map[string][][][]interface{}, 0)
		req.ExtFilter[param.Field] = param.Filter
	}
	return
}
