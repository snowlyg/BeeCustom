package controllers

import (
	"encoding/json"

	"BeeCustom/enums"
	"BeeCustom/models"
	"BeeCustom/transforms"
	"github.com/snowlyg/gotransform"
)

type OrderFileController struct {
	BaseController
}

func (c *OrderFileController) Prepare() {
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
func (c *OrderFileController) DataGrid() {
	// 直接获取参数 getDataGridData()
	params := models.NewOrderFileQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	// 获取数据列表和总数
	data, total := models.OrderFilePageList(&params)
	c.ResponseList(c.TransformOrderFileList(data), total)
	c.ServeJSON()
}

// TransformOrderFile 格式化列表数据
func (c *OrderFileController) TransformOrderFileList(ms []*models.OrderFile) []*transforms.OrderFile {

	var uts []*transforms.OrderFile
	for _, v := range ms {
		ut := &transforms.OrderFile{}
		g := gotransform.NewTransform(ut, v, enums.BaseDateTimeFormat)
		_ = g.Transformer()

		uts = append(uts, ut)
	}
	return uts

}
