package controllers

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"reflect"
	"strconv"
	"strings"
	"time"

	"BeeCustom/enums"
	"BeeCustom/models"
	"BeeCustom/utils"
)

type HandBookController struct {
	BaseController
}

func (c *HandBookController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	//默认认证 "Index", "Create", "Edit", "Delete"
	c.checkAuthor()

	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()

}

func (c *HandBookController) Index() {

	params := models.NewCompanyQueryParam()
	limit, err := c.GetInt64("limit", 10)
	offset, err := c.GetInt64("offset", 1)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "关联关系获取失败", nil)
	}

	searchWork := c.GetString("searchWork", "")
	params.SearchWork = searchWork
	params.Limit = limit
	params.Offset = offset

	companies, count := models.CompanyPageList(&params)

	cs, err := models.CompanyGetRelations(companies, "HandBooks")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "关联关系获取失败", nil)
	}
	//页面模板设置
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "handbook/index_footerjs.html"
	c.Data["m"] = cs
	c.Data["count"] = count
	c.Data["searchWork"] = searchWork

	//页面里按钮权限控制
	c.getActionData("Delete", "Import")

	c.GetXSRFToken()
}

//删除
func (c *HandBookController) Delete() {
	id, _ := c.GetInt64(":id")
	if num, err := models.HandBookDelete(id); err == nil {
		c.SetLastUpdteTime("handBookLastUpdateTime", time.Now().Format(enums.BaseFormat))
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), "")
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", err)
	}
}

//导入
func (c *HandBookController) Import() {

	handBookType, err := c.GetInt8(":type", -1)
	if err != nil || handBookType == -1 {
		utils.LogDebug(fmt.Sprintf("GetInt8:%v", err))
		c.jsonResult(enums.JRCodeFailed, "参数错误", nil)
	}

	xmlTitle := c.GetString("xmlTitle", "")
	if len(xmlTitle) == 0 {
		c.jsonResult(enums.JRCodeFailed, "请设置表头", nil)
	}

	_, err = models.HandBookDeleteAll(handBookType)
	if err != nil || handBookType == -1 {
		utils.LogDebug(fmt.Sprintf("HandBookDeleteAll:%v", err))
		c.jsonResult(enums.JRCodeFailed, "清空数据报错", nil)
	}

	fileType := "handBook/" + strconv.FormatInt(c.curUser.Id, 10) + "/"
	fileNamePath, err := c.BaseUpload(fileType)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("BaseUpload:%v", err))
		c.jsonResult(enums.JRCodeFailed, "上传失败", nil)
	}

	cDatas := make([]*models.HandBook, 0)
	handBook := models.HandBook{}

	info := c.ImportHandBookXlsx(handBook, handBookType, fileNamePath, xmlTitle)
	cDatas, err = c.GetXlsxContent(info, cDatas, &handBook)

	if err != nil {
		utils.LogDebug(fmt.Sprintf("GetXlsxContent:%v", err))
		c.jsonResult(enums.JRCodeFailed, "上传失败", nil)
	}

	mun, err := models.InsertHandBookMulti(cDatas)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("InsertMulti:%v", err))
		c.jsonResult(enums.JRCodeFailed, "上传失败", nil)
	}

	c.SetLastUpdteTime("handBookLastUpdateTime", time.Now().Format(enums.BaseFormat))
	c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("上传成功%d项基础参数", mun), mun)

}

