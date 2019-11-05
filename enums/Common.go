package enums

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"BeeCustom/utils"
)

type JsonResultCode int

const (
	JRCodeFailed JsonResultCode = iota //接口返回状态 0
	JRCodeSucc                         //接口返回状态 1
	JRCode302    = 302                 //跳转至地址
	JRCode401    = 401                 //未授权访问
)

const (
	Deleted  = -1
	Disabled = false
	Enabled  = true
)

const BaseFormat = "2006-01-02 15:04:05"

//设置值
func SetObjValue(objName, v string, t reflect.Value) {
	switch t.FieldByName(objName).Kind() {
	case reflect.String:
		t.FieldByName(objName).Set(reflect.ValueOf(v))
	case reflect.Float64:
		handBookV, err := strconv.ParseFloat(v, 64)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("ParseFloat:%v", err))
		}
		t.FieldByName(objName).Set(reflect.ValueOf(handBookV))
	case reflect.Uint64:
		reflect.ValueOf(v)
		handBookV, err := strconv.ParseUint(v, 0, 64)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("ParseUint:%v", err))
		}
		t.FieldByName(objName).Set(reflect.ValueOf(handBookV))
	case reflect.Struct:
		handBookV, err := time.Parse("2006-01-02", v)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("Parse:%v", err))
		}
		t.FieldByName(objName).Set(reflect.ValueOf(handBookV))
	default:
		//utils.LogDebug("未知类型")
	}
}

//设置值
func FilpValueString(obj map[string]string) map[string]string {
	for i, v := range obj {
		obj[v] = i
	}

	return obj
}
