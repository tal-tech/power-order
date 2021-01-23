/*
 * @Author: xiaotao@tal.com
 * @Date: 2020-05-13 16:20:17
 * @Description:
 */
package utils

import (
	"context"
)

const (
	APP_SEA_MAXSIZE = 50 // 大海业务线需求

	TOC_QUERY_MAX_DAY = 366 * 86400
	TOC_QUERY_MIN_DAY = 1

	SQL_WHERE_IN_MAX = 100 // mysql in查询最大个数
)
const (
	DB_SHADOW = "shadow"
)

// 查询压测标识
func GetPts(ctx context.Context) bool {
	var pts bool
	pts = ctx.Value("pts").(bool)
	return pts
}

// 影子库 读写替换
func DbShadowHandler(ctx context.Context, handler string) string {
	var pts bool
	pts = ctx.Value("pts").(bool)
	if pts == true {
		return handler + "_" + DB_SHADOW
	}
	return handler
}
