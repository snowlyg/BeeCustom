package controllers

import (
	"BeeCustom/enums"
	"BeeCustom/models"
	"BeeCustom/utils"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/validation"
	"strconv"
	"strings"
)

type BackendUserController struct {
	BaseController
}

func (c *BackendUserController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare

	perms := []string{
		"Index",
		"Create",
		"Edit",
		"Delete",
		"Freeze",
	}
	c.checkAuthor(perms)

	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()

}
func (c *BackendUserController) Index() {
	//页面模板设置
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "backenduser/index_footerjs.html"

	//页面里按钮权限控制
	c.getActionData("", "Edit", "Delete", "Create", "Freeze")

	c.GetXSRFToken()
}

//列表数据
func (c *BackendUserController) DataGrid() {
	//直接获取参数 getDataGridData()
	params := models.NewBackendUserQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	//获取数据列表和总数
	data, total := models.BackendUserPageList(&params)
	ms, err := models.BackendUsersGetRelations(data)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "关联关系获取失败", nil)
	}

	c.ResponseList(ms, total)
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
	roleIdStrings := c.GetString("RoleIds")
	roleIds := strings.Split(roleIdStrings, ",")
	m := models.NewBackendUser(0)
	//获取form里的值
	if err := c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}

	c.validRequestData(m)

	valid := validation.Validation{}
	valid.Required(m.UserPwd, "密码")
	valid.MinSize(m.UserPwd, 6, "密码")
	valid.MaxSize(m.UserPwd, 18, "密码")

	if valid.HasErrors() {
		// 如果有错误信息，证明验证没通过
		// 打印错误信息
		for _, err := range valid.Errors {
			c.jsonResult(enums.JRCodeFailed, err.Key+err.Message, m)
		}
	}

	if _, err := models.BackendUserSave(&m, roleIds); err != nil {
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

	c.getRolesName(m)

	params := models.NewRoleQueryParam()
	roles := models.RoleDataList(&params)
	c.Data["roles"] = roles

	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "backenduser/edit_footerjs.html"
	c.GetXSRFToken()
}

func (c *BackendUserController) getRolesName(m *models.BackendUser) {
	roleIds, err := utils.E.GetRolesForUser(strconv.FormatInt(m.Id, 10))
	if err != nil {
		utils.LogDebug(fmt.Sprintf("GetRolesForUser error:%v", err))
	}
	var roleIds64 []interface{}
	for _, roleId := range roleIds {
		roleId64, _ := strconv.ParseInt(roleId, 10, 64)
		roleIds64 = append(roleIds64, roleId64)
	}
	m.RoleIds = roleIds64
	c.Data["m"] = m
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
	roleIdStrings := c.GetString("RoleIds")
	roleIds := strings.Split(roleIdStrings, ",")
	m := models.NewBackendUser(Id)

	//获取form里的值
	if err := c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}

	c.validRequestData(m)

	valid := validation.Validation{}
	if len(m.UserPwd) > 0 {
		valid.MinSize(m.UserPwd, 6, "密码")
		valid.MaxSize(m.UserPwd, 18, "密码")
	}

	if valid.HasErrors() {
		// 如果有错误信息，证明验证没通过
		// 打印错误信息
		for _, err := range valid.Errors {
			c.jsonResult(enums.JRCodeFailed, err.Key+err.Message, m)
		}
	}

	if _, err := models.BackendUserSave(&m, roleIds); err != nil {
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

// Edit 添加 编辑 页面
func (c *BackendUserController) Profile() {

	m, err := models.BackendUserOne(c.curUser.Id)
	if m != nil && c.curUser.Id > 0 {
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
		//添加用户时默认状态为启用
		m.Status = enums.Enabled
	}

	_ = models.BackendUserGetRelations(m)

	c.Data["m"] = m
	params := models.NewRoleQueryParam()
	roles := models.RoleDataList(&params)
	c.Data["roles"] = roles

	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "backenduser/profile_footerjs.html"
	c.GetXSRFToken()
}
