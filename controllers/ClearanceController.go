package controllers

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"

	"BeeCustom/enums"
	"BeeCustom/models"
	"BeeCustom/utils"
	"github.com/astaxie/beego"
)

type ClearanceController struct {
	BaseController

	clearanceType map[string]string
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
	clearanceType, err := beego.AppConfig.GetSection("clearance_type")
	if err != nil {
		utils.LogDebug(fmt.Sprintf("clearance_type:%v", err))
	}

	c.clearanceType = clearanceType

}

func (c *ClearanceController) Index() {

	c.Data["type"] = c.clearanceType
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
	c.Data["type"] = c.clearanceType
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
	c.Data["type"] = c.clearanceType
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

	cIP := models.ClearanceImportParam{}
	cIP.Clearance = models.NewClearance(0)
	cIP.Obj = make([]*models.Clearance, 0)

	clearanceType, err := c.GetInt8(":type", -1)
	if err != nil || clearanceType == -1 {
		utils.LogDebug(fmt.Sprintf("GetInt8:%v", err))
		c.jsonResult(enums.JRCodeFailed, "参数错误", nil)
	}

	cIP.Clearance.Type = clearanceType

	xmlTitle := c.GetString("xmlTitle", "")
	_, err = models.ClearanceDeleteAll(cIP.Clearance.Type)

	if err != nil || cIP.Clearance.Type == -1 {
		utils.LogDebug(fmt.Sprintf("ClearanceDeleteAll:%v", err))
		c.jsonResult(enums.JRCodeFailed, "清空数据报错", nil)
	}

	cIP.XmlTitle = xmlTitle

	fileType := "clearance/" + strconv.FormatInt(c.curUser.Id, 10) + "/"
	fileNamePath, err := c.BaseUpload(fileType)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("BaseUpload:%v", err))
		c.jsonResult(enums.JRCodeFailed, "上传失败", nil)
	}

	cIP.FileNamePath = fileNamePath

	c.ImportClearanceXlsx(&cIP)
	err = c.GetXlsxContent(&cIP)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("GetXlsxContent:%v", err))
		c.jsonResult(enums.JRCodeFailed, "上传失败", nil)
	}

	mun, err := models.InsertClearanceMulti(cIP.Obj)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("InsertMulti:%v", err))
		c.jsonResult(enums.JRCodeFailed, "上传失败", nil)
	}

	c.SetLastUpdteTime("clearanceLastUpdateTime", time.Now().Format(enums.BaseFormat))
	c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("上传成功%d项基础参数", mun), mun)

}

//获取 xlsx 文件内容
func (c *ClearanceController) GetXlsxContent(cIP *models.ClearanceImportParam) error {
	//忽略标题行
	for i := 1; i < len(cIP.Info); i++ {
		t := reflect.ValueOf(&cIP.Clearance).Elem()
		for k, v := range cIP.Info[i] {
			enums.SetObjValue(k, v, t)
		}
		cIP.Obj = append(cIP.Obj, &cIP.Clearance)
	}

	return nil
}

//导入基础参数 xlsx 文件内容
func (c *ClearanceController) ImportClearanceXlsx(cIP *models.ClearanceImportParam) {

	rXmlTitles, _ := utils.GetRXmlTitles(cIP.XmlTitle, "clearance_excel_tile")

	f, err := excelize.OpenFile(cIP.FileNamePath)
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

		if err != nil {
			utils.LogDebug(fmt.Sprintf("GetSection:%v", err))
			c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
		}

		for _, row := range rows {
			//将数组  转成对应的 map
			var info = make(map[string]string)
			// 模型前两个字段是 BaseModel ，Type 不需要赋值
			for i := 0; i < reflect.ValueOf(cIP.Clearance).NumField(); i++ {
				obj := reflect.TypeOf(cIP.Clearance).Field(i)
				for _, iw := range rXmlTitles {
					if iw == obj.Name {
						rI := funcName(rXmlTitles, iw)
						// 模板字段数量定义
						if rI != -1 && rI <= len(row) {
							info[obj.Name] = row[rI]
						}
					}

				}
			}

			cIP.Info = append(cIP.Info, info)
		}

	} else {
		utils.LogDebug(fmt.Sprintf("导入失败:%v", err))
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}

}
