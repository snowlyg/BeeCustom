package enums

import (
	"errors"
	"fmt"
	"strconv"

	"BeeCustom/utils"
	"github.com/astaxie/beego"
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

//根据中文查询对应参数
func GetSectionWithString(wordCh, configSection string) (int8, error) {
	sections, err := beego.AppConfig.GetSection(configSection)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("GetSection:%v", err))
	}

	for i, v := range sections {
		if v == wordCh {

			sectionI, err := strconv.Atoi(i)
			if err != nil {
				utils.LogDebug(fmt.Sprintf("ParseInt:%v", err))
				return -1, err
			}

			return int8(sectionI), nil

		}
	}

	return -1, errors.New("查询参数错误")
}

//根据参数查询对应中文
func GetSectionWithInt(wordInt int8, configSection string) (string, error) {
	sections, err := beego.AppConfig.GetSection(configSection)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("GetSection:%v", err))
	}

	for i, v := range sections {
		sectionI, err := strconv.Atoi(i)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("ParseInt:%v", err))
			return "", err
		}

		if int8(sectionI) == wordInt {
			return v, nil

		}
	}

	return "", errors.New("查询参数错误")
}