//获取 xlsx 文件内容
func (c *HandBookController) GetXlsxContent(info []map[string]string, obj []*models.HandBook, handBook *models.HandBook) ([]*models.HandBook, error) {
	//忽略标题行
	for i := 1; i < len(info); i++ {
		t := reflect.ValueOf(handBook).Elem()
		for k, v := range info[i] {
			switch t.FieldByName(k).Kind() {
			case reflect.String:
				t.FieldByName(k).Set(reflect.ValueOf(v))
			case reflect.Float64:
				handBookV, err := strconv.ParseFloat(v, 64)
				if err != nil {
					utils.LogDebug(fmt.Sprintf("ParseFloat:%v", err))
					return nil, err
				}
				t.FieldByName(k).Set(reflect.ValueOf(handBookV))
			case reflect.Uint64:
				reflect.ValueOf(v)
				handBookV, err := strconv.ParseUint(v, 0, 64)
				if err != nil {
					utils.LogDebug(fmt.Sprintf("ParseUint:%v", err))
					return nil, err
				}
				t.FieldByName(k).Set(reflect.ValueOf(handBookV))
			case reflect.Struct:
				handBookV, err := time.Parse("2006-01-02", v)
				if err != nil {
					utils.LogDebug(fmt.Sprintf("Parse:%v", err))
					return nil, err
				}
				t.FieldByName(k).Set(reflect.ValueOf(handBookV))
			default:
				utils.LogDebug("未知类型")
			}
		}

		obj = append(obj, handBook)

	}

	return obj, nil
}

//导入基础参数 xlsx 文件内容
func (c *HandBookController) ImportHandBookXlsx(handBook models.HandBook, handBookType int8, fileNamePath, xmlTitle string) []map[string]string {

	xmlTitles := strings.Split(xmlTitle, "/")
	rXmlTitles := map[string]int{}
	for k, v := range xmlTitles {
		rXmlTitles[v] = k
	}

	f, err := excelize.OpenFile(fileNamePath)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("导入失败:%v", err))
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}

	if f != nil {
		// Get all the rows in the Sheet1.
		rows, err := f.GetRows("Sheet1")

		if err != nil {
			utils.LogDebug(fmt.Sprintf("导入失败:%v", err))
			c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
		}

		var Info []map[string]string
		for _, row := range rows {
			//将数组  转成对应的 map
			var info = make(map[string]string)
			// 模型前两个字段是 BaseModel ，Type 不需要赋值
			for i := 0; i < reflect.ValueOf(handBook).NumField(); i++ {
				obj := reflect.TypeOf(handBook).Field(i)
				switch obj.Name {
				case "Type":
					info[obj.Name] = string(handBookType)
				case "CustomsCode":
					funcName(rXmlTitles, info, obj, row, "CustomsCode")
				case "Name":
					funcName(rXmlTitles, info, obj, row, "Name")
				case "ShortName":
					funcName(rXmlTitles, info, obj, row, "ShortName")
				case "EnName":
					funcName(rXmlTitles, info, obj, row, "EnName")
				case "InspectionCode":
					funcName(rXmlTitles, info, obj, row, "InspectionCode")
				case "ShortEnName":
					funcName(rXmlTitles, info, obj, row, "ShortEnName")
				case "MandatoryLevel":
					funcName(rXmlTitles, info, obj, row, "MandatoryLevel")
				case "CertificateType":
					funcName(rXmlTitles, info, obj, row, "CertificateType")
				case "StatisticalUnitCode":
					funcName(rXmlTitles, info, obj, row, "StatisticalUnitCode")
				case "ConversionRate":
					funcName(rXmlTitles, info, obj, row, "ConversionRate")
				case "NatureMark":
					funcName(rXmlTitles, info, obj, row, "NatureMark")
				case "Iso2":
					funcName(rXmlTitles, info, obj, row, "Iso2")
				case "Iso3":
					funcName(rXmlTitles, info, obj, row, "Iso3")
				case "TypeCode":
					funcName(rXmlTitles, info, obj, row, "TypeCode")
				case "OldCustomCode":
					funcName(rXmlTitles, info, obj, row, "OldCustomCode")
				case "OldCustomName":
					funcName(rXmlTitles, info, obj, row, "OldCustomName")
				case "OldCiqCode":
					funcName(rXmlTitles, info, obj, row, "OldCiqCode")
				case "OldCiqName":
					funcName(rXmlTitles, info, obj, row, "OldCiqName")
				case "Remark":
					funcName(rXmlTitles, info, obj, row, "Remark")
				}

			}

			Info = append(Info, info)
		}

		return Info

	} else {
		utils.LogDebug(fmt.Sprintf("导入失败:%v", err))
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}

	return nil

}
