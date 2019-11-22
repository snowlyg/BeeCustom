package enums

import (
	"fmt"

	"BeeCustom/utils"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func NewPDFGenerator() {

	// Create new PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		utils.LogDebug(fmt.Sprintf(" NewPDFGenerator error:%v", err))
	}

	// Set global options
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationLandscape)
	pdfg.Grayscale.Set(true)

	// Create a new input page from an URL
	page := wkhtmltopdf.NewPage("https://godoc.org/github.com/SebastiaanKlippert/go-wkhtmltopdf")

	// Set options for this page
	page.FooterRight.Set("[page]")
	page.FooterFontSize.Set(10)
	page.Zoom.Set(0.95)

	// Add to document
	pdfg.AddPage(page)

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		utils.LogDebug(fmt.Sprintf(" pdfg.Create error:%v", err))
	}

	//// Write buffer contents to file on disk
	//err = pdfg.WriteFile("/static/generate/annotation/simplesample.pdf")
	//if err != nil {
	//	utils.LogDebug(fmt.Sprintf("  pdfg.WriteFile error:%v",err))
	//}

	fmt.Println("Done")
}
