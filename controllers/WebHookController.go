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
	signature := c.Ctx.Request.Header.Get("X-Coding-Signature")

	utils.LogDebug(fmt.Sprintf(" c.Ctx.Input.RequestBody :%v", c.Ctx.Input.RequestBody))
	sha1 := enums.Hmac(SECRETTOKEN, c.Ctx.Input.RequestBody)
	calculateSignature := "sha1=" + sha1
	if calculateSignature == signature {
		enums.Cmd("cd", "", []string{"/root/go/src/BeeCustom"})
		enums.Cmd("git", "", []string{"pull"})
		enums.Cmd("go build", "", []string{""})
		enums.Cmd("supervisorctl", "", []string{"restart", "beepkg"})
	}

	c.Data["json"] = "ok"
	c.ServeJSON()
}
