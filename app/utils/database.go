package utils

import (
	"fmt"
	logger "github.com/tal-tech/loggerX"
	"powerorder/app/constant"
	"strconv"
)

/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-02-16 13:18:25
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-03-12 00:14:56
 * @Description:
 */
var shardingOrderInfo, _ = NewSharding(constant.DatabaseNum, constant.TableOrderInfoNum, 0)
var shardingOrderDetail, _ = NewSharding(constant.DatabaseNum, constant.TableOrderDetailNum, 0)

/**
 * @description:
 * @params {type}
 * @return:
 */
func GenDatabaseName(iUserId uint64) string {
	var uUserId = uint(iUserId % 10000)

	return fmt.Sprintf("%s%02d", constant.DatabaseNamePrx, shardingOrderInfo.DatabaseNo(uUserId))
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func GenDatabaseNameByOrderId(strOrderId string) string {
	strTailNumber := strOrderId[len(strOrderId)-4:]
	iUserId, _ := strconv.Atoi(strTailNumber)

	return GenDatabaseName(uint64(iUserId))
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func GenOrderInfoTableName(iUserId uint64) string {
	var uUserId = uint(iUserId % 10000)

	return fmt.Sprintf("%s%02d", constant.TableOrderInfoNamePrx, shardingOrderInfo.TableNo(uUserId))
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func GenOrderInfoTableNameByOrderId(strOrderId string) string {
	strTailNumber := strOrderId[len(strOrderId)-4:]
	iUserId, _ := strconv.Atoi(strTailNumber)

	return GenOrderInfoTableName(uint64(iUserId))
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func GenOrderDetailTableName(iUserId uint64) string {
	var uUserId = uint(iUserId % 10000)

	return fmt.Sprintf("%s%02d", constant.TableOrderDetailNamePrx, shardingOrderDetail.TableNo(uUserId))
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func GenOrderExtensionTableName(uUserId uint64, uAppId uint, strExtension string) (tableName string, err error) {
	var AppInfo map[string]uint
	var ok bool
	var TableNum uint
	if AppInfo, ok = ExtensionTableShardingInfo[uAppId]; !ok {
		err = logger.NewError("error appid")
		return
	}
	if TableNum, ok = AppInfo[strExtension]; !ok {
		err = logger.NewError(fmt.Sprintf("error extension table name:%s", strExtension))
		return
	}

	shardingObject, err := NewSharding(constant.DatabaseNum, TableNum, 0)

	if err != nil {
		return
	}

	return fmt.Sprintf("%s_%03d_%02d", strExtension, uAppId, shardingObject.TableNo((uint)(uUserId%10000))), nil

}

/**
 * @description:
 * @params {type}
 * @return:
 */
func GenDatabaseIdxName() string {
	return constant.DatabaseIdxName
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func GenOrderInfoIdxTableName(uAppId uint) string {
	return fmt.Sprintf("%s%03d", constant.TableOrderInfoIdxNamePrx, uAppId)
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func GenOrderDetailIdxTableName(uAppId uint) string {
	return fmt.Sprintf("%s%03d", constant.TableOrderDetailIdxNamePrx, uAppId)
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func GenOrderExtensionIdxTableName(uAppId uint, strExtension string) (tableName string, err error) {
	var AppInfo map[string]uint
	var ok bool
	if AppInfo, ok = ExtensionTableShardingInfo[uAppId]; !ok {
		err = logger.NewError("error appid")
		return
	}
	if _, ok = AppInfo[strExtension]; !ok {
		err = logger.NewError(fmt.Sprintf("error extension table name:%s", strExtension))
		return
	}
	tableName = fmt.Sprintf("%s_idx_%03d", strExtension, uAppId)
	return

}
