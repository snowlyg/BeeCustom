package controllers

import (
	"fmt"

	"BeeCustom/utils"
)

const SECRETTOKEN = "bee_custom_auto_pull"

// WebHookController handles WebSocket requests.
type WebHookController struct {
	BaseController
}

func (c *WebHookController) Get() {

	event := c.GetString("X-Coding-Event")

	utils.LogDebug(fmt.Sprintf("X-Coding-Event:%v", event))
	//delivery := request.get_header("X-Coding-Delivery")
	//webHook_version := request.get_header("X-Coding-WebHook-Version")
	//signature := request.get_header("X-Coding-Signature")
	//content := request.body.read()
	//sha1 := hmac.new(bytes(SECRET_TOKEN, encoding := "utf8"), content, "sha1")
	//sha1 := sha1.hexdigest()
	//calculate_signature := "sha1=" + sha1

	c.ServeJSON()
}
