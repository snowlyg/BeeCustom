package controllers

import (
	"encoding/json"
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
	"github.com/snowlyg/gotransform"
)

type HandBookController struct {
	BaseController
}

func (c *HandBookController) Prepare() {
	//  先执行
	c.BaseController.Prepare()
	//  如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	perms := []string{
		"Index",
		"Create",
		"Edit",
		"Delete",
	}
	c.checkAuthor(perms)

	//  如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//  权限控制里会进行登录验证，因此这里不用再作登录验证
	//  c.checkLogin()

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
	// 页面模板设置
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "handbook/index_footerjs.html"
	c.Data["m"] = cs
	c.Data["count"] = count
	c.Data["searchWord"] = searchWord

	// 页面里按钮权限控制
	c.getActionData("", "Delete", "Import")

	c.GetXSRFToken()
}

// handbookgoods 列表数据
func (c *HandBookController) GoodDataGrid() {
	// 直接获取参数 GoodDataGrid()
	params := models.NewHandBookGoodQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	// 获取数据列表和总数
	data, total := models.HandBookGoodPageList(&params)
	c.ResponseList(data, total)
	c.ServeJSON()
}

//  根据 handbookid 获取 handbookgoods
func (c *HandBookController) GetHandBookGoodByHandBookId() {
	params := models.NewHandBookGoodQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	data, _ := models.GetHandBookGoodById(&params)

	handBookGoodsList := c.TransformHandBookGood(data)

	c.Data["json"] = handBookGoodsList
	c.ServeJSON()
}

// HandBook 列表数据
func (c *HandBookController) DataGrid() {
	params := models.NewHandBookQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	data, total := models.HandBookPageList(&params)
	c.ResponseList(data, total)
	c.ServeJSON()
}

//  Ullage 列表数据
func (c *HandBookController) UllageDataGrid() {
	params := models.NewHandBookUllageQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	data, total := models.HandBookUllagePageList(&params)
	c.ResponseList(data, total)
	c.ServeJSON()
}

//  Edit 添加 编辑 页面
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

	handBookTypeS, err := c.getHandBookTypes()
	chandBookType, err, done := enums.TransformCnToInt(handBookTypeS, "手册")
	if !done {
		c.jsonResult(enums.JRCodeFailed, fmt.Sprintf("手账册类型获取失败:%v", err), nil)
	}

	ahandBookType, err, done := enums.TransformCnToInt(handBookTypeS, "账册")
	if !done {
		c.jsonResult(enums.JRCodeFailed, fmt.Sprintf("手账册类型获取失败:%v", err), nil)
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

func (c *HandBookController) getHandBookTypes() (map[string]string, error) {
	return models.GetSettingRValueByKey("handBookType", false)
}

// 删除
func (c *HandBookController) Delete() {
	id, _ := c.GetInt64(":id")
	if num, err := models.HandBookDelete(id); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), "")
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", err)
	}
}

// 导入
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

	var sheetName, sheet1Title, sheet2Name, sheet2Title, sheet3Name, sheet3Title, sheet4Name, sheet4Title string
	handBookTypeS, err := c.getHandBookTypes()
	importManualType, err, _ := enums.TransformCnToInt(handBookTypeS, "手册")

	importAccountType, err, _ := enums.TransformCnToInt(handBookTypeS, "手册")

	if importType == importManualType {
		sheetName = "handbookManualExcelSheetName"
		sheet1Title = "handbookManualExcelSheet1Title"
		sheet2Name = "handbookManualExcelSheet2Name"
		sheet2Title = "handbookManualExcelSheet2Title"
		sheet3Name = "handbookManualExcelSheet3Name"
		sheet3Title = "handbookManualExcelSheet3Title"
		sheet4Name = "handbookManualExcelSheet4Name"
		sheet4Title = "handbookManualExcelSheet4Title"
	} else if importType == importAccountType {
		sheetName = "handbookAccountExcelSheetName"
		sheet1Title = "handbookAccountExcelSheet1Title"
		sheet2Name = "handbookAccountExcelSheet2Name"
		sheet2Title = "handbookAccountExcelSheet2Title"
		sheet3Name = "handbookAccountExcelSheet3Name"
		sheet3Title = "handbookAccountExcelSheet3Title"
		sheet4Name = "handbookAccountExcelSheet4Name"
		sheet4Title = "handbookAccountExcelSheet4Title"
	}

	handBookName, err := models.GetSettingRValueByKey(sheetName, false)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("handBook1Name:%v", err))
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}
	sheet2Name = handBookName["2"]
	sheet3Name = handBookName["3"]
	sheet4Name = handBookName["4"]

	handBook1Title, err := models.GetSettingRValueByKey(sheet1Title, false)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("handBook1Title:%v", err))
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}

	hIP := models.HandBookImportParam{
		BaseImportParam: xlsx.BaseImportParam{
			ExcelTitle:   handBook1Title,
			ExcelName:    handBookName["1"],
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

// 导入账册表体
func (c *HandBookController) InsertHandBookGoods(hIP *models.HandBookImportParam, hBGIP *models.HandBookGoodImportParam) {
	handBookSheetName, err := models.GetSettingRValueByKey(hBGIP.ExcelName, true)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}

	handBookSheetTitle, err := models.GetSettingRValueByKey(hBGIP.ExcelTitleString, false)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}

	hIP.ExcelName = handBookSheetName["0"]
	hIP.ExcelTitle = handBookSheetTitle

	c.ImportHandBookXlsxByRow(hIP, hBGIP.HandBookTypeString)

}

