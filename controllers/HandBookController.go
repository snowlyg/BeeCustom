package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"BeeCustom/enums"
	"BeeCustom/models"
	"BeeCustom/utils"
	"BeeCustom/xlsx"
)

type HandBookController struct {
	BaseController
}

func (c *HandBookController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	perms := []string{
		"Index",
		"Create",
		"Edit",
		"Delete",
	}
	c.checkAuthor(perms)

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
	c.getActionData("", "Delete", "Import")

	c.GetXSRFToken()
}

//列表数据
func (c *HandBookController) GoodDataGrid() {
	//直接获取参数 GoodDataGrid()
	params := models.NewHandBookGoodQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	//获取数据列表和总数
	data, total := models.HandBookGoodPageList(&params)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	result["code"] = 0
	c.Data["json"] = result

	c.ServeJSON()
}

//列表数据
func (c *HandBookController) GetHandBookGoodByHandBookId() {

	params := models.NewHandBookGoodQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	//获取数据列表和总数
	data, _ := models.HandBookGoodPageList(&params)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["rows"] = data
	result["status"] = 1
	c.Data["json"] = result

	c.ServeJSON()
}

//列表数据
func (c *HandBookController) DataGrid() {
	//直接获取参数 GoodDataGrid()
	params := models.NewHandBookQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	//获取数据列表和总数
	data, total := models.HandBookPageList(&params)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	result["code"] = 0
	c.Data["json"] = result

	c.ServeJSON()
}

//列表数据
func (c *HandBookController) UllageDataGrid() {
	//直接获取参数 getDataGridData()
	params := models.NewHandBookUllageQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	//获取数据列表和总数
	data, total := models.HandBookUllagePageList(&params)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	result["code"] = 0
	c.Data["json"] = result

	c.ServeJSON()
}

// Edit 添加 编辑 页面
func (c *HandBookController) Show() {
	Id, _ := c.GetInt64(":id", 0)
	m, err := models.HandBookOne(Id, "Company")
	if m != nil && Id > 0 {
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}

	c.Data["m"] = m

	var html, showFooterjs string
	chandBookType, err := enums.GetSectionWithString("手册", "hand_book_type")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, fmt.Sprintf("账册类型获取失败:%v", err), nil)
	}

	ahandBookType, err := enums.GetSectionWithString("账册", "hand_book_type")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, fmt.Sprintf("账册类型获取失败:%v", err), nil)
	}
	if m.Type == chandBookType {
		html = "handbook/manual/show.html"
		showFooterjs = "handbook/manual/show_footerjs.html"
	} else if m.Type == ahandBookType {
		html = "handbook/account/show.html"
		showFooterjs = "handbook/account/show_footerjs.html"
	}

	c.setTpl(html)
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = showFooterjs
	c.GetXSRFToken()
}

//删除
func (c *HandBookController) Delete() {
	id, _ := c.GetInt64(":id")
	if num, err := models.HandBookDelete(id); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), "")
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", err)
	}
}

