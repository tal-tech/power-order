package params

type ReqDSearch struct {
	AppId  uint              `form:"app_id" json:"app_id"`
	Filter [][][]interface{} `form:"filter" json:"filter"`
	Field  string            `form:"field" json:"field"`
	Sorter []string          `form:"sorter" json:"sorter"`
	Start  uint              `form:"start" json:"start"`
	End    uint              `form:"end" json:"end"`
}


type TobReqSearch struct {
	ReqSearch
	ExtFilter map[string][][][]interface{} `form:"ext_filter" json:"ext_filter"`
}

