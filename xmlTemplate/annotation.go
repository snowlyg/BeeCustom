package xmlTemplate

/**清单报文结构*/

import "encoding/xml"

type DigestMethod struct {
	Algorithm string `xml:"Algorithm,attr"`
}

type Reference struct {
	URI          string       `xml:"URI,attr"`
	DigestMethod DigestMethod `xml:"DigestMethod"`
	DigestValue  string       `xml:"DigestValue"`
}

type CanonicalizationMethod struct {
	Algorithm string `xml:"Algorithm,attr"`
}

type SignatureMethod struct {
	Algorithm string `xml:"Algorithm,attr"`
}

type SignedInfo struct {
	CanonicalizationMethod CanonicalizationMethod `xml:"CanonicalizationMethod"`
	SignatureMethod        SignatureMethod        `xml:"SignatureMethod"`
	Reference              Reference              `xml:"Reference"`
}

type EnvelopInfo struct {
	Version     string `xml:"version"`
	BusinessId  string `xml:"business_id"`
	MessageId   string `xml:"message_id"`
	FileName    string `xml:"file_name"`
	MessageType string `xml:"message_type"`
	SenderId    string `xml:"sender_id"`
	ReceiverId  string `xml:"receiver_id"`
	SendTime    string `xml:"send_time"`
}

type PocketInfo struct {
	PocketId       string `xml:"pocket_id"`
	TotalPocketQty string `xml:"total_pocket_qty"`
	CurPocketNo    string `xml:"cur_pocket_no"`
	IsUnstructured string `xml:"is_unstructured"`
}

type InvtHeadType struct {
	SeqNo                  string `xml:"SeqNo"`
	BondInvtNo             string `xml:"BondInvtNo"`
	ChgTmsCnt              string `xml:"ChgTmsCnt"`
	PutrecNo               string `xml:"PutrecNo"`
	InvtType               string `xml:"InvtType"`
	EtpsInnerInvtNo        string `xml:"EtpsInnerInvtNo"`
	BizopEtpsno            string `xml:"BizopEtpsno"`
	BizopEtpsSccd          string `xml:"BizopEtpsSccd"`
	BizopEtpsNm            string `xml:"BizopEtpsNm"`
	RcvgdEtpsno            string `xml:"RcvgdEtpsno"`
	RvsngdEtpsSccd         string `xml:"RvsngdEtpsSccd"`
	RcvgdEtpsNm            string `xml:"RcvgdEtpsNm"`
	DclEtpsno              string `xml:"DclEtpsno"`
	DclEtpsSccd            string `xml:"DclEtpsSccd"`
	DclEtpsNm              string `xml:"DclEtpsNm"`
	InputCode              string `xml:"InputCode"`
	InputCreditCode        string `xml:"InputCreditCode"`
	InputName              string `xml:"InputName"`
	InputTime              string `xml:"InputTime"`
	RltInvtNo              string `xml:"RltInvtNo"`
	RltPutrecNo            string `xml:"RltPutrecNo"`
	CorrEntryDclEtpsNo     string `xml:"CorrEntryDclEtpsNo"`
	CorrEntryDclEtpsSccd   string `xml:"CorrEntryDclEtpsSccd"`
	CorrEntryDclEtpsNm     string `xml:"CorrEntryDclEtpsNm"`
	RltEntryBizopEtpsno    string `xml:"RltEntryBizopEtpsno"`
	RltEntryBizopEtpsSccd  string `xml:"RltEntryBizopEtpsSccd"`
	RltEntryBizopEtpsNm    string `xml:"RltEntryBizopEtpsNm"`
	RltEntryRcvgdEtpsno    string `xml:"RltEntryRcvgdEtpsno"`
	RltEntryRvsngdEtpsSccd string `xml:"RltEntryRvsngdEtpsSccd"`
	RltEntryRcvgdEtpsNm    string `xml:"RltEntryRcvgdEtpsNm"`
	RltEntryDclEtpsno      string `xml:"RltEntryDclEtpsno"`
	RltEntryDclEtpsSccd    string `xml:"RltEntryDclEtpsSccd"`
	RltEntryDclEtpsNm      string `xml:"RltEntryDclEtpsNm"`
	ImpexpPortcd           string `xml:"ImpexpPortcd"`
	DclPlcCuscd            string `xml:"DclPlcCuscd"`
	ImpexpMarkcd           string `xml:"ImpexpMarkcd"`
	MtpckEndprdMarkcd      string `xml:"MtpckEndprdMarkcd"`
	SupvModecd             string `xml:"SupvModecd"`
	TrspModecd             string `xml:"TrspModecd"`
	ApplyNo                string `xml:"ApplyNo"`
	ListType               string `xml:"ListType"`
	DclcusFlag             string `xml:"DclcusFlag"`
	DclcusTypecd           string `xml:"DclcusTypecd"`
	IcCardNo               string `xml:"IcCardNo"`
	DecType                string `xml:"DecType"`
	Rmk                    Cdata  `xml:"Rmk"`
	StshipTrsarvNatcd      string `xml:"StshipTrsarvNatcd"`
	EntryNo                string `xml:"EntryNo"`
	RltEntryNo             string `xml:"RltEntryNo"`
	DclTypecd              string `xml:"DclTypecd"`
	GenDecFlag             string `xml:"GenDecFlag"`
}

