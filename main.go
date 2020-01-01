package main

import (
	_ "BeeCustom/routers"
	_ "BeeCustom/sysinit"
	"github.com/astaxie/beego"
)

func main() {

	beego.BConfig.WebConfig.TemplateLeft = "@{{"
	beego.BConfig.WebConfig.TemplateRight = "}}"

	beego.Run()
}
