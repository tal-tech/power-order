/*
 * @Author: lichanglin@tal.com
 * @Date: 2019-12-31 14:20:01
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2019-12-31 14:26:59
 * @Description:
 */
package tool

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"powerorder/app/utils"
)

//HealthCheck 健康检测
func HealthCheck(ctx *gin.Context) {
	resp := utils.Success(os.Getpid())
	ctx.JSON(http.StatusOK, resp)
	return
}
