package controllers

import (
	"fmt"

	"BeeCustom/mysoap"
	"BeeCustom/utils"
	"github.com/hooklift/gowsdl/soap"
)

type SoapController struct {
	BaseController
}

func (c *SoapController) Soap() {
	client := soap.NewClient("http://www.cusdectrans.com:8014/BGCDWebService/services/InBoundsService?wsdl")
	header := `<tns:Authentication xmlns:tns="authentication">DHBG<tns:Username xmlns:tns="authentication"></tns:Username>
<tns:Password xmlns:tns="authentication">Bg888888</tns:Password></tns:Authentication>`
	client.AddHeader(header)

	boundsServicePortType := mysoap.NewInBoundsServicePortType(client)
	inBounds, err := boundsServicePortType.InBounds(&mysoap.InBoundsRequest{MessageType: "SDATE"})
	if err != nil {
		utils.LogDebug(fmt.Sprintf("InBounds error: %v", err))
	}

	c.Data["json"] = inBounds
	c.ServeJSON()
}
