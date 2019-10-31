package controllers

import (
	"fmt"
	"strconv"

	"BeeCustom/enums"
	"BeeCustom/models"
	"BeeCustom/utils"
)

type CompanySealController struct {
	BaseController
}

func (c *CompanySealController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	//默认认证 "Index", "Create", "Edit", "Delete"
	c.checkAuthor()

	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()

}

// Create 添加 新建 页面
func (c *CompanySealController) Create() {
	Id, _ := c.GetInt64(":cid", 0)
	c.Data["companyId"] = Id

	c.setTpl("company/seal/create.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "company/seal/create_footerjs.html"
	c.GetXSRFToken()
}

// Store 添加 新建 页面
func (c *CompanySealController) Store() {
	m := models.NewCompanySeal(0)
	//获取form里的值
	if err := c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}

	c.checkSealNameCount(m.Id, m.CompanyId, m.SealName)

	c.validRequestData(m)

	if _, err := models.CompanySealSave(&m); err != nil {
		utils.LogDebug(err)
		c.jsonResult(enums.JRCodeFailed, "添加失败", m)
	} else {
		c.jsonResult(enums.JRCodeSucc, "添加成功", m)
	}
}

// Edit 添加 编辑 页面
func (c *CompanySealController) Edit() {
	Id, _ := c.GetInt64(":id", 0)
	m, err := models.CompanySealOne(Id)
	if m != nil && Id > 0 {
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}
	c.Data["m"] = m

	c.LayoutSections = make(map[string]string)
	c.setTpl("company/seal/create.html")
	c.LayoutSections["footerjs"] = "company/seal/edit_footerjs.html"
	c.GetXSRFToken()
}

// Update 添加 编辑 页面
func (c *CompanySealController) Update() {
	Id, _ := c.GetInt64(":id", 0)
	m := models.NewCompanySeal(Id)

	//获取form里的值
	if err := c.ParseForm(&m); err != nil {
		utils.LogDebug(fmt.Sprintf("获取数据失败:%v", err))
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}

	c.checkSealNameCount(m.Id, m.CompanyId, m.SealName)

	c.validRequestData(m)

	if _, err := models.CompanySealSave(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "编辑失败", m)
	} else {
		c.jsonResult(enums.JRCodeSucc, "编辑成功", m)
	}
}

//删除
func (c *CompanySealController) Delete() {
	id, _ := c.GetInt64(":id")
	if num, err := models.CompanySealDelete(id); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), "")
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", err)
	}
}

//判断公章名称，是否只有一个
func (c *CompanySealController) checkSealNameCount(id, companyId int64, sealName string) {

	params := models.NewCompanySealQueryParam()
	params.SealName = sealName
	params.CompanyId = strconv.FormatInt(companyId, 10)
	cs, count := models.CompanySealPageList(&params)

	if count == 1 {
		for _, v := range cs {
			if v.SealName == sealName {
				if v.Id != id {
					c.jsonResult(enums.JRCodeFailed, "只能有一个"+sealName, nil)
				}
			}
		}

	} else if count > 1 {
		c.jsonResult(enums.JRCodeFailed, sealName+"已经超过限制，请修改签章列表", nil)
	}

}
