package controllers

import (
	"encoding/json"
	"fmt"

	"BeeCustom/enums"
	"BeeCustom/models"
	"BeeCustom/utils"
)

type CompanyForeignController struct {
	BaseController
}

func (c *CompanyForeignController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	//默认认证 "Index", "Create", "Edit", "Delete"
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

//列表数据
func (c *CompanyForeignController) DataGrid() {
	//直接获取参数 getDataGridData()
	params := models.NewCompanyForeignQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	//获取数据列表和总数
	data, total := models.CompanyForeignPageList(&params)
	c.ResponseList(data, total)
	c.ServeJSON()
}

// Create 添加 新建 页面
func (c *CompanyForeignController) Create() {
	Id, _ := c.GetInt64(":cid", 0)
	c.Data["companyId"] = Id

	c.setTpl("company/foreign/create.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "company/foreign/create_footerjs.html"
	c.GetXSRFToken()
}

// Store 添加 新建 页面
func (c *CompanyForeignController) Store() {
	m := models.NewCompanyForeign(0)
	//获取form里的值
	if err := c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}

	c.validRequestData(m)

	if _, err := models.CompanyForeignSave(&m); err != nil {
		utils.LogDebug(err)
		c.jsonResult(enums.JRCodeFailed, "添加失败", m)
	} else {
		c.jsonResult(enums.JRCodeSucc, "添加成功", m)
	}
}

// Edit 添加 编辑 页面
func (c *CompanyForeignController) Edit() {
	Id, _ := c.GetInt64(":id", 0)
	m, err := models.CompanyForeignOne(Id)
	if m != nil && Id > 0 {
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}
	c.Data["m"] = m

	c.LayoutSections = make(map[string]string)
	c.setTpl("company/foreign/create.html")
	c.LayoutSections["footerjs"] = "company/foreign/edit_footerjs.html"
	c.GetXSRFToken()
}

// Update 添加 编辑 页面
func (c *CompanyForeignController) Update() {
	Id, _ := c.GetInt64(":id", 0)
	m := models.NewCompanyForeign(Id)

	//获取form里的值
	if err := c.ParseForm(&m); err != nil {
		utils.LogDebug(fmt.Sprintf("获取数据失败:%v", err))
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}

	c.validRequestData(m)

	if _, err := models.CompanyForeignSave(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "编辑失败", m)
	} else {
		c.jsonResult(enums.JRCodeSucc, "编辑成功", m)
	}
}

//删除
func (c *CompanyForeignController) Delete() {
	id, _ := c.GetInt64(":id")
	if num, err := models.CompanyForeignDelete(id); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), "")
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", err)
	}
}
