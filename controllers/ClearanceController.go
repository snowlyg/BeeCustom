package controllers

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"BeeCustom/enums"
	"BeeCustom/models"
	"BeeCustom/utils"
	"BeeCustom/xlsx"
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
	perms := []string{
		"Index",
		"Create",
		"Edit",
		"Delete",
		"Import",
	}
	c.checkAuthor(perms)

	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()

}

func (c *ClearanceController) Index() {

	clearanceType, err := beego.AppConfig.GetSection("clearance_type")
	if err != nil {
		utils.LogDebug(fmt.Sprintf("clearance_type:%v", err))
	}
	c.Data["type"] = clearanceType

	//页面模板设置
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "clearance/index_footerjs.html"

	//页面里按钮权限控制
	c.getActionData("", "Edit", "Delete", "Create", "Import")

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

//通关参数更新时间
func (c *ClearanceController) GetClearanceUpdateTime() {

	lastUpdateTime := c.getLastUpdteTime()
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = 0
	result["rows"] = lastUpdateTime
	result["code"] = 0
	c.Data["json"] = result

	c.ServeJSON()
}

//通关参数更新时间
func (c *ClearanceController) GetClearanceUpdateTimeByType() {
	clearanceType, _ := c.GetInt8(":type")
	format := "超过一个月未更新"
	lastUpdateTime, err := models.GetLastUpdteTimeByClearanceType(clearanceType)
	if err == nil {
		if lastUpdateTime != nil {
			format = lastUpdateTime.LastUpdatedAt.Format("2006-01-02 15:04:05")
		}
	}

	//定义返回的数据结构
	result := make(map[string]interface{})
	result["data"] = format
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
		c.setLastUpdteTime(m.Type)
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
		c.setLastUpdteTime(m.Type)
		c.jsonResult(enums.JRCodeSucc, "编辑成功", m)
	}
}

//删除
func (c *ClearanceController) Delete() {
	id, _ := c.GetInt64(":id")
	m, _ := models.ClearanceOne(id)
	if num, err := models.ClearanceDelete(id); err == nil {
		c.setLastUpdteTime(m.Type)
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), "")
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", err)
	}
}

//导入
func (c *ClearanceController) Import() {

	clearanceType, err := c.GetInt8(":type", -1)
	if err != nil || clearanceType == -1 {
		c.jsonResult(enums.JRCodeFailed, "参数错误", nil)
	}

	_, err = models.ClearanceDeleteAll(clearanceType) //清空对应基础参数
	if err != nil || clearanceType == -1 {
		c.jsonResult(enums.JRCodeFailed, "清空数据报错", nil)
	}

	fileType := "clearance/" + strconv.FormatInt(c.curUser.Id, 10) + "/"
	fileNamePath, err := c.BaseUpload(fileType)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "上传失败", nil)
	}

	clearances := make([]*models.Clearance, 0)
	param := xlsx.BaseImportParam{
		ExcelName:    "Sheet1",
		FileNamePath: fileNamePath,
	}
	cIP := models.ClearanceImportParam{
		Obj:             clearances,
		ClearanceType:   clearanceType,
		BaseImportParam: param,
	}

	c.ImportClearanceXlsx(&cIP)

	mun, err := models.InsertClearanceMulti(cIP.Obj)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}

	c.setLastUpdteTime(clearanceType)
	c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("导入成功 %d 项基础参数", mun), mun)

}

//导入基础参数 xlsx 文件内容
func (c *ClearanceController) ImportClearanceXlsx(cIP *models.ClearanceImportParam) {

	rXmlTitles, err := xlsx.GetExcelTitles("", "clearance_excel_tile")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}

	rows, err := xlsx.GetExcelRows(cIP.FileNamePath, cIP.ExcelName)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}

	//提取 excel 数据
	var Info []map[string]string
	for roI, row := range rows {
		if roI > 0 {
			info := make(map[string]string)
			//将数组  转成对应的 map
			inObj := models.NewClearance(0)
			for i := 0; i < reflect.ValueOf(inObj).NumField(); i++ {
				obj := reflect.TypeOf(inObj).Field(i)
				for _, iw := range rXmlTitles {
					if iw == obj.Name {
						rI := xlsx.ObjIsExists(rXmlTitles, iw)
						// 模板字段数量定义
						if rI != -1 && rI <= len(row) {
							info[obj.Name] = row[rI]
						} else {
							continue
						}
					}
				}
			}

			Info = append(Info, info)
		}

	}

	//转换 excel 数据
	//忽略标题行
	for i := 0; i < len(Info); i++ {
		inObj := models.NewClearance(0)
		inObj.Type = cIP.ClearanceType
		t := reflect.ValueOf(&inObj).Elem()
		for k, v := range Info[i] {
			xlsx.SetObjValue(k, v, t)
		}
		cIP.Obj = append(cIP.Obj, &inObj)
	}

}

//获取最后更新时间
func (c *ClearanceController) getLastUpdteTime() []*models.ClearanceUpdateTime {

	//直接获取参数 getDataGridData()
	params := models.NewClearanceUpdateTimeQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	//获取数据列表和总数
	data, _ := models.ClearanceUpdateTimePageList(&params)

	return data
}

//设置最后更新时间
func (c *ClearanceController) setLastUpdteTime(cType int8) {

	oldLastUpdateTime, err := models.GetLastUpdteTimeByClearanceType(cType)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "设置最后更新时间失败", nil)
	}

	if oldLastUpdateTime == nil {
		lastUpdateTime := models.NewClearanceUpdateTime(0)
		lastUpdateTime.LastUpdatedAt = time.Now()
		lastUpdateTime.Type = cType
		_, err = models.ClearanceUpdateTimeSave(&lastUpdateTime)
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, "设置最后更新时间失败", nil)
		}
	} else {
		oldLastUpdateTime.LastUpdatedAt = time.Now()
		_, err = models.ClearanceUpdateTimeSave(oldLastUpdateTime)
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, "设置最后更新时间失败", nil)
		}
	}

}
