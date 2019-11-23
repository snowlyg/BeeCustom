package utils

import (
	"strings"

	"github.com/astaxie/beego"
)

//自定义末模板方法
func InitFunc() {
	_ = beego.AddFuncMap("inArray", inArray)
	_ = beego.AddFuncMap("inArrayStr", inArrayStr)
	_ = beego.AddFuncMap("inArrayStrSlice", inArrayStrSlice)
	_ = beego.AddFuncMap("canArray", canArray)
}

//is Array
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

//is inArrayStr
func inArrayStr(str string, s []string) bool {

	if len(s) == 0 {
		return false
	}

	if len(str) == 0 {
		return false
	}

	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

//is inArrayStrSlice
func inArrayStrSlice(str string, index int, s []map[int][]string) bool {
	if len(s) == 0 {
		return false
	}

	if len(str) == 0 {
		return false
	}

	for _, v := range s {
		for _, iv := range v[index] {
			if iv == str {
				return true
			}
		}
	}

	return false
}

//canArray 数组权限 清单，货物进出口
//{
//	"I" : {
//		"canCreate" :true
//		},
// 	"E" : {
//		"canCreate" :true
//		},
//}
func canArray(s map[string]map[string]bool, index, perm string) bool {
	//LogDebug(s)
	//LogDebug(index)
	//LogDebug(perm)
	if len(s) == 0 {
		return false
	}
	for sI, v := range s {
		if sI == strings.ToUpper(index) {
			return v[perm]
		}
	}
	return false
}
