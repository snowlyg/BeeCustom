package controllers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"BeeCustom/enums"
	"BeeCustom/file"
	"BeeCustom/models"
	"BeeCustom/utils"
	"BeeCustom/xmlTemplate"
)

type BaseOrderController struct {
	BaseController
}

// 列表数据
func (c *BaseOrderController) bDataGrid(ieFlag string) {
	// 直接获取参数 getDataGridData()
	params := models.NewOrderQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	params.IEFlag = ieFlag
	// 获取数据列表和总数
	data, total, err := models.OrderPageList(&params)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据列表和总数失败", nil)
	}
	err = models.OrderGetRelations(data, "Company,BackendUsers")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "关联关系获取失败", nil)
	}
	// 格式化数据
	orderList := c.TransformOrderList(data)
	c.ResponseList(orderList, total)
	c.ServeJSON()
}

func (c *BaseOrderController) bIndex(ieFlag string) {
	// 页面模板设置
	c.setTpl("order/index.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "order/index_footerjs.html"

	// 页面里按钮权限控制
	c.getActionData(ieFlag, "Index", "Create", "Edit", "Make", "ReMake", "Audit", "Delete", "Distribute", "Recheck", "Push", "PushXml", "StoreError", "Change", "Restart", "Cancel", "Copy")

	// 获取制单人
	backendUsers := models.GetCreateBackendUsers("OrderController.Make")
	c.Data["BackendUsers"] = backendUsers
	c.Data["IEFlag"] = ieFlag
	c.Data["IEFlagName"] = enums.GetImpexpMarkcdCNName(ieFlag)
	c.GetXSRFToken()
}

// 数据统计
func (c *BaseOrderController) bStatusCount(ieflag string) {
	// 直接获取参数 getDataGridData()
	params := models.NewOrderQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	params.IEFlag = ieflag
	// 获取数据列表和总数
	data, err := models.OrderStatusCount(&params)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据列表和总数出错", nil)
	}
	// 定义返回的数据结构
	result := make(map[string]interface{})
	result["rows"] = data
	result["status"] = 1
	c.Data["json"] = result
	c.ServeJSON()
}

// Create 添加 新建 页面
func (c *BaseOrderController) bCreate(ieflag string) {
	c.Data["canStore"] = c.getCanStore(nil, ieflag)
	c.getResponses(ieflag)
}

// Store 添加 新建 页面
func (c *BaseOrderController) bStore(iEFlag string) {
	m := models.NewOrder(0)
	// 获取form里的值
	if err := c.ParseForm(&m); err != nil {
		utils.LogDebug(fmt.Sprintf("ParseForm:%v", err))
		c.jsonResult(enums.JRCodeFailed, "获取数据出错", m)
	}

	aplDateString := c.GetString("AplDate")
	aplDate, err := time.Parse(enums.BaseDateFormat, aplDateString)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "时间格式出错", nil)
	}
	m.AplDate = aplDate
	m.IEFlag = iEFlag
	m.ContactSignDate = aplDate.AddDate(0, -1, 0)

	company, err := models.CompanyByManageCode(m.TradeCode)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取客户出错", nil)
	}
	if err = UpdateOrderStatus(&m, "待审核", false); err != nil {
		c.jsonResult(enums.JRCodeFailed, "添加失败", nil)
	}

	m.Company = company
	m.ClientSeqNo = c.getClientSeqNo(iEFlag, m.CustomMaster)
	decOtherPacks := c.GetStrings("DecOtherPacks")
	m.DecOtherPacks = strings.Join(decOtherPacks, ",")

	c.validRequestData(m)
	if err := models.OrderUpdateOrSave(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "添加失败", m)
	} else {
		if err := c.setAnnotaionUserRelType(&m, nil, "创建人"); err != nil {
			c.jsonResult(enums.JRCodeFailed, "添加失败", m)
		}
		orderRecord := c.newOrderRecord(&m, "创建订单")
		if err := models.OrderRecordSave(orderRecord); err != nil {
			c.jsonResult(enums.JRCodeFailed, "添加失败", m)
		}
		c.jsonResult(enums.JRCodeSucc, "添加成功", m)
	}
}

