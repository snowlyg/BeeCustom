package controllers

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"BeeCustom/enums"
	"BeeCustom/models"
	"BeeCustom/utils"

	"github.com/astaxie/beego/orm"
)

type BackendUserController struct {
	BaseController
	Roles []*models.Role //当前用户信息
}

func (c *BackendUserController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor("DataGrid")
	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()

	var params = models.RoleQueryParam{}
	c.Roles = models.RoleDataList(&params)

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
	c.Data["canEdit"] = c.checkActionAuthor("BackendUserController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("BackendUserController", "Delete")
}

//列表数据
func (c *BackendUserController) DataGrid() {
	//直接获取参数 getDataGridData()
	var params models.BackendUserQueryParam
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

	c.Data["Roles"] = c.Roles
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "backenduser/create_footerjs.html"
}

// Store 添加 新建 页面
func (c *BackendUserController) Store() {
	c.Save(0)
}

// Edit 添加 编辑 页面
func (c *BackendUserController) Edit() {

	Id, _ := c.GetInt(":id", 0)
	m := &models.BackendUser{}
	var err error
	if Id > 0 {
		m, err = models.BackendUserOne(Id)
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}

	} else {
		//添加用户时默认状态为启用
		m.Status = enums.Enabled
	}
	c.Data["m"] = m
	c.Data["Roles"] = c.Roles
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "backenduser/edit_footerjs.html"

}

// Update 添加 编辑 页面
func (c *BackendUserController) Update() {
	Id, _ := c.GetInt(":id", 0)
	c.Save(Id)
}

//保存数据
func (c *BackendUserController) Save(id int) {
	m := models.BackendUser{
		Id: id,
	}
	o := orm.NewOrm()
	var err error

	//获取form里的值
	if err = c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
	}

	if m.Id == 0 {
		//对密码进行加密
		m.UserPwd = utils.String2md5(m.UserPwd)
		m.CreatedAt = time.Now()
		m.UpdatedAt = time.Now()

		if oR, err := models.RoleOne(m.RoleId); err != nil {
			c.jsonResult(enums.JRCodeFailed, "数据无效，请刷新后重试", m.Id)
		} else {
			m.Role = oR
		}

		if _, err := o.Insert(&m); err != nil {
			c.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
		} else {
			c.jsonResult(enums.JRCodeSucc, "添加成功", m.Id)
		}

	} else {

		if oM, err := models.BackendUserOne(m.Id); err != nil {
			c.jsonResult(enums.JRCodeFailed, "数据无效，请刷新后重试", m.Id)
		} else {
			m.UserPwd = strings.TrimSpace(m.UserPwd)
			m.CreatedAt = oM.CreatedAt
			m.UpdatedAt = time.Now()
			if len(m.UserPwd) == 0 {
				//如果密码为空则不修改
				m.UserPwd = oM.UserPwd
			} else {
				m.UserPwd = utils.String2md5(m.UserPwd)
			}
			//本页面不修改头像和密码，直接将值附给新m
			m.Avatar = oM.Avatar
		}

		if oR, err := models.RoleOne(m.RoleId); err != nil {
			c.jsonResult(enums.JRCodeFailed, "数据无效，请刷新后重试", m.Id)
		} else {
			m.Role = oR
		}

		if _, err := o.Update(&m); err != nil {
			c.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
		} else {
			c.jsonResult(enums.JRCodeSucc, "编辑成功", m.Id)
		}
	}

}

//删除
func (c *BackendUserController) Delete() {
	id, _ := c.GetInt(":id")

	o := orm.NewOrm()
	if num, err := o.Delete(&models.BackendUser{Id: id}); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), "")
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", err)
	}
}
