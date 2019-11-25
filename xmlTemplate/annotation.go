package xmlTemplate

/**清单报文结构*/

import "encoding/xml"

type DigestMethod struct {
	//DigestMethod xml.Name `xml:"DigestMethod"`
	Algorithm string `xml:"Algorithm,attr"`
}

type Reference struct {
	//Reference    xml.Name     `xml:"Reference"`
	URI          string       `xml:"URI,attr"`
	DigestMethod DigestMethod `xml:"DigestMethod"`
	DigestValue  string       `xml:"DigestValue"`
}

type CanonicalizationMethod struct {
	//CanonicalizationMethod xml.Name `xml:"CanonicalizationMethod"`
	Algorithm string `xml:"Algorithm,attr"`
}

type SignatureMethod struct {
	//SignatureMethod xml.Name `xml:"SignatureMethod"`
	Algorithm string `xml:"Algorithm,attr"`
}

type SignedInfo struct {
	//SignedInfo             xml.Name               `xml:"SignedInfo"`
	CanonicalizationMethod CanonicalizationMethod `xml:"CanonicalizationMethod"`
	SignatureMethod        SignatureMethod        `xml:"SignatureMethod"`
	Reference              Reference              `xml:"Reference"`
}

type EnvelopInfo struct {
	//EnvelopInfo xml.Name `xml:"EnvelopInfo"`
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
	//PocketInfo     xml.Name `xml:"PocketInfo"`
	PocketId       string `xml:"pocket_id"`
	TotalPocketQty string `xml:"total_pocket_qty"`
	CurPocketNo    string `xml:"cur_pocket_no"`
	IsUnstructured string `xml:"is_unstructured"`
}

type InvtHeadType struct {
	//InvtHeadType                 xml.Name `xml:"InvtHeadType"`
	SeqNo                        string `xml:"SeqNo"`
	BondInvtNo                   string `xml:"BondInvtNo"`
	ChgTmsCntstring              string `xml:"ChgTmsCnt"`
	PutrecNostring               string `xml:"PutrecNo"`
	InvtTypestring               string `xml:"InvtType"`
	EtpsInnerInvtNostring        string `xml:"EtpsInnerInvtNo"`
	BizopEtpsnostring            string `xml:"BizopEtpsno"`
	BizopEtpsSccdstring          string `xml:"BizopEtpsSccd"`
	BizopEtpsNmstring            string `xml:"BizopEtpsNm"`
	RcvgdEtpsnostring            string `xml:"RcvgdEtpsno"`
	RvsngdEtpsSccdstring         string `xml:"RvsngdEtpsSccd"`
	RcvgdEtpsNmstring            string `xml:"RcvgdEtpsNm"`
	DclEtpsnostring              string `xml:"DclEtpsno"`
	DclEtpsSccdstring            string `xml:"DclEtpsSccd"`
	DclEtpsNmstring              string `xml:"DclEtpsNm"`
	InputCodestring              string `xml:"InputCode"`
	InputCreditCodestring        string `xml:"InputCreditCode"`
	InputNamestring              string `xml:"InputName"`
	InputTimestring              string `xml:"InputTime"`
	RltInvtNostring              string `xml:"RltInvtNo"`
	RltPutrecNostring            string `xml:"RltPutrecNo"`
	CorrEntryDclEtpsNostring     string `xml:"CorrEntryDclEtpsNo"`
	CorrEntryDclEtpsSccdstring   string `xml:"CorrEntryDclEtpsSccd"`
	CorrEntryDclEtpsNmstring     string `xml:"CorrEntryDclEtpsNm"`
	RltEntryBizopEtpsnostring    string `xml:"RltEntryBizopEtpsno"`
	RltEntryBizopEtpsSccdstring  string `xml:"RltEntryBizopEtpsSccd"`
	RltEntryBizopEtpsNmstring    string `xml:"RltEntryBizopEtpsNm"`
	RltEntryRcvgdEtpsnostring    string `xml:"RltEntryRcvgdEtpsno"`
	RltEntryRvsngdEtpsSccdstring string `xml:"RltEntryRvsngdEtpsSccd"`
	RltEntryRcvgdEtpsNmstring    string `xml:"RltEntryRcvgdEtpsNm"`
	RltEntryDclEtpsnostring      string `xml:"RltEntryDclEtpsno"`
	RltEntryDclEtpsSccdstring    string `xml:"RltEntryDclEtpsSccd"`
	RltEntryDclEtpsNmstring      string `xml:"RltEntryDclEtpsNm"`
	ImpexpPortcdstring           string `xml:"ImpexpPortcd"`
	DclPlcCuscdstring            string `xml:"DclPlcCuscd"`
	ImpexpMarkcdstring           string `xml:"ImpexpMarkcd"`
	MtpckEndprdMarkcdstring      string `xml:"MtpckEndprdMarkcd"`
	SupvModecdstring             string `xml:"SupvModecd"`
	TrspModecdstring             string `xml:"TrspModecd"`
	ApplyNostring                string `xml:"ApplyNo"`
	ListTypestring               string `xml:"ListType"`
	DclcusFlagstring             string `xml:"DclcusFlag"`
	DclcusTypecdstring           string `xml:"DclcusTypecd"`
	IcCardNostring               string `xml:"IcCardNo"`
	DecTypestring                string `xml:"DecType"`
	Rmkstring                    string `xml:"Rmk"`
	StshipTrsarvNatcdstring      string `xml:"StshipTrsarvNatcd"`
	EntryNostring                string `xml:"EntryNo"`
	RltEntryNostring             string `xml:"RltEntryNo"`
	DclTypecdstring              string `xml:"DclTypecd"`
	GenDecFlagstring             string `xml:"GenDecFlag"`
}

