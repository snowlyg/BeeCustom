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

func getAnnotationXmlNames(pathConfig string, pathNames []string) []string {
	path := beego.AppConfig.String(pathConfig)
	fullpath := path + time.Now().Format(enums.BaseDateFormatN)
	if file.IsExist(fullpath) {
		pathCfiles, err := ioutil.ReadDir(fullpath)
		if err != nil {
			utils.LogError(fmt.Sprintf("ioutil.ReadDir :%v", err))
			return nil
		}
		for _, f := range pathCfiles {
			name := getAnnotationXmlName(f)
			if len(name) > 0 {
				pathNames = append(pathNames, name)
			}
		}

	}
	return pathNames
}

func getAnnotationXmlName(f os.FileInfo) string {
	names := strings.Split(f.Name(), `__`)
	if len(names) >= 2 && len(names[1]) > 0 {
		namesS := strings.Split(strings.Replace(names[1], `.xml`, "", -1), `_`)
		return namesS[0]
	}
	return ""
}
