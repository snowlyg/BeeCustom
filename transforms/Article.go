package transforms

type Article struct {
	Id        int64
	Type      string `gtf:"Func.GetTypeName()"`
	Title     string
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
