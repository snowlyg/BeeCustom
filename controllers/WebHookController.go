package controllers

import (
	"io/ioutil"

	"BeeCustom/enums"
)

const SECRETTOKEN = "bee_custom_auto_pull"

// WebHookController handles WebSocket requests.
type WebHookController struct {
	BaseController
}

func (c *WebHookController) Prepare() {
	c.EnableXSRF = false
}

func (c *WebHookController) Get() {
	signature := c.Ctx.Request.Header.Get("X-Coding-Signature") //获取加密签名
	contentType := c.Ctx.Request.Header.Get("Content-Type")     //获取加密签名

	var res []byte
	if contentType == "application/json; charset=UTF-8" {
		res, _ = ioutil.ReadAll(c.Ctx.Request.Body) // for application/json
	} else if contentType == "application/x-www-form-urlencoded; charset=UTF-8" { // for application/x-www-form-urlencoded
		res = []byte(c.GetString(":payload"))
	}

	sha1 := enums.Hmac(SECRETTOKEN, res) // for application/x-www-form-urlencoded
	calculateSignature := "sha1=" + sha1 // 重新加密内容
	if calculateSignature == signature {
		enums.Cmd("cd", "", []string{"/root/go/src/BeeCustom"})
		enums.Cmd("git", "", []string{"pull"})
		enums.Cmd("/usr/local/go/bin/go ", "", []string{"build"})
		enums.Cmd("supervisorctl", "", []string{"restart", "beepkg"})
	}
	data := struct {
		Status      bool
		Payload     string
		ContentType string
	}{
		Status:      calculateSignature == signature,
		Payload:     string(res),
		ContentType: contentType,
	}
	c.Data["json"] = data
	c.ServeJSON()
}
