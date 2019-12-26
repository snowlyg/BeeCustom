package controllers

import (
	"encoding/json"

	"BeeCustom/enums"
	"BeeCustom/models"
	"BeeCustom/transforms"
	"github.com/snowlyg/gotransform"
)

type AnnotationReturnController struct {
	BaseController
}

func (c *AnnotationReturnController) Prepare() {
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
func (c *AnnotationReturnController) DataGrid() {
	// 直接获取参数 getDataGridData()
	params := models.NewAnnotationReturnQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	// 获取数据列表和总数
	data, total := models.AnnotationReturnPageList(&params)
	c.ResponseList(c.transformAnnotationReturnList(data), total)
	c.ServeJSON()
}

//  格式化列表数据
func (c *AnnotationReturnController) transformAnnotationReturnList(ms []*models.AnnotationReturn) []*transforms.AnnotationReturn {
	var annotationReturnList []*transforms.AnnotationReturn
	for _, v := range ms {
		annotationReturnT := transforms.AnnotationReturn{}
		g := gotransform.NewTransform(&annotationReturnT, v, enums.BaseDateTimeFormat)
		_ = g.Transformer()
		annotationReturnList = append(annotationReturnList, &annotationReturnT)
	}

	return annotationReturnList
}
