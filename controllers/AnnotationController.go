package controllers

import (
	"BeeCustom/enums"
	"BeeCustom/models"
	"encoding/json"
	"fmt"
)

type AnnotationController struct {
	BaseController
}

func (c *AnnotationController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	//默认认证 "Index", "Create", "Edit", "Delete"
	c.checkAuthor()

	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()

}
func (c *AnnotationController) Index() {
	//页面模板设置
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "annotation/index_footerjs.html"

	//页面里按钮权限控制
	c.getActionData("Edit", "Delete", "Create", "Freeze")

	c.GetXSRFToken()
}

//列表数据
func (c *AnnotationController) DataGrid() {
	//直接获取参数 getDataGridData()
	params := models.NewAnnotationQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	//获取数据列表和总数
	data, total := models.AnnotationPageList(&params)
	ms, err := models.AnnotationGetRelations(data, "Role")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "关联关系获取失败", nil)
	}
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = ms
	result["code"] = 0
	c.Data["json"] = result

	c.ServeJSON()
}

// Create 添加 新建 页面
func (c *AnnotationController) Create() {
	params := models.NewRoleQueryParam()
	roles := models.RoleDataList(&params)

	c.Data["roles"] = roles
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "annotation/create_footerjs.html"
	c.GetXSRFToken()
}

// Store 添加 新建 页面
func (c *AnnotationController) Store() {
	m := models.NewAnnotation(0)
	//获取form里的值
	if err := c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}

	c.validRequestData(m)

	//valid := validation.Validation{}
	//valid.Required(m.UserPwd, "密码")
	//valid.MinSize(m.UserPwd, 6, "密码")
	//valid.MaxSize(m.UserPwd, 18, "密码")
	//
	//if valid.HasErrors() {
	//	// 如果有错误信息，证明验证没通过
	//	// 打印错误信息
	//	for _, err := range valid.Errors {
	//		c.jsonResult(enums.JRCodeFailed, err.Key+err.Message, m)
	//	}
	//}

	if _, err := models.AnnotationSave(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "添加失败", m)
	} else {
		c.jsonResult(enums.JRCodeSucc, "添加成功", m)
	}
}

// Edit 添加 编辑 页面
func (c *AnnotationController) Edit() {
	Id, _ := c.GetInt64(":id", 0)
	m, err := models.AnnotationOne(Id)
	if m != nil && Id > 0 {
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}

	c.Data["m"] = m
	params := models.NewRoleQueryParam()
	roles := models.RoleDataList(&params)
	c.Data["roles"] = roles

	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "annotation/edit_footerjs.html"
	c.GetXSRFToken()
}

// Update 添加 编辑 页面
func (c *AnnotationController) Update() {
	Id, _ := c.GetInt64(":id", 0)
	m := models.NewAnnotation(Id)

	//获取form里的值
	if err := c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}

	c.validRequestData(m)

	//valid := validation.Validation{}
	//if len(m.UserPwd) > 0 {
	//	valid.MinSize(m.UserPwd, 6, "密码")
	//	valid.MaxSize(m.UserPwd, 18, "密码")
	//}
	//
	//if valid.HasErrors() {
	//	// 如果有错误信息，证明验证没通过
	//	// 打印错误信息
	//	for _, err := range valid.Errors {
	//		c.jsonResult(enums.JRCodeFailed, err.Key+err.Message, m)
	//	}
	//}

	if _, err := models.AnnotationSave(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "编辑失败", m)
	} else {
		c.jsonResult(enums.JRCodeSucc, "编辑成功", m)
	}
}

//删除
func (c *AnnotationController) Delete() {
	id, _ := c.GetInt64(":id")
	if num, err := models.AnnotationDelete(id); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), "")
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", err)
	}
}

// Edit 添加 编辑 页面
func (c *AnnotationController) Profile() {

	m, err := models.AnnotationOne(c.curUser.Id)
	if m != nil && c.curUser.Id > 0 {
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}

	c.Data["m"] = m
	params := models.NewRoleQueryParam()
	roles := models.RoleDataList(&params)
	c.Data["roles"] = roles

	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "annotation/profile_footerjs.html"
	c.GetXSRFToken()
}
