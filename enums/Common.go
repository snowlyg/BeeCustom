package enums

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

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

const BaseDateTimeFormat = "2006-01-02 15:04:05"
const BaseDateFormatN = "2006-01-02"
const BaseDateTimeSecondFormat = "20060102150405"
const BaseDateFormat = "20060102"

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

//获取4位随机数
func CreateCaptcha() string {
	return fmt.Sprintf("%04v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000))
}

//判断时间，格式时间
func GetDateTimeString(v *time.Time, format string) string {
	if v.IsZero() {
		return ""
	} else {
		return v.Format(format)
	}
}

//返回进出口中文
func GetImpexpMarkcdCNName(impexpMarkcd string) string {

	if impexpMarkcd == "I" {
		return "进口"
	} else if impexpMarkcd == "E" {
		return "出口"
	} else {
		return ""
	}
}

//获取时间段
func GetOrderAnnotationDateTime(timeString, filedName string) string {
	var sql string
	switch timeString {
	case "今天":
		sql = " WHERE TO_DAYS(" + filedName + ") = TO_DAYS(NOW()) "
	case "昨天":
		sql = " WHERE  DATEDIFF(" + filedName + ",NOW()) = -1 "
	case "最近三天":
		sql = " WHERE   DATEDIFF(" + filedName + ",NOW()) <= 0 AND DATEDIFF(" + filedName + ",NOW()) > -3 "
	case "本周":
		sql = " WHERE YEARWEEK(DATE_FORMAT(" + filedName + ",'%Y-%m-%d')) = YEARWEEK(NOW()) "
	case "本月":
		sql = " WHERE DATE_FORMAT(" + filedName + ",'%Y%m') = DATE_FORMAT(CURDATE(),'%Y%m') "
	case "上月":
		sql = " WHERE PERIOD_DIFF(DATE_FORMAT(NOW(),'%Y%m'),DATE_FORMAT(" + filedName + ",'%Y%m')) = 1 "
	case "本季度":
		sql = " WHERE QUARTER(" + filedName + ") = QUARTER(NOW()) "
	case "上季度":
		sql = " WHERE QUARTER(" + filedName + ") = QUARTER(DATE_SUB(NOW(),INTERVAL 1 QUARTER)) "
	case "今年":
		sql = " WHERE YEAR(" + filedName + ")=YEAR(NOW()) "
	case "去年":
		sql = " WHERE YEAR(" + filedName + ") = YEAR(DATE_SUB(NOW(),INTERVAL 1 YEAR)) "
	default:
		sql = " WHERE TO_DAYS(" + filedName + ") = TO_DAYS(NOW()) "
	}

	return sql
}
