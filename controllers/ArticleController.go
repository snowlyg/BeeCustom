package controllers

import (
	"encoding/json"

	"BeeCustom/enums"
	"BeeCustom/transforms"
	"BeeCustom/utils"
	gtf "github.com/snowlyg/gotransformer"

	"BeeCustom/models"
)

type ArticleController struct {
	BaseController
}

func (c *ArticleController) Prepare() {
	//先执行
	c.BaseController.Prepare()

	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()

}

//列表数据
func (c *ArticleController) DataGrid() {
	//直接获取参数 getDataGridData()
	params := models.NewArticleQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	utils.LogDebug(params.Type)

	//获取数据列表和总数
	data, total := models.ArticlePageList(&params)
	c.ResponseList(c.transformArticleList(data), total)
	c.ServeJSON()
}

//  格式化列表数据
func (c *ArticleController) transformArticleList(ms []*models.Article) []*transforms.Article {
	var uts []*transforms.Article
	for _, v := range ms {
		ut := transforms.Article{}
		g := gtf.NewTransform(&ut, v, enums.BaseDateTimeFormat)
		_ = g.Transformer()

		uts = append(uts, &ut)
	}

	return uts
}
