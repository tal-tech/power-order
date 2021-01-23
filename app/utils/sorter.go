/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-26 06:37:14
 * @LastEditors  : lichanglin@tal.com
 * @LastEditTime : 2020-02-04 01:38:30
 * @Description:
 */

package utils

import (
	"fmt"
	logger "github.com/tal-tech/loggerX"
	"powerorder/app/constant"
)

type Sorter struct {
	//承载以任意结构体为元素构成的Slice
	Slice *[]interface{}
	//排序规则，比如{"status", "desc", "type", "asc"},以order by status desc, type asc排序
	Rules []string
}

/**
 * @description: 校验排序方式参数
 * @params {rules}
 * @return: err
 */
func CheckRules(rules []string) (ret bool) {
	if len(rules)%2 != 0 {
		logger.I("CheckRules", "len of rules is not even number")
		return false
	}
	for i := 0; i < len(rules); i += 2 {
		if rules[i+1] != constant.Desc && rules[i+1] != constant.Asc {
			logger.E("CheckRules", "rules[%d] is wrong(not desc and not asc)", i)
			return false
		}
	}
	return true
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func NewSorter(data *[]interface{}, rules []string) (object *Sorter, err error) {

	ret := CheckRules(rules)
	if !ret {
		err = logger.NewError("checkrules failed")
		return nil, err
	}
	for i := 0; i < len(*data); i++ {
		switch (*data)[i].(type) {
		case map[string]interface{}:
			element, _ := (*data)[i].(map[string]interface{})
			for j := 0; j < len(rules); j += 2 {

				key := rules[j]
				if _, ok := element[key]; !ok {
					return nil, logger.NewError(fmt.Sprintf("data [%d] has not %s", i, key))

				}
			}
		default:
			return nil, logger.NewError(fmt.Sprintf("data [%d] type is not map[string]interface{}", i))
		}
	}

	object = new(Sorter)

	object.Slice = data
	object.Rules = rules

	return
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this Sorter) Len() int {
	return len(*this.Slice)
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this Sorter) Swap(i, j int) {
	(*this.Slice)[i], (*this.Slice)[j] = (*this.Slice)[j], (*this.Slice)[i]
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this Sorter) Less(i, j int) bool {

	elementi, _ := (*this.Slice)[i].(map[string]interface{})
	elementj, _ := (*this.Slice)[j].(map[string]interface{})

	for k := 0; k < len(this.Rules); k += 2 {

		key := this.Rules[k]
		valuei, _ := elementi[key]
		valuej, _ := elementj[key]

		if valuei == valuej {
			continue
		}

		switch valuei.(type) {
		case float64:
			vali, _ := JsonInterface2Int(valuei)
			valj, _ := JsonInterface2Int(valuej)
			if this.Rules[k+1] == constant.Asc {

				return vali < valj
			} else {

				return vali > valj
			}
		default:
			vali, _ := JsonInterface2String(valuei)
			valj, _ := JsonInterface2String(valuej)
			if this.Rules[k+1] == constant.Asc {

				return vali < valj
			} else {

				return vali > valj
			}
		}

	}
	return false
}
