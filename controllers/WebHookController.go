package controllers

import (
	"fmt"
	"io/ioutil"

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
	content, err := ioutil.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("c.Ctx.Request.Body.Read:%v", err))
	}

	sha1 := enums.Hmac(SECRETTOKEN, content)
	calculateSignature := "sha1=" + sha1

	utils.LogDebug(fmt.Sprintf("calculateSignature == signature:%v -%v -%v", calculateSignature == signature, calculateSignature, signature))
	//if calculateSignature == signature {
	enums.Cmd("cd", "", []string{"/root/go/src/BeeCustom"})
	enums.Cmd("git", "", []string{"pull"})
	enums.Cmd("go", "", []string{"build"})
	enums.Cmd("supervisorctl", "", []string{"restart", "beepkg"})
	//}

	c.Data["json"] = "ok"
	c.ServeJSON()
}
