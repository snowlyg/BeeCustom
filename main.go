package main

import (
	_ "BeeCustom/routers"
	_ "BeeCustom/sysinit"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
