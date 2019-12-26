package transforms

type Resource struct {
	Id         int64
	Rtype      int
	Name       string
	Icon       string
	UrlFor     string
	ParentName string `gtf:"Parent.Name"`
	Sons       []*Resource
	CreatedAt  string
	UpdatedAt  string
}
