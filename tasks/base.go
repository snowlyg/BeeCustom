package tasks

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"BeeCustom/enums"
	"BeeCustom/file"
	"BeeCustom/utils"
	"github.com/astaxie/beego"
)

func getXmlNames(pathConfig string, pathNames []string) []string {
	path := beego.AppConfig.String(pathConfig)
	filepath := path + time.Now().Format(enums.BaseDateFormatN)
	if file.IsExist(filepath) {
		pathChides, err := ioutil.ReadDir(filepath)
		if err != nil {
			utils.LogError(fmt.Sprintf("ioutil.ReadDir :%v", err))
			return nil
		}
		for _, f := range pathChides {
			name := getXmlName(f)
			if len(name) > 0 {
				pathNames = append(pathNames, name)
			}
		}

	}
	return pathNames
}

func getXmlName(f os.FileInfo) string {
	names := strings.Split(f.Name(), `__`)
	if len(names) >= 2 && len(names[1]) > 0 {
		namesS := strings.Split(strings.Replace(names[1], `.xml`, "", -1), `_`)
		return namesS[0]
	}
	return ""
}

// 打开文件
func openFile(fullPath string) (*os.File, error, []byte) {
	xmlFile, err := os.Open(fullPath)
	if err != nil {
		utils.LogError(fmt.Sprintf("os.Open :%v", err))
	}
	data, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		utils.LogError(fmt.Sprintf(" ioutil.ReadAll :%v", err))
	}
	return xmlFile, err, data
}

// ws 自动更新
func wsPush() {
	msg := utils.Message{Message: "清单状态更新", IsUpdated: true}
	utils.Broadcast <- msg
}

// 移动文件
func moveFile(historyPath, v, fullPath string, f os.FileInfo) error {
	path := historyPath + time.Now().Format(enums.BaseDateFormat) + "/" + v + "/"
	if err := file.CreateFile(path); err != nil {
		utils.LogError(fmt.Sprintf("文件夹创建失败:%v", err))
	}

	err := os.Rename(fullPath, path+f.Name())
	if err != nil {
		return err
	}

	return nil
}
