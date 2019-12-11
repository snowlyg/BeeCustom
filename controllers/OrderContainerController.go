package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"BeeCustom/enums"
	"BeeCustom/models"
	"BeeCustom/utils"
)

type OrderContainerController struct {
	BaseController
}

func (c *OrderContainerController) Prepare() {
	// 先执行
	c.BaseController.Prepare()
	// 如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	// 默认认证 "Index", "Create", "Edit", "Delete"
	var perms []string
	c.checkAuthor(perms)

	// 如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	// 权限控制里会进行登录验证，因此这里不用再作登录验证
	// c.checkLogin()

}

// Store 添加 新建 页面
func (c *OrderContainerController) Store() {
	Id, _ := c.GetInt64(":aid", 0)
	m := models.NewOrderContainer(0)

	c.saveOrUpdate(&m, Id)
}

// Update 添加 编辑 页面
func (c *OrderContainerController) Update() {
	Id, _ := c.GetInt64(":id", 0)
	m := models.NewOrderContainer(Id)

	c.saveOrUpdate(&m, 0)
}

// Update 添加 编辑 页面
func (c *OrderContainerController) saveOrUpdate(m *models.OrderContainer, aId int64) {
	// 获取form里的值
	if err := c.ParseForm(m); err != nil {
		utils.LogDebug(fmt.Sprintf("获取数据失败:%v", err))
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}

	c.validRequestData(m)

	if aId == 0 {
		aId = m.OrderId
	}

	annotation, err := models.OrderOne(aId, "")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取表头数据失败", m)
	}

	m.Order = annotation

	if err := models.OrderContainerSave(m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	} else {
		c.jsonResult(enums.JRCodeSucc, "操作成功", m)
	}
}

// 删除
func (c *OrderContainerController) Delete() {
	idsString := c.GetString("Ids")
	Ids := strings.Split(idsString, ",")
	utils.LogDebug(idsString)
	utils.LogDebug(Ids)
	for _, i := range Ids {
		id, err := strconv.ParseInt(i, 10, 64)
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, "删除失败", err)
		}
		if _, err := models.OrderContainerDelete(id); err != nil {
			c.jsonResult(enums.JRCodeFailed, "删除失败", err)
		}
	}

	c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", 0), "")

}
