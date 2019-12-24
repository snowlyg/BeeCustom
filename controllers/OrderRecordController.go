package controllers

import (
	"encoding/json"
	"strconv"

	"BeeCustom/enums"
	"BeeCustom/models"
)

type OrderRecordController struct {
	BaseController
}

func (c *OrderRecordController) Prepare() {
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
func (c *OrderRecordController) DataGrid() {
	// 直接获取参数 getDataGridData()
	params := models.NewOrderRecordQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	// 获取数据列表和总数
	data, total := models.OrderRecordPageList(&params)
	c.ResponseList(c.transformOrderRecordList(data), total)
	c.ServeJSON()
}

// TransformOrderList 格式化列表数据
func (c *OrderRecordController) transformOrderRecordList(ms []*models.OrderRecord) []*map[string]interface{} {
	var orderRecordList []*map[string]interface{}
	for _, v := range ms {
		OrderRecord := make(map[string]interface{})
		OrderRecord["Id"] = strconv.FormatInt(v.Id, 10)
		OrderRecord["Content"] = v.Content
		OrderRecord["Remark"] = v.Remark
		OrderRecord["CreatedAt"] = v.CreatedAt.Format(enums.BaseDateTimeFormat)
		if v.BackendUser != nil {
			OrderRecord["BackendUserMobile"] = v.BackendUser.Mobile
			OrderRecord["BackendUserName"] = v.BackendUser.RealName
		} else {
			OrderRecord["BackendUserMobile"] = "系统自动操作"
			OrderRecord["BackendUserName"] = "系统自动操作"
		}

		orderRecordList = append(orderRecordList, &OrderRecord)
	}

	return orderRecordList
}
