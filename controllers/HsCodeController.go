package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"

	"BeeCustom/models"
)

type HsCodeController struct {
	BaseController
}

func (c *HsCodeController) Prepare() {
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

func (c *HsCodeController) Index() {

	//页面模板设置
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "hscode/index_footerjs.html"

	//页面里按钮权限控制
	c.getActionData("", "Import")

	c.GetXSRFToken()
}

//列表数据
func (c *HsCodeController) DataGrid() {
	//直接获取参数 getDataGridData()
	params := models.NewHsCodeQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	//获取数据列表和总数
	data, total := models.HsCodePageList(&params)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	result["code"] = 0
	c.Data["json"] = result

	c.ServeJSON()
}

//导入
func (c *HsCodeController) Import() {
	f, err := excelize.OpenFile("./Book1.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	// Get value from cell by given worksheet name and axis.
	cell, err := f.GetCellValue("Sheet1", "B2")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cell)
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet1")
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
}
