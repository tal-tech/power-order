/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-22 15:37:32
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-04-13 15:57:16
 * @Description:
 */
package utils

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	logger "github.com/tal-tech/loggerX"
	"powerorder/app/constant"
	"strconv"
	"time"
)

/**
 * @description: 类型转换
 * @params {from}
 * @return:
 */
func JsonInterface2Int(from interface{}) (to int, ret bool) {
	to = 0
	ret = true
	var err error
	switch from.(type) {
	case float64:
		to, err = strconv.Atoi(strconv.FormatFloat(from.(float64), 'f', -1, 64))
		if err != nil {
			logger.E("error from type", "str to float64 error")
			ret = false
			return
		}
	case int64:
		to = int(from.(int64))
		return
	default:
		logger.E("error from type", "")
		ret = false
		return
	}
	return
}

/**
 * @description: 类型转换
 * @params {from}
 * @return:
 */
func JsonInterface2UInt(from interface{}) (to uint, ret bool) {
	to = 0
	ret = true
	var temp int
	temp, ret = JsonInterface2Int(from)

	if ret == false {
		return
	}

	to = uint(temp)
	return
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func JsonInterface2IntArray(from interface{}) (tos []int, ret bool) {
	tos = []int{}
	ret = true
	var ok bool
	var tmpArr []interface{}
	var tmpInt int
	switch from.(type) {
	case []interface{}:
		if tmpArr, ok = from.([]interface{}); ok {

			for i := 0; i < len(tmpArr); i++ {
				tmpInt, ret = JsonInterface2Int(tmpArr[i])
				if ret == false {
					return
				}
				tos = append(tos, tmpInt)
			}
		} else {
			logger.E("error from type", "")
			ret = false

		}
	default:
		logger.E("error from type", "")
		ret = false
		return
	}
	return
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func JsonInterface2String(from interface{}) (to string, ret bool) {
	to = ""
	ret = true
	switch from.(type) {
	case string:
		to = from.(string)
		return
	default:
		ret = false
		return
	}
	return
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func JsonInterface2StringArray(from interface{}) (tos []string, ret bool) {
	tos = []string{}
	ret = true
	var ok bool
	var tmpArr []interface{}
	var tmpStr string
	switch from.(type) {
	case []interface{}:
		if tmpArr, ok = from.([]interface{}); ok {

			for i := 0; i < len(tmpArr); i++ {
				tmpStr, ret = JsonInterface2String(tmpArr[i])
				if ret == false {
					return
				}
				tos = append(tos, tmpStr)
			}
		} else {
			logger.E("error from type", "")
			ret = false
		}
	default:
		logger.E("error from type", "")
		ret = false
		return
	}
	return
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func Struct2Map(from interface{}) (to map[string]interface{}, err error) {
	var data []byte
	data, err = json.Marshal(from)
	if err != nil {
		return
	}
	to = make(map[string]interface{})
	err = json.Unmarshal(data, &to)
	return

}

/**
 * @description:
 * @params {type}
 * @return:
 */
func StructArray2MapArray(from []interface{}) (to []map[string]interface{}, err error) {
	to = make([]map[string]interface{}, len(from))
	for i := 0; i < len(from); i++ {
		to[i], err = Struct2Map(from[i])
		if err != nil {
			return
		}
	}
	return
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func JsonInterface2Timestamp(from interface{}) (to int64, ret bool) {
	strTime, ret := JsonInterface2String(from)
	if ret == false {
		return
	}
	var stamp time.Time
	var err error
	stamp, err = time.ParseInLocation("2006-01-02 15:04:05", strTime, time.Local) //使用parseInLocation将字符串格式化返回本地时区时间
	if err != nil {
		ret = false
		return
	}
	to = (stamp.Unix()) //输出：1546926630
	return
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func JsonInterface2TimeString(from interface{}) (to string, ret bool) {
	var t int
	t, ret = JsonInterface2Int(from)
	if ret == false {
		return
	}
	to = time.Unix(int64(t), 0).Format("2006-01-02 15:04:05")
	return

}

/**
 * @description:  解析filter成为mysql可使用的
 * @params {type}
 * @return:
 */
func And2Where(from [][]interface{}) (whereStr string, whereValue []interface{}) {
	whereValue = make([]interface{}, 0)
	whereStr = " 1=1 "
	for i := 0; i < len(from); i++ {
		field, _ := JsonInterface2String(from[i][0])
		comparer, _ := JsonInterface2UInt(from[i][1])
		value := from[i][2]
		switch comparer {
		case constant.IntGreater:
			whereStr = fmt.Sprintf("%s and %s > ?", whereStr, field)
			whereValue = append(whereValue, value)
		case constant.IntGreaterOrEqual:
			whereStr = fmt.Sprintf("%s and %s >= ?", whereStr, field)
			whereValue = append(whereValue, value)
		case constant.IntLess:
			whereStr = fmt.Sprintf("%s and %s < ?", whereStr, field)
			whereValue = append(whereValue, value)
		case constant.IntLessOrEqual:
			whereStr = fmt.Sprintf("%s and %s <= ?", whereStr, field)
			whereValue = append(whereValue, value)
		case constant.IntEqual:
			whereStr = fmt.Sprintf("%s and %s = ?", whereStr, field)
			whereValue = append(whereValue, value)
		case constant.IntNotEqual:
			whereStr = fmt.Sprintf("%s and %s != ?", whereStr, field)
			whereValue = append(whereValue, value)
		}
	}
	return
}

func ToStringStringMap(data interface{}) (map[string]string, error) {
	if valueArr, ok := data.([]interface{}); ok {
		back := make(map[string]string)
		for i := 0; i < len(valueArr); i = i + 2 {
			back[cast.ToString(valueArr[i])] = cast.ToString(valueArr[i+1])
		}
		return back, nil
	}
	return cast.ToStringMapString(data), nil
}
