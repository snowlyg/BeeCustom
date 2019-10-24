package controllers

import (
	"encoding/json"
	"fmt"

	"BeeCustom/enums"
	"BeeCustom/models"
)

type BackendUserController struct {
	BaseController
}

func (c *BackendUserController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	//默认认证 "Index", "Create", "Edit", "Delete"
	c.checkAuthor("Freeze")

	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()

}
func (c *BackendUserController) Index() {
	//是否显示更多查询条件的按钮弃用，前端自动判断
	//c.Data["showMoreQuery"] = true
	//将页面左边菜单的某项激活
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	//页面模板设置
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "backenduser/index_footerjs.html"

	//页面里按钮权限控制
	c.getActionData("Edit", "Delete", "Create", "Freeze")

	c.GetXSRFToken()
}

//列表数据
func (c *BackendUserController) DataGrid() {
	//直接获取参数 getDataGridData()
	params := models.NewBackendUserQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	//获取数据列表和总数
	data, total := models.BackendUserPageList(&params)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	result["code"] = 0
	c.Data["json"] = result

	c.ServeJSON()
}

// Create 添加 新建 页面
func (c *BackendUserController) Create() {
	params := models.NewRoleQueryParam()
	roles := models.RoleDataList(&params)

	c.Data["roles"] = roles
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "backenduser/create_footerjs.html"
	c.GetXSRFToken()
}

// Store 添加 新建 页面
func (c *BackendUserController) Store() {
	m := models.NewBackendUser(0)
	//获取form里的值
	if err := c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}

	c.validRequestData(m)

	if _, err := models.BackendUserSave(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "添加失败", m)
	} else {
		c.jsonResult(enums.JRCodeSucc, "添加成功", m)
	}
}

// Edit 添加 编辑 页面
func (c *BackendUserController) Edit() {
	Id, _ := c.GetInt64(":id", 0)
	m, err := models.BackendUserOne(Id)
	if m != nil && Id > 0 {
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
		//添加用户时默认状态为启用
		m.Status = enums.Enabled
	}

	c.Data["m"] = m
	params := models.NewRoleQueryParam()
	roles := models.RoleDataList(&params)
	c.Data["roles"] = roles

	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "backenduser/edit_footerjs.html"
	c.GetXSRFToken()
}

// 禁用 启用
func (c *BackendUserController) Freeze() {
	Id, _ := c.GetInt64(":id", 0)
	m, err := models.BackendUserOne(Id)
	if m != nil && Id > 0 {
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
		m.Status = !m.Status
	}

	if _, err := models.BackendUserFreeze(m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	} else {
		c.jsonResult(enums.JRCodeSucc, "操作成功", m)
	}
}

// Update 添加 编辑 页面
func (c *BackendUserController) Update() {
	Id, _ := c.GetInt64(":id", 0)
	m := models.NewBackendUser(Id)

	//获取form里的值
	if err := c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}

	c.validRequestData(m)

	if _, err := models.BackendUserSave(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "编辑失败", m)
	} else {
		c.jsonResult(enums.JRCodeSucc, "编辑成功", m)
	}
}

//删除
func (c *BackendUserController) Delete() {
	id, _ := c.GetInt64(":id")
	if num, err := models.BackendUserDelete(id); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), "")
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", err)
	}
}
