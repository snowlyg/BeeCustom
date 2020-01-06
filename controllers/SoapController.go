package controllers

import (
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
