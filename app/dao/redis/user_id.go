/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-23 20:14:46
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-02-22 19:51:40
 * @Description:
 */

package redis

import (
	"context"
	"fmt"
	"github.com/tal-tech/xredis"
	"powerorder/app/constant"
	"powerorder/app/utils"
	"strconv"
)

type UserId struct {
	Dao
	ctx context.Context
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func NewUserIdDao(ctx context.Context, uAppId uint, instance string) *UserId {

	object := new(UserId)
	object.uAppId = uAppId
	object.ctx = ctx

	return object
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *UserId) GenKVKey(OrderId string) string {
	if utils.GetPts(this.ctx) {
		return fmt.Sprintf("pts_%s%03d_%v", constant.UserId_KVKeyPrx, this.uAppId, OrderId)
	}
	return fmt.Sprintf("%s%03d_%v", constant.UserId_KVKeyPrx, this.uAppId, OrderId)
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *UserId) MGet(OrderIds []string) (ret map[string]uint64, err error) {
	OtherKeys := make([]interface{}, 0)
	for i := 1; i < len(OrderIds); i++ {
		OtherKeys = append(OtherKeys, this.GenKVKey(OrderIds[i]))
		OtherKeys = append(OtherKeys, nil)
	}
	Key1 := this.GenKVKey(OrderIds[0])

	var back []string
	back, err = xredis.NewSimpleXesRedis(this.ctx, constant.Order_Redis_Cluster).MGet(Key1, []interface{}{}, OtherKeys...)

	var UserId uint64
	var OrderId string

	ret = make(map[string]uint64, 0)
	for i := 0; i < len(back); i++ {
		OrderId = OrderIds[i]
		UserId, err = strconv.ParseUint(back[i], 10, 64)

		if err != nil {
			return
		}

		ret[OrderId] = UserId

	}

	return
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *UserId) Set(OrderId string, UserId uint64) (back string, err error) {

	Key := this.GenKVKey(OrderId)

	back, err = xredis.NewSimpleXesRedis(this.ctx, constant.Order_Redis_Cluster).Set(Key, []interface{}{}, UserId, int64(constant.UserId_ExpiredTimeDefault))

	return

}
