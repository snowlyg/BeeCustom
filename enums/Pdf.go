package enums

import (
	"fmt"
	"strconv"

	"BeeCustom/file"
	"BeeCustom/utils"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/astaxie/beego"
)

func NewPDFGenerator(Id int64, etpsInnerInvtNo, url string) error {

	// Create new PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		utils.LogDebug(fmt.Sprintf(" NewPDFGenerator error:%v", err))
		return err
	}

	// Set global options
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationLandscape)
	pdfg.Grayscale.Set(false) //彩色 false， 灰色 true

	httpaddr := beego.AppConfig.String("httpaddr")
	httpport := beego.AppConfig.String("httpport")

	// basic auth 认证用户名和密码
	username := beego.AppConfig.String("pdf_username")
	password := beego.AppConfig.String("pdf_password")

	// Create a new input page from an URL
	page := wkhtmltopdf.NewPage(httpaddr + ":" + httpport + "/pdf/" + url + "/" + strconv.FormatInt(Id, 10))

	// Set options for this page
	page.FooterCenter.Set("第[page]页 共[topage]页")
	page.FooterFontSize.Set(12)
	page.Zoom.Set(0.95)
	page.DisableJavascript.Set(false)
	page.DebugJavascript.Set(true)
	page.MinimumFontSize.Set(12)

	page.Username.Set(username)
	page.Password.Set(password)

	// // Add to document
	pdfg.AddPage(page)

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		utils.LogDebug(fmt.Sprintf(" pdfg.Create error:%v", err))
		return err
	}

	path := "./static/generate/annotation/" + strconv.FormatInt(Id, 10)
	if err := file.CreateFile(path); err != nil {
		utils.LogDebug(fmt.Sprintf("文件夹创建失败:%v", err))
		return err
	}

	// Write buffer contents to file on disk
	err = pdfg.WriteFile(path + "/" + etpsInnerInvtNo + ".pdf")
	if err != nil {
		utils.LogDebug(fmt.Sprintf("  pdfg.WriteFile error:%v", err))
		return err
	}

	return nil
}
