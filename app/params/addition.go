package params

import "powerorder/app/dao/mysql"

type ReqAddition struct {
	Info       mysql.OrderInfo                     `form:"info" json:"info" validate:"required"`
	Detail     []mysql.OrderDetail                 `form:"detail" json:"detail" validate:"required"`
	Extensions map[string][]map[string]interface{} `form:"extensions" json:"extensions" validate:"required"`
}

type ReqBegin struct {
	AppId     uint       `form:"app_id" json:"app_id" validate:"required"`
	UserId    uint64     `form:"user_id" json:"user_id" validate:"required"`
	Additions []ReqAddition `form:"additions" json:"additions" validate:"required"`
}

type ReqRollback struct {
	AppId    uint   `form:"app_id" json:"app_id" validate:"required"`
	UserId   uint64 `form:"user_id" json:"user_id" validate:"required"`
	TxId     string `form:"tx_id" json:"tx_id" validate:"required"`
	TxStatus uint   `form:"tx_status" json:"tx_status" validate:"required"`
}

type ReqCommit struct {
	ReqRollback
}

type extensions map[string]interface{}

type ReqOrder struct {
	Info       mysql.OrderInfo                     `json:"info"`
	Detail     []mysql.OrderDetail                 `json:"detail"`
	Extensions map[string][]extensions 			   `json:"extensions"`
}
