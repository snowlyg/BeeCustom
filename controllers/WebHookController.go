package controllers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"

	"BeeCustom/enums"
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

	sha1 := enums.Hmac(SECRETTOKEN, content)
	CalculateSignature := "sha1=" + sha1
	utils.LogDebug(fmt.Sprintf("CalculateSignature:%v", CalculateSignature))
	if CalculateSignature == signature {
		utils.LogDebug(fmt.Sprintf("calculate_signature:%v", CalculateSignature))
		cmd := exec.Command("cd", "/root/go/src/BeeCustom")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			utils.LogDebug(fmt.Sprintf("calculate_signature:%v", err))
		}
	}

	c.ServeJSON()
}
