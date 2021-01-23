package params

import "powerorder/app/dao/mysql"

type ReqUpdate struct {
	AppId   uint    `json:"app_id"`
	UserId  uint64  `json:"user_id"`
	Orders  []order `json:"orders"`
	Consist bool    `json:"consist"`
}

// 更新的数据结构
type order struct {
	OrderId    string               `json:"order_id"`
	Info       info                 `form:"info" json:"info"`
	Detail     []detail             `form:"details" json:"details"`
	Extensions map[string]extension `form:"extensions" json:"extensions"`
}

type info struct {
	Version int             `json:"version"`
	Info    mysql.OrderInfo `json:"info"`
}

type detail struct {
	Version int               `json:"version"`
	Detail  mysql.OrderDetail `json:"detail"`
}

type extension struct {
	Updates    []extUpdate              `json:"updates"`
	Insertions []map[string]interface{} `json:"insertions"`
}

type extUpdate struct {
	Version int                    `json:"version"`
	Update  map[string]interface{} `json:"update"`
}