package xmlTemplate

import "encoding/xml"

type CommonResponeMessage struct {
	CommonResponeMessage xml.Name `xml:"CommonResponeMessage"`

	SeqNo        string `xml:"SeqNo"`
	EtpsPreentNo string `xml:"EtpsPreentNo"`
	CheckInfo    string `xml:"CheckInfo"`
	DealFlag     string `xml:"DealFlag"`
}

type ReturnPackage struct {
	Package xml.Name `xml:"Package"`

	Version               string `xml:"EnvelopInfo>version"`
	EnvelopInfoBusinessId string `xml:"EnvelopInfo>business_id"`
	MessageId             string `xml:"EnvelopInfo>message_id"`
	FileName              string `xml:"EnvelopInfo>file_name"`
	MessageType           string `xml:"EnvelopInfo>message_type"`
	SenderId              string `xml:"EnvelopInfo>sender_id"`
	ReceiverId            string `xml:"EnvelopInfo>receiver_id"`
	SendTime              string `xml:"EnvelopInfo>send_time"`

	InvPreentNo        string `xml:"DataInfo>BussinessData>INV202>InvApprResult>invPreentNo"`
	DataInfoBusinessId string `xml:"DataInfo>BussinessData>INV202>InvApprResult>businessId"`
	EntrySeqNo         string `xml:"DataInfo>BussinessData>INV202>InvApprResult>entrySeqNo"`
	ManageResult       string `xml:"DataInfo>BussinessData>INV202>InvApprResult>manageResult"`
	CreateDate         string `xml:"DataInfo>BussinessData>INV202>InvApprResult>createDate"`
	Reason             string `xml:"DataInfo>BussinessData>INV202>InvApprResult>reason"`
}
