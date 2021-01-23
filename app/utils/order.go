/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-07 18:47:22
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-03-13 01:02:50
 * @Description:
 */
package utils

import (
	"fmt"
	"math/rand"
	"time"
)

/**
 * @description: 生成订单Id，生成格式:时间秒（12位）+ appid（3位）+ 随机数（5位）+ 用户Id后四位（4位）=24
 * @params {uAppId} 业务线Id
 * @params {time} 订单创建时间
 * @params {uUserId} 用户Id
 * @return: 生成的订单Id
 */
func GenOrderId(uAppId uint, time time.Time, uUserId uint64) string {
	return fmt.Sprintf("%s%03d%05d%04d", time.Format("060102150405"), uAppId%1000, rand.Int63n(100000), uUserId%10000)
}
