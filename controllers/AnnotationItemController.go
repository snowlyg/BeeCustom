package controllers

import (
	"encoding/json"
	"fmt"

	"BeeCustom/enums"
	"BeeCustom/models"
	"BeeCustom/utils"
)

type AnnotationItemController struct {
	BaseController
}

func (c *AnnotationItemController) Prepare() {
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

// 列表数据
func (c *AnnotationItemController) DataGrid() {
	// 直接获取参数 getDataGridData()
	params := models.NewAnnotationItemQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	// 获取数据列表和总数
	data, total := models.AnnotationItemPageList(&params)

	// 定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	result["code"] = 0
	c.Data["json"] = result

	c.ServeJSON()
}

// Store 添加 新建 页面
func (c *AnnotationItemController) Store() {
	Id, _ := c.GetInt64(":aid", 0)
	m := models.NewAnnotationItem(0)

	c.saveOrUpdate(&m, Id)
}

// Update 添加 编辑 页面
func (c *AnnotationItemController) Update() {
	Id, _ := c.GetInt64(":id", 0)
	m := models.NewAnnotationItem(Id)

	c.saveOrUpdate(&m, 0)

}

// Update 添加 编辑 页面
func (c *AnnotationItemController) saveOrUpdate(m *models.AnnotationItem, aId int64) {
	// 获取form里的值
	if err := c.ParseForm(m); err != nil {
		utils.LogDebug(fmt.Sprintf("获取数据失败:%v", err))
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}

	c.validRequestData(m)

	if aId == 0 {
		aId = m.AnnotationId
	}

	annotation, err := models.AnnotationOne(aId, "")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取表头数据失败", m)
	}

	m.Annotation = annotation

	if err := models.AnnotationItemSave(m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	} else {
		c.jsonResult(enums.JRCodeSucc, "操作成功", m)
	}
}

// 删除
func (c *AnnotationItemController) Delete() {
	id, _ := c.GetInt64(":id")

	annotationItem, err := models.AnnotationItemOne(id)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "删除失败", err)
	}
	aId := annotationItem.Annotation.Id

	if _, err := models.AnnotationItemDelete(id); err != nil {
		c.jsonResult(enums.JRCodeFailed, "删除失败", err)
	}

	annotationItems, err := models.AnnotationItemsByAnnotationId(aId)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "删除失败", err)
	}

	/*修改商品序号*/
	for i, _ := range annotationItems {
		item := annotationItems[i]
		if item.GdsSeqno != i+1 {
			item.GdsSeqno = i + 1
			err = models.AnnotationItemSave(item)
			if err != nil {
				c.jsonResult(enums.JRCodeFailed, "删除失败", err)
			}
		}

	}

	c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", 0), "")

}
