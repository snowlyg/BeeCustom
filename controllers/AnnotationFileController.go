package controllers

import (
	"encoding/json"
	"strconv"

	"BeeCustom/enums"
	"BeeCustom/models"
)

type AnnotationFileController struct {
	BaseController
}

func (c *AnnotationFileController) Prepare() {
	// 先执行
	c.BaseController.Prepare()
	// 如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	// 默认认证 "Index", "Create", "Edit", "Delete"
	var perms []string
	c.checkAuthor(perms)

	// 如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	// 权限控制里会进行登录验证，因此这里不用再作登录验证
	// c.checkLogin()

}

// 列表数据
func (c *AnnotationFileController) DataGrid() {
	// 直接获取参数 getDataGridData()
	params := models.NewAnnotationFileQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	// 获取数据列表和总数
	data, total := models.AnnotationFilePageList(&params)
	c.ResponseList(c.transformAnnotationFileList(data), total)
	c.ServeJSON()
}

// TransformAnnotationList 格式化列表数据
func (c *AnnotationFileController) transformAnnotationFileList(ms []*models.AnnotationFile) []*map[string]interface{} {
	var annotationFileList []*map[string]interface{}
	for _, v := range ms {
		AnnotationFile := make(map[string]interface{})
		AnnotationFile["Id"] = strconv.FormatInt(v.Id, 10)
		AnnotationFile["Type"] = v.Type
		AnnotationFile["Name"] = v.Name
		AnnotationFile["Url"] = v.Url
		AnnotationFile["Creator"] = v.Creator
		AnnotationFile["Version"] = v.Version
		AnnotationFile["CreatedAt"] = v.CreatedAt.Format(enums.BaseDateTimeFormat)

		annotationFileList = append(annotationFileList, &AnnotationFile)
	}

	return annotationFileList
}
