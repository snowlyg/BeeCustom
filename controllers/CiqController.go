package controllers

import (
	"BeeCustom/enums"
	"BeeCustom/transforms"
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/snowlyg/gotransform"

	"BeeCustom/models"
)

type CiqController struct {
	BaseController
}

func (c *CiqController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	//默认认证 "Index", "Create", "Edit", "Delete"
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

func (c *CiqController) Index() {

	//页面模板设置
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "ciq/index_footerjs.html"

	//页面里按钮权限控制
	c.getActionData("", "Import")
	c.GetXSRFToken()
}

//列表数据
func (c *CiqController) DataGrid() {
	//直接获取参数 getDataGridData()
	params := models.NewCiqQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	//获取数据列表和总数
	data, total := models.CiqPageList(&params)
	c.ResponseList(c.transformCiqList(data), total)
	c.ServeJSON()
}

//导入
func (c *CiqController) Import() {
	f, err := excelize.OpenFile("./Book1.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	// Get value from cell by given worksheet name and axis.
	cell := f.GetCellValue("Sheet1", "B2")

	fmt.Println(cell)
	// Get all the rows in the Sheet1.
	rows := f.GetRows("Sheet1")
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
}

//  格式化列表数据
func (c *CiqController) transformCiqList(ms []*models.Ciq) []*transforms.Ciq {
	var uts []*transforms.Ciq
	for _, v := range ms {
		ut := transforms.Ciq{}
		g := gotransform.NewTransform(&ut, v, enums.BaseDateTimeFormat)
		_ = g.Transformer()

		uts = append(uts, &ut)
	}

	return uts
}
