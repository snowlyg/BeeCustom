package enums

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	"BeeCustom/utils"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/astaxie/beego"
)

//var (
//	htmlTplEngine  *template.Template
//	htmlTplEngineErr error
//)

func NewPDFGenerator(m interface{}) error {

	// Create new PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		utils.LogDebug(fmt.Sprintf(" NewPDFGenerator error:%v", err))
		return err
	}

	// Set global options
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationLandscape)
	pdfg.Grayscale.Set(true)

	httpaddr := beego.AppConfig.String("httpaddr")
	httpport := beego.AppConfig.String("httpport")
	// Create a new input page from an URL
	page := wkhtmltopdf.NewPage(httpaddr + ":" + httpport + "/backenduser/index")

	tmpl := template.Must(template.ParseFiles("./views/annotation/pdf/index.html"))
	// Error checking elided
	data := struct {
		M interface{}
	}{M: m}

	err = tmpl.Execute(os.Stdin, data)
	if err != nil {
		utils.LogDebug(fmt.Sprintf(" Execute error:%v", err))
	}

	//_ = htmlTplEngine.ExecuteTemplate(
	//	w,
	//	"index/index",
	//	map[string]interface{}{"PageTitle": "首页", "Name": "sqrt_cat", "Age": 25},
	//)

	// Set options for this page
	page.FooterRight.Set("[page]")
	page.FooterFontSize.Set(10)
	page.Zoom.Set(0.95)

	//// Add to document
	//pdfg.AddPage(page)

	html := "<html>Hi</html>"
	pdfg.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(html)))

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		utils.LogDebug(fmt.Sprintf(" pdfg.Create error:%v", err))
		return err
	}
	//
	//path := "./static/generate/annotation/" + strconv.FormatInt(m.Id, 10)
	//if err := file.CreateFile(path); err != nil {
	//	utils.LogDebug(fmt.Sprintf("文件夹创建失败:%v", err))
	//	return err
	//}

	//// Write buffer contents to file on disk
	//err = pdfg.WriteFile(path + "/simplesample.pdf")
	//if err != nil {
	//	utils.LogDebug(fmt.Sprintf("  pdfg.WriteFile error:%v", err))
	//	return err
	//}

	fmt.Println("Done")
	return nil
}
