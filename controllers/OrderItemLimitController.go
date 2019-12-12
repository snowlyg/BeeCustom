package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"BeeCustom/enums"
	"BeeCustom/models"
	"BeeCustom/utils"
)

type OrderItemLimitController struct {
	BaseController
}

func (c *OrderItemLimitController) Prepare() {
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
func (c *OrderItemLimitController) Store() {
	Id, _ := c.GetInt64(":aid", 0)
	m := models.NewOrderItemLimit(0)

	c.saveOrUpdate(&m, Id)
}

// Update 添加 编辑 页面
func (c *OrderItemLimitController) Update() {
	Id, _ := c.GetInt64(":id", 0)
	m := models.NewOrderItemLimit(Id)

	c.saveOrUpdate(&m, 0)
}

// Update 添加 编辑 页面
func (c *OrderItemLimitController) saveOrUpdate(m *models.OrderItemLimit, aId int64) {
	// 获取form里的值
	if err := c.ParseForm(m); err != nil {
		utils.LogDebug(fmt.Sprintf("获取数据失败:%v", err))
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}

	c.validRequestData(m)

	if aId == 0 {
		aId = m.OrderItemId
	}

	orderItem, err := models.OrderItemOne(aId)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取表头数据失败", m)
	}

	m.OrderItem = orderItem

	if err := models.OrderItemLimitSave(m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	} else {
		c.jsonResult(enums.JRCodeSucc, "操作成功", m)
	}
}

// 删除
func (c *OrderItemLimitController) Delete() {
	idsString := c.GetString("Ids")
	Ids := strings.Split(idsString, ",")

	for _, i := range Ids {
		id, err := strconv.ParseInt(i, 10, 64)
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, "删除失败", err)
		}
		if _, err := models.OrderItemLimitDelete(id); err != nil {
			c.jsonResult(enums.JRCodeFailed, "删除失败", err)
		}
	}

	c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", len(Ids)), "")

}