type InvtListType struct {
	SeqNo            string `xml:"SeqNo"`
	GdsSeqno         string `xml:"GdsSeqno"`
	PutrecSeqno      string `xml:"PutrecSeqno"`
	GdsMtno          string `xml:"GdsMtno"`
	Gdecd            string `xml:"Gdecd"`
	GdsNm            Cdata  `xml:"GdsNm"`
	GdsSpcfModelDesc string `xml:"GdsSpcfModelDesc"`
	DclUnitcd        string `xml:"DclUnitcd"`
	LawfUnitcd       string `xml:"LawfUnitcd"`
	SecdLawfUnitcd   string `xml:"SecdLawfUnitcd"`
	Natcd            string `xml:"Natcd"`
	DclUprcAmt       string `xml:"DclUprcAmt"`
	DclTotalAmt      string `xml:"DclTotalAmt"`
	UsdStatTotalAmt  string `xml:"UsdStatTotalAmt"`
	DclCurrcd        string `xml:"DclCurrcd"`
	LawfQty          string `xml:"LawfQty"`
	SecdLawfQty      string `xml:"SecdLawfQty"`
	WtSfVal          string `xml:"WtSfVal"`
	FstSfVal         string `xml:"FstSfVal"`
	SecdSfVal        string `xml:"SecdSfVal"`
	DclQty           string `xml:"DclQty"`
	GrossWt          string `xml:"GrossWt"`
	NetWt            string `xml:"NetWt"`
	UseCd            string `xml:"UseCd"`
	LvyrlfModecd     string `xml:"LvyrlfModecd"`
	UcnsVerno        string `xml:"UcnsVerno"`
	ClyMarkcd        string `xml:"ClyMarkcd"`
	EntryGdsSeqno    string `xml:"EntryGdsSeqno"`
	ApplyTbSeqno     string `xml:"ApplyTbSeqno"`
	DestinationNatcd string `xml:"DestinationNatcd"`
	ModfMarkcd       string `xml:"ModfMarkcd"`
	Rmk              Cdata  `xml:"Rmk"`
}

type InvtDecHeadType struct {
	SeqNo                  string `xml:"SeqNo"`
	DecSeqNo               string `xml:"DecSeqNo"`
	PutrecNo               string `xml:"PutrecNo"`
	BizopEtpsSccd          string `xml:"BizopEtpsSccd"`
	BizopEtpsno            string `xml:"BizopEtpsno"`
	BizopEtpsNm            string `xml:"BizopEtpsNm"`
	RvsngdEtpsSccd         string `xml:"RvsngdEtpsSccd"`
	RcvgdEtpsno            string `xml:"RcvgdEtpsno"`
	RcvgdEtpsNm            string `xml:"RcvgdEtpsNm"`
	DclEtpsSccd            string `xml:"DclEtpsSccd"`
	DclEtpsno              string `xml:"DclEtpsno"`
	DclEtpsNm              string `xml:"DclEtpsNm"`
	InputCode              string `xml:"InputCode"`
	InputCreditCode        string `xml:"InputCreditCode"`
	InputName              string `xml:"InputName"`
	ImpexpPortcd           string `xml:"ImpexpPortcd"`
	DclPlcCuscd            string `xml:"DclPlcCuscd"`
	ImpexpMarkcd           string `xml:"ImpexpMarkcd"`
	SupvModecd             string `xml:"SupvModecd"`
	TrspModecd             string `xml:"TrspModecd"`
	TradeCountry           string `xml:"TradeCountry"`
	DecType                string `xml:"DecType"`
	Rmk                    Cdata  `xml:"Rmk"`
	CreateFlag             string `xml:"CreateFlag"`
	BillNo                 string `xml:"BillNo"`
	ContrNo                string `xml:"ContrNo"`
	CutMode                string `xml:"CutMode"`
	DistinatePort          string `xml:"DistinatePort"`
	FeeCurr                string `xml:"FeeCurr"`
	FeeMark                string `xml:"FeeMark"`
	FeeRate                string `xml:"FeeRate"`
	GrossWet               string `xml:"GrossWet"`
	InsurCurr              string `xml:"InsurCurr"`
	InsurMark              string `xml:"InsurMark"`
	InsurRate              string `xml:"InsurRate"`
	LicenseNo              string `xml:"LicenseNo"`
	NetWt                  string `xml:"NetWt"`
	OtherCurr              string `xml:"OtherCurr"`
	OtherMark              string `xml:"OtherMark"`
	OtherRate              string `xml:"OtherRate"`
	PackNo                 string `xml:"PackNo"`
	TrafName               string `xml:"TrafName"`
	TransMode              string `xml:"TransMode"`
	Type                   string `xml:"Type"`
	WrapType               string `xml:"WrapType"`
	PromiseItems           string `xml:"PromiseItems"`
	TradeAreaCode          string `xml:"TradeAreaCode"`
	DespPortCode           string `xml:"DespPortCode"`
	EntryPortCode          string `xml:"EntryPortCode"`
	GoodsPlace             string `xml:"GoodsPlace"`
	OverseasConsignorCode  string `xml:"OverseasConsignorCode"`
	OverseasConsignorCname string `xml:"OverseasConsignorCname"`
	OverseasConsignorEname string `xml:"OverseasConsignorEname"`
	OverseasConsignorAddr  string `xml:"OverseasConsignorAddr"`
	OverseasConsigneeCode  string `xml:"OverseasConsigneeCode"`
	OverseasConsigneeEname string `xml:"OverseasConsigneeEname"`
}

