package controllers

import (
	"encoding/json"

	"BeeCustom/enums"
	"BeeCustom/models"
	"github.com/snowlyg/gotransform"
)

type OrderReturnController struct {
	BaseController
}

func (c *OrderReturnController) Prepare() {
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
func (c *OrderReturnController) DataGrid() {
	// 直接获取参数 getDataGridData()
	params := models.NewOrderReturnQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	// 获取数据列表和总数
	data, total := models.OrderReturnPageList(&params)
	c.ResponseList(c.transformOrderReturnList(data), total)
	c.ServeJSON()
}

// TransformOrderList 格式化列表数据
func (c *OrderReturnController) transformOrderReturnList(ms []*models.OrderReturn) []*map[string]interface{} {
	var orderReturnList []*map[string]interface{}
	for _, v := range ms {
		OrderReturn := make(map[string]interface{})
		//	OrderReturn["Id"] = strconv.FormatInt(v.Id, 10)
		//	//OrderReturn["CheckInfo"] = v.CheckInfo
		//	//OrderReturn["DealFlag"] = v.DealFlag
		//	//OrderReturn["EtpsPreentNo"] = v.EtpsPreentNo
		//	//OrderReturn["ManageResult"] = v.ManageResult
		//	//OrderReturn["BusinessId"] = v.BusinessId
		//	//OrderReturn["Reason"] = v.Reason
		//	//OrderReturn["SeqNo"] = v.SeqNo
		//	//OrderReturn["Rmk"] = v.Rmk
		//	//OrderReturn["CreateDate"] = v.CreateDate.Format(enums.BaseDateTimeFormat)
		//	OrderReturn["CreatedAt"] = v.CreatedAt.Format(enums.BaseDateTimeFormat)
		//

		g := gotransform.NewTransform(&OrderReturn, v, enums.BaseDateTimeFormat)
		_ = g.Transformer()
		orderReturnList = append(orderReturnList, &OrderReturn)
	}

	return orderReturnList
}
