package controllers

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"BeeCustom/xlsx"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/astaxie/beego"

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

	searchWord := c.GetString("searchWord", "")
	params.SearchWord = searchWord
	params.Limit = limit
	params.Offset = offset

	companies, count := models.CompanyPageList(&params)

	cs, err := models.CompaniesGetRelations(companies, "HandBooks")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "关联关系获取失败", nil)
	}
	//页面模板设置
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "handbook/index_footerjs.html"
	c.Data["m"] = cs
	c.Data["count"] = count
	c.Data["searchWord"] = searchWord

	//页面里按钮权限控制
	c.getActionData("Delete", "Import")

	c.GetXSRFToken()
}

// Edit 添加 编辑 页面
func (c *HandBookController) Show() {
	Id, _ := c.GetInt64(":id", 0)
	m, err := models.HandBookOne(Id, "Company,HandBookGoods")
	if m != nil && Id > 0 {
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}

	c.Data["m"] = m

	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "handbook/show_footerjs.html"
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

	hIP := models.HandBookImportParam{}
	hIP.HandBook = models.NewHandBook(0)

	fileType := "handBook/" + strconv.FormatInt(c.curUser.Id, 10) + "/"
	fileNamePath, err := c.BaseUpload(fileType)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("BaseUpload:%v", err))
		c.jsonResult(enums.JRCodeFailed, "上传失败", nil)
	}

	hIP.FileNamePath = fileNamePath

	fileNamePaths := strings.Split(fileNamePath, ".")
	fileExt := fileNamePaths[len(fileNamePaths)-1]
	if fileExt != "xlsx" {
		c.jsonResult(enums.JRCodeFailed, "文件格式错误，只能导入 xlsx 文件", nil)
	}

	c.ImportHandBookXlsxByCell(&hIP)
	c.ImportHandBookXlsxByRow(&hIP)
	//cDatas, err = c.GetXlsxContent(info, cDatas, &handBook)

	if err != nil {
		utils.LogDebug(fmt.Sprintf("GetXlsxContent:%v", err))
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}

	m, err := models.HandBookSave(&hIP.HandBook)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("InsertMulti:%v", err))
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}

	c.SetLastUpdteTime("handBookLastUpdateTime", time.Now().Format(enums.BaseFormat))
	c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("导入成功"), m.Id)

}

//获取 xlsx 文件内容
func (c *HandBookController) GetXlsxContent(info []map[string]string, obj []*models.HandBook, handBook *models.HandBook) ([]*models.HandBook, error) {
	//忽略标题行
	for i := 1; i < len(info); i++ {
		t := reflect.ValueOf(handBook).Elem()
		for k, v := range info[i] {
			xlsx.SetObjValue(k, v, t)
		}

		obj = append(obj, handBook)

	}

	return obj, nil
}

//导入基础参数 xlsx 文件内容
func (c *HandBookController) ImportHandBookXlsxByCell(hIP *models.HandBookImportParam) {

	f, err := excelize.OpenFile(hIP.FileNamePath)

	if err != nil {
		utils.LogDebug(fmt.Sprintf("OpenFile:%v ,fileNamePath;%s", err, hIP.FileNamePath))
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}

	if f != nil {
		accountSheet1Name, _ := beego.AppConfig.GetSection("handbook_account_excel_sheet1_name")
		if err != nil {
			utils.LogDebug(fmt.Sprintf("GetSection:%v", err))
			c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
		}

		accountSheet1Title, _ := beego.AppConfig.GetSection("handbook_account_excel_sheet1_title")
		if err != nil {
			utils.LogDebug(fmt.Sprintf("GetSection:%v", err))
			c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
		}

		handBookTypeString := models.GetHandBookTypeWithString("账册")
		if len(handBookTypeString) == 0 {
			c.jsonResult(enums.JRCodeFailed, "账册类型获取失败", nil)
		}

		handBookType, err := strconv.Atoi(handBookTypeString)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("ParseInt:%v", err))
			c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
		}

		hIP.HandBook.Type = int8(handBookType)
		t := reflect.ValueOf(&hIP.HandBook).Elem()
		for i := 0; i < reflect.ValueOf(hIP.HandBook).NumField(); i++ {
			obj := reflect.TypeOf(hIP.HandBook).Field(i)
			for iw, v := range accountSheet1Title {
				// Get value from cell by given worksheet name and axis.
				if iw == strings.ToLower(obj.Name) {
					cell, err := f.GetCellValue(accountSheet1Name["name"], v)
					if err != nil {
						utils.LogDebug(fmt.Sprintf("GetCellValue:%v", err))
						c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
					}
					xlsx.SetObjValue(obj.Name, cell, t)
				}
			}
		}

		CompanyManageCode := hIP.HandBook.CompanyManageCode // 经营单位代码
		company, err := models.CompanyByManageCode(CompanyManageCode)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("CompanyByManageCode:%v", err))
			c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
		}

		hIP.HandBook.Company = company

	} else {
		utils.LogDebug(fmt.Sprintf("导入失败:%v", err))
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}

}

//导入基础参数 xlsx 文件内容
func (c *HandBookController) ImportHandBookXlsxByRow(hIP *models.HandBookImportParam) {

	f, err := excelize.OpenFile(hIP.FileNamePath)

	if err != nil {
		utils.LogDebug(fmt.Sprintf("OpenFile:%v ,fileNamePath;%s", err, hIP.FileNamePath))
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}

	if f != nil {
		accountSheet1Name, err := xlsx.GetExcelName("handbook_account_excel_sheet1_name")
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
		}

		accountSheet1Title, err := xlsx.GetExcelTitles("", "handbook_account_excel_sheet1_title")
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
		}

		handBookTypeString := models.GetHandBookTypeWithString("账册")
		if len(handBookTypeString) == 0 {
			c.jsonResult(enums.JRCodeFailed, "账册类型获取失败", nil)
		}

		handBookType, err := strconv.Atoi(handBookTypeString)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("ParseInt:%v", err))
			c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
		}

		hIP.HandBook.Type = int8(handBookType)
		t := reflect.ValueOf(&hIP.HandBook).Elem()
		for i := 0; i < reflect.ValueOf(hIP.HandBook).NumField(); i++ {
			obj := reflect.TypeOf(hIP.HandBook).Field(i)
			for iw, v := range accountSheet1Title {
				// Get value from cell by given worksheet name and axis.
				if iw == strings.ToLower(obj.Name) {
					cell, err := f.GetCellValue(accountSheet1Name, v)
					if err != nil {
						utils.LogDebug(fmt.Sprintf("GetCellValue:%v", err))
						c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
					}
					xlsx.SetObjValue(obj.Name, cell, t)
				}
			}
		}

		CompanyManageCode := hIP.HandBook.CompanyManageCode // 经营单位代码
		company, err := models.CompanyByManageCode(CompanyManageCode)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("CompanyByManageCode:%v", err))
			c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
		}

		hIP.HandBook.Company = company

	} else {
		utils.LogDebug(fmt.Sprintf("导入失败:%v", err))
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}

}
