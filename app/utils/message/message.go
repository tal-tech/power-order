/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-27 18:40:38
 * @LastEditors  : lichanglin@tal.com
 * @LastEditTime : 2020-02-04 01:41:41
 * @Description:
 */
package message

import (
	"fmt"
	"github.com/spf13/cast"
	"github.com/tal-tech/xtools/kafkautil"
	"powerorder/app/constant"
)

type Message struct {
	strTopic string
	uUserId  uint64
	uAppId   uint
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func NewMessage(uUserId uint64, uAppId uint) *Message {
	object := new(Message)
	object.uUserId = uUserId
	object.uAppId = uAppId
	object.InitTopic()
	return object
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *Message) InitTopic() (err error) {

	this.strTopic = fmt.Sprintf("%s%03d", constant.TopicPrex, this.uAppId)
	err = nil
	return
}

/**
 * @description:
 * @params {type}
 * @return:
 */
func (this *Message) Send(param interface{}) (err error) {
	err = nil
	err = kafkautil.Send2Proxy(this.strTopic, []byte("kafka "+cast.ToString(param)))

	return
}
