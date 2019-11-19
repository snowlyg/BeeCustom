package utils

import (
	"github.com/astaxie/beego"
)

//自定义末模板方法
func InitFunc() {
	_ = beego.AddFuncMap("inArray", inArray)
}

func inArray(in int64, s []interface{}) bool {

	if len(s) == 0 {
		return false
	}

	for _, v := range s {
		if v == in {
			return true
		}
	}

	return false
}
