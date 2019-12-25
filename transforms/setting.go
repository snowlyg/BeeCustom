package transforms

import "html"

type Setting struct {
	Id        int64
	Key       string
	Value     string `gtf:"Func.GetValueEnd"`
	Rmk       string
	DeletedAt string
	CreatedAt string
	UpdatedAt string
}

func (s *Setting) GetValueEnd(v string) string {
	value := html.UnescapeString(v)
	valueEnd := value[:len(value)-1]
	if len(value) > 30 {
		valueEnd = value[:30] + "..."
	}

	return valueEnd
}
