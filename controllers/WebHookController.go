package controllers

import (
	"fmt"
	"io/ioutil"

	"BeeCustom/enums"
	"BeeCustom/file"
	"BeeCustom/utils"
	"github.com/astaxie/beego"
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
		dbType := beego.AppConfig.String("db_type")
		// 数据库名称
		dbName := beego.AppConfig.String(dbType + "::db_name")
		// 数据库连接用户名
		dbUser := beego.AppConfig.String(dbType + "::db_user")
		// 数据库连接用户名
		dbPwd := beego.AppConfig.String(dbType + "::db_pwd")
		// 数据库IP（域名）
		dbHost := beego.AppConfig.String(dbType + "::db_host")
		// 数据库端口
		dbPort := beego.AppConfig.String(dbType + "::db_port")
		arv := []string{"-driver=mysql", fmt.Sprintf(`-conn="%s:%s@tcp(%s:%s)/%s"`, dbUser, dbPwd, dbHost, dbPort, dbName)}
		enums.Cmd("cd", []string{"/root/go/src/BeeCustom"})
		enums.Cmd("git", []string{"pull"})
		enums.Cmd("bee", []string{"pack"})
		enums.Cmd("bee", arv)
		if file.IsExist("/root/go/src/BeeCustom/BeeCustom.tar.gz") {
			enums.Cmd("mv", []string{"BeeCustom.tar.gz", "/root/back"})
			utils.LogDebug("mv BeeCustom.tar.gz")
		} else {
			utils.LogDebug("mv BeeCustom.tar.gz error")
		}

		if !file.IsExist("/root/go/src/BeeCustom/BeeCustom.tar.gz") && file.IsExist("/root/back/BeeCustom.tar.gz") {
			enums.Cmd("cd", []string{"/root/back"})
			enums.Cmd("tar", []string{"-zxvf", "BeeCustom.tar.gz", "BeeCustom"})
			enums.Cmd("rm", []string{"BeeCustom.tar.gz"})
			utils.LogDebug("tar BeeCustom.tar.gz")
		} else {
			utils.LogDebug("tar BeeCustom.tar.gz error")
		}

		if file.IsExist("/root/back/BeeCustom") && !file.IsExist("/root/back/BeeCustom.tar.gz") {
			enums.Cmd("mv", []string{"BeeCustom", "/root/go/src/BeeCustom"})
			utils.LogDebug("mv BeeCustom")
		} else {
			utils.LogDebug("mv BeeCustom error")
		}

		enums.Cmd("cd", []string{"/etc/supervisord.conf.d"})
		enums.Cmd("supervisorctl", []string{"restart", "beepkg"})
	}

	c.ServeJSON()
}
