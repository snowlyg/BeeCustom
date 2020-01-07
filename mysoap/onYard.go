package mysoap

import (
	"encoding/xml"
	"time"

	"github.com/hooklift/gowsdl/soap"
)

// against "unused imports"
var _ time.Time
var _ xml.Name

type OnYard_conta struct {
	XMLName xml.Name `xml:"http://webservice.bgcd.hwt.com onYard_conta"`

	Conta_no string `xml:"conta_no,omitempty"`

	Dock_code string `xml:"dock_code,omitempty"`
}

type OnYard_contaResponse struct {
	XMLName xml.Name `xml:"http://webservice.bgcd.hwt.com onYard_contaResponse"`

	Return_ string `xml:"return,omitempty"`
}

type Main struct {
	XMLName xml.Name `xml:"http://webservice.bgcd.hwt.com main"`

	Args []string `xml:"args,omitempty"`
}

type GetDCWShipInfo struct {
	XMLName xml.Name `xml:"http://webservice.bgcd.hwt.com getDCWShipInfo"`

	Conta_no string `xml:"conta_no,omitempty"`
}

type GetDCWShipInfoResponse struct {
	XMLName xml.Name `xml:"http://webservice.bgcd.hwt.com getDCWShipInfoResponse"`

	Return_ *Map1 `xml:"return,omitempty"`
}

type CreateDCXDoc struct {
	XMLName xml.Name `xml:"http://webservice.bgcd.hwt.com createDCXDoc"`

	Map_ *Map2 `xml:"map,omitempty"`
}

type CreateDCXDocResponse struct {
	XMLName xml.Name `xml:"http://webservice.bgcd.hwt.com createDCXDocResponse"`

	Return_ string `xml:"return,omitempty"`
}

type Entry1 struct {
	XMLName xml.Name `xml:"http://ws.apache.org/namespaces/axis2/map entry1"`

	Key string `xml:"key,omitempty"`

	Value string `xml:"value,omitempty"`
}

type Map1 struct {
	XMLName xml.Name `xml:"http://ws.apache.org/namespaces/axis2/map map1"`

	Entry []*Entry1 `xml:"entry,omitempty"`
}

type Entry2 struct {
	XMLName xml.Name `xml:"http://ws.apache.org/namespaces/axis2/map entry2"`

	Key interface{} `xml:"key,omitempty"`

	Value interface{} `xml:"value,omitempty"`
}

type Map2 struct {
	XMLName xml.Name `xml:"http://ws.apache.org/namespaces/axis2/map map2"`

	Entry []*Entry2 `xml:"entry,omitempty"`
}

type QueryContaStatusPortType interface {
	GetDCWShipInfo(request *GetDCWShipInfo) (*GetDCWShipInfoResponse, error)

	CreateDCXDoc(request *CreateDCXDoc) (*CreateDCXDocResponse, error)

	OnYard_conta(request *OnYard_conta) (*OnYard_contaResponse, error)
}

type queryContaStatusPortType struct {
	client *soap.Client
}

func NewQueryContaStatusPortType(client *soap.Client) QueryContaStatusPortType {
	return &queryContaStatusPortType{
		client: client,
	}
}

func (service *queryContaStatusPortType) GetDCWShipInfo(request *GetDCWShipInfo) (*GetDCWShipInfoResponse, error) {
	response := new(GetDCWShipInfoResponse)
	err := service.client.Call("urn:getDCWShipInfo", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *queryContaStatusPortType) CreateDCXDoc(request *CreateDCXDoc) (*CreateDCXDocResponse, error) {
	response := new(CreateDCXDocResponse)
	err := service.client.Call("urn:createDCXDoc", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *queryContaStatusPortType) OnYard_conta(request *OnYard_conta) (*OnYard_contaResponse, error) {
	response := new(OnYard_contaResponse)
	err := service.client.Call("urn:onYard_conta", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
