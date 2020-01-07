package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"BeeCustom/enums"
	"BeeCustom/transforms"
	"BeeCustom/xlsx"
	gtf "github.com/snowlyg/gotransformer"

	"BeeCustom/models"
)

type HsCodeController struct {
	BaseController
}

func (c *HsCodeController) Prepare() {
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

}

func (c *HsCodeController) Index() {

	// 页面模板设置
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "hscode/index_footerjs.html"

	// 页面里按钮权限控制
	c.getActionData("", "Import")

	c.GetXSRFToken()
}

// 列表数据
func (c *HsCodeController) DataGrid() {
	// 直接获取参数 getDataGridData()
	params := models.NewHsCodeQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	// 获取数据列表和总数
	data, total := models.HsCodePageList(&params)
	c.ResponseList(c.transformHsCodeList(data), total)
	c.ServeJSON()
}

// 列表数据
func (c *HsCodeController) Get() {
	hsCodeS := c.GetString(":hs_code")
	hsCode, _ := models.GetHsCodeByCode(hsCodeS)
	c.Data["json"] = hsCode
	c.ServeJSON()
}

// 导入
func (c *HsCodeController) Import() {
	fileType := "hs_code/" + strconv.FormatInt(c.curUser.Id, 10) + "/"
	fileNamePath, err := c.BaseUpload(fileType)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "上传失败", err)
	}

	_, err = models.HsCodeDeleteAll() // 清空对应基础参数
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "清空数据报错", err)
	}

	cIP := xlsx.BaseImport{
		ExcelName:    "Sheet1",
		FileNamePath: fileNamePath,
	}

	titles, err := models.GetSettingRValueByKey("hsCodeExcelTile", false)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "上传失败", err)
	}

	rows, err := xlsx.GetExcelRows(cIP.FileNamePath, cIP.ExcelName)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "上传失败", err)
	}

	// 提取 excel 数据
	hsCodes := make([]*models.HsCode, 0)
	//cs := models.GetClearancesByTypes("计量单位代码", false)
	for roI, row := range rows {
		if roI > 0 {
			// 将数组  转成对应的 map
			c := models.NewHsCode(0)
			x := gtf.NewXlxsTransform(&c, titles, row, "", "", nil)
			err := x.XlxsTransformer()
			if err != nil {
				//c.jsonResult(enums.JRCodeFailed, "上传失败", err)
			}
			//for _, cc := range cs {
			//	if cc[1] == c.Unit1 {
			//		c.Unit1Name = cc[0].(string)
			//	}
			//	if cc[1] == c.Unit2 {
			//		c.Unit2Name = cc[0].(string)
			//	}
			//}
			hsCodes = append(hsCodes, &c)
		}
	}

	mun, err := models.InsertHsCodeMulti(hsCodes)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "导入失败", nil)
	}

	c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("导入成功 %d 项基础参数", mun), mun)
}

//  格式化列表数据
func (c *HsCodeController) transformHsCodeList(ms []*models.HsCode) []*transforms.HsCode {
	var uts []*transforms.HsCode
	for _, v := range ms {
		ut := transforms.HsCode{}
		g := gtf.NewTransform(&ut, v, enums.BaseDateTimeFormat)
		_ = g.Transformer()

		uts = append(uts, &ut)
	}

	return uts
}
