package controllers

import (
	"fmt"
	"log"

	"BeeCustom/enums"
	"BeeCustom/models"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
)

// SingeController handles WebSocket requests.
type SingeController struct {
	BaseController
}

func (s *SingeController) Get() {
	var aType int8

	articles := make([]*models.Article, 0, 200)

	// 1.准备收集器实例
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36"),
		colly.AllowURLRevisit(),
		// 开启本机debug
		colly.Debugger(&debug.LogDebugger{}),
		colly.AllowedDomains("singlewindow.cn", "www.singlewindow.cn"),
		// 防止页面重复下载
		//colly.CacheDir("./Single_cache"),
	)

	// Create another collector to scrape article details
	detailCollector := c.Clone()

	//// On every a element which has href attribute call callback
	//c.OnHTML("div.lsli > a", func(e *colly.HTMLElement) {
	//	link := e.Attr("href")
	//	if len(e.Text) == 0 {
	//		return
	//	}
	//	switch e.Text {
	//	case "新闻动态":
	//		aType = 1
	//	case "通知公告":
	//		aType = 2
	//	default:
	//		aType = 1
	//	}
	//	e.Request.Visit(e.Request.AbsoluteURL(link))
	//})

	// On every a element which has href attribute call callback
	c.OnHTML("div.lsbt > a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		a, _ := models.GetArticleByTitle(e.Text)
		if a != nil && a.Id != 0 {
			return
		}
		detailCollector.Visit(e.Request.AbsoluteURL(link))
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		//r.Headers.Set("User-Agent", RandomString())
		fmt.Println("Visiting", r.URL.String())
	})

	// Extract details of the course
	detailCollector.OnHTML(`#indetmain`, func(e *colly.HTMLElement) {
		log.Println("Course found", e.Request.URL)
		title := e.ChildText("#dettitle")
		if title == "" {
			log.Println("No title found", e.Request.URL)
		}
		article := &models.Article{
			Type:    aType,
			Title:   title,
			Origin:  e.Request.URL.String(),
			Content: e.ChildText("div.wnp"),
			NewTime: e.ChildText("span#sj"),
		}
		articles = append(articles, article)
	})

	// 下一页
	c.OnHTML("a[href].page-link", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	// 启动
	c.Visit("http://www.singlewindow.cn/xwdt/index.jhtml")

	num := int64(0)
	if len(articles) > 0 {
		num, _ = models.InsertArticleMulti(articles)
	}

	s.jsonResult(enums.JRCodeSucc, fmt.Sprintf("采集 %v 个新闻", num), nil)
}
