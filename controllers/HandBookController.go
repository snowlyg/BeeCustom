package controllers

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"BeeCustom/enums"
	"BeeCustom/models"
	"BeeCustom/utils"
	"BeeCustom/xlsx"
	"github.com/360EntSecGroup-Skylar/excelize"
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

	accountSheet1Name, _ := xlsx.GetExcelName("handbook_account_excel_sheet1_name")
	if err != nil {
		utils.LogDebug(fmt.Sprintf("GetSection:%v", err))
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}
	hIP.ExcelName = accountSheet1Name

	accountSheet1Title, _ := xlsx.GetExcelTitles("", "handbook_account_excel_sheet1_title")
	if err != nil {
		utils.LogDebug(fmt.Sprintf("GetSection:%v", err))
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}
	hIP.ExcelTitle = accountSheet1Title

	c.ImportHandBookXlsxByCell(&hIP)

	m, err := models.HandBookSave(&hIP.HandBook)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("InsertMulti:%v", err))
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}

	hIP.HandBook = *m

	hBGIP := models.HandBookGoodImportParam{
		"handbook_account_excel_sheet2_name",
		"handbook_account_excel_sheet2_title",
		"料件",
	}
	c.InsertHandBookGoods(&hIP, &hBGIP)

	hBGIP = models.HandBookGoodImportParam{
		"handbook_account_excel_sheet3_name",
		"handbook_account_excel_sheet3_title",
		"成品",
	}
	c.InsertHandBookGoods(&hIP, &hBGIP)

	hBGIP = models.HandBookGoodImportParam{
		"handbook_account_excel_sheet4_name",
		"handbook_account_excel_sheet4_title",
		"",
	}

	c.InsertHandBookGoods(&hIP, &hBGIP)

	c.SetLastUpdteTime("handBookLastUpdateTime", time.Now().Format(enums.BaseFormat))
	c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("导入成功"), m.Id)

}

//导入账册表体
func (c *HandBookController) InsertHandBookGoods(hIP *models.HandBookImportParam, hBGIP *models.HandBookGoodImportParam) {
	accountSheetName, err := xlsx.GetExcelName(hBGIP.ExcelNameString)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}
	hIP.ExcelName = accountSheetName

	accountSheetTitle, err := xlsx.GetExcelTitles("", hBGIP.ExcelTitleString)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}
	hIP.ExcelTitle2 = accountSheetTitle

	c.ImportHandBookXlsxByRow(hIP)

	if len(hBGIP.HandBookTypeString) > 0 { //表体
		handBookGoodType, err := models.GetHandBookTypeWithString(hBGIP.HandBookTypeString, "hand_book_good_type")
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, fmt.Sprintf("账册类型获取失败:%v", err), nil)
		}
		hIP.HandBookGoodType = handBookGoodType

		err = c.InsertHandBookGood(hIP)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("InsertHandBookGood:%v", err))
			c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
		}
	} else { //单损
		err = c.InsertHandBookUllage(hIP)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("InsertHandBookUllage:%v", err))
			c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
		}
	}

}

//获取 xlsx 文件内容
func (c *HandBookController) InsertHandBookGood(hIP *models.HandBookImportParam) error {
	//忽略标题行
	for i := 1; i < len(hIP.Info); i++ {
		t := reflect.ValueOf(&hIP.HandBookGood).Elem()
		for k, v := range hIP.Info[i] {
			xlsx.SetObjValue(k, v, t)
		}

		hIP.HandBookGood.HandBook = &hIP.HandBook
		hIP.HandBookGood.Type = hIP.HandBookGoodType
		hIP.HandBookGoods = append(hIP.HandBookGoods, &hIP.HandBookGood)
	}

	num, err := models.InsertHandBookGoodMulti(hIP.HandBookGoods)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("InsertHandBookGoodMulti:%v ", err))
		return err
	}

	if num == 0 {
		return errors.New("InsertHandBookGoodMulti:导入失败")
	}

	return nil
}

//获取 xlsx 文件内容
func (c *HandBookController) InsertHandBookUllage(hIP *models.HandBookImportParam) error {
	//忽略标题行
	for i := 1; i < len(hIP.Info); i++ {
		t := reflect.ValueOf(&hIP.HandBookUllage).Elem()
		for k, v := range hIP.Info[i] {
			xlsx.SetObjValue(k, v, t)
		}
		handBookGood, err := models.GetHandBookGoodBySerial(hIP.HandBookUllage.FinishProNo)
		if err != nil && err.Error() != "<QuerySeter> no row found" {
			return err
		}

		hIP.HandBookUllage.HandBookGood = handBookGood
		hIP.HandBookUllages = append(hIP.HandBookUllages, &hIP.HandBookUllage)
	}

	num, err := models.InsertHandBookUllageMulti(hIP.HandBookUllages)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("InsertHandBookGoodMulti:%v ", err))
		return err
	}

	if num == 0 {
		return errors.New("InsertHandBookGoodMulti:导入失败")
	}

	return nil
}

//导入基础参数 xlsx 文件内容
func (c *HandBookController) ImportHandBookXlsxByCell(hIP *models.HandBookImportParam) {

	f, err := excelize.OpenFile(hIP.FileNamePath)

	if err != nil {
		utils.LogDebug(fmt.Sprintf("OpenFile:%v ,fileNamePath;%s", err, hIP.FileNamePath))
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}

	if f != nil {
		handBookType, err := models.GetHandBookTypeWithString("账册", "hand_book_type")
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, fmt.Sprintf("账册类型获取失败:%v", err), nil)
		}

		hIP.HandBook.Type = handBookType
		t := reflect.ValueOf(&hIP.HandBook).Elem()
		for i := 0; i < reflect.ValueOf(hIP.HandBook).NumField(); i++ {
			obj := reflect.TypeOf(hIP.HandBook).Field(i)
			for iw, v := range hIP.ExcelTitle {
				// Get value from cell by given worksheet name and axis.
				if iw == strings.ToLower(obj.Name) {
					cell, err := f.GetCellValue(hIP.ExcelName, v)
					if err != nil {
						utils.LogDebug(fmt.Sprintf("GetCellValue:%v", err))
						c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
					}
					xlsx.SetObjValue(obj.Name, cell, t)
				}
			}
		}

		hB, err := models.GetHandBookByContractNumber(hIP.HandBook.ContractNumber)
		if err != nil && err.Error() != "<QuerySeter> no row found" {
			c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
		}

		if hB != nil && hB.Id != 0 {
			c.jsonResult(enums.JRCodeFailed, "账册已存在", nil)
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
		rows, err := f.GetRows(hIP.ExcelName)
		// Get all the rows in the Sheet1.
		if err != nil {
			utils.LogDebug(fmt.Sprintf("导入失败:%v", err))
			c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
		}

		for _, row := range rows {
			//将数组  转成对应的 map
			var info = make(map[string]string)
			// 模型前两个字段是 BaseModel ，Type 不需要赋值
			for i := 0; i < reflect.ValueOf(hIP.HandBookGood).NumField(); i++ {
				obj := reflect.TypeOf(hIP.HandBookGood).Field(i)
				for _, iw := range hIP.ExcelTitle2 {
					if iw == obj.Name {
						rI := xlsx.ObjIsExists(hIP.ExcelTitle2, iw)
						// 模板字段数量定义
						if rI != -1 && rI <= len(row) {
							info[obj.Name] = row[rI]
						}
					}
				}
			}

			hIP.Info = append(hIP.Info, info)
		}

	} else {
		utils.LogDebug(fmt.Sprintf("导入失败:%v", err))
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}

}
