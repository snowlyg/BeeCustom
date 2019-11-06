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
