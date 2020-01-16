package transforms

import "html"

type Article struct {
	Id        int64
	Type      string `gtf:"Func.GetTypeName()"`
	Title     string `gtf:"Func.GetValueEnd()"`
	Content   string
	Overview  string
	Origin    string
	NewTime   string
	CreatedAt string
	UpdatedAt string
}

func (s *Article) GetTypeName(v string) string {
	if v == "1" {
		return "新闻动态"
	} else if v == "2" {
		return "通知公告"
	}

	return "新闻动态"
}

func (s *Article) GetValueEnd(v string) string {
	value := []rune(html.UnescapeString(v))

	valueEnd := string(value[:len(value)-1])
	if len(value) > 40 {
		valueEnd = string(value[:40]) + `...`
	}

	return valueEnd
}