// copy 复制
func (c *BaseOrderController) bCopy(id int64) {
	m, err := models.OrderOne(id, "Company,OrderItems,OrderContainers,OrderDocuments")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据出错", nil)
	}
	if err := UpdateOrderStatus(m, "待审核", true); err != nil {
		c.jsonResult(enums.JRCodeFailed, "添加失败", nil)
	}
	// 重置数据
	if m != nil {
		m.Id = 0
		m.AplDate = time.Now()
		m.ClientSeqNo = c.getClientSeqNo(m.IEFlag, m.CustomMaster)
		m.SeqNo = ""
		m.EntryId = ""
	}
	if err := models.OrderUpdateOrSave(m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "添加失败", m)
	} else {
		if err := c.setAnnotaionUserRelType(m, nil, "创建人"); err != nil {
			c.jsonResult(enums.JRCodeFailed, "添加失败", m)
		}
		orderRecord := c.newOrderRecord(m, "创建订单")
		if err := models.OrderRecordSave(orderRecord); err != nil {
			c.jsonResult(enums.JRCodeFailed, "添加失败", m)
		}
		c.jsonResult(enums.JRCodeSucc, "添加成功", m)
	}
}

// Edit 添加 编辑 页面
func (c *BaseOrderController) bEdit(id int64) {
	m, err := models.OrderOne(id, "OrderItems,OrderContainers,OrderDocuments")
	if err != nil {
		c.pageError("数据无效，请刷新后重试")
	}
	c.setStatusOnly(m, "审核中", false)
	// 获取制单人
	backendUsers := models.GetCreateBackendUsers("OrderController.Make")
	c.Data["BackendUsers"] = backendUsers
	c.Data["m"] = models.TransformOrder(id, "OrderItems,OrderContainers,OrderDocuments", false)
	c.Data["canStore"] = c.getCanStore(m, "")
	if m != nil {
		c.getResponses(m.IEFlag)
	}
}

// bMake 制单
func (c *BaseOrderController) bMake(id int64) {
	m, err := models.OrderOne(id, "")
	if m != nil && id > 0 {
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}
	c.setStatusOnly(m, "制单中", false)
	c.Data["m"] = models.TransformOrder(id, "OrderItems,OrderContainers,OrderDocuments", false)
	c.Data["canStore"] = c.getCanStore(m, "")
	if m != nil {
		c.getResponses(m.IEFlag)
	}
}

// bReMake 驳回修改
func (c *BaseOrderController) bReMake(id int64) {
	m, err := models.OrderOne(id, "")
	if m != nil && id > 0 {
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}

	c.Data["m"] = models.TransformOrder(id, "OrderItems,OrderContainers,OrderDocuments,OrderRecords", false)
	c.Data["canStore"] = c.getCanStore(m, "")
	if m != nil {
		c.getResponses(m.IEFlag)
	}
}

// 编辑相关页面返回
func (c *BaseOrderController) getResponses(ieflag string) {
	// 页面里按钮权限控制
	c.getActionData(ieflag, "Audit", "Distribute", "ForRecheck", "Print", "Remark", "ReForRecheck")
	c.Data["IEFlagName"] = enums.GetImpexpMarkcdCNName(ieflag)
	c.Data["IEFlag"] = ieflag
	c.setTpl("order/change_create_edit_show.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "order/create_footerjs.html"
	c.GetXSRFToken()
}

// Cancel 取消订单
func (c *BaseOrderController) bCancel(id int64) {
	m, err := models.OrderOne(id, "")
	if m != nil && id > 0 {
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}
	if err = UpdateOrderStatus(m, "订单关闭", false); err != nil {
		c.jsonResult(enums.JRCodeFailed, "取消失败", m)
	}
	if err := models.OrderUpdateStatus(m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "取消失败", m)
	}
	orderRecord := c.newOrderRecord(m, "取消订单")
	if err := models.OrderRecordSave(orderRecord); err != nil {
		c.jsonResult(enums.JRCodeFailed, "取消失败", m)
	}
	c.jsonResult(enums.JRCodeSucc, "取消成功", m)
}

