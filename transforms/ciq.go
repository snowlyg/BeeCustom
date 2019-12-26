package transforms

type Ciq struct {
	Id          int64
	Hs          string
	Name        string
	CiqCode     string
	CiqName     string
	Version     string
	VersionDate string `gtf:"Func.FormatTime()"`
	Mark        int8
	Status      string
	CreatedAt   string
	UpdatedAt   string
}