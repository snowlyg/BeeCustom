package xlsx

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"BeeCustom/utils"
	"github.com/astaxie/beego"
)

//获取导入文件表头
func GetExcelTitles(xmlTitle, configSection string) (map[string]string, error) {
	rXmlTitles := map[string]string{}
	if len(xmlTitle) == 0 {
		importWord, err := beego.AppConfig.GetSection(configSection)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("GetSection:%v", err))
			return nil, err
		}
		rXmlTitles = importWord
	} else {
		xmlTitles := strings.Split(xmlTitle, "/")
		for k, v := range xmlTitles {
			rXmlTitles[v] = strconv.Itoa(k)
		}
	}

	return rXmlTitles, nil
}

//获取导入文件表名称
func GetExcelName(configSection string) (string, error) {
	nameMap, err := beego.AppConfig.GetSection(configSection)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("GetSection:%v", err))
		return configSection, err
	}
	name := nameMap["name"]
	return name, nil
}

//设置值
func SetObjValue(objName, v string, t reflect.Value) {
	switch t.FieldByName(objName).Kind() {
	case reflect.String:
		t.FieldByName(objName).Set(reflect.ValueOf(v))
	case reflect.Float64:
		handBookV, err := strconv.ParseFloat(v, 64)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("ParseFloat:%v,%v", err, v))
		}
		t.FieldByName(objName).Set(reflect.ValueOf(handBookV))
	case reflect.Uint64:
		reflect.ValueOf(v)
		handBookV, err := strconv.ParseUint(v, 0, 64)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("ParseUint:%v,%v", err, v))
		}
		t.FieldByName(objName).Set(reflect.ValueOf(handBookV))
	case reflect.Struct:
		handBookV, err := time.Parse("20060102", v)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("Parse:%v,%v", err, v))
		}
		t.FieldByName(objName).Set(reflect.ValueOf(handBookV))
	default:
		utils.LogDebug("未知类型")
	}
}

//设置值
func FilpValueString(obj map[string]string) map[string]string {
	for i, v := range obj {
		obj[v] = i
	}

	return obj
}
