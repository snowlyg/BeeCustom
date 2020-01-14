package controllers

import (
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
	//res, err := ioutil.ReadAll(c.Ctx.Request.Body) // for application/json
	palyload := c.GetString(":payload")
	sha1 := enums.Hmac(SECRETTOKEN, []byte(palyload)) // for application/x-www-form-urlencoded
	calculateSignature := "sha1=" + sha1              // 重新加密内容
	if calculateSignature == signature {
		enums.Cmd("cd", "", []string{"/root/go/src/BeeCustom"})
		enums.Cmd("git", "", []string{"pull"})
		enums.Cmd("go build", "", []string{""})
		enums.Cmd("supervisorctl", "", []string{"restart", "beepkg"})
	}
	data := struct {
		Status  bool
		Payload string
	}{
		Status:  calculateSignature == signature,
		Payload: palyload,
	}
	c.Data["json"] = data
	c.ServeJSON()
}
