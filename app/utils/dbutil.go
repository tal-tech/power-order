package utils

import logger "github.com/tal-tech/loggerX"

/*
 * @Author: lichanglin@tal.com
 * @Date: 2019-12-31 15:58:23
 * @LastEditors  : lichanglin@tal.com
 * @LastEditTime : 2020-01-06 19:24:22
 * @Description:
 以示例说明分库分表原理
 假设分成2库2表（每个库中有2张表），精度为1024
0------------------256------------------512------------------768------------------1024
0---------------Database 0--------------512---------------Database 1--------------1024
0-----Table 0------256------Table 1-----512-----Table 0------768-----Table 1------1024
*/

const (
	iDefaultPrecision = 1024
)

type Sharding struct {
	iDatabaseNum       uint //分库的数目，只分表的话此值为1
	iTableNum          uint //单库中分表的数目
	iDatabasePrecision uint //分库的精度 默认为1024
	iTablePrecision    uint //分表的精度 为分库的精度/分库的数目
}

/**
 * @description: 生成分库分表的对象
 * @param {iDatabaseNum 分库的个数，只分表不分库，此值为1}
 * @param {iTableNum 分表的个数（每个库中表的个数）}
 * @param {iPrecision 计算的精度}
 * @return: error 是否有异常
 * @return: *Sharding 生成的对象
 */
func NewSharding(iDatabaseNum, iTableNum, iPrecision uint) (*Sharding, error) {
	if iPrecision == 0 {
		iPrecision = iDefaultPrecision
	}
	if iDatabaseNum == 0 || iTableNum == 0 || iPrecision == 0 {
		return nil, logger.NewError("参数异常")
	}
	if iPrecision%iDatabaseNum != 0 || (iPrecision/iDatabaseNum)%iTableNum != 0 {
		return nil, logger.NewError("参数异常")
	}

	sharding := new(Sharding)
	sharding.iDatabaseNum = iDatabaseNum
	sharding.iTableNum = iTableNum
	sharding.iDatabasePrecision = iPrecision
	sharding.iTablePrecision = sharding.iDatabasePrecision / iDatabaseNum
	return sharding, nil

}

/**
 * @description: 计算分库后的库编号（从0开始）
 * @param {iBreakPoint 计算的奇点}
 * @return: 库编号
 */
func (this *Sharding) DatabaseNo(iBreakPoint uint) uint {
	return (iBreakPoint % this.iDatabasePrecision) / (this.iDatabasePrecision / this.iDatabaseNum)
}

/**
 * @description: 计算分表后的表编号（从0开始）
 * @param {iBreakPoing 计算的奇点}
 * @return: 表编号
 */
func (this *Sharding) TableNo(iBreakPoint uint) uint {
	return (iBreakPoint % this.iTablePrecision) / (this.iTablePrecision / this.iTableNum)
}
