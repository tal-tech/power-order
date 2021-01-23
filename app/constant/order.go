/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-06 22:56:33
 * @LastEditors  : lichanglin@tal.com
 * @LastEditTime : 2020-01-11 22:40:06
 * @Description: 订单相关的常量定义
 */
package constant

//业务线AppId
const (
	AppBus1 uint = 1
	AppBus2 uint = 8
)

// 下单tx_status 状态
const (
	TxStatusUnCommitted uint = 0
	TxStatusCommitted   uint = 1
)

const (
	OrderInfo   string = "info"
	OrderDetail string = "detail"
)

// 订单支付状态
const (
	StatusDefault     int = 1
	StatusPaying      int = 2
	StatusPaySuccess  int = 3
	StatusCancel      int = 3
	StatusCancelBySys int = 3
)

const (
	Addition = 1 // 添加订单
	Sync     = 2 // 订单同步
	ZeroTime = -62135596800
)

const (
	MaxOrderCountAtInsertionTime uint = 20
	MaxOrderCountAtSearchTime    uint = 20

	MaxBatchOrderCount     = 20
	MaxBatchUserCount      = 5
	MaxBatchUserOrderCount = 10
)

const (
	TopicPrex = "pwr_order_"
)

const (
	IndexPrex = "pwr_order_"

	OrderInfoType   = OrderInfo
	OrderDetailType = OrderDetail
)

const (
	Desc = "desc"
	Asc  = "asc"
)

const (
	Info   = "info"
	Detail = "detail"
)

// 自定义类型
const (
	TypeString      = 1
	TypeInt         = 2
	TypeTimestamp   = 3
	TypeTimestamp2 = 101 //类型为string，但要转成int
)
