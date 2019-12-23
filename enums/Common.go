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
	"reflect"
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

	i, err, done := TransformCnToInt(sections, wordCh)
	if done {
		return i, err
	}

	return -1, errors.New("查询参数错误")
}

func TransformCnToInt(sections map[string]string, wordCh string) (int8, error, bool) {
	for i, v := range sections {
		if v == wordCh {

			sectionI, err := strconv.Atoi(i)
			if err != nil {
				utils.LogDebug(fmt.Sprintf("ParseInt:%v", err))
				return -1, err, true
			}

			return int8(sectionI), nil, true

		}
	}
	return 0, nil, false
}

// 根据参数查询对应中文
func GetSectionWithInt(wordInt int8, configSection string) (string, error) {
	sections, err := beego.AppConfig.GetSection(configSection)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("GetSection:%v", err))
	}

	s, err, done := TransformIntToCn(sections, wordInt)
	if done {
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

// 设置值 slice
func SetObjValueFromSlice(inObj interface{}, Info []map[string]string) {
	for i := 0; i < len(Info); i++ {
		SetObjValue(inObj, Info[i])
	}
}

// 设置值
func SetObjValue(inObj interface{}, Info map[string]string) {
	t := reflect.ValueOf(inObj).Elem()
	for k, v := range Info {
		SetObjValueIn(k, v, t)
	}
}

// 设置值
func SetObjValueIn(objName, v string, t reflect.Value) {
	switch t.FieldByName(objName).Kind() {
	case reflect.String:
		t.FieldByName(objName).Set(reflect.ValueOf(v))
	case reflect.Float64:
		if len(v) > 0 {
			objV, err := strconv.ParseFloat(v, 64)
			if err != nil {
				utils.LogDebug(fmt.Sprintf("Parse:%v,%v,%v", err, v, objName))
			}
			t.FieldByName(objName).Set(reflect.ValueOf(objV))
		}
	case reflect.Int8:
		if len(v) > 0 {
			objV, err := strconv.Atoi(v)
			if err != nil {
				utils.LogDebug(fmt.Sprintf("Parse:%v,%v,%v", err, v, objName))
			}
			t.FieldByName(objName).Set(reflect.ValueOf(int8(objV)))
		}
	case reflect.Uint64:
		reflect.ValueOf(v)
		objV, err := strconv.ParseUint(v, 0, 64)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("Parse:%v,%v,%v", err, v, objName))
		}
		t.FieldByName(objName).Set(reflect.ValueOf(objV))
	case reflect.Struct:
		if len(v) > 0 {
			objV, err := time.Parse("20060102", v)
			if err != nil {
				utils.LogDebug(fmt.Sprintf("Parse:%v,%v,%v", err, v, objName))
			}
			t.FieldByName(objName).Set(reflect.ValueOf(objV))
		}

	default:
		utils.LogDebug(fmt.Sprintf("未知类型:%v,%v", v, objName))
	}
}

// 设置值
func SetObjValueFromObj(outObj interface{}, inObj interface{}) {

	outObjE := reflect.ValueOf(outObj).Elem()
	outObjET := outObjE.Type()

	inObjE := reflect.ValueOf(inObj).Elem()
	inObjET := inObjE.Type()

	for i := 0; i < outObjE.NumField(); i++ {

		outObjEF := outObjE.Field(i)

		for iI := 0; iI < inObjE.NumField(); iI++ {

			inObjEF := inObjE.Field(iI)

			if outObjET.Field(i).Name == inObjET.Field(iI).Name && outObjEF.Type() == inObjEF.Type() {
				if outObjEF.CanSet() {
					switch inObjEF.Kind() {
					case reflect.String:
						outObjEF.SetString(inObjEF.String())
					case reflect.Bool:
						outObjEF.SetBool(inObjEF.Bool())
					case reflect.Float64, reflect.Float32:
						outObjEF.SetFloat(inObjEF.Float())
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						outObjEF.SetInt(inObjEF.Int())
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						outObjEF.SetUint(inObjEF.Uint())
					case reflect.Struct:
						SetObjValueFromObj(outObjEF, inObjEF)
					default:
						utils.LogDebug(fmt.Sprintf("未知类型:%v,%v", inObjEF.Kind(), inObjEF))
					}
				}
			}
		}

	}

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
