package controllers

import (
	"encoding/json"
	"strconv"

	"BeeCustom/enums"
	"BeeCustom/models"
)

type OrderFileController struct {
	BaseController
}

func (c *OrderFileController) Prepare() {
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
func (c *OrderFileController) DataGrid() {
	// 直接获取参数 getDataGridData()
	params := models.NewOrderFileQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	// 获取数据列表和总数
	data, total := models.OrderFilePageList(&params)
	c.ResponseList(c.transformOrderFileList(data), total)
	c.ServeJSON()
}

// TransformOrderList 格式化列表数据
func (c *OrderFileController) transformOrderFileList(ms []*models.OrderFile) []*map[string]interface{} {
	var orderFileList []*map[string]interface{}
	for _, v := range ms {
		OrderFile := make(map[string]interface{})
		OrderFile["Id"] = strconv.FormatInt(v.Id, 10)
		//OrderFile["Type"] = v.Type
		//OrderFile["Name"] = v.Name
		//OrderFile["Url"] = v.Url
		OrderFile["Creator"] = v.Creator
		OrderFile["Version"] = v.Version
		OrderFile["CreatedAt"] = v.CreatedAt.Format(enums.BaseDateTimeFormat)

		orderFileList = append(orderFileList, &OrderFile)
	}

	return orderFileList
}
