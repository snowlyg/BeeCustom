package mysoap

import (
	"encoding/xml"
	"github.com/hooklift/gowsdl/soap"
	"time"
)

// against "unused imports"
var _ time.Time
var _ xml.Name

type Out_bounds struct {
	XMLName xml.Name `xml:"http://webservice.bgcd.hwt.com out_bounds"`

	XmlStr string `xml:"xmlStr,omitempty"`
}

type Out_boundsResponse struct {
	XMLName xml.Name `xml:"http://webservice.bgcd.hwt.com out_boundsResponse"`

	Return_ string `xml:"return,omitempty"`
}

type OutBoundsServicePortType interface {
	Out_bounds(request *Out_bounds) (*Out_boundsResponse, error)
}

type outBoundsServicePortType struct {
	client *soap.Client
}

func NewOutBoundsServicePortType(client *soap.Client) OutBoundsServicePortType {
	return &outBoundsServicePortType{
		client: client,
	}
}

func (service *outBoundsServicePortType) Out_bounds(request *Out_bounds) (*Out_boundsResponse, error) {
	response := new(Out_boundsResponse)
	err := service.client.Call("urn:out_bounds", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
