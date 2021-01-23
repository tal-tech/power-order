/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-17 16:58:44
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-02-23 01:51:54
 * @Description:
 */
package order

import (
	"context"
	"powerorder/app/constant"
	"powerorder/app/dao/mysql"
)

type Extension struct {
	uAppId     uint
	uUserId    uint64
	strExtName string
	ctx        context.Context
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func NewExtensionModel(ctx context.Context, uAppId uint, uUserId uint64, strExtName string) *Extension {
	object := new(Extension)
	object.uAppId = uAppId
	object.uUserId = uUserId
	object.strExtName = strExtName
	object.ctx = ctx
	return object
}

/**
 * @description:  组装返回数据，根据redis缓存的订单信息【extensions】
 * @params {OrderIds} 想要获得Extension的OrderId
 * @params {OrderCached} 可能已从cache中读取到的cache
 * @return: 返回Extensions集合
 * @return: 返回需要种缓存的Extensions集合
 */
func (this *Extension) Get(OrderIds []string, OrderCached map[string]map[string]interface{}) (ret map[string][]map[string]interface{}, PipeData map[string]map[string]interface{}, err error) {
	ret = make(map[string][]map[string]interface{})
	PipeData = make(map[string]map[string]interface{})
	var OrderIdsWithoutExtension []string

	for i := 0; i < len(OrderIds); i++ {
		PipeData[OrderIds[i]] = make(map[string]interface{})
		var ExtValue map[string]interface{}
		var Value interface{}
		var ok bool
		//如果某个OrderId在Cache中存在
		if ExtValue, ok = OrderCached[OrderIds[i]]; ok {
			if Value, ok = ExtValue[this.strExtName]; ok {
				var ExtData []map[string]interface{}
				if ExtData, ok = Value.([]map[string]interface{}); ok {
					ret[OrderIds[i]] = ExtData
				}
				continue
			}

		}

		OrderIdsWithoutExtension = append(OrderIdsWithoutExtension, OrderIds[i])
		PipeData[OrderIds[i]][this.strExtName] = constant.Order_HashFieldDefault
	}

	if len(OrderIdsWithoutExtension) == 0 {
		return
	}
	ExtensionDao := mysql.NewExtensionDao(this.ctx, this.uAppId, this.uUserId, this.strExtName, constant.DBReader)
	ExtDataFromMysql, err := ExtensionDao.Get(OrderIdsWithoutExtension)

	if err != nil {
		return
	}

	for i := 0; i < len(ExtDataFromMysql); i++ {
		OrderId := ExtDataFromMysql[i]["order_id"].(string)
		ret[OrderId] = append(ret[OrderId], ExtDataFromMysql[i])
		//种缓存
		if Pipe, ok := PipeData[OrderId][this.strExtName]; ok {

			// todo  考虑改成if
			switch Pipe.(type) {
			case string:
				PipeData[OrderId][this.strExtName] = make([]map[string]interface{}, 0)
			}

		} else {
			PipeData[OrderId][this.strExtName] = make([]map[string]interface{}, 0)

		}
		PipeData[OrderId][this.strExtName] = append(PipeData[OrderId][this.strExtName].([]map[string]interface{}), ExtDataFromMysql[i])

	}

	return

}
