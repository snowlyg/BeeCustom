package xlsx

import (
	"errors"
	"fmt"

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