//导入
func (c *HandBookController) Import() {
	importType, _ := c.GetInt8(":type")

	fileType := "handBook/" + strconv.FormatInt(c.curUser.Id, 10) + "/"
	fileNamePath, err := c.BaseUpload(fileType)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("BaseUpload:%v", err))
		c.jsonResult(enums.JRCodeFailed, "上传失败", nil)
	}

	fileNamePaths := strings.Split(fileNamePath, ".")
	fileExt := fileNamePaths[len(fileNamePaths)-1]
	if fileExt != "xlsx" {
		c.jsonResult(enums.JRCodeFailed, "文件格式错误，只能导入 xlsx 文件", nil)
	}

	var sheet1Name, sheet1Title, sheet2Name, sheet2Title, sheet3Name, sheet3Title, sheet4Name, sheet4Title string
	importManualType, _ := enums.GetSectionWithString("手册", "hand_book_type")
	importAccountType, _ := enums.GetSectionWithString("账册", "hand_book_type")
	if importType == importManualType {
		sheet1Name = "handbook_manual_excel_sheet1_name"
		sheet1Title = "handbook_manual_excel_sheet1_title"
		sheet2Name = "handbook_manual_excel_sheet2_name"
		sheet2Title = "handbook_manual_excel_sheet2_title"
		sheet3Name = "handbook_manual_excel_sheet3_name"
		sheet3Title = "handbook_manual_excel_sheet3_title"
		sheet4Name = "handbook_manual_excel_sheet4_name"
		sheet4Title = "handbook_manual_excel_sheet4_title"
	} else if importType == importAccountType {
		sheet1Name = "handbook_account_excel_sheet1_name"
		sheet1Title = "handbook_account_excel_sheet1_title"
		sheet2Name = "handbook_account_excel_sheet2_name"
		sheet2Title = "handbook_account_excel_sheet2_title"
		sheet3Name = "handbook_account_excel_sheet3_name"
		sheet3Title = "handbook_account_excel_sheet3_title"
		sheet4Name = "handbook_account_excel_sheet4_name"
		sheet4Title = "handbook_account_excel_sheet4_title"
	}

	accountSheet1Name, _ := xlsx.GetExcelName(sheet1Name)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("GetSection:%v", err))
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}

	accountSheet1Title, _ := xlsx.GetExcelTitles("", sheet1Title)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("GetSection:%v", err))
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}

	hIP := models.HandBookImportParam{
		BaseImportParam: xlsx.BaseImportParam{
			ExcelTitle:   accountSheet1Title,
			ExcelName:    accountSheet1Name,
			FileNamePath: fileNamePath,
		},
		HandBook: models.NewHandBook(0),
	}

	hIP.HandBook.Type = importType

	c.ImportHandBookXlsxByCell(&hIP)

	m, err := models.HandBookSave(&hIP.HandBook)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("InsertMulti:%v", err))
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}
	hIP.HandBook = *m

	hBGIP := models.HandBookGoodImportParam{
		sheet2Name,
		sheet2Title,
		"料件",
	}

	c.InsertHandBookGoods(&hIP, &hBGIP)

	hBGIP = models.HandBookGoodImportParam{
		sheet3Name,
		sheet3Title,
		"成品",
	}

	c.InsertHandBookGoods(&hIP, &hBGIP)

	hBGIP = models.HandBookGoodImportParam{
		sheet4Name,
		sheet4Title,
		"",
	}

	c.InsertHandBookGoods(&hIP, &hBGIP)

	c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("导入成功"), m.Id)

}

//导入账册表体
func (c *HandBookController) InsertHandBookGoods(hIP *models.HandBookImportParam, hBGIP *models.HandBookGoodImportParam) {
	accountSheetName, err := xlsx.GetExcelName(hBGIP.ExcelNameString)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}

	accountSheetTitle, err := xlsx.GetExcelTitles("", hBGIP.ExcelTitleString)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}

	hIP.ExcelName = accountSheetName
	hIP.ExcelTitle = accountSheetTitle

	c.ImportHandBookXlsxByRow(hIP, hBGIP.HandBookTypeString)

}