// Audit 审核通过订单
func (c *BaseOrderController) bAudit(id int64) {
	m, err := models.OrderOne(id, "")
	if m != nil && id > 0 {
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}
	if err := c.setAnnotaionUserRelType(m, nil, "审单人"); err != nil {
		c.jsonResult(enums.JRCodeFailed, "审核失败", m)
	}
	c.setStatusOnly(m, "审核通过", false)
	orderRecord := c.newOrderRecord(m, "审核订单")
	if err := models.OrderRecordSave(orderRecord); err != nil {
		c.jsonResult(enums.JRCodeFailed, "审核失败", m)
	}
	c.jsonResult(enums.JRCodeSucc, "审核通过", m)
}

// Distribute 分配
func (c *BaseOrderController) bDistribute(backendUserId, id int64) {
	bu, err := models.BackendUserOne(backendUserId)
	if bu != nil && backendUserId > 0 {
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}
	m, err := models.OrderOne(id, "")
	if m != nil && id > 0 {
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}
	if err = c.setAnnotaionUserRelType(m, bu, "制单人,派单人"); err != nil {
		c.jsonResult(enums.JRCodeFailed, "派单失败", m)
	}
	if err = UpdateOrderStatus(m, "待制单", false); err != nil {
		c.jsonResult(enums.JRCodeFailed, "派单失败", m)
	}
	if err := models.OrderUpdateOrSave(m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "派单失败", m)
	}
	orderRecord := c.newOrderRecord(m, "派单："+bu.RealName)
	if err := models.OrderRecordSave(orderRecord); err != nil {
		c.jsonResult(enums.JRCodeFailed, "派单失败", m)
	}
	c.jsonResult(enums.JRCodeSucc, "派单通过", m)
}

// Update 添加 编辑 页面
func (c *BaseOrderController) bUpdate(id int64) {
	m, err := models.OrderOne(id, "")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}
	// 获取form里的值
	if err := c.ParseForm(m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "ParseForm", m)
	}

	c.validRequestData(m)
	if err := models.OrderUpdateOrSave(m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "编辑失败", m)
	}
	orderRecord := c.newOrderRecord(m, "保存数据")
	if err := models.OrderRecordSave(orderRecord); err != nil {
		c.jsonResult(enums.JRCodeFailed, "编辑失败", m)
	}
	c.jsonResult(enums.JRCodeSucc, "编辑成功", m)
}

// bForRecheck 申请复核
func (c *BaseOrderController) bForRecheck(id int64) {
	m, err := models.OrderOne(id, "")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}
	if err = UpdateOrderStatus(m, "待复核", true); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}
	if err := models.OrderUpdateStatus(m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}
	orderRecord := c.newOrderRecord(m, "复核")
	if err := models.OrderRecordSave(orderRecord); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}
	c.jsonResult(enums.JRCodeSucc, "操作成功", m)
}

// bRestart 重新开启
func (c *BaseOrderController) bRestart(id int64) {
	m, err := models.OrderOne(id, "")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}
	if err = UpdateOrderStatus(m, "待审核", true); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}
	if err := models.OrderUpdateStatus(m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}
	orderRecord := c.newOrderRecord(m, "重新开启订单")
	if err := models.OrderRecordSave(orderRecord); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}
	c.jsonResult(enums.JRCodeSucc, "操作成功", m)
}

// bReForRecheck 重新申请复核
func (c *BaseOrderController) bReForRecheck(id int64) {
	m, err := models.OrderOne(id, "")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}
	if err = UpdateOrderStatus(m, "待复核", true); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}
	if err := models.OrderUpdateStatus(m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}
	orderRecord := c.newOrderRecord(m, "复核")
	if err := models.OrderRecordSave(orderRecord); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}
	c.jsonResult(enums.JRCodeSucc, "操作成功", m)
}

