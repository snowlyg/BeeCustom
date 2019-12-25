package controllers

import (
	"encoding/json"
	"fmt"

	"BeeCustom/enums"
	"BeeCustom/models"
	"BeeCustom/transforms"
	"BeeCustom/utils"
	"github.com/snowlyg/gotransform"
)

type ResourceController struct {
	BaseController
}

func (c *ResourceController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的少数Action需要权限控制，则将验证放到需要控制的Action里

	perms := []string{
		"Index",
		"Create",
		"Edit",
		"Delete",
	}
	c.checkAuthor(perms)
	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//这里注释了权限控制，因此这里需要登录验证
	c.checkLogin()
}

func (c *ResourceController) Index() {
	//将页面左边菜单的某项激活

	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "resource/index_footerjs.html"
	//页面里按钮权限控制
	c.getActionData("", "Edit", "Delete", "Create")
	c.GetXSRFToken()
}

// Create 添加 新建 页面
func (c *ResourceController) Create() {
	//直接反序化获取json格式的requestbody里的值
	params := models.NewResourceQueryParam()
	params.IsParent = true

	//获取数据列表和总数
	data, _ := models.ResourceTreeGrid(&params)
	c.Data["parent_perms"] = data
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "resource/create_footerjs.html"
	c.GetXSRFToken()
}

// Store 添加 新建 页面
func (c *ResourceController) Store() {
	m := models.NewResource(0)

	//获取form里的值
	if err := c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}

	c.validRequestData(m)

	if _, err := models.ResourceSave(&m); err == nil {
		c.jsonResult(enums.JRCodeSucc, "添加成功", m)
	} else {
		c.jsonResult(enums.JRCodeFailed, "添加失败", m)
	}
}

//TreeGrid 获取所有资源的列表
func (c *ResourceController) TreeGrid() {
	//直接反序化获取json格式的requestbody里的值
	params := models.NewResourceQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	//获取数据列表和总数
	data, total := models.ResourceTreeGrid(&params)
	c.ResponseList(c.transformResourceList(data), total)
	c.ServeJSON()
}

//Edit 资源编辑页面
func (c *ResourceController) Edit() {
	//直接反序化获取json格式的requestbody里的值
	params := models.NewResourceQueryParam()
	params.IsParent = true

	//获取数据列表和总数
	data, _ := models.ResourceTreeGrid(&params)
	c.Data["parent_perms"] = data

	Id, _ := c.GetInt64(":id", 0)
	if Id > 0 {
		m, err := models.ResourceOne(Id)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("数据无效出错：%v", err))
			c.pageError("数据无效，请刷新后重试")
		}
		c.Data["m"] = m
	}

	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "resource/edit_footerjs.html"
	c.GetXSRFToken()
}

//Update 添加、编辑角色界面
func (c *ResourceController) Update() {
	id, _ := c.GetInt64(":id", 0)
	m := models.NewResource(id)

	//获取form里的值
	if err := c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}

	c.validRequestData(m)

	if _, err := models.ResourceSave(&m); err == nil {
		c.jsonResult(enums.JRCodeSucc, "编辑成功", m)
	} else {
		c.jsonResult(enums.JRCodeFailed, "编辑失败", m)
	}
}

// Delete 删除
func (c *ResourceController) Delete() {
	Id, _ := c.GetInt64(":id", 0)
	if Id == 0 {
		c.jsonResult(enums.JRCodeFailed, "选择的数据无效", 0)
	}

	if _, err := models.ResourceDelete(Id); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("删除成功"), 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
	}
}

// TransformAnnotationList 格式化列表数据
func (c *ResourceController) transformResourceList(ms []*models.Resource) []*transforms.Resource {
	var uts []*transforms.Resource
	for _, v := range ms {
		ut := transforms.Resource{}
		g := gotransform.NewTransform(&ut, v, enums.BaseDateTimeFormat)
		_ = g.Transformer()

		uts = append(uts, &ut)
	}

	return uts
}
