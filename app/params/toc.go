package params


type ReqSearch struct {
	AppId        uint              `form:"app_id" json:"app_id" validate:"required"`
	UserId       uint64            `form:"user_id" json:"user_id" validate:"required"`
	InfoFilter   [][][]interface{} `form:"info_filter" json:"info_filter"`
	DetailFilter [][][]interface{} `form:"detail_filter" json:"detail_filter"`
	Fields       []string          `form:"fields" json:"fields"`
	Sorter       []string          `form:"sorter" json:"sorter"`
	Start        uint              `form:"start" json:"start"`
	End          uint              `form:"end" json:"end"`
}


type ReqGet struct {
	AppId    uint     `form:"app_id" json:"app_id"`
	UserId   uint64   `form:"user_id" json:"user_id"`
	OrderIds []string `form:"order_ids" json:"order_ids"`
	Fields   []string `form:"fields" json:"fields"`
}

type ReqGenerate struct {
	AppId  uint   `form:"app_id" json:"app_id"`
	UserId uint64 `form:"user_id" json:"user_id"`
	Num    uint   `form:"num" json:"num"`
}

type ReqOrderIdSearch struct {
	AppId  uint              `form:"app_id" json:"app_id"`
	Filter [][][]interface{} `form:"filter" json:"filter"`
	Table  string            `form:"table" json:"table"`
	Num    int64             `form:"num" json:"num"`
}

type ReqBGet struct {
	AppId  uint     `form:"app_id" json:"app_id"`
	Get    []ReqGet    `form:"get" json:"get"`
	Fields []string `form:"fields" json:"fields"`
}

type ReqQuery struct {
	AppId        uint              `form:"app_id" json:"app_id"`
	StartDate    int64              `form:"start_date" json:"start_date"`
	EndDate      int64              `form:"end_date" json:"end_date"`
	UserId       uint64            `form:"user_id" json:"user_id"`
	InfoFilter   [][][]interface{} `form:"info_filter" json:"info_filter"`
	DetailFilter [][][]interface{} `form:"detail_filter" json:"detail_filter"`
	Fields       []string          `form:"fields" json:"fields"`
	Sorter       []string          `form:"sorter" json:"sorter"`
	Start        uint              `form:"start" json:"start"`
	End          uint              `form:"end" json:"end"`
}

