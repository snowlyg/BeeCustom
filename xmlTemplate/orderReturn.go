package xmlTemplate

import "encoding/xml"

// successed / failed：没有 SeqNo
type DecImportResponse struct {
	DecImportResponse xml.Name `xml:"DecImportResponse"`

	ResponseCode string `xml:"ResponseCode"`
	ErrorMessage string `xml:"ErrorMessage"`
	ClientSeqNo  string `xml:"ClientSeqNo"`
	SeqNo        string `xml:"SeqNo"`
}

// receipt
type DecResult struct {
	DecResult xml.Name `xml:"DEC_RESULT"`

	CusCiqNo     string `xml:"CUS_CIQ_NO"`
	EntryId      string `xml:"ENTRY_ID"`
	NoticeDate   string `xml:"NOTICE_DATE"`
	Channel      string `xml:"CHANNEL"`
	Note         string `xml:"NOTE"`
	CustomMaster string `xml:"CUSTOM_MASTER"`
	IEDate       string `xml:"I_E_DATE"`
	DDate        string `xml:"D_DATE"`
}
