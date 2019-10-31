package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"BeeCustom/enums"
	"BeeCustom/models"
	"BeeCustom/utils"
)

type CompanyContactController struct {
	BaseController
}

func (c *CompanyContactController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	//默认认证 "Index", "Create", "Edit", "Delete"
	c.checkAuthor()

	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()

}

func (c *CompanyContactController) Index() {
	//是否显示更多查询条件的按钮弃用，前端自动判断
	//c.Data["showMoreQuery"] = true
	//将页面左边菜单的某项激活
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)

	//页面模板设置
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "companycontact/index_footerjs.html"

	//页面里按钮权限控制
	c.getActionData("Edit", "Delete", "Create")

	c.GetXSRFToken()
}

//列表数据
func (c *CompanyContactController) DataGrid() {
	//直接获取参数 getDataGridData()
	params := models.NewCompanyContactQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	//获取数据列表和总数
	data, total := models.CompanyContactPageList(&params)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	result["code"] = 0
	c.Data["json"] = result

	c.ServeJSON()
}

// Create 添加 新建 页面
func (c *CompanyContactController) Create() {
	Id, _ := c.GetInt64(":cid", 0)
	c.Data["companyId"] = Id

	c.setTpl("company/contact/create.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "company/contact/create_footerjs.html"
	c.GetXSRFToken()
}

// Store 添加 新建 页面
func (c *CompanyContactController) Store() {
	m := models.NewCompanyContact(0)
	//获取form里的值
	if err := c.ParseForm(&m); err != nil {
		utils.LogDebug(fmt.Sprintf("获取数据失败:%v", err))
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}

	c.checkAdminContactCount(m.CompanyId, m.IsAdmin)

	c.validRequestData(m)

	if _, err := models.CompanyContactSave(&m); err != nil {
		utils.LogDebug(fmt.Sprintf("添加失败:%v", err))
		c.jsonResult(enums.JRCodeFailed, "添加失败", m)
	} else {
		c.jsonResult(enums.JRCodeSucc, "添加成功", m)
	}
}

// Edit 添加 编辑 页面
func (c *CompanyContactController) Edit() {
	Id, _ := c.GetInt64(":id", 0)
	m, err := models.CompanyContactOne(Id)
	if m != nil && Id > 0 {
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}

	c.Data["m"] = m

	c.LayoutSections = make(map[string]string)
	c.setTpl("company/contact/create.html")
	c.LayoutSections["footerjs"] = "company/contact/edit_footerjs.html"
	c.GetXSRFToken()
}

// Update 添加 编辑 页面
func (c *CompanyContactController) Update() {
	Id, _ := c.GetInt64(":id", 0)
	m := models.NewCompanyContact(Id)

	//获取form里的值
	if err := c.ParseForm(&m); err != nil {
		utils.LogDebug(fmt.Sprintf("获取数据失败:%v", err))
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}

	c.checkAdminContactCount(m.CompanyId, m.IsAdmin)

	c.validRequestData(m)

	if _, err := models.CompanyContactSave(&m); err != nil {
		utils.LogDebug(err)
		c.jsonResult(enums.JRCodeFailed, "编辑失败", m)
	} else {
		c.jsonResult(enums.JRCodeSucc, "编辑成功", m)
	}
}

//删除
func (c *CompanyContactController) Delete() {
	id, _ := c.GetInt64(":id")
	if num, err := models.CompanyContactDelete(id); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), "")
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", err)
	}
}

//导入
func (c *CompanyContactController) Import() {
	id, _ := c.GetInt64(":id")
	if num, err := models.CompanyContactDelete(id); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), "")
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", err)
	}
}

func (c *CompanyContactController) checkAdminContactCount(companyId int64, isAdmin int8) {
	params := models.NewCompanyContactQueryParam()
	params.IsAdmin = true
	params.CompanyId = strconv.FormatInt(companyId, 10)

	_, count := models.CompanyContactPageList(&params)

	if count == 1 {
		if isAdmin == 1 {
			c.jsonResult(enums.JRCodeFailed, "只能有一个管理员", nil)
		}
	} else if count > 1 {
		c.jsonResult(enums.JRCodeFailed, "管理员已经超过限制，请修改联系人列表", nil)
	}

}
