package mysoap

import (
	"encoding/xml"
)

type Manifest struct {
	XMLName     xml.Name    `xml:"Manifest"`
	Head        Head        `xml:"Head"`
	Declaration Declaration `xml:"Declaration"`
}

type Head struct {
	MessageID    string `xml:"MessageID"`
	FunctionCode string `xml:"FunctionCode"`
	MessageType  string `xml:"MessageType"`
	SenderID     string `xml:"SenderID"`
	ReceiverID   string `xml:"ReceiverID"`
	SendTime     string `xml:"SendTime"`
	Version      string `xml:"Version"`
}

type Declaration struct {
	RepresentativePerson RepresentativePerson `xml:"RepresentativePerson"`
	ExitCustomsOffice    ExitCustomsOffice    `xml:"ExitCustomsOffice"`
	Carrier              Carrier              `xml:"Carrier"`
	OptType              string               `xml:"OptType"`
	TradeMode            string               `xml:"TradeMode"`
	DeclType             string               `xml:"DeclType"`
	BorderTransportMeans BorderTransportMeans `xml:"BorderTransportMeans"`
	Consignment          Consignment          `xml:"Consignment"`
}

type RepresentativePerson struct {
	Name string `xml:"Name"`
}
type ExitCustomsOffice struct {
	ID string `xml:"ID"`
}
type Carrier struct {
	ID string `xml:"ID"`
}
type BorderTransportMeans struct {
	JourneyID string `xml:"JourneyID"`
	ID        string `xml:"ID"`
	Name      string `xml:"Name"`
}
type Consignment struct {
	TransportContractDocument TransportContractDocument `xml:"TransportContractDocument"`
	LoadingLocation           LoadingLocation           `xml:"LoadingLocation"`
	UnloadingLocation         UnloadingLocation         `xml:"UnloadingLocation"`
	CustomsStatusCode         string                    `xml:"CustomsStatusCode"`
	FreightPayment            FreightPayment            `xml:"FreightPayment"`
	ConsignmentPackaging      ConsignmentPackaging      `xml:"ConsignmentPackaging"`
	TotalGrossMassMeasure     string                    `xml:"TotalGrossMassMeasure"`
	Consignee                 Consignee                 `xml:"Consignee"`
	Consignor                 Consignor                 `xml:"Consignor"`
	UNDGContact               UNDGContact               `xml:"UNDGContact"`
	TransportEquipment        []TransportEquipment      `xml:"TransportEquipment"`
	ConsignmentItem           []ConsignmentItem         `xml:"ConsignmentItem"`
}
type TransportContractDocument struct {
	ConditionCode string `xml:"ConditionCode"`
}

type LoadingLocation struct {
	ID          string `xml:"ID"`
	LoadingDate string `xml:"LoadingDate"`
}
type UnloadingLocation struct {
	ID string `xml:"ID"`
}
type FreightPayment struct {
	MethodCode string `xml:"MethodCode"`
}

type ConsignmentPackaging struct {
	QuantityQuantity string `xml:"QuantityQuantity"`
	TypeCode         string `xml:"TypeCode"`
}

type Consignee struct {
	ID            string        `xml:"ID"`
	Name          string        `xml:"Name"`
	Address       Address       `xml:"Address"`
	Communication Communication `xml:"Communication"`
}

type Address struct {
	Line        string `xml:"Line"`
	CountryCode string `xml:"CountryCode"`
}

type Communication struct {
	ID     string `xml:"ID"`
	TypeID string `xml:"TypeID"`
}

type Consignor struct {
	ID            string        `xml:"ID"`
	Name          string        `xml:"Name"`
	Address       Address       `xml:"Address"`
	Communication Communication `xml:"Communication"`
}

type UNDGContact struct {
	Name          string        `xml:"Name"`
	Communication Communication `xml:"Communication"`
}

type TransportEquipment struct {
	EquipmentIdentification EquipmentIdentification `xml:"EquipmentIdentification"`
	CharacteristicCode      string                  `xml:"CharacteristicCode"`
	SupplierPartyTypeCode   string                  `xml:"SupplierPartyTypeCode"`
	FullnessCode            string                  `xml:"FullnessCode"`
	SealID                  string                  `xml:"SealID"`
	IsLcl                   string                  `xml:"IsLcl"`
	LclNum                  string                  `xml:"LclNum"`
	BookingNumber           string                  `xml:"BookingNumber"`
}

type EquipmentIdentification struct {
	ID string `xml:"ID"`
}

type ConsignmentItem struct {
	SequenceNumeric          string                   `xml:"SequenceNumeric"`
	ConsignmentItemPackaging ConsignmentItemPackaging `xml:"ConsignmentItemPackaging"`
	Commodity                Commodity                `xml:"Commodity"`
	GoodsMeasure             GoodsMeasure             `xml:"GoodsMeasure"`
	EquipmentIdentification  EquipmentIdentification  `xml:"EquipmentIdentification"`
}

type ConsignmentItemPackaging struct {
	QuantityQuantity string `xml:"QuantityQuantity"`
	TypeCode         string `xml:"TypeCode"`
	MarksNumbers     string `xml:"MarksNumbers"`
}

type Commodity struct {
	CargoDescription string `xml:"CargoDescription"`
	UNDGCode         string `xml:"UNDGCode"`
}

type GoodsMeasure struct {
	GrossMassMeasure string `xml:"GrossMassMeasure"`
}