type InvtListType struct {
	//InvtListType     xml.Name `xml:"InvtListType"`
	SeqNo            string `xml:"SeqNo"`
	GdsSeqno         string `xml:"GdsSeqno"`
	PutrecSeqno      string `xml:"PutrecSeqno"`
	GdsMtno          string `xml:"GdsMtno"`
	Gdecd            string `xml:"Gdecd"`
	GdsNm            string `xml:"GdsNm"`
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
	Rmk              string `xml:"Rmk"`
}

type InvtDecHeadType struct {
	//InvtDecHeadType        xml.Name `xml:"InvtDecHeadType"`
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
	Rmk                    string `xml:"Rmk"`
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
	//InvtDecListType        xml.Name `xml:"InvtDecListType"`
	SeqNostring            string `xml:"SeqNo"`
	DecSeqNostring         string `xml:"DecSeqNo"`
	EntryGdsSeqnostring    string `xml:"EntryGdsSeqno"`
	PutrecSeqnostring      string `xml:"PutrecSeqno"`
	Gdecdstring            string `xml:"Gdecd"`
	GdsNmstring            string `xml:"GdsNm"`
	GdsSpcfModelDescstring string `xml:"GdsSpcfModelDesc"`
	DclUnitcdstring        string `xml:"DclUnitcd"`
	LawfUnitcdstring       string `xml:"LawfUnitcd"`
	SecdLawfUnitcdstring   string `xml:"SecdLawfUnitcd"`
	DclUprcAmtstring       string `xml:"DclUprcAmt"`
	DclTotalAmtstring      string `xml:"DclTotalAmt"`
	DclCurrCdstring        string `xml:"DclCurrCd"`
	NatCdstring            string `xml:"NatCd"`
	DestinationNatcdstring string `xml:"DestinationNatcd"`
	LawfQtystring          string `xml:"LawfQty"`
	SecdLawfQtystring      string `xml:"SecdLawfQty"`
	DclQtystring           string `xml:"DclQty"`
	UseCdstring            string `xml:"UseCd"`
	Rmkstring              string `xml:"Rmk"`
	LvyrlfModecdstring     string `xml:"LvyrlfModecd"`
	CiqCodestring          string `xml:"CiqCode"`
	DeclGoodsEnamestring   string `xml:"DeclGoodsEname"`
	OrigPlaceCodestring    string `xml:"OrigPlaceCode"`
	Purposestring          string `xml:"Purpose"`
	ProdValidDtstring      string `xml:"ProdValidDt"`
	ProdQgpstring          string `xml:"ProdQgp"`
	GoodsAttrstring        string `xml:"GoodsAttr"`
	Stuffstring            string `xml:"Stuff"`
	UnCodestring           string `xml:"UnCode"`
	DangNamestring         string `xml:"DangName"`
	DangPackTypestring     string `xml:"DangPackType"`
	DangPackSpecstring     string `xml:"DangPackSpec"`
	EngManEntCnmstring     string `xml:"EngManEntCnm"`
	NoDangFlagstring       string `xml:"NoDangFlag"`
	DestCodestring         string `xml:"DestCode"`
	GoodsSpecstring        string `xml:"GoodsSpec"`
	GoodsModelstring       string `xml:"GoodsModel"`
	GoodsBrandstring       string `xml:"GoodsBrand"`
	ProduceDatestring      string `xml:"ProduceDate"`
	ProdBatchNostring      string `xml:"ProdBatchNo"`
	DistrictCodestring     string `xml:"DistrictCode"`
	CiqNamestring          string `xml:"CiqName"`
	MnufctrRegnostring     string `xml:"MnufctrRegno"`
	MnufctrRegNamestring   string `xml:"MnufctrRegName"`
}

type InvtMessage struct {
	//InvtMessage     xml.Name          `xml:"InvtMessage"`
	InvtHeadType    InvtHeadType      `xml:"InvtHeadType"`
	InvtListType    []InvtListType    `xml:"InvtListType"`
	InvtDecHeadType InvtDecHeadType   `xml:"InvtDecHeadType"`
	InvtDecListType []InvtDecListType `xml:"InvtDecListType"`
	OperCusRegCode  string            `xml:"OperCusRegCode"`
	SysId           string            `xml:"SysId"`
}

type BussinessData struct {
	//BussinessData xml.Name    `xml:"BussinessData"`
	InvtMessage InvtMessage `xml:"InvtMessage"`
	DelcareFlag string      `xml:"DelcareFlag"`
}

type DataInfo struct {
	//DataInfo      xml.Name      `xml:"DataInfo"`
	PocketInfo    PocketInfo    `xml:"PocketInfo"`
	BussinessData BussinessData `xml:"BussinessData"`
}

type Package struct {
	//Package     xml.Name    `xml:"Package"`
	EnvelopInfo EnvelopInfo `xml:"EnvelopInfo"`
	DataInfo    DataInfo    `xml:"DataInfo"`
}

type Object struct {
	//Object  xml.Name `xml:"Object"`
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
