package controllers

import (
	"fmt"
	"strings"

	"BeeCustom/enums"
	"BeeCustom/models"
	"BeeCustom/utils"
	"BeeCustom/validations"
)

type HomeController struct {
	BaseController
}

func (c *HomeController) Index() {
	//判断是否登录
	c.checkLogin()
	c.setTpl("home/index.html", "shared/layout_base.html")
}

func (c *HomeController) Control() {
	//判断是否登录
	c.checkLogin()
	c.setTpl()
}

func (c *HomeController) Page404() {
	c.setTpl()
}

func (c *HomeController) Error() {
	c.Data["error"] = c.GetString(":error")
	c.setTpl("home/error.html", "shared/layout_app.html")
}

func (c *HomeController) Login() {
	c.LayoutSections = make(map[string]string)
	c.setTpl("auth/login.html", "auth/layout_base.html")
	c.GetXSRFToken()
}

func (c *HomeController) DoLogin() {
	username := strings.TrimSpace(c.GetString("UserName"))
	userpwd := strings.TrimSpace(c.GetString("UserPwd"))

	errMsg := validations.LoginValid(username, userpwd)
	if len(errMsg) > 0 {
		c.jsonResult(enums.JRCodeFailed, errMsg, "")
	}

	if len(username) == 0 || len(userpwd) == 0 {
		c.jsonResult(enums.JRCodeFailed, "用户名和密码不正确", "")
	}

	userpwd = utils.String2md5(userpwd)
	user, err := models.BackendUserOneByUserName(username, userpwd)

	if user != nil && err == nil {
		if user.Status == enums.Disabled {
			c.jsonResult(enums.JRCodeFailed, "用户被禁用，请联系管理员", "")
		}
		//保存用户信息到session
		if err = c.setBackendUser2Session(user.Id); err != nil {
			utils.LogDebug(fmt.Sprintf("用户sessions失败:%v", err))
			c.jsonResult(enums.JRCodeFailed, "用户sessions失败", "")
		}

		//获取用户信息
		c.jsonResult(enums.JRCodeSucc, "登录成功", "")
	} else {
		c.jsonResult(enums.JRCodeFailed, "用户名或者密码错误", "")
	}
}

//
func (c *HomeController) Logout() {
	user := models.NewBackendUser(0)
	c.SetSession("backenduser", user)
	c.pageLogin()
}

//初始化数据
func (c *HomeController) DataReset() {
	if ok, err := models.DataReset(); ok {
		c.jsonResult(enums.JRCodeSucc, "初始化成功", "")
	} else {
		c.jsonResult(enums.JRCodeFailed, "初始化失败,可能原因:"+err.Error(), "")
	}

}
