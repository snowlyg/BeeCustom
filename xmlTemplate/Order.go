package xmlTemplate

/**清单报文结构*/

import "encoding/xml"

type DecHead struct {
	SeqNo                  string `xml:"SeqNo"`
	IEFlag                 string `xml:"IEFlag"`
	Type                   string `xml:"Type"`
	AgentCode              string `xml:"AgentCode"`
	AgentName              string `xml:"AgentName"`
	ApprNo                 string `xml:"ApprNo"`
	BillNo                 string `xml:"BillNo"`
	ContrNo                string `xml:"ContrNo"`
	CustomMaster           string `xml:"CustomMaster"`
	CutMode                string `xml:"CutMode"`
	DistinatePort          string `xml:"DistinatePort"`
	FeeCurr                string `xml:"FeeCurr"`
	FeeMark                string `xml:"FeeMark"`
	FeeRate                string `xml:"FeeRate"`
	GrossWet               string `xml:"GrossWet"`
	IEDate                 string `xml:"IEDate"`
	IEPort                 string `xml:"IEPort"`
	InsurCurr              string `xml:"InsurCurr"`
	InsurMark              string `xml:"InsurMark"`
	InsurRate              string `xml:"InsurRate"`
	LicenseNo              string `xml:"LicenseNo"`
	ManualNo               string `xml:"ManualNo"`
	NetWt                  string `xml:"NetWt"`
	NoteS                  string `xml:"NoteS"`
	OtherCurr              string `xml:"OtherCurr"`
	OtherMark              string `xml:"OtherMark"`
	OtherRate              string `xml:"OtherRate"`
	OwnerCode              string `xml:"OwnerCode"`
	OwnerName              string `xml:"OwnerName"`
	PackNo                 string `xml:"PackNo"`
	TradeCode              string `xml:"TradeCode"`
	TradeCountry           string `xml:"TradeCountry"`
	TradeMode              string `xml:"TradeMode"`
	TradeName              string `xml:"TradeName"`
	TrafMode               string `xml:"TrafMode"`
	TrafName               string `xml:"TrafName"`
	TransMode              string `xml:"TransMode"`
	WrapType               string `xml:"WrapType"`
	EntryId                string `xml:"EntryId"`
	PreEntryId             string `xml:"PreEntryId"`
	EdiId                  string `xml:"EdiId"`
	Risk                   string `xml:"Risk"`
	CopName                string `xml:"CopName"`
	CopCode                string `xml:"CopCode"`
	EntryType              string `xml:"EntryType"`
	PDate                  string `xml:"PDate"`
	TypistNo               string `xml:"TypistNo"`
	InputerName            string `xml:"InputerName"`
	PartenerID             string `xml:"PartenerID"`
	TgdNo                  string `xml:"TgdNo"`
	DataSource             string `xml:"DataSource"`
	DeclTrnRel             string `xml:"DeclTrnRel"`
	ChkSurety              string `xml:"ChkSurety"`
	BillType               string `xml:"BillType"`
	CopCodeScc             string `xml:"CopCodeScc"`
	OwnerCodeScc           string `xml:"OwnerCodeScc"`
	AgentCodeScc           string `xml:"AgentCodeScc"`
	TradeCoScc             string `xml:"TradeCoScc"`
	PromiseItmes           string `xml:"PromiseItmes"`
	TradeAreaCode          string `xml:"TradeAreaCode"`
	CheckFlow              string `xml:"CheckFlow"`
	TaxAaminMark           string `xml:"TaxAaminMark"`
	MarkNo                 string `xml:"MarkNo"`
	DespPortCode           string `xml:"DespPortCode"`
	EntyPortCode           string `xml:"EntyPortCode"`
	GoodsPlace             string `xml:"GoodsPlace"`
	BLNo                   string `xml:"BLNo"`
	InspOrgCode            string `xml:"InspOrgCode"`
	SpecDeclFlag           string `xml:"SpecDeclFlag"`
	PurpOrgCode            string `xml:"PurpOrgCode"`
	DespDate               string `xml:"DespDate"`
	CmplDschrgDt           string `xml:"CmplDschrgDt"`
	CorrelationReasonFlag  string `xml:"CorrelationReasonFlag"`
	VsaOrgCode             string `xml:"VsaOrgCode"`
	OrigBoxFlag            string `xml:"OrigBoxFlag"`
	DeclareName            string `xml:"DeclareName"`
	NoOtherPack            string `xml:"NoOtherPack"`
	OrgCode                string `xml:"OrgCode"`
	OverseasConsignorCode  string `xml:"OverseasConsignorCode"`
	OverseasConsignorCname string `xml:"OverseasConsignorCname"`
	OverseasConsignorEname string `xml:"OverseasConsignorEname"`
	OverseasConsignorAddr  string `xml:"OverseasConsignorAddr"`
	OverseasConsigneeCode  string `xml:"OverseasConsigneeCode"`
	OverseasConsigneeEname string `xml:"OverseasConsigneeEname"`
	DomesticConsigneeEname string `xml:"DomesticConsigneeEname"`
	CorrelationNo          string `xml:"CorrelationNo"`
	EdiRemark              string `xml:"EdiRemark"`
	EdiRemark2             string `xml:"EdiRemark2"`
	TradeCiqCode           string `xml:"TradeCiqCode"`
	OwnerCiqCode           string `xml:"OwnerCiqCode"`
	DeclCiqCode            string `xml:"DeclCiqCode"`
}

type DecLists struct {
	DecList DecList `xml:"DecList"`
}

