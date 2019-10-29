package controllers

import (
	"encoding/json"

	"BeeCustom/models"
)

type CiqController struct {
	BaseController
}

func (c *CiqController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	//默认认证 "Index", "Create", "Edit", "Delete"
	c.checkAuthor()

	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()

}

func (c *CiqController) Index() {

	//是否显示更多查询条件的按钮弃用，前端自动判断
	//c.Data["showMoreQuery"] = true
	//将页面左边菜单的某项激活
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.Data["lastUpdateTime"] = c.GetLastUpdteTime("ciqLastUpdteTime")

	//页面模板设置
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "ciq/index_footerjs.html"

	//页面里按钮权限控制
	c.getActionData("Edit")

	c.GetXSRFToken()
}

//列表数据
func (c *CiqController) DataGrid() {
	//直接获取参数 getDataGridData()
	params := models.NewCiqQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	//获取数据列表和总数
	data, total := models.CiqPageList(&params)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	result["code"] = 0
	c.Data["json"] = result

	c.ServeJSON()
}
