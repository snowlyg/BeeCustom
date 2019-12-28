package enums

import (
	"fmt"
	"strconv"

	"BeeCustom/file"
	"BeeCustom/utils"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/astaxie/beego"
)

type PdfData struct {
	Id              int64
	EtpsInnerInvtNo string
	Url             string
	Action          string
	ModelName       string
	Header          string
	Username        string
	Password        string
	MarginTop       uint
}

func NewPDFGenerator(pdfData *PdfData) (string, error) {

	// Create new PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		utils.LogDebug(fmt.Sprintf(" NewPDFGenerator error:%v", err))
		return "", err
	}

	// Set global options
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationLandscape)
	pdfg.Grayscale.Set(false)             // 彩色 false， 灰色 true
	pdfg.MarginTop.Set(pdfData.MarginTop) // 上边距

	httpaddr := beego.AppConfig.String("httpaddr")
	httpport := beego.AppConfig.String("httpport")

	// Create a new input page from an URL
	formatInt := strconv.FormatInt(pdfData.Id, 10)
	page := wkhtmltopdf.NewPage(httpaddr + ":" + httpport + "/pdf/" + pdfData.Url + "/" + formatInt)

	// Set options for this page
	if pdfData.ModelName == "annotation" {
		page.FooterCenter.Set("第[page]页 共[topage]页")
		page.FooterFontSize.Set(12)
	}

	if pdfData.ModelName == "order" {
		page.HeaderHTML.Set(httpaddr + ":" + httpport + "/pdf/" + pdfData.Header + "/" + formatInt)
	}

	// page.Zoom.Set(0.95)
	page.DisableJavascript.Set(false)
	//page.EnableJavascript.Set(true)
	page.DebugJavascript.Set(true)
	page.MinimumFontSize.Set(12)

	page.Username.Set(pdfData.Username)
	page.Password.Set(pdfData.Password)

	// // Add to document
	pdfg.AddPage(page)

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		utils.LogDebug(fmt.Sprintf(" pdfg.Create error:%v", err))
		return "", err
	}

	path := "./static/generate/" + pdfData.ModelName + "/" + formatInt + "/" + pdfData.Action
	if err := file.CreateFile(path); err != nil {
		utils.LogDebug(fmt.Sprintf("文件夹创建失败:%v", err))
		return "", err
	}

	// Write buffer contents to file on disk
	fileFullPath := path + "/" + pdfData.EtpsInnerInvtNo + ".pdf"
	err = pdfg.WriteFile(fileFullPath)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("  pdfg.WriteFile error:%v", err))
		return "", err
	}

	return fileFullPath, nil
}