// bRecheckPass 通过复核、驳回
func (c *BaseOrderController) bRecheckPassReject(statusString, action, actionName string) {
	Id, _ := c.GetInt64(":id", 0)
	m, err := models.OrderOne(Id, "")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}
	c.validRequestData(m)
	if err = UpdateOrderStatus(m, statusString, false); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}
	recheckErrorInputIds := c.GetString("RecheckErrorInputIds")
	itemRecheckErrorInputIds := c.GetString("ItemRecheckErrorInputIds")
	m.RecheckErrorInputIds = recheckErrorInputIds
	m.ItemRecheckErrorInputIds = itemRecheckErrorInputIds
	if err := models.OrderUpdateStatusRecheckErrorInputIds(m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}
	content := c.GetString("Content")
	remark := c.GetString("Remark")
	if len(content) > 0 {
		statusString += ":" + content
	}
	orderRecord := c.newOrderRecord(m, statusString)
	orderRecord.Remark = remark
	if err := models.OrderRecordSave(orderRecord); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}
	if err := c.setAnnotaionUserRelType(m, nil, "复核人"); err != nil {
		c.jsonResult(enums.JRCodeFailed, "添加失败", m)
	}
	// 生成 pdf 凭证
	// basic auth 认证用户名和密码
	username, _ := models.GetSettingValueByKey("pdf_username")
	password, _ := models.GetSettingValueByKey("pdf_password")
	pdfData := enums.PdfData{
		m.Id,
		m.ClientSeqNo,
		"order_recheck_pdf",
		action,
		"order",
		"order_recheck_pdf_header",
		username,
		password,
		30,
	}
	if ffp, err := enums.NewPDFGenerator(&pdfData); err != nil {
		c.jsonResult(enums.JRCodeFailed, "添加失败", m)
	} else {
		aFile := models.NewOrderFile(0)
		aFile.EdocCopUrl = strings.Replace(ffp, ".", "", 1)
		aFile.EdocCode = actionName
		aFile.EdocCodeName = actionName
		aFile.Creator = c.curUser.RealName
		aFile.Order = m
		aFile.Version = 1.0
		err = models.OrderFileSaveOrUpdate(&aFile)
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, "添加失败", err)
		}
	}
	c.jsonResult(enums.JRCodeSucc, "操作成功", m)
}

