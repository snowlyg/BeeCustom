package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"BeeCustom/enums"
	"BeeCustom/models"

	"github.com/astaxie/beego/orm"
)

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
	c.Save(0)
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

//Edit 添加、编辑角色界面
func (c *RoleController) Edit() {
	Id, _ := c.GetInt(":id", 0)
	m := models.Role{BaseModel: models.BaseModel{Id: Id}, Seq: 100}
	if Id > 0 {
		o := orm.NewOrm()
		err := o.Read(&m)
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}
	c.Data["m"] = m
	c.setTpl("role/edit.html", "shared/layout_app.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "role/edit_footerjs.html"
}

//Update 添加、编辑角色界面
func (c *RoleController) Update() {
	id, _ := c.GetInt(":id", 0)

	c.Save(id)
}

//Save 添加、编辑页面 保存
func (c *RoleController) Save(id int) {
	var err error
	m := models.Role{BaseModel: models.BaseModel{id, time.Now(), time.Now()}}

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
		if _, err = o.Update(&m, "Name", "UpdatedAt"); err == nil {
			c.jsonResult(enums.JRCodeSucc, "编辑成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
		}
	}

}

//Delete 批量删除
func (c *RoleController) Delete() {

	id, _ := c.GetInt(":id")

	o := orm.NewOrm()
	if num, err := o.Delete(&models.Role{BaseModel: models.BaseModel{Id: id}}); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), "")
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", id)
	}
}

//Allocate 给角色分配资源界面
func (c *RoleController) Allocate() {

	roleId, _ := c.GetInt("id", 0)
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

func (c *RoleController) UpdateSeq() {

	Id, _ := c.GetInt("pk", 0)
	oM, err := models.RoleOne(Id)
	if err != nil || oM == nil {
		c.jsonResult(enums.JRCodeFailed, "选择的数据无效", Id)
	}
	value, _ := c.GetInt("value", 0)
	oM.Seq = value
	o := orm.NewOrm()
	if _, err := o.Update(oM); err == nil {
		c.jsonResult(enums.JRCodeSucc, "修改成功", oM.Id)
	} else {
		c.jsonResult(enums.JRCodeFailed, "修改失败", oM.Id)
	}

}
