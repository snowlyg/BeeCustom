package mysoap

import (
	"encoding/xml"
	"fmt"
	"time"

	"BeeCustom/utils"
	"github.com/hooklift/gowsdl/soap"
)

// against "unused imports"
var _ time.Time
var _ xml.Name

type Authentication struct {
	XMLName xml.Name `xml:"Authentication"`

	MessageType string `xml:"messageType,omitempty"`
}
type InBoundsRequest struct {
	XMLName xml.Name `xml:"http://webservice.bgcd.hwt.com in_bounds"`

	MessageType string `xml:"messageType,omitempty"`
}

type InBoundsResponse struct {
	XMLName xml.Name `xml:"http://webservice.bgcd.hwt.com in_boundsResponse"`

	Return_ string `xml:"return,omitempty"`
}

type FormateDate struct {
	XMLName xml.Name `xml:"http://webservice.bgcd.hwt.com formateDate"`

	Format string `xml:"format,omitempty"`
}

type FormateDateResponse struct {
	XMLName xml.Name `xml:"http://webservice.bgcd.hwt.com formateDateResponse"`

	Return_ string `xml:"return,omitempty"`
}

type InBoundsServicePortType interface {
	FormateDate(request *FormateDate) (*FormateDateResponse, error)

	InBounds(request *InBoundsRequest) (*InBoundsResponse, error)
}

type inBoundsServicePortType struct {
	client *soap.Client
}

func NewInBoundsServicePortType(client *soap.Client) InBoundsServicePortType {
	return &inBoundsServicePortType{
		client: client,
	}
}

func (service *inBoundsServicePortType) FormateDate(request *FormateDate) (*FormateDateResponse, error) {
	response := new(FormateDateResponse)
	err := service.client.Call("urn:formateDate", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *inBoundsServicePortType) InBounds(request *InBoundsRequest) (*InBoundsResponse, error) {
	response := new(InBoundsResponse)
	err := service.client.Call("urn:in_bounds", request, response)
	if err != nil {
		utils.LogDebug(fmt.Sprintf(" service.client.Call  error: %v", err))
		return nil, err
	}

	return response, nil
}