// bRecheck 复核
func (c *BaseOrderController) bRecheck(id int64) {
	m, err := models.OrderOne(id, "")
	if err != nil {
		c.pageError("数据无效，请刷新后重试")
	}
	if m != nil {
		c.getActionData(m.IEFlag, "RecheckPass", "RecheckReject")
	}
	c.setStatusOnly(m, "复核中", false)
	order := models.TransformOrder(id, "OrderItems,OrderContainers,OrderDocuments", true)
	c.Data["m"] = order
	c.Data["IEFlagName"] = enums.GetImpexpMarkcdCNName(m.IEFlag)
	c.setTpl("order/recheck.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "order/recheck_footerjs.html"
	// 页面里按钮权限控制
	c.GetXSRFToken()
}

// bPrint 打印
func (c *BaseOrderController) bPrint(id int64) {
	m, err := models.OrderOne(id, "OrderItems,OrderContainers,OrderDocuments")
	if err != nil {
		c.pageError("数据无效，请刷新后重试")
	}
	if m != nil {
		// 生成
		username, _ := models.GetSettingValueByKey("pdf_username")
		password, _ := models.GetSettingValueByKey("pdf_password")
		pdfData := enums.PdfData{
			m.Id,
			m.ClientSeqNo,
			"order_pdf",
			"report",
			"order",
			"order_pdf_header",
			username,
			password,
			30,
		}
		if ffp, err := enums.NewPDFGenerator(&pdfData); err != nil {
			c.jsonResult(enums.JRCodeFailed, "添加失败", m)
		} else {
			c.Data["json"] = strings.Replace(ffp, ".", "", 1)
		}
	}
	c.ServeJSON()
}

// bPushXml 提交单一
func (c *BaseOrderController) bPushXml(id int64) {
	m, err := models.OrderOne(id, "OrderItems,OrderContainers,OrderDocuments")
	if err != nil || m == nil {
		c.pageError("数据无效，请刷新后重试")
	} else {
		/*报文对象*/
		decMessage := &xmlTemplate.DecMessage{
			Version: "3.1",
			Xmlns:   "http://www.chinaport.gov.cn/dec",
		}

		decHead := xmlTemplate.DecHead{}
		enums.SetObjValueFromObj(&decHead, m) // 设置数据到 xml 结构体

		gName := xmlTemplate.Cdata{Value: m.NoteS}
		decHead.NoteS = gName

		decLists := xmlTemplate.DecLists{}
		var decListsl []xmlTemplate.DecList
		for _, dl := range m.OrderItems {
			decList := xmlTemplate.DecList{}
			enums.SetObjValueFromObj(&decList, dl) // 设置数据到 xml 结构体

			gName := xmlTemplate.Cdata{Value: dl.GName}
			decList.GName = gName

			gModel := xmlTemplate.Cdata{Value: dl.GModel}
			decList.GModel = gModel

			ciqName := xmlTemplate.Cdata{Value: dl.CiqName}
			decList.CiqName = ciqName

			decGoodsLimits := xmlTemplate.DecGoodsLimits{}
			var decGoodsLimitsl []xmlTemplate.DecGoodsLimit
			for _, oil := range dl.OrderItemLimits {
				decGoodsLimit := xmlTemplate.DecGoodsLimit{}
				enums.SetObjValueFromObj(&decGoodsLimit, oil) // 设置数据到 xml 结构体

				var decGoodsLimitVins []xmlTemplate.DecGoodsLimitVin
				for _, oilv := range oil.OrderItemLimitVins {
					decGoodsLimitVin := xmlTemplate.DecGoodsLimitVin{}
					enums.SetObjValueFromObj(&decGoodsLimitVin, oilv) // 设置数据到 xml 结构体
					decGoodsLimitVins = append(decGoodsLimitVins, decGoodsLimitVin)
				}
				decGoodsLimit.DecGoodsLimitVin = decGoodsLimitVins
				decGoodsLimitsl = append(decGoodsLimitsl, decGoodsLimit)
			}

			decGoodsLimits.DecGoodsLimit = decGoodsLimitsl
			decList.DecGoodsLimits = decGoodsLimits
			decListsl = append(decListsl, decList)
		}

		decLists.DecList = decListsl

		decLicenseDocus := xmlTemplate.DecLicenseDocus{}
		var licenseDocusl []xmlTemplate.LicenseDocu
		for _, odec := range m.OrderDocuments {
			decLicenseDocu := xmlTemplate.LicenseDocu{}
			enums.SetObjValueFromObj(&decLicenseDocu, odec) // 设置数据到 xml 结构体
			licenseDocusl = append(licenseDocusl, decLicenseDocu)
		}
		decLicenseDocus.LicenseDocu = licenseDocusl

		decContainers := xmlTemplate.DecContainers{}
		var decContainersl []xmlTemplate.DecContainer
		for _, oc := range m.OrderContainers {
			decContainer := xmlTemplate.DecContainer{}
			enums.SetObjValueFromObj(&decContainer, oc) // 设置数据到 xml 结构体
			decContainersl = append(decContainersl, decContainer)
		}
		decContainers.DecContainer = decContainersl

		decSign := xmlTemplate.DecSign{}
		decFreeTxt := xmlTemplate.DecFreeTxt{}

		var ecoRelations []xmlTemplate.EcoRelation

		decRequestCerts := xmlTemplate.DecRequestCerts{}
		var decRequestCertsl []xmlTemplate.DecRequestCert
		for _, odrc := range m.DecRequestCerts {
			decRequestCert := xmlTemplate.DecRequestCert{}
			enums.SetObjValueFromObj(&decRequestCert, odrc) // 设置数据到 xml 结构体
			decRequestCertsl = append(decRequestCertsl, decRequestCert)
		}
		decRequestCerts.DecRequestCert = decRequestCertsl

		decOtherPacks := xmlTemplate.DecOtherPacks{}
		var decOtherPacksl []xmlTemplate.DecOtherPack
		for _, odop := range m.DecOtherPacks {
			decOtherPack := xmlTemplate.DecOtherPack{}
			enums.SetObjValueFromObj(&decOtherPack, odop) // 设置数据到 xml 结构体
			decOtherPacksl = append(decOtherPacksl, decOtherPack)
		}
		decOtherPacks.DecOtherPack = decOtherPacksl

		decCopLimits := xmlTemplate.DecCopLimits{}
		var decCopLimitsl []xmlTemplate.DecCopLimit
		decCopLimits.DecCopLimit = decCopLimitsl

		decUsers := xmlTemplate.DecUsers{}
		var decUsersl []xmlTemplate.DecUser
		for _, odu := range m.DecUsers {
			decUser := xmlTemplate.DecUser{}
			enums.SetObjValueFromObj(&decUser, odu) // 设置数据到 xml 结构体
			decUsersl = append(decUsersl, decUser)
		}
		decUsers.DecUser = decUsersl

		decCopPromises := xmlTemplate.DecCopPromises{}
		decCopPromise := xmlTemplate.DecCopPromise{}
		decCopPromises.DecCopPromise = decCopPromise

		var edocRealations []xmlTemplate.EdocRealation
		for _, odr := range m.OrderFiles {
			edocCodes, _ := models.GetSettingRValueByKey("sendEdocCodes", false)
			if enums.InStringMap(odr.EdocCode, edocCodes) {
				edocRealation := xmlTemplate.EdocRealation{}
				enums.SetObjValueFromObj(&edocRealation, odr) // 设置数据到 xml 结构体
				edocRealations = append(edocRealations, edocRealation)
			}
		}

		decMessage.DecHead = decHead
		decMessage.DecLists = decLists
		decMessage.DecLicenseDocus = decLicenseDocus
		decMessage.DecContainers = decContainers
		decMessage.DecSign = decSign
		decMessage.DecFreeTxt = decFreeTxt
		decMessage.EcoRelation = ecoRelations
		decMessage.DecRequestCerts = decRequestCerts
		decMessage.DecOtherPacks = decOtherPacks
		decMessage.DecCopLimits = decCopLimits
		decMessage.DecUsers = decUsers
		decMessage.DecCopPromises = decCopPromises
		decMessage.EdocRealation = edocRealations

		path, _ := models.GetSettingValueByKey("order_xml_path")
		pathTemp := "./static/generate/order/" + strconv.FormatInt(id, 10) + "/temp/"
		// 报文名称
		mName := time.Now().Format(enums.BaseDateTimeSecondFormat) + "_" + m.ClientSeqNo
		fileName := mName + ".xml"

		output, err := xml.MarshalIndent(decMessage, "", "")
		if err != nil {
			utils.LogDebug(fmt.Sprintf("MarshalIndent error:%v", err))
			c.jsonResult(enums.JRCodeFailed, "操作失败", nil)
		}

		if err := file.CreateFile(pathTemp); err != nil {
			utils.LogDebug(fmt.Sprintf("文件夹创建失败:%v", err))
			c.jsonResult(enums.JRCodeFailed, "操作失败", nil)
		}

		err = file.WriteFile(pathTemp+fileName, []byte(xml.Header))
		if err != nil {
			utils.LogDebug(fmt.Sprintf("WriteFile error:%v", err))
			c.jsonResult(enums.JRCodeFailed, "操作失败", nil)
		}

		err = file.AppendFile(pathTemp+fileName, output)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("WriteFile error:%v", err))
			c.jsonResult(enums.JRCodeFailed, "操作失败", nil)
		}

		f1, err := os.Open(pathTemp + fileName)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("os.Open error:%v", err))
			c.jsonResult(enums.JRCodeFailed, "操作失败", nil)
		}
		defer f1.Close()

		var files = []*os.File{f1}
		err = file.Compress(files, path+mName+".zip")
		if err != nil {
			utils.LogDebug(fmt.Sprintf("file.Compress error:%v", err))
			c.jsonResult(enums.JRCodeFailed, "操作失败", nil)
		}

		err = os.Remove(pathTemp + fileName)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("os.Remove error:%v", err))
			c.jsonResult(enums.JRCodeFailed, "操作失败", nil)
		}

		c.jsonResult(enums.JRCodeSucc, "操作成功", nil)
	}
}

