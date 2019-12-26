package transforms

import (
	"BeeCustom/models"
)

type Company struct {
	Id                  int64
	Number              string
	AdminName           string `gtf:"Func.GetAdminName(CompanyContacts)"`
	Name                string
	Short               string
	Registration        string
	Address             string
	DeclareType         int8
	RegistrationCode    string
	Phone               string
	Fax                 string
	CreditCode          string
	BusinessName        string
	BusinessCode        string
	Bank                string
	CustomCertification int8
	CompanyType         int8
	CompanyKind         int8
	ControlRating       int8
	Remark              string
	IsOpenSubEmail      int8
	IsOpenSubPhone      int8
	SubPhone            string
	SubEmail            string
	SubContentCheck     int8
	SubContentSubmit    int8
	SubContentReject    int8
	SubContentPass      int8
	StatementDate       int8
	IsTrade             int8
	IsOwner             int8
	Business            string
	BusinessAuditStatus int8
	BusinessAuditAt     string `gtf:"Func.FormatTime()"`
	Tax                 int8
	CreatedAt           string
	UpdatedAt           string
}

func (c *Company) GetAdminName(vs []*models.CompanyContact) string {
	for _, v := range vs {
		if v.IsAdmin == 1 {
			return v.Name
		}
	}

	return ""
}
