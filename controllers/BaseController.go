package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

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
	controllerName string             // 当前控制名称
	actionName     string             // 当前action名称
	curUser        models.BackendUser // 当前用户信息
}

func (c *BaseController) Prepare() {
	// 附值
	c.controllerName, c.actionName = c.GetControllerAndAction()
	// 从Session里获取数据 设置用户信息
	c.adapterUserInfo()
}

func (c *BaseController) GetXSRFToken() {
	c.Data["xsrf_token"] = c.XSRFToken()
}

func (c *BaseController) ResponseList(data interface{}, total int64) {
	// 定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	result["code"] = 0
	c.Data["json"] = result
}

// checkLogin判断用户是否登录，未登录则跳转至登录页面
// 一定要在BaseController.Prepare()后执行
func (c *BaseController) checkLogin() {
	if c.curUser.Id == 0 {
		// 登录页面地址
		urlstr := c.URLFor("HomeController.Login") + "?url="
		// 登录成功后返回的址为当前
		returnURL := c.Ctx.Request.URL.Path
		// 如果ajax请求则返回相应的错码和跳转的地址
		if c.Ctx.Input.IsAjax() {
			// 由于是ajax请求，因此地址是header里的Referer
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

	// 从session获取用户信息
	user := c.GetSession("backenduser")
	// 类型断言
	bu, ok := user.(models.BackendUser)
	if ok {
		// 如果是超级管理员，则直接通过
		if bu.IsSuper {
			return true
		}
		roleIds, err := utils.E.GetRolesForUser(strconv.FormatInt(bu.Id, 10))
		if err != nil {
			return false
		}

		for _, roleId := range roleIds {
			return utils.E.HasPermissionForUser(roleId, ctrlName+"."+ActName)
		}
	}

	return false
}

// 判断某 Controller.Action 当前用户是否有权访问
func (c *BaseController) getActionData(upper string, actionNames ...string) {

	// 清单，货物进出口分离权限
	if len(upper) > 0 {
		actions := make(map[string]map[string]bool)
		action := make(map[string]bool)
		for _, v := range actionNames {
			isAction := c.checkActionAuthor(c.controllerName, upper+v)
			action["can"+v] = isAction
		}
		actions[upper] = action
		c.Data["actions"] = actions
	} else {
		for _, v := range actionNames {
			c.Data["can"+v] = c.checkActionAuthor(c.controllerName, v)
		}
	}

}

// checkLogin判断用户是否有权访问某地址，无权则会跳转到错误页面
// 一定要在BaseController.Prepare()后执行
// 会调用checkLogin
// 传入的参数为需要权限控制的Action
func (c *BaseController) checkAuthor(actionNames []string) {
	// 先判断是否登录
	c.checkLogin()
	// 如果Action在忽略列表里，则直接通用
	for _, actionName := range actionNames {
		if actionName == c.actionName {
			hasAuthor := c.checkActionAuthor(c.controllerName, c.actionName)
			if !hasAuthor {
				utils.LogDebug(fmt.Sprintf("author control: path=%s.%s userid=%v  无权访问", c.controllerName, c.actionName, c.curUser.Id))
				// 如果没有权限
				if !hasAuthor {
					if c.Ctx.Input.IsAjax() {
						c.jsonResult(enums.JRCode401, "无权访问", "")
					} else {
						c.pageError("无权访问")
					}
				}
			}
		} else {
			continue
		}
	}

}

// 从session里取用户信息
func (c *BaseController) adapterUserInfo() {
	a := c.GetSession("backenduser")
	if a != nil {
		c.curUser = a.(models.BackendUser)
		c.Data["backenduser"] = a
	}
}

// SetBackendUser2Session 获取用户信息（包括资源UrlFor）保存至Session
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
		// 不要Controller这个10个字母
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

// 上传文件
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

// 格式时间
func (c *BaseController) GetDateTime(timeString, timeFormatString string) (*time.Time, error) {
	tS := c.GetString(timeString)
	fTimeFormat, err := time.ParseInLocation(timeFormatString, tS, time.Local)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("ParseInLocation:%v", err))

		return nil, err
	}

	return &fTimeFormat, nil
}

// 更新状态和状态更新时间
func UpdateAnnotationStatus(m *models.Annotation, StatusString string, isRestart bool) error {
	aStatusS, err := models.GetSettingRValueByKey("annotationStatus", false)
	aStatus, err, done := enums.TransformCnToInt(aStatusS, StatusString)
	if done {
		return err
	}

	if err != nil {
		utils.LogDebug(fmt.Sprintf("转换清单状态出错:%v", err))
		return err
	}

	// 禁止状态回退
	if m.Status < aStatus || isRestart {
		m.Status = aStatus
		m.StatusUpdatedAt = time.Now()
	}

	return nil
}

// 更新状态和状态更新时间
func UpdateOrderStatus(m *models.Order, StatusString string, isRestart bool) error {
	aStatusS, err := models.GetSettingRValueByKey("orderStatus", false)
	aStatus, err, done := enums.TransformCnToInt(aStatusS, StatusString)
	if done {
		return err
	}
	if err != nil {
		utils.LogDebug(fmt.Sprintf("转换状态出错:%v", err))
		return err
	}

	// 禁止状态回退
	if m.Status < aStatus || isRestart {
		m.Status = aStatus
		m.StatusUpdatedAt = time.Now()
	}

	return nil
}

// TransformHandBookGood 格式化列表数据
func (c *HandBookController) TransformHandBookGood(v *models.HandBookGood) map[string]interface{} {

	clearances1 := models.GetClearancesByTypes("货币代码", true)
	clearances2 := models.GetClearancesByTypes("计量单位代码", false)
	var unitOneCode interface{}
	var unitTwoCode interface{}
	var moneyunitCode interface{}
	for _, c := range clearances2 {
		if c[0] == v.UnitOne {
			unitOneCode = c[1]
		}

		if c[0] == v.UnitTwo {
			unitTwoCode = c[1]
		}
	}

	for _, c := range clearances1 {
		if c[0] == v.Moneyunit {
			moneyunitCode = c[1]
		}

	}
	handBook := make(map[string]interface{})
	handBook["Id"] = strconv.FormatInt(v.Id, 10)
	handBook["RecordNo"] = v.RecordNo
	handBook["HsCode"] = v.HsCode
	handBook["Name"] = v.Name
	handBook["Special"] = v.Special
	handBook["UnitOne"] = v.UnitOne
	handBook["UnitOneCode"] = unitOneCode
	handBook["UnitTwo"] = v.UnitTwo
	handBook["UnitTwoCode"] = unitTwoCode
	handBook["Price"] = v.Price
	handBook["Moneyunit"] = v.Moneyunit
	handBook["MoneyunitCode"] = moneyunitCode

	return handBook
}
