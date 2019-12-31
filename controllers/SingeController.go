package controllers

import (
	"fmt"
	"net/http"

	"BeeCustom/utils"
	"github.com/PuerkitoBio/goquery"
)

// SingeController handles WebSocket requests.
type SingeController struct {
	BaseController
}

func (c *SingeController) Get() {
	c.EnableRender = false
	//req := httplib.Get("https://www.singlewindow.cn")
	////req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	//req.SetTimeout(100*time.Second, 30*time.Second)
	//req.Param("username", "13800138000")
	//req.Param("password", "DongHua@22898086")
	//req.Header("Accept-Encoding", "gzip,deflate,sdch")
	//req.Header("Host", "www.singlewindow.cn")
	//req.Header("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.57 Safari/537.36")
	//str, err := req.String()
	//if err != nil {
	//	utils.LogDebug(fmt.Sprintf("req.String err %v \n", err))
	//}
	//utils.LogDebug(fmt.Sprintf("req.String %v \n", str))
	//
	//res, err := req.Response()
	//if err != nil {
	//	utils.LogDebug(fmt.Sprintf("req.Response err %v \n", err))
	//}
	//utils.LogDebug(fmt.Sprintf("req.Response %v \n", res))
	//
	//req.Debug(true)

	// Request the HTML page.
	res, err := http.Get("https://beego.me/docs/module/httplib.md")
	if err != nil {
		utils.LogDebug(fmt.Sprintf("req.Get %v \n", err))
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		utils.LogDebug(fmt.Sprintf("req.StatusCode %v \n", err))
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("req.NewDocumentFromReader %v \n", err))
	}
	utils.LogDebug(fmt.Sprintf("Cookies %v \n", res.Cookies()))
	utils.LogDebug(fmt.Sprintf("StatusCode %v \n", res.StatusCode))
	utils.LogDebug(fmt.Sprintf("Url %v \n", doc.Url))

	// Find the review items
	doc.Find(".nav.navbar-nav").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		band := s.Find("a").Text()
		utils.LogDebug(fmt.Sprintf("Review %d: %s \n", i, band))
	})

}
