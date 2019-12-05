package controllers

import (
	"strconv"

	"BeeCustom/enums"
)

//FileController 文件上传
type FileController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *FileController) Prepare() {
	//先执行
	c.BaseController.Prepare()

	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	c.checkLogin() //权限控制里会进行登录验证，因此这里不用再作登录验证
}

//image 文件上传
func (c *FileController) Upload() {
	fileType := "img/" + strconv.FormatInt(c.curUser.Id, 10) + "/"
	if fileNamePath, err := c.BaseUpload(fileType); err != nil {
		c.jsonResult(enums.JRCodeFailed, "上传失败", nil)
	} else {
		c.jsonResult(enums.JRCodeFailed, "上传成功", "/"+fileNamePath)
	}
}

//pdf 文件上传
func (c *FileController) OrderDataUpload() {
	orderId := c.GetString(":id", "0")
	if orderId == "0" {
		c.jsonResult(enums.JRCodeFailed, "参数错误", nil)
	}

	fileType := "order/" + orderId + "/" + strconv.FormatInt(c.curUser.Id, 10) + "/"
	if fileNamePath, err := c.BaseUpload(fileType); err != nil {
		c.jsonResult(enums.JRCodeFailed, "上传失败", nil)
	} else {
		c.jsonResult(enums.JRCodeFailed, "上传成功", fileNamePath)
	}
}
