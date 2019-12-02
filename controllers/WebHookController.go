package controllers

import (
	"encoding/json"
	"fmt"

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
	var ob interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &ob)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("c.Ctx.Request.Body.Read:%v", err))
	}

	//mac := hmac.New(sha1.New,[]byte(SECRETTOKEN))
	//mac.Write([]byte(ob))
	//
	//CalculateSignature := "sha1=" + mac.Sum(nil)
	//if CalculateSignature == signature {
	//	utils.LogDebug(fmt.Sprintf("calculate_signature:%v", CalculateSignature))
	//}

	c.ServeJSON()
}
