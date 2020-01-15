package tasks

import (
	"fmt"

	"BeeCustom/utils"
	"github.com/astaxie/beego/toolbox"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
)

//const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

//func RandomString() string {
//	b := make([]byte, rand.Intn(10)+10)
//	for i := range b {
//		b[i] = letterBytes[rand.Intn(len(letterBytes))]
//	}
//	return string(b)
//}

func getArticles() {
	// 1.准备收集器实例
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36"),
		colly.AllowURLRevisit(),
		// 开启本机debug
		colly.Debugger(&debug.LogDebugger{}),
		colly.AllowedDomains("singlewindow.cn"),
		// 防止页面重复下载
		colly.CacheDir("./learnku_cache"),
	)

	//c.OnRequest(func(r *colly.Request) {
	//	r.Headers.Set("User-Agent", RandomString())
	//})

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		c.Visit(e.Request.AbsoluteURL(link))
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	//// 2.分析页面数据
	//c.OnHTML("div.blog-article-list > .event", func(e *colly.HTMLElement) {
	//	article := models.Article{
	//		Type:     1,
	//		Title:    e.ChildText("div.content > div.summary"),
	//		Origin:   e.ChildAttr("div.content a.title", "href"),
	//		Overview: e.ChildAttr("div.content a.title", "href"),
	//		Content:  e.ChildText("div.item-meta > a:first-child"),
	//	}
	//	// 查找同一集合不同子项
	//	e.ForEach("div.content > div.meta > div.date>a", func(i int, el *colly.HTMLElement) {
	//		switch i {
	//		case 1:
	//			article.CreatedAt = time.Now()
	//		}
	//	})
	//})
	//
	//// 下一页
	//c.OnHTML("a[href].page-link", func(e *colly.HTMLElement) {
	//	e.Request.Visit(e.Attr("href"))
	//})

	// 启动
	c.Visit("http://www.singlewindow.cn/xwdt/index.jhtml")
}

// 回执解释
func getSingleArticles() *toolbox.Task {

	task := toolbox.NewTask("task", "0 0 1 * * *", func() error {
		getArticles()
		return nil
	})

	err := task.Run()
	if err != nil {
		utils.LogError(fmt.Sprintf("tk.Run error :%v", err))
		return nil
	}

	return task
}
