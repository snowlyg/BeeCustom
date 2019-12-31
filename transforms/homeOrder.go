package transforms

type HomeOrder struct {
	Country string `gtf:"Func.GetCountry()"`
	Value   int
	Year    string
}

func (s *HomeOrder) GetCountry(v string) string {
	if v == "I" {
		return "进口"
	} else if v == "E" {
		return "出口"
	}
	return ""
}