// bRemark 附注
func (c *BaseOrderController) bRemark(id int64, remark string) {
	m, err := models.OrderOne(id, "")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", err)
	}
	m.Remark = remark
	if err = models.OrderUpdateRemark(m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", err)
	}
	c.jsonResult(enums.JRCodeSucc, "操作成功", err)
}

// bAuditFirstRejectLog 驳回原因
func (c *BaseOrderController) bAuditFirstRejectLog(id int64) {
	m, err := models.OrderOne(id, "")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", err)
	}

	aRecord := models.NewOrderRecord(0)
	aRecord.Order = m
	aStatusS, err := c.getOrderStatus("orderStatus")
	sSting, err, _ := enums.TransformIntToCn(aStatusS, m.Status)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", err)
	}

	if err = models.OrderRecordOneByStatusAndOrderId(&aRecord, sSting); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", err)
	}

	c.jsonResult(enums.JRCodeSucc, "操作成功", aRecord)
}

// 删除
func (c *BaseOrderController) bDelete(id int64) {
	m, err := models.OrderOne(id, "")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "删除失败", err)
	}
	if _, err := models.OrderDelete(id); err != nil {
		c.jsonResult(enums.JRCodeFailed, "删除失败", err)
	}
	orderRecord := c.newOrderRecord(m, "删除订单")
	if err := models.OrderRecordSave(orderRecord); err != nil {
		c.jsonResult(enums.JRCodeFailed, "删除失败", m)
	}

	c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", 1), "")
}

