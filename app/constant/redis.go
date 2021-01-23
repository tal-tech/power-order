/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-11 22:39:53
 * @LastEditors  : lichanglin@tal.com
 * @LastEditTime : 2020-02-12 13:01:36
 * @Description:
 */
package constant

const (
	Order_HashKeyPrx         = "power_order_"
	Order_HashSubKey_UserId  = "user_id"
	Order_HashSubKey_OrderId = "order_id"
	Order_HashSubKey_Info    = Info
	Order_HashSubKey_Detail  = Detail

	OrderId_SetKeyPrx = "power_order_id_"

	UserId_KVKeyPrx = "power_order_user_id_"
)
const (
	OrderId_MemberDefault        = "-1"
	Order_HashFieldDefault       = "-2"
	Order_HashFieldDefaultDecode = "\"-2\""
)
const (
	Order_SecondsPassedFromCreatedOneDay   = 60 * 60 * 24
	Order_SecondsPassedFromCreatedOneMonth = Order_SecondsPassedFromCreatedOneDay * 30
)

const (
	Order_ExpiredTimeNoOperation = -2
	Order_ExpiredTimeNoCache     = -1
	Order_ExpiredTimeTwoMinutes  = 2 * 60
	Order_ExpiredTimeOneDay      = 24 * 60 * 60

	OrderId_ExpiredTimeDefault = 24 * 60 * 60
	OrderId_ExpiredTime30Days  = 24 * 60 * 60 * 30

	//UserId_ExpiredTimeDefault = 5 * 30 * 24 * 60 * 60
	//设置为不过期
	UserId_ExpiredTimeDefault = -1 //5 * 30 * 24 * 60 * 60
)

// redis集群名称配置
const (
	Order_Redis_Cluster = "order"
)

var Order_HashSubKey map[uint][]string

func init() {
	Order_HashSubKey = make(map[uint][]string)
	Order_HashSubKey[AppBus1] = []string{Order_HashSubKey_UserId, Order_HashSubKey_Info, Order_HashSubKey_Detail, BUS1_GrouponInfoTableName, BUS1_GrouponDetailTableName, BUS2_ExtensionPromotionInfoTableName}
	Order_HashSubKey[AppBus2] = []string{Order_HashSubKey_UserId, Order_HashSubKey_Info, Order_HashSubKey_Detail, BUS2_ExtensionTableName, BUS2_ExtensionPromotionInfoTableName}

}
