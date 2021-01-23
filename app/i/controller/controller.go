/*
 * @Author: lichanglin@tal.com
 * @Date: 2020-01-10 14:59:46
 * @LastEditors  : lichanglin@tal.com
 * @LastEditTime : 2020-01-11 20:52:15
 * @Description:
 */
package controller

type Controller interface {
	InitParam() error
	VldParam() error
	Run() error
	Output()
}

//func Index(instance Controller) {
//
//	err := instance.InitParam()
//
//	if err != nil {
//		return
//	}
//
//	err = instance.VldParam()
//
//
//	if err != nil {
//		return
//	}
//
//	instance.Run()
//	instance.Output()
//
//}
