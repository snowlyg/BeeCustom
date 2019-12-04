package controllers

import (
	"fmt"
	"io/ioutil"

	"BeeCustom/enums"
	"BeeCustom/file"
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

	if calculateSignature == signature {
		enums.Cmd("cd", []string{"/root/go/src/BeeCustom"})
		enums.Cmd("git", []string{"pull"})
		enums.Cmd("bee", []string{"pack"})
		if file.IsExist("/root/go/src/BeeCustom/Beecustom.tar.gz") {
			enums.Cmd("mv", []string{"Beecustom.tar.gz", "/root/back"})
		}

		if !file.IsExist("/root/go/src/BeeCustom/Beecustom.tar.gz") && file.IsExist("/root/back/Beecustom.tar.gz") {
			enums.Cmd("cd", []string{"/root/back"})
			enums.Cmd("tar", []string{"-zxvf", "Beecustom.tar.gz", "Beecustom"})
		}

		if file.IsExist("/root/back/BeeCustom") {
			enums.Cmd("mv", []string{"BeeCustom", "/root/go/src/BeeCustom"})
		}

		enums.Cmd("cd", []string{"/etc/supervisord.conf.d"})
		enums.Cmd("supervisorctl", []string{"restart", "beepkg"})
	}

	c.ServeJSON()
}
