package middlewares

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	logger "github.com/tal-tech/loggerX"
	"net/http"
	"powerorder/app/constant"
	"powerorder/app/utils"
)

// 接口验签
func VerifySignMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 签名验证
		ok := verifySign("1","app_001")
		if !ok {
			outError(c, constant.ORDER_SIGNATURE_VERIFIED_ERROR)
			return
		}
		c.Next()
	}
}

func outError(c *gin.Context, code int) {
	msg := "order system error"
	if value, ok:=constant.ERROR_MSG_MAP[code]; ok{
		msg = value
	}
	err := logger.NewError(msg)
	err.Code = code
	resp := utils.Error(err)
	c.AbortWithStatusJSON(http.StatusOK, resp)
}

func verifySign(appId, sign string) bool {
	bodyMd5 := _md5(appId)
	return bodyMd5 == sign
}

func _md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
