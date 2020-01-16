package controllers

import (
	"fmt"
	"log"

	"BeeCustom/enums"
	"BeeCustom/models"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
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
		//colly.UserAgent("Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36"),
		colly.AllowURLRevisit(),
		// 开启本机debug
		//colly.Debugger(&debug.LogDebugger{}),
		colly.AllowedDomains("singlewindow.cn", "www.singlewindow.cn"),
		// 防止页面重复下载
		colly.CacheDir("./Single_cache"),
	)

	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	// Create another collector to scrape article details
	listCollector := c.Clone()
	detailCollector := c.Clone()

	// On every a element which has href attribute call callback
	c.OnHTML("li.lsli > a", func(e *colly.HTMLElement) {
		//log.Println("Course found", e.Request.URL)
		link := e.Attr("href")
		if len(e.Text) == 0 {
			return
		}
		switch e.Text {
		case "新闻动态":
			aType = 1
		case "通知公告":
			aType = 2
		default:
			aType = 1
		}

		_ = listCollector.Visit(e.Request.AbsoluteURL(link))
	})

	// On every a element which has href attribute call callback
	listCollector.OnHTML("div.lsbt > a", func(e *colly.HTMLElement) {
		//log.Println("listCollector found", e.Request.URL)
		link := e.Attr("href")
		a, _ := models.GetArticleByTitle(e.Text)
		if a != nil && a.Id != 0 {
			return
		}
		_ = detailCollector.Visit(e.Request.AbsoluteURL(link))
	})

	// Extract details of the course
	detailCollector.OnHTML(`#indetmain`, func(e *colly.HTMLElement) {
		//log.Println("detailCollector found", e.Request.URL)
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
	//listCollector.OnHTML("a[href].page-link", func(e *colly.HTMLElement) {
	//	_ = e.Request.Visit(e.Attr("href"))
	//})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	// Before making a request print "Visiting ..."
	listCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("listCollector - Visiting", r.URL.String())
	})
	// Before making a request print "Visiting ..."
	detailCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("detailCollector - Visiting", r.URL.String())
	})

	// 启动
	_ = c.Visit("http://www.singlewindow.cn/xwdt/index.jhtml")

	num := int64(0)
	if len(articles) > 0 {
		num, _ = models.InsertArticleMulti(articles)
	}

	s.jsonResult(enums.JRCodeSucc, fmt.Sprintf("采集 %v 个新闻", num), nil)
}
