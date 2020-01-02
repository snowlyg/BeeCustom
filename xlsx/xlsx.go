package xlsx

import (
	"errors"
	"fmt"
	"strconv"

	"BeeCustom/utils"
	"github.com/360EntSecGroup-Skylar/excelize"
)

// BaseImport
type BaseImport struct {
	FileNamePath string
	ExcelTitle   map[string]string
	ExcelName    string
}

// 导入基础参数 rows 文件内容
func GetExcelRows(fileNamePath, excelName string) ([][]string, error) {

	f, err := excelize.OpenFile(fileNamePath)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("GetExcelRows.OpenFile:%v", err))
		return nil, err
	}

	if f == nil {
		return nil, errors.New("excelize.OpenFile 出错")
	}

	rows := f.GetRows(excelName)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("GetExcelRows.OpenFile:%v", err))
		return nil, err
	}

	return rows, nil

}

// 导入基础参数 Cell 文件内容
func GetExcel(fileNamePath string) (*excelize.File, error) {

	f, err := excelize.OpenFile(fileNamePath)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("GetExcelCell.OpenFile:%v", err))
		return nil, err
	}

	if f == nil {
		return nil, errors.New("excelize.OpenFile 出错")
	}

	return f, nil
}

// 设置值
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
