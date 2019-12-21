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
	c.Data["M"] = order
	c.Data["Now"] = time.Now()
	c.setTpl("order/pdf/recheck.html", "shared/layout_app.html")
}

func (c *PdfController) OrderPdf() {
	id, _ := c.GetInt64(":id")
	order := models.TransformOrder(id, "OrderItems,OrderContainers,OrderDocuments", true)
	c.Data["M"] = order
	c.Data["Now"] = time.Now()
	c.setTpl("order/pdf/index.html", "shared/layout_app.html")
}