// 清单订单号
func (c *BaseOrderController) getClientSeqNo(iEFlag, customMasterName string) string {
	return iEFlag + customMasterName + time.Now().Format(enums.BaseDateTimeSecondFormat) + enums.CreateCaptcha()
}

// 更新状态和状态更新时间
func (c *BaseOrderController) setAnnotaionUserRelType(m *models.Order, bu *models.BackendUser, userTypes string) error {
	rs := strings.Split(userTypes, ",")
	for _, v := range rs {
		aur := models.NewOrderUserRel(0)
		aStatusS, err := c.getOrderStatus("orderUserType")
		aStatus, err, _ := enums.TransformCnToInt(aStatusS, v)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("转换制单人类型出错:%v", err))
			return err
		}
		// 除了制单人，其他人都是当前用户
		if v == "制单人" && bu != nil {
			aur.BackendUser = bu
		} else {
			aur.BackendUser = &c.curUser
		}
		aur.Order = m
		aur.UserType = aStatus
		if err = models.OrderUserRelSave(&aur); err != nil {
			return err
		}
	}
	return nil
}

// TransformOrderList 格式化列表数据
func (c *BaseOrderController) TransformOrderList(ms []*models.Order) []*map[string]interface{} {

	var orderList []*map[string]interface{}
	for _, v := range ms {
		orderCreatorName := "" // 制单人
		orderItem := make(map[string]interface{})
		aStatusS, err := c.getOrderStatus("orderStatus")
		aStatus, err, _ := enums.TransformIntToCn(aStatusS, v.Status)
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, "获取状态转中文出错", nil)
		}
		userTypeS, err := c.getOrderStatus("orderUserType")
		userType, err, _ := enums.TransformCnToInt(userTypeS, "制单人")
		if err != nil {
			utils.LogDebug(fmt.Sprintf("转换制单人类型出错:%v", err))
		}
		if len(v.BackendUsers) > 0 {
			for _, bu := range v.BackendUsers {
				abur, err := models.OrderUserRelByUserIdAndOrderId(bu.Id, v.Id, userType)
				if err != nil {
					c.jsonResult(enums.JRCodeFailed, "获取制单人出错", nil)
				}
				if abur != nil && abur.Id != 0 {
					orderCreatorName = bu.RealName
				}
			}
		}

		// 集装箱ids
		containerIds := ""
		if len(v.OrderContainers) > 0 {
			for _, oc := range v.OrderContainers {
				containerIds += oc.ContainerId + ","
			}
		}

		// 清单编号
		annotationBondInvtNo := ""
		annotationId := 0
		if v.Annotation != nil {
			annotationBondInvtNo = v.Annotation.BondInvtNo
			annotationId = int(v.Annotation.Id)
		}
		orderItem["Id"] = strconv.FormatInt(v.Id, 10)
		orderItem["StatusString"] = aStatus
		orderItem["ManualNo"] = v.ManualNo
		orderItem["ContrNo"] = v.ContrNo
		orderItem["IEPortName"] = v.IEPortName
		orderItem["IEPort"] = v.IEPort
		orderItem["EntryId"] = v.EntryId
		orderItem["AnnotationBondInvtNo"] = annotationBondInvtNo
		orderItem["AnnotationId"] = annotationId
		orderItem["ClientSeqNo"] = v.ClientSeqNo
		orderItem["BillNo"] = v.BillNo
		orderItem["ContainerId"] = containerIds // 集装箱ids
		orderItem["CompanyName"] = v.Company.Name
		orderItem["DeclareName"] = orderCreatorName
		orderItem["AplDate"] = v.AplDate.Format(enums.BaseDateTimeFormat)
		orderItem["CreatedAt"] = v.CreatedAt.Format(enums.BaseDateTimeFormat)
		orderList = append(orderList, &orderItem)
	}
	return orderList
}

