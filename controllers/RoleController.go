package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"BeeCustom/enums"
	"BeeCustom/models"
	"BeeCustom/transforms"
	"BeeCustom/utils"
	"github.com/snowlyg/gotransform"
)

type PermList struct {
	Title   string      `json:"title"`
	Value   string      `json:"value"`
	Checked bool        `json:"checked"`
	Data    []*PermList `json:"data"`
}

type PermTreeList struct {
	Data []*PermList `json:"data"`
}

//RoleController 角色管理
type RoleController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *RoleController) Prepare() {
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
	//c.checkLogin()//权限控制里会进行登录验证，因此这里不用再作登录验证
}

//Index 角色管理首页
func (c *RoleController) Index() {
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "role/index_footerjs.html"

	//页面里按钮权限控制
	c.getActionData("", "Edit", "Delete", "Create")
	c.GetXSRFToken()
}

// Create 添加 新建 页面
func (c *RoleController) Create() {
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "role/create_footerjs.html"
	c.GetXSRFToken()
}

// Store 添加 新建 页面
func (c *RoleController) Store() {
	m := models.NewRole(0)

	//获取form里的值
	if err := c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}

	c.validRequestData(m)

	err := models.RoleSave(&m)
	if err == nil {
		c.jsonResult(enums.JRCodeSucc, "添加成功", m)
	} else {
		c.jsonResult(enums.JRCodeFailed, "添加失败", m)
	}
}

// DataGrid 角色管理首页 表格获取数据
func (c *RoleController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值
	params := models.NewRoleQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	//获取数据列表和总数
	data, total := models.RolePageList(&params)

	c.ResponseList(c.transformRoleList(data), total)
	c.ServeJSON()
}

//DataList 角色列表
func (c *RoleController) DataList() {
	params := models.NewRoleQueryParam()
	//获取数据列表和总数
	data := models.RoleDataList(&params)
	//定义返回的数据结构
	c.jsonResult(enums.JRCodeSucc, "", data)
}

//PermLists 权限列表
func (c *RoleController) PermLists() {
	var ptl PermTreeList
	var m *models.Role
	var err error

	Id, _ := c.GetInt64(":id", 0)
	if Id > 0 {
		m, err = models.RoleOne(Id, true)
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}

	//直接反序化获取json格式的requestbody里的值
	params := models.NewResourceQueryParam()
	params.IsParent = true
	//获取数据列表和总数
	datas := models.ResourceDataList(&params)

	for _, v := range datas {
		getSonsPerms(&ptl, v, m) //生成子权限树结构
	}

	//定义返回的数据结构
	c.jsonResult(enums.JRCodeSucc, "", ptl)
}

//生成子权限树结构
func getSonsPerms(ptl *PermTreeList, v *models.Resource, m *models.Role) {
	pl := PermList{}
	pl.Title = v.Name
	pl.Value = v.UrlFor
	pl.Checked = getChecked(v, m) //是否有权限

	if v.Sons != nil {
		for _, sv := range v.Sons {
			pls := PermList{}
			pls.Title = sv.Name
			pls.Value = sv.UrlFor
			pls.Checked = getChecked(sv, m) //是否有权限
			pl.Data = append(pl.Data, &pls)
		}
	}

	ptl.Data = append(ptl.Data, &pl)
}

//是否有权限
func getChecked(v *models.Resource, m *models.Role) bool {
	if m == nil {
		return false
	}

	hasPermission := utils.E.HasPermissionForUser(strconv.FormatInt(m.Id, 10), v.UrlFor)

	return hasPermission

}

//Edit 添加、编辑角色界面
func (c *RoleController) Edit() {
	Id, _ := c.GetInt64(":id", 0)
	if Id > 0 {

		m, err := models.RoleOne(Id, false)
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}

		c.Data["m"] = m
	}

	c.setTpl("role/edit.html", "shared/layout_app.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "role/edit_footerjs.html"
	c.GetXSRFToken()
}

//Update 添加、编辑角色界面
func (c *RoleController) Update() {

	id, _ := c.GetInt64(":id", 0)
	//获取form里的值
	if id == 0 {
		c.jsonResult(enums.JRCodeFailed, "参数错误", id)
	}

	m := models.NewRole(id)

	//获取form里的值
	if err := c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}

	c.validRequestData(m)

	_, err := models.RoleUpdate(&m)
	if err == nil {
		c.jsonResult(enums.JRCodeSucc, "编辑成功", m)
	} else {
		c.jsonResult(enums.JRCodeFailed, "编辑失败", m)
	}
}

//Delete 批量删除
func (c *RoleController) Delete() {
	id, _ := c.GetInt64(":id")
	if num, err := models.RoleDelete(id); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), "")
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", id)
	}
}

//  格式化列表数据
func (c *RoleController) transformRoleList(ms []*models.Role) []*transforms.Role {
	var uts []*transforms.Role
	for _, v := range ms {
		ut := transforms.Role{}
		g := gotransform.NewTransform(&ut, v, enums.BaseDateTimeFormat)
		_ = g.Transformer()

		uts = append(uts, &ut)
	}

	return uts
}
