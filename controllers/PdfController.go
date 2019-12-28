package controllers

import (
	"time"

	"BeeCustom/models"
)

type PdfController struct {
	BaseController
}

func (c *PdfController) AnnotationRecheckPdf() {
	id, _ := c.GetInt64(":id")
	annotation := models.TransformAnnotation(id, "AnnotationItems")
	c.Data["M"] = annotation
	c.Data["Now"] = time.Now()
	c.setTpl("annotation/pdf/recheck.html", "shared/layout_app.html")
}

func (c *PdfController) AnnotationPdf() {
	id, _ := c.GetInt64(":id")
	annotation := models.TransformAnnotation(id, "AnnotationItems")
	c.Data["M"] = annotation
	c.Data["Now"] = time.Now()
	c.setTpl("annotation/pdf/index.html", "shared/layout_app.html")
}

func (c *PdfController) OrderRecheckPdf() {
	id, _ := c.GetInt64(":id")
	order := models.TransformOrder(id, "OrderItems,OrderContainers,OrderDocuments", true)
	c.Data["m"] = order
	c.setTpl("order/pdf/recheck/index.html", "")
}

func (c *PdfController) OrderPdf() {
	id, _ := c.GetInt64(":id")
	order := models.TransformOrder(id, "OrderItems,OrderContainers,OrderDocuments", true)
	c.Data["m"] = order
	c.setTpl("order/pdf/declaration/index.html", "")
}

func (c *PdfController) OrderPdfHeader() {
	id, _ := c.GetInt64(":id")
	page, _ := c.GetInt64("page")
	topage, _ := c.GetInt64("topage")
	order := models.TransformOrder(id, "", false)
	c.Data["m"] = order
	c.Data["IsRecheck"] = false
	c.Data["topage"] = topage
	c.Data["page"] = page
	c.setTpl("order/pdf/header/header.html", "shared/layout_app.html")
}

func (c *PdfController) OrderRecheckPdfHeader() {
	id, _ := c.GetInt64(":id")
	page, _ := c.GetInt64("page")
	topage, _ := c.GetInt64("topage")
	order := models.TransformOrder(id, "", false)
	c.Data["m"] = order
	c.Data["IsRecheck"] = true
	c.Data["topage"] = topage
	c.Data["page"] = page
	c.setTpl("order/pdf/header/header.html", "shared/layout_app.html")
}
