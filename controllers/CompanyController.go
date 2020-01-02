package controllers

import (
	"encoding/json"
	"fmt"

	"BeeCustom/transforms"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/snowlyg/gotransform"

	"BeeCustom/enums"
	"BeeCustom/models"
	"BeeCustom/utils"
)

type CompanyController struct {
	BaseController
}

func (c *CompanyController) Prepare() {
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

func (c *CompanyController) Index() {

	//页面模板设置
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "company/index_footerjs.html"

	//页面里按钮权限控制
	c.getActionData("", "Edit", "Delete", "Create")

	c.GetXSRFToken()
}

//列表数据
func (c *CompanyController) DataGrid() {
	//直接获取参数 getDataGridData()
	params := models.NewCompanyQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	//获取数据列表和总数
	data, total := models.CompanyPageList(&params)
	data, _ = models.CompaniesGetRelations(data, "CompanyContacts")
	c.ResponseList(c.transformCompanyList(data), total)
	c.ServeJSON()
}

// Create 添加 新建 页面
func (c *CompanyController) Create() {

	params := models.NewBackendUserQueryParam()
	backendUser := models.BackenduserDataList(&params)
	c.Data["backendUsers"] = backendUser

	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "company/create_footerjs.html"
	c.GetXSRFToken()
}

// Store 添加 新建 页面
func (c *CompanyController) Store() {
	m := models.NewCompany(0)
	//获取form里的值
	if err := c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}

	c.validRequestData(m)

	if _, err := models.CompanySave(&m); err != nil {
		utils.LogDebug(err)
		c.jsonResult(enums.JRCodeFailed, "添加失败", m)
	} else {
		c.jsonResult(enums.JRCodeSucc, "添加成功", m)
	}
}

// Edit 添加 编辑 页面
func (c *CompanyController) Edit() {
	Id, _ := c.GetInt64(":id", 0)
	m, err := models.CompanyOne(Id, "CompanySeals")

	if m != nil && Id > 0 {
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}

	c.Data["m"] = m

	params := models.NewBackendUserQueryParam()
	backendUser := models.BackenduserDataList(&params)
	c.Data["backendUsers"] = backendUser

	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.setTpl("company/create.html")
	c.LayoutSections["footerjs"] = "company/create_footerjs.html"
	c.GetXSRFToken()
}

// Update 添加 编辑 页面
func (c *CompanyController) Update() {
	Id, _ := c.GetInt64(":id", 0)
	m := models.NewCompany(Id)

	//获取form里的值
	if err := c.ParseForm(&m); err != nil {
		utils.LogDebug(fmt.Sprintf("获取数据失败:%v", err))
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}

	c.validRequestData(m)

	if _, err := models.CompanySave(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "编辑失败", m)
	} else {
		c.jsonResult(enums.JRCodeSucc, "编辑成功", m)
	}
}

//删除
func (c *CompanyController) Delete() {
	id, _ := c.GetInt64(":id")
	if num, err := models.CompanyDelete(id); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), "")
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", err)
	}
}

//导入
func (c *CompanyController) Import() {
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
func (c *CompanyController) transformCompanyList(ms []*models.Company) []*transforms.Company {
	var uts []*transforms.Company
	for _, v := range ms {
		ut := transforms.Company{}
		g := gotransform.NewTransform(&ut, v, enums.BaseDateTimeFormat)
		_ = g.Transformer()
		//ut.AdminName = c.getAdminName(v)
		uts = append(uts, &ut)
	}

	return uts
}

func (c *CompanyController) getAdminName(v *models.Company) string {
	for _, cc := range v.CompanyContacts {
		if cc.IsAdmin == 1 {
			return cc.Name
		}
	}
	return ""
}
