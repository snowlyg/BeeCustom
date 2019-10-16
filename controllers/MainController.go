package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "BeeCustom"
	c.Data["Email"] = "569616226@qq.com"
	c.TplName = "index.tpl"
}
