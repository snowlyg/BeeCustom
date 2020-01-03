package controllers

import (
	"BeeCustom/mysoap"
	"github.com/hooklift/gowsdl/soap"
)

type SoapController struct {
	BaseController
}

func (c *SoapController) Soap() {
	client := soap.NewClient("http://www.cusdectrans.com:8014/BGCDWebService/services/InBoundsService?wsdl")
	header := `<tns:Authentication xmlns:tns="authentication">DHBG-IT<tns:Username xmlns:tns="authentication"></tns:Username><tns:Password xmlns:tns="authentication">88888888</tns:Password></tns:Authentication>`
	client.AddHeader(header)
	boundsServicePortType := mysoap.NewInBoundsServicePortType(client)
	inBounds, err := boundsServicePortType.InBounds(&mysoap.InBoundsRequest{MessageType: "SDATE"})
	if err != nil {
		c.Data["json"] = err
	} else {
		c.Data["json"] = inBounds
	}

	c.ServeJSON()
}
