package constant

// 通用错误码
const(
	SYSTEM_ERROR = 50000
	ORDER_SYSTEM_ERROR  = 60000
	ORDER_SIGNATURE_ERROR = 60001
	ORDER_SIGNATURE_REQUIRED_ERROR = 60002
	ORDER_SIGNATURE_TIMEFORMAT_ERROR = 60003
	ORDER_SIGNATURE_EXPIRED_ERROR = 60004
	ORDER_SIGNATURE_VERIFIED_ERROR = 60005

	ORDER_TOC_ORDER_INFO_ERROR = 61000
)

var  ERROR_MSG_MAP map[int]string = map[int]string{
	50000 : "system error",
	60000 : "order system error",
	60001 : "order signature error",
	60002 : "order signature required",
	60003 : "order signature time format error",
	60004 : "order system signature expired",
	60005 : "order system signature verified error",
	61000 : "save order info error",

}