// 获取 xlsx 文件内容
func (c *HandBookController) InsertHandBookGood(hIP *models.HandBookImportParam, Info []map[string]string) error {
	var handBookGoods []*models.HandBookGood

	handBookGood := models.NewHandBookGood(0)
	gf := gotransform.NewTransform(&Info, &handBookGood, enums.BaseDateTimeFormat)
	err := gf.Transformer()
	if err != nil {
		return err
	}
	//enums.SetObjValueFromSlice(&handBookGood, Info)

	handBookGood.HandBook = &hIP.HandBook
	handBookGood.Type = hIP.HandBookGoodType
	handBookGoods = append(handBookGoods, &handBookGood)

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

// 获取 xlsx 文件内容
func (c *HandBookController) InsertHandBookUllage(hIP *models.HandBookImportParam, Info []map[string]string) error {

	var handBookUllages []*models.HandBookUllage

	handBookUllage := models.NewHandBookUllage(0)
	gf := gotransform.NewTransform(&Info, &handBookUllage, enums.BaseDateTimeFormat)
	err := gf.Transformer()
	if err != nil {
		return err
	}
	//enums.SetObjValueFromSlice(&handBookUllage, Info)
	handBookGood, err := models.GetHandBookGoodBySerial(handBookUllage.FinishProNo)
	if err != nil {
		return err
	}

	handBookUllage.HandBookGood = handBookGood
	handBookUllages = append(handBookUllages, &handBookUllage)

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

// 导入基础参数 xlsx 文件内容
func (c *HandBookController) ImportHandBookXlsxByCell(hIP *models.HandBookImportParam) {

	t := reflect.ValueOf(&hIP.HandBook).Elem()
	for i := 0; i < t.NumField(); i++ {
		tf := t.Field(i)
		hb := t.Type().Field(i)

		for iw, v := range hIP.ExcelTitle {
			//  Get value from cell by given worksheet name and axis.
			if iw == hb.Name {
				cell, err := xlsx.GetExcelCell(hIP.FileNamePath, hIP.ExcelName, v)
				if err != nil {
					utils.LogDebug(fmt.Sprintf("err:%v", err))
					c.jsonResult(enums.JRCodeFailed, "导入失败", err)
				}

				switch tf.Kind() {
				case reflect.String:
					tf.SetString(cell)
				case reflect.Float64:
					if len(cell) > 0 {
						objV, err := strconv.ParseFloat(cell, 64)
						if err != nil {
							utils.LogDebug(fmt.Sprintf("Parse:%v,%v,%v", err, cell, hb.Name))
						}
						tf.SetFloat(objV)
					}
				case reflect.Int8:
					if len(cell) > 0 {
						objV, err := strconv.Atoi(cell)
						if err != nil {
							utils.LogDebug(fmt.Sprintf("Parse:%v,%v,%v", err, cell, hb.Name))
						}
						tf.SetInt(int64(objV))
					}
				case reflect.Uint64:
					reflect.ValueOf(cell)
					objV, err := strconv.ParseUint(v, 0, 64)
					if err != nil {
						utils.LogDebug(fmt.Sprintf("Parse:%v,%v,%v", err, cell, hb.Name))
					}
					tf.SetUint(objV)
				case reflect.Struct:
					if len(cell) > 0 {
						objV, err := time.Parse("20060102", cell)
						if err != nil {
							utils.LogDebug(fmt.Sprintf("Parse:%v,%v,%v", err, cell, hb.Name))
						}
						tf.Set(reflect.ValueOf(objV))
					}

				default:
					utils.LogDebug(fmt.Sprintf("未知类型:%v,%v", cell, hb.Name))
				}
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

	CompanyManageCode := hIP.HandBook.CompanyManageCode //  经营单位代码
	company, err := models.CompanyByManageCode(CompanyManageCode)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}

	hIP.HandBook.Company = company

}

// 导入基础参数 xlsx 文件内容
func (c *HandBookController) ImportHandBookXlsxByRow(hIP *models.HandBookImportParam, handBookTypeString string) {
	rows, err := xlsx.GetExcelRows(hIP.FileNamePath, hIP.ExcelName)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}

	var Info []map[string]string
	if len(handBookTypeString) > 0 { // 表体
		obj := models.NewHandBookGood(0)
		for roI, row := range rows {
			if roI > 1 { // 忽略标题和表头 2 行
				// 将数组  转成对应的 map
				var info = make(map[string]string)
				//  模型前两个字段是 BaseModel ，Type 不需要赋值
				for i := 0; i < reflect.ValueOf(obj).NumField(); i++ {
					obj := reflect.TypeOf(obj).Field(i)
					for _, iw := range hIP.ExcelTitle {
						if iw == obj.Name {
							rI := xlsx.ObjIsExists(hIP.ExcelTitle, iw)
							//  模板字段数量定义
							if rI != -1 && rI <= len(row)-1 {
								info[obj.Name] = row[rI]
							}
						}
					}
				}
				Info = append(Info, info)
			}

		}

		handBookGoodTypeS, err := c.getHandBookTypes()
		handBookGoodType, err, _ := enums.TransformCnToInt(handBookGoodTypeS, handBookTypeString)
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, fmt.Sprintf("账册类型获取失败:%v", err), nil)
		}

		hIP.HandBookGoodType = handBookGoodType

		err = c.InsertHandBookGood(hIP, Info)
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
		}

	} else { // 单损

		obj := models.NewHandBookUllage(0)
		for roI, row := range rows {
			if roI > 1 { // 忽略标题行
				// 将数组  转成对应的 map
				var info = make(map[string]string)
				//  模型前两个字段是 BaseModel ，Type 不需要赋值
				for i := 0; i < reflect.ValueOf(obj).NumField(); i++ {
					obj := reflect.TypeOf(obj).Field(i)
					for _, iw := range hIP.ExcelTitle {
						if iw == obj.Name {
							rI := xlsx.ObjIsExists(hIP.ExcelTitle, iw)
							//  模板字段数量定义
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

//  TransformHandBookGoodsList 格式化列表数据
func (c *HandBookController) TransformHandBookGoodsList(ms []*models.HandBookGood) []*map[string]interface{} {
	var handBookList []*map[string]interface{}
	clearances1 := models.GetClearancesByTypes("货币代码", true)
	clearances2 := models.GetClearancesByTypes("计量单位代码", false)
	for _, v := range ms {
		var unitOneCode interface{}
		var unitTwoCode interface{}
		var moneyunitCode interface{}
		for _, c := range clearances2 {
			if c[0] == v.UnitOne {
				unitOneCode = c[1]
			}

			if c[0] == v.UnitTwo {
				unitTwoCode = c[1]
			}
		}
		for _, c := range clearances1 {
			if c[0] == v.Moneyunit {
				moneyunitCode = c[1]
			}

		}
		handBook := make(map[string]interface{})
		handBook["Id"] = strconv.FormatInt(v.Id, 10)
		handBook["RecordNo"] = v.RecordNo
		handBook["HsCode"] = v.HsCode
		handBook["Name"] = v.Name
		handBook["Special"] = v.Special
		handBook["UnitOne"] = v.UnitOne
		handBook["UnitOneCode"] = unitOneCode
		handBook["UnitTwo"] = v.UnitTwo
		handBook["UnitTwoCode"] = unitTwoCode
		handBook["Price"] = v.Price
		handBook["Moneyunit"] = v.Moneyunit
		handBook["MoneyunitCode"] = moneyunitCode

		handBookList = append(handBookList, &handBook)
	}

	return handBookList
}
