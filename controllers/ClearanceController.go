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
)

type ClearanceController struct {
	BaseController

	clearanceType map[string]string
}

func (c *ClearanceController) Prepare() {
	// 先执行
	c.BaseController.Prepare()
	// 如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	perms := []string{
		"Index",
		"Create",
		"Edit",
		"Delete",
		"Import",
	}
	c.checkAuthor(perms)

	// 如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	// 权限控制里会进行登录验证，因此这里不用再作登录验证
	// c.checkLogin()
	c.clearanceType, _ = models.GetSettingRValueByKey("ClearanceTypes")

}

func (c *ClearanceController) Index() {
	c.Data["type"] = c.clearanceType
	// 页面模板设置
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "clearance/index_footerjs.html"

	// 页面里按钮权限控制
	c.getActionData("", "Edit", "Delete", "Create", "Import")

	c.GetXSRFToken()
}

// 列表数据
func (c *ClearanceController) DataGrid() {
	// 直接获取参数 getDataGridData()
	params := models.NewClearanceQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	// 获取数据列表和总数
	data, total := models.ClearancePageList(&params)
	c.ResponseList(data, total)
	c.ServeJSON()
}

// 通关参数更新时间
func (c *ClearanceController) GetClearanceUpdateTimeByType() {
	clearanceType, _ := c.GetInt8(":type")
	format := "超过一个月未更新"
	lastUpdateTime, err := models.GetLastUpdteTimeByClearanceType(clearanceType)
	if err == nil {
		if lastUpdateTime != nil {
			format = lastUpdateTime.LastUpdatedAt.Format("2006-01-02 15:04:05")
		}
	}

	c.Data["json"] = format

	c.ServeJSON()
}

//  Create 添加 新建 页面
func (c *ClearanceController) Create() {
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "clearance/create_footerjs.html"
	c.Data["type"] = c.clearanceType
	c.GetXSRFToken()
}

//  Store 添加 新建 页面
func (c *ClearanceController) Store() {
	m := models.NewClearance(0)
	// 获取form里的值
	if err := c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}
	c.validRequestData(m)
	if nm, err := models.ClearanceSave(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "添加失败", nil)
	} else {
		c.setLastUpdteTime(nm.Type)
		c.jsonResult(enums.JRCodeSucc, "添加成功", nil)
	}
}

//  Edit 添加 编辑 页面
func (c *ClearanceController) Edit() {
	Id, _ := c.GetInt64(":id", 0)
	m, err := models.ClearanceOne(Id)
	if m != nil && Id > 0 {
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}

	c.Data["m"] = m
	c.Data["cType"] = strconv.Itoa(int(m.Type))
	c.Data["type"] = c.clearanceType
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "clearance/edit_footerjs.html"
	c.GetXSRFToken()
}

//  commonClearance
func (c *ClearanceController) CommonClearance() {

	arg := []int8{1, 3, 4, 5, 9, 19, 23}
	data := models.ClearancePageListInTypes(arg)

	jsonData := c.transforClearance(data, arg)
	c.Data["json"] = jsonData
	c.ServeJSON()
}

// 格式基础参数
func (c *ClearanceController) transforClearance(data []*models.Clearance, arg []int8) map[int8][]map[string]string {
	jsonData := map[int8][]map[string]string{}
	for _, i := range arg {
		var clearances []map[string]string
		for _, v := range data {
			if i == v.Type {
				clearance := map[string]string{}
				clearance["Name"] = v.Name
				clearance["CustomsCode"] = v.CustomsCode
				clearance["OldCustomName"] = v.OldCustomName
				clearance["OldCustomCode"] = v.OldCustomCode
				clearances = append(clearances, clearance)
			}
			jsonData[i] = clearances
		}
	}

	return jsonData
}

//  annotationClearance
func (c *ClearanceController) AnnotationClearance() {
	arg := []int8{33, 34, 35, 36, 37, 38, 40, 41}
	data := models.ClearancePageListInTypes(arg)
	jsonData := c.transforClearance(data, arg)
	c.Data["json"] = jsonData
	c.ServeJSON()
}

//  orderClearance
func (c *ClearanceController) OrderClearance() {
	arg := []int8{2, 6, 7, 8, 10, 11, 12, 13, 14, 15, 16, 17, 18, 20, 21, 22, 24, 25, 26, 27, 28, 29, 30, 31, 32, 39, 42}
	data := models.ClearancePageListInTypes(arg)
	jsonData := c.transforClearance(data, arg)
	c.Data["json"] = jsonData
	c.ServeJSON()
}

//  Update 添加 编辑 页面
func (c *ClearanceController) Update() {
	Id, _ := c.GetInt64(":id", 0)
	m := models.NewClearance(Id)

	// 获取form里的值
	if err := c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}

	c.validRequestData(m)

	if _, err := models.ClearanceSave(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "编辑失败", nil)
	} else {
		nm, _ := models.ClearanceOne(Id)
		utils.LogDebug(nm)
		c.setLastUpdteTime(nm.Type)
		c.jsonResult(enums.JRCodeSucc, "编辑成功", nil)
	}
}

// 删除
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

// 导入
func (c *ClearanceController) Import() {

	clearanceType, err := c.GetInt8(":type", -1)
	if err != nil || clearanceType == -1 {
		c.jsonResult(enums.JRCodeFailed, "参数错误", nil)
	}

	_, err = models.ClearanceDeleteAll(clearanceType) // 清空对应基础参数
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

// 导入基础参数 xlsx 文件内容
func (c *ClearanceController) ImportClearanceXlsx(cIP *models.ClearanceImportParam) {

	rXmlTitles, err := xlsx.GetExcelTitles("", "clearance_excel_tile")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}

	rows, err := xlsx.GetExcelRows(cIP.FileNamePath, cIP.ExcelName)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}

	// 提取 excel 数据
	var Info []map[string]string
	for roI, row := range rows {
		if roI > 0 {
			info := make(map[string]string)
			// 将数组  转成对应的 map
			inObj := models.NewClearance(0)
			for i := 0; i < reflect.ValueOf(inObj).NumField(); i++ {
				obj := reflect.TypeOf(inObj).Field(i)
				for _, iw := range rXmlTitles {
					if iw == obj.Name {
						rI := xlsx.ObjIsExists(rXmlTitles, iw)
						//  模板字段数量定义
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

	// 转换 excel 数据
	// 忽略标题行
	for i := 0; i < len(Info); i++ {
		inObj := models.NewClearance(0)
		inObj.Type = cIP.ClearanceType
		enums.SetObjValue(&inObj, Info, i)
		cIP.Obj = append(cIP.Obj, &inObj)
	}

}

// 获取最后更新时间
func (c *ClearanceController) getLastUpdteTime() []*models.ClearanceUpdateTime {
	// 直接获取参数 getDataGridData()
	params := models.NewClearanceUpdateTimeQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	// 获取数据列表和总数
	data, _ := models.ClearanceUpdateTimePageList(&params)

	return data
}

// 设置最后更新时间
func (c *ClearanceController) setLastUpdteTime(cType int8) {
	if cType == 0 {
		c.jsonResult(enums.JRCodeFailed, "类型错误", nil)
	}

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
		oldLastUpdateTime.Type = cType
		_, err = models.ClearanceUpdateTimeSave(oldLastUpdateTime)
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, "设置最后更新时间失败", nil)
		}
	}

}
