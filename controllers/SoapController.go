package controllers

import (
	"bytes"
	"encoding/xml"
	"time"

	"BeeCustom/enums"
	"BeeCustom/file"
	"BeeCustom/mysoap"
	"github.com/hooklift/gowsdl/soap"
)

type SoapController struct {
	BaseController
}

//获取回执、船期、特殊费用
func (c *SoapController) ACKMsg() {
	client := soap.NewClient("http://www.cusdectrans.com:8014/BGCDWebService/services/InBoundsService?wsdl")
	header := mysoap.Authentication{Username: "DHBG-IT", Password: "88888888"}
	client.AddHeader(header)
	m := mysoap.NewInBoundsServicePortType(client)
	r, err := m.InBounds(&mysoap.InBoundsRequest{MessageType: "SDATE"})
	if err != nil {
		c.Data["json"] = err
	} else {
		c.Data["json"] = r
	}

	c.ServeJSON()
}

//获取在场箱信息
func (c *SoapController) OnYard() {
	client := soap.NewClient("http://www.cusdectrans.com:8014/BGCDWebService/services/QueryContaStatus?wsdl")
	header := mysoap.Authentication{Username: "DHBG-IT", Password: "88888888"}
	client.AddHeader(header)
	m := mysoap.NewQueryContaStatusPortType(client)
	r, err := m.OnYard_conta(&mysoap.OnYard_conta{Conta_no: "HLBU2358636", Dock_code: "CNYTN"})
	if err != nil {
		c.Data["json"] = err
	} else {
		c.Data["json"] = r
	}

	c.ServeJSON()
}

//发送报文
func (c *SoapController) SendDt() {
	client := soap.NewClient("http://www.cusdectrans.com:8014/BGCDWebService/services/OutBoundsService?wsdl")
	header := mysoap.Authentication{Username: "DHBG-IT", Password: "88888888"}
	client.AddHeader(header)
	m := mysoap.NewOutBoundsServicePortType(client)
	xmlStr := c.getXmlStr()
	r, err := m.Out_bounds(&mysoap.Out_bounds{XmlStr: string(xmlStr)})
	if err != nil {
		c.Data["json"] = err
	} else {
		c.Data["json"] = r
	}

	c.ServeJSON()
}

