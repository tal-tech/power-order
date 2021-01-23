/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-21 21:07:19
 * @LastEditors  : lichanglin@tal.com
 * @LastEditTime : 2020-01-25 19:08:19
 * @Description:
 */
package utils

import (
	logger "github.com/tal-tech/loggerX"
	"powerorder/app/constant"
)

type Filter struct {
	Rules [][][]interface{}
}

/**
 * @description: 过滤条件初始化，校验
 * @params {Rules}
 * @return:
 */
func NewFilter(Rules [][][]interface{}) *Filter {
	object := new(Filter)
	object.Rules = Rules
	return object
}

/**
 * @description:
 * @params {data}
 * @return:
 */
func (this *Filter) Execute(data map[string]interface{}) (real bool, err error) {
	real = false
	err = nil

	for i := 0; i < len(this.Rules); i++ {
		AndRules := this.Rules[i]

		real = true
		for j := 0; j < len(AndRules); j++ {
			var key string
			var field interface{}
			var comparer uint
			var ok bool
			Rule := AndRules[j]

			if len(Rule) != 3 {
				logger.E("filter error", "rule params len != 3 ")
				return
			}
			if key, ok = Rule[0].(string); !ok {
				logger.E("filter error", "error rule params[0]")
				return
			}

			if field, ok = data[key]; !ok {
				logger.E("filter error", "error filed, key:%s", key)
				return
			}

			comparer, real = JsonInterface2UInt(Rule[1])
			if real == false {
				return
			}

			switch field.(type) {
			case float64:
			default:
				comparer += 100
			}
			if comparer < constant.StrGreater {
				real, err = this.IntCompare(field, comparer, Rule[2])
			} else {
				real, err = this.StrCompare(field, comparer, Rule[2])
			}

			if err != nil {
				return
			}
			if real == false {
				break
			}
		}

		if real == true {
			return
		}
	}
	real = false
	err = nil
	return
}

/**
 * @description:
 * @params {field} 当前字段的值
 * @params {comparer} 运算符
 * @params {value} 参数对应的值
 * @return:
 */
func (this *Filter) IntCompare(field interface{}, comparer uint, value interface{}) (real bool, err error) {
	var c1 int
	var c2 []int
	var c3 int

	c1, real = JsonInterface2Int(field)
	if real == false {
		return
	}
	if comparer == constant.IntWithIn || comparer == constant.IntNotWithIn {
		c2, real = JsonInterface2IntArray(value)
		if real == false {
			return
		}
	} else {
		c3, real = JsonInterface2Int(value)

		if real == false {
			return
		}

	}
	switch comparer {
	case constant.IntGreater:
		return c1 > c3, nil
	case constant.IntGreaterOrEqual:
		return c1 >= c3, nil
	case constant.IntLess:
		return c1 < c3, nil
	case constant.IntLessOrEqual:
		return c1 <= c3, nil
	case constant.IntEqual:
		return c1 == c3, nil
	case constant.IntNotEqual:
		return c1 != c3, nil
	case constant.IntNotWithIn:
		for i := 0; i < len(c2); i++ {
			if c1 == c2[i] {
				return false, nil
			}
		}
		return true, nil
	case constant.IntWithIn:
		for i := 0; i < len(c2); i++ {
			if c1 == c2[i] {
				return true, nil
			}
		}
		return false, nil
	default:
		logger.E("comparer error", "%d", comparer)
		return false, nil
	}
	return false, nil
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *Filter) StrCompare(field interface{}, comparer uint, value interface{}) (real bool, err error) {
	var c1 string
	var c2 []string
	var c3 string
	var ok bool
	if c1, ok = field.(string); !ok {
		logger.E("error field type", "")
		real = false
		return
	}
	if comparer == constant.StrWithIn || comparer == constant.StrNotWithIn {
		if c2, real = JsonInterface2StringArray(value); real == false {
			return
		}
	} else {
		if c3, real = JsonInterface2String(value); real == false {
			return
		}
	}
	switch comparer {
	case constant.StrGreater:
		return c1 > c3, nil
	case constant.StrGreaterOrEqual:
		return c1 >= c3, nil
	case constant.StrLess:
		return c1 < c3, nil
	case constant.StrLessOrEqual:
		return c1 <= c3, nil
	case constant.StrEqual:
		return c1 == c3, nil
	case constant.StrNotEqual:
		return c1 != c3, nil
	case constant.StrNotWithIn:
		for i := 0; i < len(c2); i++ {
			if c1 == c2[i] {
				return false, nil
			}
		}
		return true, nil
	case constant.StrWithIn:
		for i := 0; i < len(c2); i++ {
			if c1 == c2[i] {
				return true, nil
			}
		}
		return false, nil
	default:
		logger.E("error comparer", "%d", comparer)
		return false, nil
	}
	return false, nil
}
