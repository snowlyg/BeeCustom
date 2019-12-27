package enums

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"BeeCustom/utils"
	"github.com/astaxie/beego"
)

type JsonResultCode int

const (
	JRCodeFailed JsonResultCode = iota // 接口返回状态 0
	JRCodeSucc                         // 接口返回状态 1
	JRCode302    = 302                 // 跳转至地址
	JRCode401    = 401                 // 未授权访问
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
const RFC3339 = "2006-01-02T15:04:05"

// 根据中文查询对应参数
func GetSectionWithString(wordCh, configSection string) (int8, error) {
	sections, err := beego.AppConfig.GetSection(configSection)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("GetSection:%v", err))
	}

	i, err := TransformCnToInt(sections, wordCh)
	if err != nil {
		return i, err
	}

	return -1, errors.New("查询参数错误")
}

func TransformCnToInt(ss map[string]string, s string) (int8, error) {

	for i, v := range ss {
		if v == s {
			si, err := strconv.Atoi(i)
			if err != nil {
				utils.LogDebug(fmt.Sprintf("ParseInt:%v", err))
				return -1, err
			}

			return int8(si), nil

		}
	}

	return 0, errors.New("not in")
}

// 根据参数查询对应中文
func GetSectionWithInt(wordInt int8, configSection string) (string, error) {
	sections, err := beego.AppConfig.GetSection(configSection)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("GetSection:%v", err))
	}

	s, err, done := TransformIntToCn(sections, wordInt)
	if !done {
		return s, err
	}

	return "", errors.New("查询参数错误")
}

func TransformIntToCn(sections map[string]string, wordInt int8) (string, error, bool) {
	for i, v := range sections {
		sectionI, err := strconv.Atoi(i)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("ParseInt:%v", err))
			return "", err, true
		}

		if int8(sectionI) == wordInt {
			return v, nil, true

		}
	}
	return "", nil, false
}

// 获取4位随机数
func CreateCaptcha() string {
	return fmt.Sprintf("%04v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000))
}

// 判断时间，格式时间
func GetDateTimeString(v *time.Time, format string) string {
	if v.IsZero() {
		return ""
	} else {
		return v.Format(format)
	}
}

// 返回进出口中文
func GetImpexpMarkcdCNName(impexpMarkcd string) string {

	if impexpMarkcd == "I" {
		return "进口"
	} else if impexpMarkcd == "E" {
		return "出口"
	} else {
		return ""
	}
}

// 获取时间段
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

// string slice in
func InStringArray(s string, sS []string) bool {
	for _, v := range sS {
		if v == s {
			return true
		}
	}

	return false
}

// string map in
func InStringMap(s string, sS map[string]string) bool {
	for _, v := range sS {
		if v == s {
			return true
		}
	}

	return false
}

//  hmac 加密
func Hmac(key string, data []byte) string {
	hmacSha1 := hmac.New(sha1.New, []byte(key))
	hmacSha1.Write(data)
	return hex.EncodeToString(hmacSha1.Sum([]byte("")))
}

func Cmd(action, input string, arg []string) {
	cmd := exec.Command(action, arg...)
	if len(input) > 0 {
		cmd.Stdin = strings.NewReader(input)
	}
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		utils.LogDebug(fmt.Sprintf("cmd:%v:%v--%v --%v", err, action, arg, stderr.String()))
	}
}

// if 0 to ""
func IsFloatZore(f float64) string {
	floatString := strconv.FormatFloat(f, 'f', 0, 64)
	if floatString == "0" {
		return ""
	}
	return floatString
}

// if 0 to ""
func IsIZore(i int) string {
	if strconv.Itoa(i) == "0" {
		return ""
	}

	return strconv.Itoa(i)
}

/**
* string to slice
* s: [{"OrderItemLimitOId":"1","OrderItemLimitId":"1","GoodsNo":1,"LicTypeCode":"100","LicTypeName":"通关司类","LicenceNo":"1","LicWrtofDetailNo":"1","LicWrtofQty":"1","LicWrtofQtyUnit":"001","LicWrtofQtyUnit
* Name":"台","OrderItemId":"1","Id":"1"},{"OrderItemLimitOId":"1","OrderItemLimitId":"","GoodsNo":2,"LicTypeCode":"100","LicTypeName":"通关司类","LicenceNo":"1","LicWrtofDetailNo":"1","LicWrtofQty":"1","LicWrtofQtyUnit":"001","LicWrto
* fQtyUnitName":"台","OrderItemId":"1","Id":2},{"OrderItemLimitOId":"1","OrderItemLimitId":"","GoodsNo":3,"LicTypeCode":"000","LicTypeName":"企业产品许可类别","LicenceNo":"1","LicWrtofDetailNo":"1","LicWrtofQty":"1","LicWrtofQtyUnit":
* "001","LicWrtofQtyUnitName":"台","OrderItemId":"1","Id":3}]
* start : `[{"`
* end : `}]`
* mid : `},{"`
* imid : `","`
*
 */
func TramsformStringToSlice(s, start, mid, imid, end string) []string {
	sJ := strings.Replace(strings.Replace(strings.Replace(s, imid, "", -1), start, "", -1), end, "", -1)
	sJSlice := strings.Split(sJ, mid)
	return sJSlice
}
