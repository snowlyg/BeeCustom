package controllers

import (
	"BeeCustom/enums"
)

const SECRETTOKEN = "bee_custom_auto_pull"

// WebHookController handles WebSocket requests.
type WebHookController struct {
	BaseController
}

func (c *WebHookController) Get() {
	signature := c.Ctx.Request.Header.Get("X-Coding-Signature")
	sha1 := enums.Hmac(SECRETTOKEN, c.Ctx.Input.RequestBody)
	calculateSignature := "sha1=" + sha1
	if calculateSignature == signature {
		enums.Cmd("cd", "", []string{"/root/go/src/BeeCustom"})
		enums.Cmd("git", "", []string{"pull"})
		enums.Cmd("supervisorctl", "", []string{"restart", "beepkg"})
	}

	c.Data["json"] = "ok"
	c.ServeJSON()
}
