package controllers

import (
	"fmt"

	"BeeCustom/enums"
	"BeeCustom/utils"
)

const SECRETTOKEN = "bee_custom_auto_pull"

// WebHookController handles WebSocket requests.
type WebHookController struct {
	BaseController
}

func (c *WebHookController) Get() {

	signature := c.GetString("X-Coding-Signature")
	//content := request.body.read()
	sha1 := enums.Sha1(SECRETTOKEN)

	CalculateSignature := "sha1=" + sha1
	if CalculateSignature == signature {
		utils.LogDebug(fmt.Sprintf("calculate_signature:%v", CalculateSignature))
	}

	c.ServeJSON()
}
