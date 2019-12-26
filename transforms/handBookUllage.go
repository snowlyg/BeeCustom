package transforms

type HandBookUllage struct {
	Id                    int64
	OriginalityProNo      int8
	OriginalityProName    string
	OriginalityProSpecial string
	OriginalityProU       string
	OnlyUllage            string
	Ullage                string
	Gedition              string
	Serial                string
	FinishProNo           string
	FinishRecordNo        string
	FinishHsCode          string
	FinishName            string
	FinishSpecial         string
	FinishSpecialU        string
	OriginalityRecordNo   string
	OriginalityHsCode     string
	OnlyUllageVersion     string
	OneUllage             string
	NoUllage              string
	OnlyUllageStatus      string
	ChangeMark            string
	BondedRate            string
	CompanyExecuteFlag    string
	OnlyUllageAt          string `gtf:"Func.FormatTime()"`
	UllageFlag            string
	TalkFlag              string
	Remark                string
	CreatedAt             string
	UpdatedAt             string
}
