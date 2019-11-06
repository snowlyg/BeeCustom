package xlsx

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"reflect"
	"strconv"
	"strings"
	"time"

	"BeeCustom/utils"
	"github.com/astaxie/beego"
)

// ClearanceImportParam 用于查询的类
type BaseImportParam struct {
	Info         []map[string]string
	FileNamePath string
	ExcelTitle   map[string]string
	ExcelTitle2  map[string]string
	ExcelTitle3  map[string]string
	ExcelTitle4  map[string]string
	ExcelName    string
}

//导入基础参数 xlsx 文件内容
func GetExcelRows(fileNamePath, excelName string) ([][]string, error) {

	f, err := excelize.OpenFile(fileNamePath)
	if err != nil {
		return nil, err
	}

	if f == nil {
		return nil, errors.New("excelize.OpenFile 出错")
	}

	rows, err := f.GetRows(excelName)
	if err != nil {
		return nil, err
	}

	return rows, nil

}

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
		objV, err := strconv.ParseFloat(v, 64)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("ParseFloat:%v,%v", err, v))
		}
		t.FieldByName(objName).Set(reflect.ValueOf(objV))
	case reflect.Uint64:
		reflect.ValueOf(v)
		objV, err := strconv.ParseUint(v, 0, 64)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("ParseUint:%v,%v", err, v))
		}
		t.FieldByName(objName).Set(reflect.ValueOf(objV))
	case reflect.Struct:
		objV, err := time.Parse("20060102", v)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("Parse:%v,%v", err, v))
		}
		t.FieldByName(objName).Set(reflect.ValueOf(objV))
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

// 判断是否存在键
func ObjIsExists(rXmlTitles map[string]string, s string) int {
	fRXmlTitles := FilpValueString(rXmlTitles)
	if _, ok := fRXmlTitles[s]; ok {
		i, err := strconv.Atoi(rXmlTitles[s])
		if err != nil {
			utils.LogDebug(fmt.Sprintf("funcName=>Atoi:%v", err))
		}
		return i
	} else {
		return -1
	}
}