// 仅仅更新状态
func (c *BaseOrderController) setStatusOnly(m *models.Order, statusString string, isRestart bool) {
	if err := UpdateOrderStatus(m, statusString, isRestart); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", nil)
	}
	if err := models.OrderUpdateStatus(m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}
}

// 操作记录
func (c *BaseOrderController) newOrderRecord(m *models.Order, content string) *models.OrderRecord {
	aStatusS, _ := c.getOrderStatus("orderStatus")
	statusString, _, _ := enums.TransformIntToCn(aStatusS, m.Status)
	orderRecord := models.NewOrderRecord(0)
	orderRecord.Content = content
	orderRecord.BackendUser = &c.curUser
	orderRecord.Status = statusString
	orderRecord.Order = m
	return &orderRecord
}

// 转换订单状态和订单用户状态
func (c *BaseOrderController) getOrderStatus(status string) (map[string]string, error) {
	return models.GetSettingRValueByKey(status, false)
}

// 是否能保存
func (c *BaseOrderController) getCanStore(m *models.Order, ieflag string) bool {
	if models.IsSuperAdmin(c.curUser.Id) {
		return true
	}
	if m == nil {
		return c.checkActionAuthor(c.controllerName, ieflag+"Audit")
	} else {
		aStatusS, _ := c.getOrderStatus("orderStatus")
		aStatus, _, _ := enums.TransformIntToCn(aStatusS, m.Status)
		if c.checkActionAuthor(c.controllerName, m.IEFlag+"Audit") {
			if aStatus == "待审核" || aStatus == "审核中" {
				return true
			}
		} else if c.checkActionAuthor(c.controllerName, m.IEFlag+"Make") {
			if aStatus == "待制单" || aStatus == "制单中" {
				return true
			}
		}
	}
	return false
}
