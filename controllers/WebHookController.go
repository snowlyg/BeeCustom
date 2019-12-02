package controllers

import (
	"fmt"
	"io/ioutil"

	"BeeCustom/utils"
)

const SECRETTOKEN = "bee_custom_auto_pull"

// WebHookController handles WebSocket requests.
type WebHookController struct {
	BaseController
}

func (c *WebHookController) Get() {
	signature := c.Ctx.Request.Header.Get("X-Coding-Signature")
	utils.LogDebug(fmt.Sprintf("calculate_signature:%v", signature))
	content, err := ioutil.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("c.Ctx.Request.Body.Read:%v", err))
	}
	utils.LogDebug(fmt.Sprintf("ob:%v", content))

	//mac := hmac.New(sha1.New,[]byte(SECRETTOKEN))
	//mac.Write([]byte(ob))
	//
	//CalculateSignature := "sha1=" + mac.Sum(nil)
	//if CalculateSignature == signature {
	//	utils.LogDebug(fmt.Sprintf("calculate_signature:%v", CalculateSignature))
	//}

	c.ServeJSON()
}