//获取 xlsx 文件内容
func (c *HandBookController) InsertHandBookGood(hIP *models.HandBookImportParam, Info []map[string]string) error {
	var handBookGoods []*models.HandBookGood
	for i := 0; i < len(Info); i++ {
		handBookGood := models.NewHandBookGood(0)
		t := reflect.ValueOf(&handBookGood).Elem()
		for k, v := range Info[i] {
			xlsx.SetObjValue(k, v, t)
		}

		handBookGood.HandBook = &hIP.HandBook
		handBookGood.Type = hIP.HandBookGoodType
		handBookGoods = append(handBookGoods, &handBookGood)
	}

	num, err := models.InsertHandBookGoodMulti(handBookGoods)
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
func (c *HandBookController) InsertHandBookUllage(hIP *models.HandBookImportParam, Info []map[string]string) error {

	var handBookUllages []*models.HandBookUllage
	for i := 0; i < len(Info); i++ {
		handBookUllage := models.NewHandBookUllage(0)
		t := reflect.ValueOf(&handBookUllage).Elem()
		for k, v := range Info[i] {
			xlsx.SetObjValue(k, v, t)
		}
		handBookGood, err := models.GetHandBookGoodBySerial(handBookUllage.FinishProNo)
		if err != nil {
			return err
		}

		handBookUllage.HandBookGood = handBookGood
		handBookUllages = append(handBookUllages, &handBookUllage)
	}

	num, err := models.InsertHandBookUllageMulti(handBookUllages)
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

	t := reflect.ValueOf(&hIP.HandBook).Elem()
	for i := 0; i < reflect.ValueOf(hIP.HandBook).NumField(); i++ {
		obj := reflect.TypeOf(hIP.HandBook).Field(i)
		for iw, v := range hIP.ExcelTitle {
			// Get value from cell by given worksheet name and axis.
			if iw == strings.ToLower(obj.Name) {
				cell, err := xlsx.GetExcelCell(hIP.FileNamePath, hIP.ExcelName, v)
				if err != nil {
					c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
				}

				xlsx.SetObjValue(obj.Name, cell, t)
			}
		}
	}

	hB, err := models.GetHandBookByContractNumber(hIP.HandBook.ContractNumber)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}

	if hB != nil && hB.Id != 0 {
		var errMsg string
		if hB.Type == 1 {
			errMsg = "手册已存在"
		} else if hB.Type == 2 {
			errMsg = "账册已存在"
		}
		c.jsonResult(enums.JRCodeFailed, errMsg, nil)
	}

	CompanyManageCode := hIP.HandBook.CompanyManageCode // 经营单位代码
	company, err := models.CompanyByManageCode(CompanyManageCode)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}

	hIP.HandBook.Company = company

}

//导入基础参数 xlsx 文件内容
func (c *HandBookController) ImportHandBookXlsxByRow(hIP *models.HandBookImportParam, handBookTypeString string) {
	rows, err := xlsx.GetExcelRows(hIP.FileNamePath, hIP.ExcelName)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}

	var Info []map[string]string
	if len(handBookTypeString) > 0 { //表体
		obj := models.NewHandBookGood(0)
		for roI, row := range rows {
			if roI > 1 { //忽略标题和表头 2 行
				//将数组  转成对应的 map
				var info = make(map[string]string)
				// 模型前两个字段是 BaseModel ，Type 不需要赋值
				for i := 0; i < reflect.ValueOf(obj).NumField(); i++ {
					obj := reflect.TypeOf(obj).Field(i)
					for _, iw := range hIP.ExcelTitle {
						if iw == obj.Name {
							rI := xlsx.ObjIsExists(hIP.ExcelTitle, iw)
							// 模板字段数量定义
							if rI != -1 && rI <= len(row)-1 {
								info[obj.Name] = row[rI]
							}
						}
					}
				}
				Info = append(Info, info)
			}

		}

		handBookGoodType, err := enums.GetSectionWithString(handBookTypeString, "hand_book_good_type")
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, fmt.Sprintf("账册类型获取失败:%v", err), nil)
		}

		hIP.HandBookGoodType = handBookGoodType

		err = c.InsertHandBookGood(hIP, Info)
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
		}

	} else { //单损

		obj := models.NewHandBookUllage(0)
		for roI, row := range rows {
			if roI > 1 { //忽略标题行
				//将数组  转成对应的 map
				var info = make(map[string]string)
				// 模型前两个字段是 BaseModel ，Type 不需要赋值
				for i := 0; i < reflect.ValueOf(obj).NumField(); i++ {
					obj := reflect.TypeOf(obj).Field(i)
					for _, iw := range hIP.ExcelTitle {
						if iw == obj.Name {
							rI := xlsx.ObjIsExists(hIP.ExcelTitle, iw)
							// 模板字段数量定义
							if rI != -1 && rI <= len(row) {
								info[obj.Name] = row[rI]
							}
						}
					}
				}

				Info = append(Info, info)
			}

		}

		err = c.InsertHandBookUllage(hIP, Info)
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
		}

	}
}
