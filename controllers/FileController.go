package controllers

import (
	"BeeCustom/enums"
	"github.com/satori/go.uuid"
)

//FileController 文件上传
type FileController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *FileController) Prepare() {
	//先执行
	c.BaseController.Prepare()

	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	//默认认证 "Index", "Create", "Edit", "Delete"
	c.checkAuthor()

	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//c.checkLogin()//权限控制里会进行登录验证，因此这里不用再作登录验证
}

//Upload 文件上传
func (c *FileController) Upload() {
	f, h, err := c.GetFile("filename")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "上传失败", nil)
	}

	if f != nil {
		defer f.Close()
	} else {
		c.jsonResult(enums.JRCodeFailed, "上传失败", nil)
	}

	uid, _ := uuid.NewV4()
	var fileNamePath string
	if h != nil {
		fileNamePath = "static/upload/" + uid.String() + "_" + h.Filename
		_ = c.SaveToFile("uploadname", fileNamePath) // 保存位置在 static/upload, 没有文件夹要先创建
	} else {
		c.jsonResult(enums.JRCodeFailed, "上传失败", nil)
	}

	c.jsonResult(enums.JRCodeSucc, "添加成功", fileNamePath)

}
