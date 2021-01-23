package utils

import (
	"powerorder/app/constant"
)

/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-14 01:36:30
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-04-12 23:31:50
 * @Description:
 */

//mysql扩展字段分表信息 第一级为业务线Id 第二级为扩展表名
var ExtensionTableShardingInfo map[uint]strmap

//Idxdb集合字段类型信息 第一级为业务线Id 第二级为集合名
var IdxdbCollectionFieldInfo map[uint]map[string]map[string]int

//Idxdb集合索引信息 第一级为业务线Id 第二级为集合名
var IdxdbCollectionIndexInfo map[uint]map[string][][]string

//Info集合字段类型信息
var InfoCollectionFieldInfo = map[string]int{
	"created_time":     constant.TypeString,
	"updated_time":     constant.TypeString,
	"cancelled_time":   constant.TypeString,
	"pay_created_time": constant.TypeString,
	"paid_time":        constant.TypeString,
	"expired_time":     constant.TypeString,
	"order_id":         constant.TypeString,
	"user_id":          constant.TypeInt,
	"id":               constant.TypeInt,
}

//Info集合索引信息
var InfoCollectionIndexInfo = [][]string{
	[]string{"created_time"},
	[]string{"updated_time"},
	[]string{"cancelled_time"},
	[]string{"pay_created_time"},
	[]string{"paid_time"},
	[]string{"expired_time"},
}

//Detail集合字段类型信息
var DetailCollectionFieldInfo = map[string]int{
	"created_time":      constant.TypeString,
	"updated_time":      constant.TypeString,
	"product_name":      constant.TypeString,
	"parent_product_id": constant.TypeInt,
	"product_id":        constant.TypeInt,
	"order_id":          constant.TypeString,
	"user_id":           constant.TypeInt,
	"id":                constant.TypeInt,
}

//Detail集合索引信息
var DetailCollectionIndexInfo = [][]string{
	[]string{"created_time"},
	[]string{"updated_time"},
	[]string{"product_name"},
	[]string{"parent_product_id"},
	[]string{"product_id"},
}

type strmap map[string]uint

//业务线Id
var AppIds = [...]uint{constant.AppBus1, constant.AppBus2}

func init() {
	ExtensionTableShardingInfo = make(map[uint]strmap)
	IdxdbCollectionFieldInfo = make(map[uint]map[string]map[string]int)
	IdxdbCollectionIndexInfo = make(map[uint]map[string][][]string)

	//todo 提供一种注册模式
	initAppBus1(&ExtensionTableShardingInfo, &IdxdbCollectionFieldInfo, &IdxdbCollectionIndexInfo)
	initAppBus2(&ExtensionTableShardingInfo, &IdxdbCollectionFieldInfo, &IdxdbCollectionIndexInfo)
}

func initAppBus2(extTable *map[uint]strmap, idxTableField *map[uint]map[string]map[string]int, idx *map[uint]map[string][][]string) {
	//扩展表分表信息 扩展表
	extTableInfo := map[string]uint{
		constant.BUS2_ExtensionTableName:              1,
		constant.BUS2_ExtensionPromotionInfoTableName: 1,
	}
	(*extTable)[constant.AppBus2] = extTableInfo
	// 索引表字段信息

	//Idxdb集合字段类型信息
	idxFieldInfo := map[string]map[string]int{
		constant.Info:   InfoCollectionFieldInfo,
		constant.Detail: DetailCollectionFieldInfo,
		//ExtensionTableName : map[string]int{
		//},
	}
	(*idxTableField)[constant.AppBus2] = idxFieldInfo

	// 索引库索引表字段
	idxInfo := map[string][][]string{
		constant.Info:   InfoCollectionIndexInfo,
		constant.Detail: DetailCollectionIndexInfo,
		//ExtensionTableName : [][]string{
		//},
	}
	(*idx)[constant.AppBus2] = idxInfo
}

func initAppBus1(extTable *map[uint]strmap, idxTableField *map[uint]map[string]map[string]int, idx *map[uint]map[string][][]string) {
	//扩展表分表信息 扩展表
	extTableInfo := map[string]uint{
		constant.BUS1_GrouponInfoTableName:      1,
		constant.BUS1_GrouponDetailTableName:      1,
		constant.BUS1_PromotionInfoTableName:	4,
	}
	(*extTable)[constant.AppBus1] = extTableInfo
	// 索引表字段信息

	//Idxdb集合字段类型信息
	idxFieldInfo := map[string]map[string]int{
		constant.Info:   InfoCollectionFieldInfo,
		constant.Detail: DetailCollectionFieldInfo,
		constant.BUS1_GrouponInfoTableName: map[string]int{
			"create_time":  constant.TypeString,
			"end_time":     constant.TypeString,
			"success_time": constant.TypeString,
		},
		constant.BUS1_GrouponDetailTableName: map[string]int{
			"groupon_order_id": constant.TypeString,
		},
	}
	(*idxTableField)[constant.AppBus1] = idxFieldInfo

	// 索引库索引表字段
	idxInfo := map[string][][]string{
		constant.Info:   InfoCollectionIndexInfo,
		constant.Detail: DetailCollectionIndexInfo,
		constant.BUS1_GrouponInfoTableName: [][]string{
			[]string{"create_time"},
			[]string{"end_time"},
			[]string{"success_time"},
		},
		constant.BUS1_GrouponDetailTableName: [][]string{
			[]string{"groupon_order_id"},
		},
	}
	(*idx)[constant.AppBus1] = idxInfo
}
