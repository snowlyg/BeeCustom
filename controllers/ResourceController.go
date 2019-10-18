package controllers

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"BeeCustom/enums"
	"BeeCustom/models"

	"github.com/astaxie/beego/orm"
)

type ResourceController struct {
	BaseController
}

func (c *ResourceController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的少数Action需要权限控制，则将验证放到需要控制的Action里
	//"TreeGrid", "UserMenuTree", "ParentTreeGrid", "Select" 不用检查权限
	c.checkAuthor("TreeGrid", "UserMenuTree", "ParentTreeGrid", "Select")
	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//这里注释了权限控制，因此这里需要登录验证
	c.checkLogin()
}

func (c *ResourceController) Index() {

	//将页面左边菜单的某项激活
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "resource/index_footerjs.html"
	//页面里按钮权限控制
	c.Data["canEdit"] = c.checkActionAuthor("ResourceController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("ResourceController", "Delete")

}

// Create 添加 新建 页面
func (c *ResourceController) Create() {

	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "resource/create_footerjs.html"
}

// Store 添加 新建 页面
func (c *ResourceController) Store() {

	c.Save(0)
}

//TreeGrid 获取所有资源的列表
func (c *ResourceController) TreeGrid() {
	//直接反序化获取json格式的requestbody里的值
	var params models.ResourceQueryParam
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	//获取数据列表和总数
	data, total := models.ResourceTreeGrid(&params)

	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	result["code"] = 0
	c.Data["json"] = result
	c.ServeJSON()
}

//UserMenuTree 获取用户有权管理的菜单、区域列表
func (c *ResourceController) UserMenuTree() {
	userid := c.curUser.Id
	//获取用户有权管理的菜单列表（包括区域）
	tree := models.ResourceTreeGridByUserId(userid, 1)
	//转换UrlFor 2 LinkUrl
	c.UrlFor2Link(tree)
	c.jsonResult(enums.JRCodeSucc, "", tree)
}

//ParentTreeGrid 获取可以成为某节点的父节点列表
//func (c *ResourceController) ParentTreeGrid() {
//	Id, _ := c.GetInt("id", 0)
//	tree := models.ResourceTreeGrid4Parent(Id)
//	//转换UrlFor 2 LinkUrl
//	c.UrlFor2Link(tree)
//	c.jsonResult(enums.JRCodeSucc, "", tree)
//}

// UrlFor2LinkOne 使用URLFor方法，将资源表里的UrlFor值转成LinkUrl
func (c *ResourceController) UrlFor2LinkOne(urlfor string) string {
	if len(urlfor) == 0 {
		return ""
	}
	// ResourceController.Edit,:id,1
	strs := strings.Split(urlfor, ",")
	if len(strs) == 1 {
		return c.URLFor(strs[0])
	} else if len(strs) > 1 {
		var values []interface{}
		for _, val := range strs[1:] {
			values = append(values, val)
		}
		return c.URLFor(strs[0], values...)
	}
	return ""
}

//UrlFor2Link 使用URLFor方法，批量将资源表里的UrlFor值转成LinkUrl
func (c *ResourceController) UrlFor2Link(src []*models.Resource) {
	for _, item := range src {
		item.LinkUrl = c.UrlFor2LinkOne(item.UrlFor)
	}
}

//Edit 资源编辑页面
func (c *ResourceController) Edit() {

	Id, _ := c.GetInt(":id", 0)
	m := models.Resource{BaseModel: models.BaseModel{Id: Id}, Seq: 100}
	if Id > 0 {
		o := orm.NewOrm()
		err := o.Read(&m)
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}

	c.Data["m"] = m

	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "resource/edit_footerjs.html"
}

//Update 添加、编辑角色界面
func (c *ResourceController) Update() {

	id, _ := c.GetInt(":id", 0)

	c.Save(id)
}

//Save 资源添加编辑 保存
func (c *ResourceController) Save(id int) {
	var err error
	m := models.Resource{BaseModel: models.BaseModel{id, time.Now(), time.Now()}}

	//获取form里的值
	if err = c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
	}

	o := orm.NewOrm()
	if m.Id == 0 {
		if _, err = o.Insert(&m); err == nil {
			c.jsonResult(enums.JRCodeSucc, "添加成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
		}

	} else {
		if _, err = o.Update(&m, "Name", "Parent", "Rtype", "Seq", "Seq", "Sons", "Sons", "Icon", "UrlFor", "Roles", "UpdatedAt"); err == nil {
			c.jsonResult(enums.JRCodeSucc, "编辑成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
		}
	}
}

// Delete 删除
func (c *ResourceController) Delete() {

	Id, _ := c.GetInt(":id", 0)
	if Id == 0 {
		c.jsonResult(enums.JRCodeFailed, "选择的数据无效", 0)
	}
	query := orm.NewOrm().QueryTable(models.ResourceTBName())
	if _, err := query.Filter("id", Id).Delete(); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("删除成功"), 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
	}
}

// Select 通用选择面板
func (c *ResourceController) Select() {
	////获取调用者的类别 1表示 角色
	//desttype, _ := c.GetInt("desttype", 0)
	////获取调用者的值
	//destval, _ := c.GetInt("destval", 0)
	////返回的资源列表
	//var selectedIds []string
	//o := orm.NewOrm()
	//if desttype > 0 && destval > 0 {
	//	//如果都大于0,则获取已选择的值，例如：角色，就是获取某个角色已关联的资源列表
	//	switch desttype {
	//	case 1:
	//		{
	//			role := models.Role{Id: destval}
	//			_, _ = o.LoadRelated(&role, "RoleResourceRel")
	//			for _, item := range role.RoleResourceRel {
	//				selectedIds = append(selectedIds, strconv.Itoa(item.Resource.Id))
	//			}
	//		}
	//	}
	//}
	//c.Data["selectedIds"] = strings.Join(selectedIds, ",")
	//c.setTpl("resource/select.html", "shared/layout_app.html")
	//c.LayoutSections = make(map[string]string)
	//c.LayoutSections["headcssjs"] = "resource/select_headcssjs.html"
	//c.LayoutSections["footerjs"] = "resource/select_footerjs.html"
}

//CheckUrlFor 填写UrlFor时进行验证
func (c *ResourceController) CheckUrlFor() {
	urlfor := c.GetString("urlfor")
	link := c.UrlFor2LinkOne(urlfor)
	if len(link) > 0 {
		c.jsonResult(enums.JRCodeSucc, "解析成功", link)
	} else {
		c.jsonResult(enums.JRCodeFailed, "解析失败", link)
	}
}
func (c *ResourceController) UpdateSeq() {

	Id, _ := c.GetInt("pk", 0)
	oM, err := models.ResourceOne(Id)
	if err != nil || oM == nil {
		c.jsonResult(enums.JRCodeFailed, "选择的数据无效", 0)
	}
	value, _ := c.GetInt("value", 0)
	oM.Seq = value
	if _, err := orm.NewOrm().Update(oM); err == nil {
		c.jsonResult(enums.JRCodeSucc, "修改成功", oM.Id)
	} else {
		c.jsonResult(enums.JRCodeFailed, "修改失败", oM.Id)
	}
}
