package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"reflect"
	"strconv"
	"strings"
	"time"

	"BeeCustom/enums"
	"BeeCustom/models"
	"BeeCustom/utils"
	"github.com/astaxie/beego"
)

type ClearanceController struct {
	BaseController
}

func (c *ClearanceController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	//默认认证 "Index", "Create", "Edit", "Delete"
	c.checkAuthor()

	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()

}

func (c *ClearanceController) Index() {

	c.Data["type"] = strings.Split(beego.AppConfig.String("clearance::type"), ",")
	c.Data["lastUpdateTime"] = c.GetLastUpdteTime("clearanceLastUpdateTime")

	//页面模板设置
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "clearance/index_footerjs.html"

	//页面里按钮权限控制
	c.getActionData("Edit", "Delete", "Create", "Import")

	c.GetXSRFToken()
}

//列表数据
func (c *ClearanceController) DataGrid() {
	//直接获取参数 getDataGridData()
	params := models.NewClearanceQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	//获取数据列表和总数
	data, total := models.ClearancePageList(&params)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	result["code"] = 0
	c.Data["json"] = result

	c.ServeJSON()
}

// Create 添加 新建 页面
func (c *ClearanceController) Create() {
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "clearance/create_footerjs.html"
	c.Data["type"] = strings.Split(beego.AppConfig.String("clearance::type"), ",")
	c.GetXSRFToken()
}

// Store 添加 新建 页面
func (c *ClearanceController) Store() {
	m := models.NewClearance(0)
	//获取form里的值
	if err := c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}

	c.validRequestData(m)

	if _, err := models.ClearanceSave(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "添加失败", m)
	} else {
		c.SetLastUpdteTime("clearanceLastUpdateTime", time.Now().Format(enums.BaseFormat))
		c.jsonResult(enums.JRCodeSucc, "添加成功", m)
	}
}

// Edit 添加 编辑 页面
func (c *ClearanceController) Edit() {
	Id, _ := c.GetInt64(":id", 0)
	m, err := models.ClearanceOne(Id)
	if m != nil && Id > 0 {
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}

	c.Data["m"] = m
	c.Data["type"] = strings.Split(beego.AppConfig.String("clearance::type"), ",")
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "clearance/edit_footerjs.html"
	c.GetXSRFToken()
}

// Update 添加 编辑 页面
func (c *ClearanceController) Update() {
	Id, _ := c.GetInt64(":id", 0)
	m := models.NewClearance(Id)

	//获取form里的值
	if err := c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}

	c.validRequestData(m)

	if _, err := models.ClearanceSave(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "编辑失败", m)
	} else {
		c.SetLastUpdteTime("clearanceLastUpdateTime", time.Now().Format(enums.BaseFormat))
		c.jsonResult(enums.JRCodeSucc, "编辑成功", m)
	}
}

//删除
func (c *ClearanceController) Delete() {
	id, _ := c.GetInt64(":id")
	if num, err := models.ClearanceDelete(id); err == nil {
		c.SetLastUpdteTime("clearanceLastUpdateTime", time.Now().Format(enums.BaseFormat))
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), "")
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", err)
	}
}

//导入
func (c *ClearanceController) Import() {

	clearanceType, err := c.GetInt8(":type", -1)
	if err != nil || clearanceType == -1 {
		utils.LogDebug(fmt.Sprintf("GetInt8:%v", err))
		c.jsonResult(enums.JRCodeFailed, "参数错误", nil)
	}

	xmlTitle := c.GetString("xmlTitle", "")
	if len(xmlTitle) == 0 {
		c.jsonResult(enums.JRCodeFailed, "请设置表头", nil)
	}

	_, err = models.ClearanceDeleteAll(clearanceType)
	if err != nil || clearanceType == -1 {
		utils.LogDebug(fmt.Sprintf("ClearanceDeleteAll:%v", err))
		c.jsonResult(enums.JRCodeFailed, "清空数据报错", nil)
	}

	fileType := "clearance/" + strconv.FormatInt(c.curUser.Id, 10) + "/"
	fileNamePath, err := c.BaseUpload(fileType)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("BaseUpload:%v", err))
		c.jsonResult(enums.JRCodeFailed, "上传失败", nil)
	}

	cDatas := make([]*models.Clearance, 0)
	clearance := models.Clearance{}

	info := c.ImportClearanceXlsx(clearance, clearanceType, fileNamePath, xmlTitle)
	cDatas, err = c.GetXlsxContent(info, cDatas, &clearance)

	if err != nil {
		utils.LogDebug(fmt.Sprintf("GetXlsxContent:%v", err))
		c.jsonResult(enums.JRCodeFailed, "上传失败", nil)
	}

	mun, err := models.InsertClearanceMulti(cDatas)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("InsertMulti:%v", err))
		c.jsonResult(enums.JRCodeFailed, "上传失败", nil)
	}

	c.SetLastUpdteTime("clearanceLastUpdateTime", time.Now().Format(enums.BaseFormat))
	c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("上传成功%d项基础参数", mun), mun)

}

//获取 xlsx 文件内容
func (c *ClearanceController) GetXlsxContent(info []map[string]string, obj []*models.Clearance, clearance *models.Clearance) ([]*models.Clearance, error) {
	//忽略标题行
	for i := 1; i < len(info); i++ {
		t := reflect.ValueOf(clearance).Elem()
		for k, v := range info[i] {
			switch t.FieldByName(k).Kind() {
			case reflect.String:
				t.FieldByName(k).Set(reflect.ValueOf(v))
			case reflect.Float64:
				clearanceV, err := strconv.ParseFloat(v, 64)
				if err != nil {
					utils.LogDebug(fmt.Sprintf("ParseFloat:%v", err))
					return nil, err
				}
				t.FieldByName(k).Set(reflect.ValueOf(clearanceV))
			case reflect.Uint64:
				reflect.ValueOf(v)
				clearanceV, err := strconv.ParseUint(v, 0, 64)
				if err != nil {
					utils.LogDebug(fmt.Sprintf("ParseUint:%v", err))
					return nil, err
				}
				t.FieldByName(k).Set(reflect.ValueOf(clearanceV))
			case reflect.Struct:
				clearanceV, err := time.Parse("2006-01-02", v)
				if err != nil {
					utils.LogDebug(fmt.Sprintf("Parse:%v", err))
					return nil, err
				}
				t.FieldByName(k).Set(reflect.ValueOf(clearanceV))
			default:
				utils.LogDebug("未知类型")
			}
		}

		obj = append(obj, clearance)

	}

	return obj, nil
}

//导入基础参数 xlsx 文件内容
func (c *ClearanceController) ImportClearanceXlsx(clearance models.Clearance, clearanceType int8, fileNamePath, xmlTitle string) []map[string]string {

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
			for i := 0; i < reflect.ValueOf(clearance).NumField(); i++ {
				obj := reflect.TypeOf(clearance).Field(i)
				switch obj.Name {
				case "Type":
					info[obj.Name] = string(clearanceType)
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
