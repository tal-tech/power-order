/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-02-16 13:19:15
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2021-01-04 11:05:40
 * @Description:
 */

package constant

const (
	DBWriter = "writer"
	DBReader = "reader"
)

// 分库分表配置
const (
	// 分库数量配置
	DatabaseNum uint = 4
	// 分表order_info数量配置
	TableOrderInfoNum uint = 4
	// 分表order_detail数量配置
	TableOrderDetailNum uint = 32
)

// 分库和分表的前缀
const (
	DatabaseNamePrx         = "pwr_platform_order_"
	TableOrderInfoNamePrx   = "order_info_"
	TableOrderDetailNamePrx = "order_detail_"

	DatabaseIdxName            = "pwr_platform_order_idx"
	TableOrderInfoIdxNamePrx   = "order_info_idx_"
	TableOrderDetailIdxNamePrx = "order_detail_idx_"
)

/* 业务线1扩展表*/
const (
	BUS1_GrouponInfoTableName   = "groupon_info"
	BUS1_GrouponDetailTableName = "groupon_detail"
	BUS1_PromotionInfoTableName = "promotion_info"
)

/* 业务线2扩展表*/
const (
	BUS2_ExtensionTableName              = "order_info_ext"
	BUS2_ExtensionPromotionInfoTableName = "promotion_info"
)