type DecList struct {
	ClassMark          string         `xml:"ClassMark"`
	CodeTS             string         `xml:"CodeTS"`
	ContrItem          string         `xml:"ContrItem"`
	DeclPrice          string         `xml:"DeclPrice"`
	DutyMode           string         `xml:"DutyMode"`
	Factor             string         `xml:"Factor"`
	GModel             string         `xml:"GModel"`
	GName              string         `xml:"GName"`
	GNo                string         `xml:"GNo"`
	OriginCountry      string         `xml:"OriginCountry"`
	TradeCurr          string         `xml:"TradeCurr"`
	DeclTotal          string         `xml:"DeclTotal"`
	GQty               string         `xml:"GQty"`
	FirstQty           string         `xml:"FirstQty"`
	SecondQty          string         `xml:"SecondQty"`
	GUnit              string         `xml:"GUnit"`
	FirstUnit          string         `xml:"FirstUnit"`
	SecondUnit         string         `xml:"SecondUnit"`
	UseTo              string         `xml:"UseTo"`
	WorkUsd            string         `xml:"WorkUsd"`
	ExgNo              string         `xml:"ExgNo"`
	ExgVersion         string         `xml:"ExgVersion"`
	DestinationCountry string         `xml:"DestinationCountry"`
	CiqCode            string         `xml:"CiqCode"`
	DeclGoodsEname     string         `xml:"DeclGoodsEname"`
	OrigPlaceCode      string         `xml:"OrigPlaceCode"`
	Purpose            string         `xml:"Purpose"`
	ProdValidDt        string         `xml:"ProdValidDt"`
	ProdQgp            string         `xml:"ProdQgp"`
	GoodsAttr          string         `xml:"GoodsAttr"`
	Stuff              string         `xml:"Stuff"`
	Uncode             string         `xml:"Uncode"`
	DangName           string         `xml:"DangName"`
	DangPackType       string         `xml:"DangPackType"`
	DangPackSpec       string         `xml:"DangPackSpec"`
	EngManEntCnm       string         `xml:"EngManEntCnm"`
	NoDangFlag         string         `xml:"NoDangFlag"`
	DestCode           string         `xml:"DestCode"`
	GoodsSpec          string         `xml:"GoodsSpec"`
	GoodsModel         string         `xml:"GoodsModel"`
	GoodsBrand         string         `xml:"GoodsBrand"`
	ProduceDate        string         `xml:"ProduceDate"`
	ProdBatchNo        string         `xml:"ProdBatchNo"`
	DistrictCode       string         `xml:"DistrictCode"`
	CiqName            string         `xml:"CiqName"`
	DecGoodsLimits     DecGoodsLimits `xml:"DecGoodsLimits"`
	MnufctrRegNo       string         `xml:"MnufctrRegNo"`
	MnufctrRegName     string         `xml:"MnufctrRegName"`
}

type DecGoodsLimits struct {
	DecGoodsLimit DecGoodsLimit `xml:"DecGoodsLimit"`
}

type DecGoodsLimit struct {
}

type DecLicenseDocus struct {
}

type DecContainers struct {
}

type DecSign struct {
	OperType    string `xml:"OperType"`
	ICCode      string `xml:"ICCode"`
	CopCode     string `xml:"CopCode"`
	OperName    string `xml:"OperName"`
	ClientSeqNo string `xml:"ClientSeqNo"`
	Sign        string `xml:"Sign"`
	SignDate    string `xml:"SignDate"`
	Certificate string `xml:"Certificate"`
	HostId      string `xml:"HostId"`
	BillSeqNo   string `xml:"BillSeqNo"`
	DomainId    string `xml:"DomainId"`
	Note        string `xml:"Note"`
}

type DecFreeTxt struct {
	RelId    string `xml:"RelId"`
	RelManNo string `xml:"RelManNo"`
	BonNo    string `xml:"BonNo"`
	VoyNo    string `xml:"VoyNo"`
	DecBpNo  string `xml:"DecBpNo"`
	CusFie   string `xml:"CusFie"`
	DecNo    string `xml:"DecNo"`
}

type DecRequestCerts struct {
}

type DecOtherPacks struct {
}

type DecCopLimits struct {
}

type DecUsers struct {
}

type DecCopPromises struct {
	DecCopPromise DecCopPromise `xml:"DecCopPromise"`
}

type DecCopPromise struct {
	DeclaratioMaterialCode string `xml:"DeclaratioMaterialCode"`
}

type DecMessage struct {
	XMLName         xml.Name        `xml:"DecMessage"`   // 该XML文件的根元素为 Signature
	Version         string          `xml:"version,attr"` // 该值会作为Signature元素的属性
	Xmlns           string          `xml:"xmlns,attr"`   // 该值会作为Signature元素的属性
	DecHead         DecHead         `xml:"DecHead"`
	DecLists        DecLists        `xml:"DecLists"`
	DecLicenseDocus DecLicenseDocus `xml:"DecLicenseDocus"`
	DecContainers   DecContainers   `xml:"DecContainers"`
	DecSign         DecSign         `xml:"DecSign"`
	DecFreeTxt      DecFreeTxt      `xml:"DecFreeTxt"`
	DecRequestCerts DecRequestCerts `xml:"DecRequestCerts"`
	DecOtherPacks   DecOtherPacks   `xml:"DecOtherPacks"`
	DecCopLimits    DecCopLimits    `xml:"DecCopLimits"`
	DecUsers        DecUsers        `xml:"DecUsers"`
	DecCopPromises  DecCopPromises  `xml:"DecCopPromises"`
}
