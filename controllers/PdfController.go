package controllers

import (
	"time"

	"BeeCustom/models"
)

type PdfController struct {
	BaseController
}

func (c *PdfController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	perms := []string{}
	c.checkAuthor(perms)

	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()
}

func (c *PdfController) AnnotationPdf() {
	id, _ := c.GetInt64(":id")
	annotation := models.TransformAnnotation(id, "AnnotationItems")
	c.Data["M"] = annotation
	c.Data["Now"] = time.Now()
	c.setTpl("annotation/pdf/recheck.html", "shared/layout_app.html")
}
