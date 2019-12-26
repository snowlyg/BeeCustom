package transforms

import "html"

type Setting struct {
	Id        int64
	Key       string
	RValue    string `gtf:"Func.GetValueEnd()"`
	Rmk       string
	CreatedAt string
	UpdatedAt string
}

func (s *Setting) GetValueEnd(v string) string {
	value := []rune(html.UnescapeString(v))

	valueEnd := string(value[:len(value)-1])
	if len(value) > 30 {
		valueEnd = string(value[:30]) + `...`
	}

	return valueEnd
}
