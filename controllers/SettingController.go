package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"BeeCustom/enums"
	"BeeCustom/models"
	"BeeCustom/utils"
	"golang.org/x/net/html"
)

type SettingController struct {
	BaseController
}

func (c *SettingController) Prepare() {
	// 先执行
	c.BaseController.Prepare()
	// 如果一个Controller的少数Action需要权限控制，则将验证放到需要控制的Action里

	perms := []string{
		"Index",
		"Create",
		"Edit",
		"Delete",
	}
	c.checkAuthor(perms)
	// 如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	// 这里注释了权限控制，因此这里需要登录验证
	c.checkLogin()
}

func (c *SettingController) Index() {
	// 将页面左边菜单的某项激活

	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "setting/index_footerjs.html"
	// 页面里按钮权限控制
	c.getActionData("", "Edit", "Delete", "Create")
	c.GetXSRFToken()
}

//  Create 添加 新建 页面
func (c *SettingController) Create() {
	c.setTpl("setting/edit.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "setting/edit_footerjs.html"
	c.GetXSRFToken()
}

//  Store 添加 新建 页面
func (c *SettingController) Store() {
	m := models.NewSetting(0)
	// 获取form里的值
	if err := c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}
	c.validRequestData(m)
	if _, err := models.SettingSave(&m); err == nil {
		c.jsonResult(enums.JRCodeSucc, "添加成功", m)
	} else {
		c.jsonResult(enums.JRCodeFailed, "添加失败", m)
	}
}

// TreeGrid 获取所有资源的列表
func (c *SettingController) TreeGrid() {
	// 直接反序化获取json格式的requestbody里的值
	params := models.NewSettingQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	// 获取数据列表和总数
	data, total := models.SettingTreeGrid(&params)
	dataLists := c.TransformSettingList(data)

	c.ResponseList(dataLists, total)
	c.ServeJSON()
}

// GetOne 获取系统设置
func (c *SettingController) GetOne() {
	key := c.GetString(":key", "")
	if len(key) > 0 {
		rvalue, err := models.GetSettingRValueByKey(key)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("数据无效出错：%v", err))
			c.pageError("数据无效，请刷新后重试")
		}

		//单个设置
		if len(rvalue) == 1 {
			c.Data["json"] = rvalue["0"]
		} else {
			c.Data["json"] = rvalue
		}

	}

	c.ServeJSON()
}

// Edit 资源编辑页面
func (c *SettingController) Edit() {
	Id, _ := c.GetInt64(":id", 0)
	c.Data["m"], _ = models.SettingOne(Id)
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "setting/edit_footerjs.html"
	c.GetXSRFToken()
}

// Update 添加、编辑角色界面
func (c *SettingController) Update() {
	id, _ := c.GetInt64(":id", 0)
	m := models.NewSetting(id)

	// 获取form里的值
	if err := c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}
	c.validRequestData(m)
	if _, err := models.SettingSave(&m); err == nil {
		c.jsonResult(enums.JRCodeSucc, "编辑成功", m)
	} else {
		c.jsonResult(enums.JRCodeFailed, "编辑失败", m)
	}
}

//  Delete 删除
func (c *SettingController) Delete() {
	Id, _ := c.GetInt64(":id", 0)
	if Id == 0 {
		c.jsonResult(enums.JRCodeFailed, "选择的数据无效", 0)
	}

	if _, err := models.SettingDelete(Id); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("删除成功"), 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
	}
}

// TransformSettingList 格式化列表数据
func (c *SettingController) TransformSettingList(ms []*models.Setting) []*map[string]interface{} {
	var dataLists []*map[string]interface{}
	for _, v := range ms {
		value := html.UnescapeString(v.Value)
		valueEnd := value[:len(value)-1]
		if len(value) > 30 {
			valueEnd = value[:30] + "..."
		}
		dataList := make(map[string]interface{})
		dataList["Id"] = strconv.FormatInt(v.Id, 10)
		dataList["Key"] = v.Key
		dataList["Value"] = valueEnd
		dataList["Rmk"] = v.Rmk
		dataLists = append(dataLists, &dataList)
	}
	return dataLists
}
