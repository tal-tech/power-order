/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-14 16:20:17
 * @LastEditors  : lichanglin@tal.com
 * @LastEditTime : 2020-02-16 16:11:31
 * @Description:
 */
package utils

import (
	"errors"
	"fmt"
	logger "github.com/tal-tech/loggerX"
	"github.com/tal-tech/xredis"
	"math/rand"
	"os"
	"powerorder/app/constant"
	"strconv"
	"sync"
	"time"
)

type IdWorker struct {
	startTime             int64
	workerIdBits          uint
	datacenterIdBits      uint
	maxWorkerId           int64
	maxDatacenterId       int64
	sequenceBits          uint
	workerIdLeftShift     uint
	datacenterIdLeftShift uint
	timestampLeftShift    uint
	sequenceMask          int64
	workerId              int64
	datacenterId          int64
	sequence              int64
	lastTimestamp         int64
	signMask              int64
	idLock                *sync.Mutex
}

var idObj *IdWorker

const workerKey = "xes_order_platform_workers"

/**
 * @description:  生成64位随机int
 * @params {type}
 * @return:
 */
func GenId() (random int64) {
	currWorker, err := newIdWorker()
	if err != nil {
		logger.E("InitIdWorker failed", "err %+v", err)
		return
	}
	random, _ = currWorker.nextId()
	return
}

func newIdWorker() (obj *IdWorker, err error) {
	if idObj != nil && idObj.workerId > 0 {
		obj = idObj
		return
	}
	obj = &IdWorker{}

	// 获取workerId 和 datacenterId
	workerId, datacenterId := obj.getWorkerId()

	var baseValue int64 = -1
	obj.startTime = 1463834116272
	obj.workerIdBits = 5
	obj.datacenterIdBits = 5
	obj.maxWorkerId = baseValue ^ (baseValue << obj.workerIdBits)
	obj.maxDatacenterId = baseValue ^ (baseValue << obj.datacenterIdBits)
	obj.sequenceBits = 12
	obj.workerIdLeftShift = obj.sequenceBits
	obj.datacenterIdLeftShift = obj.workerIdBits + obj.workerIdLeftShift
	obj.timestampLeftShift = obj.datacenterIdBits + obj.datacenterIdLeftShift
	obj.sequenceMask = baseValue ^ (baseValue << obj.sequenceBits)
	obj.sequence = 0
	obj.lastTimestamp = -1
	obj.signMask = ^baseValue + 1

	obj.idLock = &sync.Mutex{}

	if obj.workerId < 0 || obj.workerId > obj.maxWorkerId {
		err = errors.New(fmt.Sprintf("workerId[%v] is less than 0 or greater than maxWorkerId[%v].", workerId, datacenterId))
		return
	}
	if obj.datacenterId < 0 || obj.datacenterId > obj.maxDatacenterId {
		err = errors.New(fmt.Sprintf("datacenterId[%d] is less than 0 or greater than maxDatacenterId[%d].", workerId, datacenterId))
		return
	}
	obj.workerId = workerId
	obj.datacenterId = datacenterId
	idObj = obj
	return
}

//方案一： 通过配置
func (this *IdWorker) getWorkerId() (workerId, datacenterId int64) {
	datacenterId = rand.Int63n(31)
	workerId = rand.Int63n(31)
	hostName, err := os.Hostname()
	if err != nil || hostName == "" {
		return
	}
	rds := xredis.NewSimpleXesRedis(nil, constant.Order_Redis_Cluster)
	workers, err := rds.HGetAll(workerKey, nil)
	if err != nil {
		return
	}
	workerNum := len(workers)
	if _, ok := workers[hostName]; !ok {
		workerNum = workerNum + 1
		rds.HSet(workerKey, nil, hostName, workerNum)
		workerId = int64(workerNum)
	} else {
		workerId, _ = strconv.ParseInt(workers[hostName], 10, 64)
	}
	return
}

func (this *IdWorker) nextId() (int64, error) {
	this.idLock.Lock()
	defer this.idLock.Unlock()
	timestamp := this.genTime()
	if timestamp < this.lastTimestamp {
		return -1, errors.New(fmt.Sprintf("Clock moved backwards.  Refusing to generate id for %d milliseconds", this.lastTimestamp-timestamp))
	}

	if timestamp == this.lastTimestamp {
		this.sequence = (this.sequence + 1) & this.sequenceMask
		if this.sequence == 0 {
			timestamp = this.tilNextMillis()
			this.sequence = 0
		}
	} else {
		this.sequence = 0
	}

	this.lastTimestamp = timestamp

	id := ((timestamp - this.startTime) << this.timestampLeftShift) |
		(this.datacenterId << this.datacenterIdLeftShift) |
		(this.workerId << this.workerIdLeftShift) |
		this.sequence

	if id < 0 {
		id = -id
	}

	return id, nil
}

func (this *IdWorker) tilNextMillis() int64 {
	timestamp := this.genTime()
	if timestamp <= this.lastTimestamp {
		timestamp = this.genTime()
	}
	return timestamp
}

func (this *IdWorker) genTime() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
