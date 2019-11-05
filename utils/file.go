package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
)

//调用os.MkdirAll递归创建文件夹
func CreateFile(filePath string) error {
	if !IsExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		return err
	}
	return nil
}

// 判断所给路径文件/文件夹是否存在(返回true是存在)
func IsExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

//获取导入文件表头
func GetRXmlTitles(xmlTitle, configSection string) (map[string]string, error) {
	rXmlTitles := map[string]string{}
	if len(xmlTitle) == 0 {
		importWord, err := beego.AppConfig.GetSection(configSection)
		if err != nil {
			LogDebug(fmt.Sprintf("GetSection:%v", err))
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
