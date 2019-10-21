package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"BeeCustom/enums"
	"BeeCustom/models"
	"github.com/astaxie/beego/orm"
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
	//"DataGrid", "DataList", "UpdateSeq" 不用检查权限
	c.checkAuthor("DataGrid", "DataList", "UpdateSeq")

	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//c.checkLogin()//权限控制里会进行登录验证，因此这里不用再作登录验证
}

//Index 角色管理首页
func (c *RoleController) Index() {

	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "role/index_footerjs.html"

	//页面里按钮权限控制
	c.Data["canEdit"] = c.checkActionAuthor("RoleController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("RoleController", "Delete")
	c.Data["canAllocate"] = c.checkActionAuthor("RoleController", "Allocate")
}

// Create 添加 新建 页面
func (c *RoleController) Create() {
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "role/create_footerjs.html"
}

// Store 添加 新建 页面
func (c *RoleController) Store() {
	m := &models.Role{BaseModel: models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()}}
	//获取form里的值
	if err := c.ParseForm(m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
	}

	_, err := models.RoleSave(m)
	if err == nil {
		c.jsonResult(enums.JRCodeSucc, "添加成功", m.Id)
	} else {
		c.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
	}
}

// DataGrid 角色管理首页 表格获取数据
func (c *RoleController) DataGrid() {

	//直接反序化获取json格式的requestbody里的值
	var params models.RoleQueryParam
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	//获取数据列表和总数
	data, total := models.RolePageList(&params)

	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	result["code"] = 0
	c.Data["json"] = result
	c.ServeJSON()
}

//DataList 角色列表
func (c *RoleController) DataList() {
	var params = models.RoleQueryParam{}
	//获取数据列表和总数
	data := models.RoleDataList(&params)
	//定义返回的数据结构
	c.jsonResult(enums.JRCodeSucc, "", data)
}

//PermLists 权限列表
func (c *RoleController) PermLists() {
	ptl := PermTreeList{}
	//直接反序化获取json格式的requestbody里的值
	var params = models.ResourceQueryParam{}
	//获取数据列表和总数
	datas := models.ResourceDataList(&params)
	for _, v := range datas {

		pl := &PermList{v.Name, v.UrlFor, false, nil}
		if v.Parent == nil {
			if v.Sons != nil {
				for _, sv := range v.Sons {
					pls := &PermList{sv.Name, sv.UrlFor, false, nil}

					pl.Data = append(pl.Data, pls)
				}
			}

			ptl.Data = append(ptl.Data, pl)
		}

	}

	//定义返回的数据结构
	c.jsonResult(enums.JRCodeSucc, "", ptl)
}

//Edit 添加、编辑角色界面
func (c *RoleController) Edit() {

	Id, _ := c.GetInt64(":id", 0)
	//m := models.Role{BaseModel: models.BaseModel{Id: Id}}
	if Id > 0 {

		m, err := models.RoleOne(Id)
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}

		c.Data["m"] = m
	}

	c.setTpl("role/edit.html", "shared/layout_app.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "role/edit_footerjs.html"
}

//Update 添加、编辑角色界面
func (c *RoleController) Update() {
	id, _ := c.GetInt64(":id", 0)
	perm_ids := c.GetString(":perm_ids")

	m := &models.Role{BaseModel: models.BaseModel{Id: id, CreatedAt: time.Now(), UpdatedAt: time.Now()}}

	//获取form里的值
	if err := c.ParseForm(m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
	}

	_, err := models.RoleSave(m, perm_ids)
	if err == nil {
		c.jsonResult(enums.JRCodeSucc, "编辑成功", m.Id)
	} else {
		c.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
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

//Allocate 给角色分配资源界面
func (c *RoleController) Allocate() {

	roleId, _ := c.GetInt64("id", 0)
	//strs := c.GetString("ids")

	o := orm.NewOrm()
	m := models.Role{BaseModel: models.BaseModel{Id: roleId}}
	if err := o.Read(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "数据无效，请刷新后重试", m.Id)
	}

	//删除已关联的历史数据
	if _, err := o.QueryTable(models.RoleResourceRelTBName()).Filter("role__id", m.Id).Delete(); err != nil {
		c.jsonResult(enums.JRCodeFailed, "删除历史关系失败", m.Id)
	}

	//var relations []models.RoleResourceRel
	//for _, str := range strings.Split(strs, ",") {
	//	if id, err := strconv.Atoi(str); err == nil {
	//		r := models.Resource{Id: id}
	//		relation := models.RoleResourceRel{Role: &m, Resource: &r}
	//		relations = append(relations, relation)
	//	}
	//}

	//if len(relations) > 0 {
	//	//批量添加
	//	if _, err := o.InsertMulti(len(relations), relations); err == nil {
	//		c.jsonResult(enums.JRCodeSucc, "保存成功", "")
	//	}
	//}

	c.jsonResult(0, "保存失败", "")
}
