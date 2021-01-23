/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-12 15:16:47
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-04-10 00:51:34
 * @Description:
 */
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tal-tech/xredis"
	"powerorder/app/constant"
	"powerorder/app/output"
	"powerorder/app/utils"
)

type TxOrder struct {
	Dao
	txId string
	ctx  context.Context
}

type TxInfo struct {
	OrderIds   []string
	Extensions map[string]bool
	OrderInfos map[string]output.Order
}

/**
 * @description:
 * @params {type}
 * @return:
 */

func NewTxOrderDao(ctx context.Context, uAppId uint, txId string, instance string) *TxOrder {
	object := new(TxOrder)
	object.txId = txId
	object.ctx = ctx
	object.Dao.uAppId = uAppId
	object.Dao.instance = instance
	return object
}

/**
 * @description: 获取order_id列表的redis无序集合的key
 * @params {}
 * @return: xes_platform_order_id_001_1
 */
func (this *TxOrder) GenSetKey() string {
	pts := utils.GetPts(this.ctx)
	if pts {
		return fmt.Sprintf("pts_%s%d_%v", "tx_order_", this.uAppId, this.txId)
	}

	return fmt.Sprintf("%s%d_%v", "tx_order_", this.uAppId, this.txId)
}

func (this *TxOrder) Set(TxInfo TxInfo) error {
	key := this.GenSetKey()
	b, err := json.Marshal(TxInfo)
	if err != nil {
		return err
	}
	_, err = xredis.NewSimpleXesRedis(this.ctx, constant.Order_Redis_Cluster).Set(key, []interface{}{}, string(b), constant.OrderId_ExpiredTimeDefault)
	return err

}

func (this *TxOrder) Get() (TxInfo, error) {

	key := this.GenSetKey()
	back, err := xredis.NewSimpleXesRedis(this.ctx, constant.Order_Redis_Cluster).Get(key, []interface{}{})
	var txInfo TxInfo

	err = json.Unmarshal([]byte(back), &txInfo)

	return txInfo, err

}
