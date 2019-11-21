package controllers

import (
	"encoding/json"
	"strconv"

	"BeeCustom/enums"
	"BeeCustom/models"
)

type AnnotationRecordController struct {
	BaseController
}

func (c *AnnotationRecordController) Prepare() {
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
func (c *AnnotationRecordController) DataGrid() {
	// 直接获取参数 getDataGridData()
	params := models.NewAnnotationRecordQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	// 获取数据列表和总数
	data, total := models.AnnotationRecordPageList(&params)
	c.ResponseList(c.transformAnnotationRecordList(data), total)
	c.ServeJSON()
}

// TransformAnnotationList 格式化列表数据
func (c *AnnotationRecordController) transformAnnotationRecordList(ms []*models.AnnotationRecord) []*map[string]interface{} {
	var annotationRecordList []*map[string]interface{}
	for _, v := range ms {
		AnnotationRecord := make(map[string]interface{})
		AnnotationRecord["Id"] = strconv.FormatInt(v.Id, 10)
		AnnotationRecord["Content"] = v.Content
		AnnotationRecord["Remark"] = v.Remark
		AnnotationRecord["CreatedAt"] = v.CreatedAt.Format(enums.BaseDateTimeFormat)
		if v.BackendUser != nil {
			AnnotationRecord["BackendUserMobile"] = v.BackendUser.Mobile
			AnnotationRecord["BackendUserName"] = v.BackendUser.RealName
		} else {
			AnnotationRecord["BackendUserMobile"] = "系统自动操作"
			AnnotationRecord["BackendUserName"] = "系统自动操作"
		}

		annotationRecordList = append(annotationRecordList, &AnnotationRecord)
	}

	return annotationRecordList
}