func (c *SoapController) getXmlStr() []byte {
	var manifest mysoap.Manifest

	declaration := mysoap.Declaration{}
	declaration.OptType = "00"
	declaration.TradeMode = "0615"
	declaration.DeclType = "Q"

	person := mysoap.RepresentativePerson{}
	person.Name = "5304192209991"
	declaration.RepresentativePerson = person

	office := mysoap.ExitCustomsOffice{}
	office.ID = "CNYTN/5316"
	declaration.ExitCustomsOffice = office

	carrier := mysoap.Carrier{}
	carrier.ID = "MSC"
	declaration.Carrier = carrier

	means := mysoap.BorderTransportMeans{}
	means.JourneyID = "FT939W"
	means.ID = "UN9469560"
	means.Name = "MSC TERESA"
	declaration.BorderTransportMeans = means

	consignment := mysoap.Consignment{}
	consignment.CustomsStatusCode = "001"
	consignment.TotalGrossMassMeasure = "19770"

	document := mysoap.TransportContractDocument{}
	document.ConditionCode = "10"
	consignment.TransportContractDocument = document

	location := mysoap.LoadingLocation{}
	location.ID = "CNYTN"
	location.LoadingDate = "201910042000"
	consignment.LoadingLocation = location

	unloadingLocation := mysoap.UnloadingLocation{}
	unloadingLocation.ID = "TRIST"
	consignment.UnloadingLocation = unloadingLocation

	payment := mysoap.FreightPayment{}
	payment.MethodCode = "CC"
	consignment.FreightPayment = payment

	packaging := mysoap.ConsignmentPackaging{}
	packaging.QuantityQuantity = "1440"
	packaging.TypeCode = "CT"
	consignment.ConsignmentPackaging = packaging

	consignee := mysoap.Consignee{}
	consignee.ID = "MERSIS NUMBER+0461039586100014"
	consignee.Name = "HEDEF TUKETIM URUNLERI SANAYI VE DIS TICARET A.S."

	address := mysoap.Address{}
	address.CountryCode = "TR"
	address.Line = "TERRA PLAZA, NO:25 KAT:3 D:9-10"
	consignee.Address = address

	communication := mysoap.Communication{}
	communication.ID = "216 6390115"
	communication.TypeID = "TE"
	consignee.Communication = communication
	consignment.Consignee = consignee

	consignor := mysoap.Consignor{}
	consignor.ID = "USCI+91441900730442249J"
	consignor.Name = "广东力王新能源股份有限公司 TEST"

	addressConsignor := mysoap.Address{}
	addressConsignor.CountryCode = "CN"
	addressConsignor.Line = "广东省东莞市塘厦镇石马管理区 TEST"
	consignor.Address = addressConsignor

	communicationConsignor := mysoap.Communication{}
	communicationConsignor.ID = "0769-87885755"
	communicationConsignor.TypeID = "TE"
	consignor.Communication = communicationConsignor
	consignment.Consignor = consignor

	contact := mysoap.UNDGContact{}
	contact.Name = "IT部测试"
	i := mysoap.Communication{}
	i.ID = "13565654852"
	i.ID = "TE"
	contact.Communication = i
	consignment.UNDGContact = contact

	var equipments []mysoap.TransportEquipment
	equipment := mysoap.TransportEquipment{}
	equipment.BookingNumber = "181AS0193642072M1"
	equipment.LclNum = "2"
	equipment.IsLcl = "1"
	equipment.SealID = "M/CNB549227"
	equipment.FullnessCode = "7"
	equipment.SupplierPartyTypeCode = "2"
	equipment.CharacteristicCode = "20GP"

	identification := mysoap.EquipmentIdentification{}
	identification.ID = "FCIU2558917"
	equipment.EquipmentIdentification = identification

	equipments = append(equipments, equipment)
	consignment.TransportEquipment = equipments

	var items []mysoap.ConsignmentItem
	item := mysoap.ConsignmentItem{}
	item.SequenceNumeric = "1"

	itemPackaging := mysoap.ConsignmentItemPackaging{}
	itemPackaging.QuantityQuantity = "1440"
	itemPackaging.TypeCode = "CT"
	itemPackaging.MarksNumbers = "N/M"
	item.ConsignmentItemPackaging = itemPackaging

	commodity := mysoap.Commodity{}
	commodity.CargoDescription = "碱性锌锰电池"
	commodity.UNDGCode = "2343"
	item.Commodity = commodity

	measure := mysoap.GoodsMeasure{}
	measure.GrossMassMeasure = "19770"
	item.GoodsMeasure = measure

	equipmentIdentification := mysoap.EquipmentIdentification{}
	equipmentIdentification.ID = "FCIU2558917"
	item.EquipmentIdentification = equipmentIdentification

	items = append(items, item)
	consignment.ConsignmentItem = items

	declaration.Consignment = consignment

	head := mysoap.Head{}
	head.MessageID = "DHBG-IT_19092708561404707"
	head.FunctionCode = "9"
	head.MessageType = "MT2101A"
	head.SenderID = "DHBG-IT"
	head.ReceiverID = "DT"
	head.SendTime = time.Now().Format(enums.BaseDateTimeSecondFormat)
	head.SendTime = "1.0"

	manifest.Declaration = declaration
	manifest.Head = head

	output, err := xml.MarshalIndent(manifest, "", "")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "xml文件解析出错", err)
	}

	if err := file.CreateFile("./"); err != nil {
		c.jsonResult(enums.JRCodeFailed, "文件创建出错", err)
	}

	bs := [][]byte{[]byte(xml.Header), output}
	moutput := bytes.Join(bs, []byte(""))

	return moutput
}
