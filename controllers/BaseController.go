package controllers

import (
	"fmt"
	"strings"

	"BeeCustom/enums"
	"BeeCustom/file"
	"BeeCustom/models"
	"BeeCustom/utils"
	"BeeCustom/validations"
	"github.com/astaxie/beego/validation"

	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
	controllerName string             //当前控制名称
	actionName     string             //当前action名称
	curUser        models.BackendUser //当前用户信息
}

func (c *BaseController) Prepare() {
	//附值
	c.controllerName, c.actionName = c.GetControllerAndAction()
	//从Session里获取数据 设置用户信息
	c.adapterUserInfo()
}

func (c *BaseController) GetXSRFToken() {
	c.Data["xsrf_token"] = c.XSRFToken()
}

// checkLogin判断用户是否登录，未登录则跳转至登录页面
// 一定要在BaseController.Prepare()后执行
func (c *BaseController) checkLogin() {
	if c.curUser.Id == 0 {
		//登录页面地址
		urlstr := c.URLFor("HomeController.Login") + "?url="
		//登录成功后返回的址为当前
		returnURL := c.Ctx.Request.URL.Path
		//如果ajax请求则返回相应的错码和跳转的地址
		if c.Ctx.Input.IsAjax() {
			//由于是ajax请求，因此地址是header里的Referer
			returnURL := c.Ctx.Input.Refer()
			c.jsonResult(enums.JRCode302, "请登录", urlstr+returnURL)
		}
		c.Redirect(urlstr+returnURL, 302)
		c.StopRun()
	}
}

// 判断某 Controller.Action 当前用户是否有权访问
func (c *BaseController) checkActionAuthor(ctrlName, ActName string) bool {
	if c.curUser.Id == 0 {
		return false
	}
	//从session获取用户信息
	user := c.GetSession("backenduser")
	//类型断言
	bu, ok := user.(models.BackendUser)
	if ok {
		//如果是超级管理员，则直接通过
		if bu.IsSuper {
			return true
		}

		//遍历用户所负责的资源列表
		for _, resource := range bu.Role.Resources {
			urlfor := resource.UrlFor
			if len(urlfor) == 0 {
				continue
			}
			if len(urlfor) > 0 && urlfor == (ctrlName+"."+ActName) {
				return true
			}
		}
	}

	return false
}

// 判断某 Controller.Action 当前用户是否有权访问
func (c *BaseController) getActionData(actionNames ...string) {
	for _, v := range actionNames {
		c.Data["can"+v] = c.checkActionAuthor(c.controllerName, v)
	}
}

// checkLogin判断用户是否有权访问某地址，无权则会跳转到错误页面
//一定要在BaseController.Prepare()后执行
// 会调用checkLogin
// 传入的参数为忽略权限控制的Action
func (c *BaseController) checkAuthor(actionNames ...string) {
	//先判断是否登录
	c.checkLogin()
	dActionNames := append(actionNames, "Index", "Create", "Edit", "Delete") //默认需要验证的权限
	//如果Action在忽略列表里，则直接通用
	for _, actionName := range dActionNames {
		if actionName == c.actionName {
			hasAuthor := c.checkActionAuthor(c.controllerName, c.actionName)
			if !hasAuthor {
				utils.LogDebug(fmt.Sprintf("author control: path=%s.%s userid=%v  无权访问", c.controllerName, c.actionName, c.curUser.Id))
				//如果没有权限
				if !hasAuthor {
					if c.Ctx.Input.IsAjax() {
						c.jsonResult(enums.JRCode401, "无权访问", "")
					} else {
						c.pageError("无权访问")
					}
				}
			}
		} else {
			return
		}
	}

}

//从session里取用户信息
func (c *BaseController) adapterUserInfo() {
	a := c.GetSession("backenduser")
	if a != nil {
		c.curUser = a.(models.BackendUser)
		c.Data["backenduser"] = a
	}
}

//SetBackendUser2Session 获取用户信息（包括资源UrlFor）保存至Session
func (c *BaseController) setBackendUser2Session(userId int64) error {
	m, err := models.BackendUserOne(userId)
	if err != nil {
		return err
	}

	c.SetSession("backenduser", *m)
	return nil
}

// 设置模板
// 第一个参数模板，第二个参数为layout
func (c *BaseController) setTpl(template ...string) {
	var tplName string
	layout := "shared/layout_app.html"
	switch {
	case len(template) == 1:
		tplName = template[0]
	case len(template) == 2:
		tplName = template[0]
		layout = template[1]
	default:
		//不要Controller这个10个字母
		ctrlName := strings.ToLower(c.controllerName[0 : len(c.controllerName)-10])
		actionName := strings.ToLower(c.actionName)
		tplName = ctrlName + "/" + actionName + ".html"
	}
	c.Layout = layout
	c.TplName = tplName
}

func (c *BaseController) jsonResult(status enums.JsonResultCode, msg string, obj interface{}) {
	r := &models.JsonResult{status, msg, obj}
	c.Data["json"] = r
	c.ServeJSON()
	c.StopRun()
}

// 重定向
func (c *BaseController) redirect(url string) {
	c.Redirect(url, 302)
	c.StopRun()
}

// 重定向 去错误页
func (c *BaseController) pageError(msg string) {
	errorurl := c.URLFor("HomeController.Error") + "/" + msg
	c.Redirect(errorurl, 302)
	c.StopRun()
}

// 重定向 去登录页
func (c *BaseController) pageLogin() {
	url := c.URLFor("HomeController.Login")
	c.Redirect(url, 302)
	c.StopRun()
}

// 验证提交数据
func (c *BaseController) validRequestData(m interface{}) {
	var errMsg string

	validation.MessageTmpls = validations.GetMessageTmpls()
	valid := validation.Validation{}

	b, err := valid.Valid(m)

	if !b {
		if len(valid.Errors) > 0 {
			errMsg = strings.Split(valid.Errors[0].Key, ".")[0] + ":" + valid.Errors[0].Message
		}
	}

	if err != nil {
		c.jsonResult(enums.JRCodeFailed, fmt.Sprintf("表单数据验证错误:%v", err), m)
	} else if len(errMsg) > 0 {
		c.jsonResult(enums.JRCodeFailed, errMsg, m)
	}

}

//上传文件
func (c *BaseController) BaseUpload(fileType string) (string, error) {
	f, h, err := c.GetFile("filename")
	if err != nil {
		utils.LogDebug(fmt.Sprintf("参数错误:%v", err))
		return "", err
	}

	if fileNamePath, err := file.GetUploadFileUPath(f, h, fileType); err != nil {
		return "", err
	} else {
		err = c.SaveToFile("filename", fileNamePath) // 保存位置在 static/upload, 没有文件夹要先创建
		if err != nil {
			utils.LogDebug(fmt.Sprintf("图片保存失败:%v", err))
			return "", err
		} else {
			return fileNamePath, nil
		}
	}
}
