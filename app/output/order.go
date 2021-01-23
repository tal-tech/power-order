/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-13 22:43:14
 * @LastEditors: lichanglin@tal.com
 * @LastEditTime: 2020-02-18 01:24:42
 * @Description:
 */
package output

import (
	"bytes"
	"encoding/json"
	"powerorder/app/dao/mysql"
)

type Order struct {
	Info       mysql.OrderInfo                     `json:"info"`
	Detail     []mysql.OrderDetail                 `json:"detail"`
	Extensions map[string][]map[string]interface{} `json:"extensions"`
}

/**
 * @description:
 * @params {type}
 * @return:
 */

func (o Order) MarshalJSON() (b []byte, err error) {
	b = make([]byte, 0)
	bs := make([][]byte, 0)
	bs = append(bs, []byte("{"))
	var str []byte
	var first = true
	if o.Info.UserId != 0 {
		str, err = json.Marshal(o.Info)
		if err != nil {
			return
		}

		bs = append(bs, []byte("\"info\":"))
		bs = append(bs, []byte(str))

		first = false
	}

	if len(o.Detail) > 0 {
		str, err = json.Marshal(o.Detail)
		if err != nil {
			return
		}

		if first == false {
			bs = append(bs, []byte(","))
		}
		bs = append(bs, []byte("\"detail\":"))
		bs = append(bs, []byte(str))

		first = false
	}

	if o.Extensions != nil {

		if first == false {
			bs = append(bs, []byte(","))
		}
		//o.Extensions = make(map[string][]map[string]interface{}, 0)
		bs = append(bs, []byte("\"extensions\":"))

		str, err = json.Marshal(o.Extensions)
		if err != nil {
			return
		}
		bs = append(bs, []byte(str))
	}
	bs = append(bs, []byte("}"))

	b = bytes.Join(bs, []byte(""))
	return b, nil
}