type InvtDecListType struct {
	SeqNo            string `xml:"SeqNo"`
	DecSeqNo         string `xml:"DecSeqNo"`
	EntryGdsSeqno    string `xml:"EntryGdsSeqno"`
	PutrecSeqno      string `xml:"PutrecSeqno"`
	Gdecd            string `xml:"Gdecd"`
	GdsNm            string `xml:"GdsNm"`
	GdsSpcfModelDesc Cdata  `xml:"GdsSpcfModelDesc"`
	DclUnitcd        string `xml:"DclUnitcd"`
	LawfUnitcd       string `xml:"LawfUnitcd"`
	SecdLawfUnitcd   string `xml:"SecdLawfUnitcd"`
	DclUprcAmt       string `xml:"DclUprcAmt"`
	DclTotalAmt      string `xml:"DclTotalAmt"`
	DclCurrCd        string `xml:"DclCurrCd"`
	NatCd            string `xml:"NatCd"`
	DestinationNatcd string `xml:"DestinationNatcd"`
	LawfQty          string `xml:"LawfQty"`
	SecdLawfQty      string `xml:"SecdLawfQty"`
	DclQty           string `xml:"DclQty"`
	UseCd            string `xml:"UseCd"`
	Rmk              Cdata  `xml:"Rmk"`
	LvyrlfModecd     string `xml:"LvyrlfModecd"`
	CiqCode          string `xml:"CiqCode"`
	DeclGoodsEname   string `xml:"DeclGoodsEname"`
	OrigPlaceCode    string `xml:"OrigPlaceCode"`
	Purpose          string `xml:"Purpose"`
	ProdValidDt      string `xml:"ProdValidDt"`
	ProdQgp          string `xml:"ProdQgp"`
	GoodsAttr        string `xml:"GoodsAttr"`
	Stuff            string `xml:"Stuff"`
	UnCode           string `xml:"UnCode"`
	DangName         string `xml:"DangName"`
	DangPackType     string `xml:"DangPackType"`
	DangPackSpec     string `xml:"DangPackSpec"`
	EngManEntCnm     string `xml:"EngManEntCnm"`
	NoDangFlag       string `xml:"NoDangFlag"`
	DestCode         string `xml:"DestCode"`
	GoodsSpec        string `xml:"GoodsSpec"`
	GoodsModel       string `xml:"GoodsModel"`
	GoodsBrand       string `xml:"GoodsBrand"`
	ProduceDate      string `xml:"ProduceDate"`
	ProdBatchNo      string `xml:"ProdBatchNo"`
	DistrictCode     string `xml:"DistrictCode"`
	CiqName          string `xml:"CiqName"`
	MnufctrRegno     string `xml:"MnufctrRegno"`
	MnufctrRegName   string `xml:"MnufctrRegName"`
}

type InvtMessage struct {
	InvtHeadType    InvtHeadType      `xml:"InvtHeadType"`
	InvtListType    []InvtListType    `xml:"InvtListType"`
	InvtDecHeadType InvtDecHeadType   `xml:"InvtDecHeadType"`
	InvtDecListType []InvtDecListType `xml:"InvtDecListType"`
	OperCusRegCode  string            `xml:"OperCusRegCode"`
	SysId           string            `xml:"SysId"`
}

type BussinessData struct {
	InvtMessage InvtMessage `xml:"InvtMessage"`
	DelcareFlag string      `xml:"DelcareFlag"`
}

type DataInfo struct {
	PocketInfo    PocketInfo    `xml:"PocketInfo"`
	BussinessData BussinessData `xml:"BussinessData"`
}

type Package struct {
	EnvelopInfo EnvelopInfo `xml:"EnvelopInfo"`
	DataInfo    DataInfo    `xml:"DataInfo"`
}

type Object struct {
	Package Package `xml:"Package"`
}

type Signature struct {
	XMLName        xml.Name   `xml:"Signature"`      //该XML文件的根元素为 Signature
	Xmlns          string     `xml:"xmlns:xsi,attr"` //该值会作为Signature元素的属性
	SignedInfo     SignedInfo `xml:"SignedInfo"`
	SignatureValue string     `xml:"SignatureValue"`
	KeyInfo        string     `xml:"KeyInfo>KeyName"`
	Object         Object     `xml:"Object"`
}
