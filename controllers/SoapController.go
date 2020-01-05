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
	header := mysoap.Authentication{Username: "DHBG-IT", Password: "88888888"}
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
