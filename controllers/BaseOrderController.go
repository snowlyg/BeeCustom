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
	"github.com/astaxie/beego"
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
	m, err := models.OrderOne(id, "Company,OrderItems")
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
	m, err := models.OrderOne(id, "")
	if err != nil {
		c.pageError("数据无效，请刷新后重试")
	}
	c.setStatusOnly(m, "审核中", false)
	// 获取制单人
	backendUsers := models.GetCreateBackendUsers("OrderController.Make")
	c.Data["BackendUsers"] = backendUsers
	c.Data["m"] = models.TransformOrder(id, "OrderItems")
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
	c.Data["m"] = models.TransformOrder(id, "OrderItems")
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

	c.Data["m"] = models.TransformOrder(id, "OrderItems,OrderRecords")
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
	m, err := models.OrderOne(id, "Company,OrderItems")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}
	// 获取form里的值
	if err := c.ParseForm(m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "ParseForm", m)
	}

	if m != nil {
		m.AplDate = time.Now()
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
	m, err := models.OrderOne(id, "Company,OrderItems")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}
	if err = UpdateOrderStatus(m, "待复核", false); err != nil {
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
	m, err := models.OrderOne(id, "Company,OrderItems")
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
	m, err := models.OrderOne(id, "Company,OrderItems")
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
	m, err := models.OrderOne(Id, "Company,OrderItems")
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
	if ffp, err := enums.NewPDFGenerator(m.Id, m.ClientSeqNo, "order_recheck_pdf", action); err != nil {
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
	m, err := models.OrderOne(id, "OrderItems")
	if err != nil {
		c.pageError("数据无效，请刷新后重试")
	}
	if m != nil {
		c.getActionData(m.IEFlag, "RecheckPass", "RecheckReject")
	}
	c.setStatusOnly(m, "复核中", false)
	order := models.TransformOrder(id, "OrderItems")
	c.Data["m"] = order
	c.setTpl("order/recheck.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "order/recheck_footerjs.html"
	// 页面里按钮权限控制
	c.GetXSRFToken()
}

// bPrint 打印
func (c *BaseOrderController) bPrint(id int64) {
	m, err := models.OrderOne(id, "OrderItems")
	if err != nil {
		c.pageError("数据无效，请刷新后重试")
	}
	if m != nil {
		// 生成 pdf 凭证
		if ffp, err := enums.NewPDFGenerator(m.Id, m.ClientSeqNo, "order_pdf", "report"); err != nil {
			c.jsonResult(enums.JRCodeFailed, "添加失败", m)
		} else {
			c.Data["json"] = strings.Replace(ffp, ".", "", 1)
		}
	}
	c.ServeJSON()
}

// bPushXml 提交单一
func (c *BaseOrderController) bPushXml(id int64) {
	m, err := models.OrderOne(id, "OrderItems")
	if err != nil || m == nil {
		c.pageError("数据无效，请刷新后重试")
	} else {
		/*清单报文对象*/
		signature := &xmlTemplate.Signature{
			Xmlns: "http://www.w3.org/2001/XMLSchema-instance",
			SignedInfo: xmlTemplate.SignedInfo{
				CanonicalizationMethod: xmlTemplate.CanonicalizationMethod{
					Algorithm: "http://www.w3.org/TR/2001/REC-xml-c14n-20010315",
				},
				SignatureMethod: xmlTemplate.SignatureMethod{
					Algorithm: "http://www.w3.org/2001/04/xmldsig-more#rsa-md5",
				},
				Reference: xmlTemplate.Reference{
					URI: "String",
					DigestMethod: xmlTemplate.DigestMethod{
						Algorithm: "http://www.w3.org/2000/09/xmldsig#sha1",
					},
				},
			},
		}
		handBook, _ := models.HandBookOne(m.HandBookId, "")
		var sysId string
		var receiverId string
		var path string
		handBookType1, _ := enums.GetSectionWithString("手册", "hand_book_type")
		handBookType2, _ := enums.GetSectionWithString("账册", "hand_book_type")
		if handBook == nil {
			c.jsonResult(enums.JRCodeFailed, "错误手账册类型", nil)
		} else {
			if handBook.Type == handBookType1 {
				receiverId = beego.AppConfig.String("OrderReceiverIdC")
				sysId = beego.AppConfig.String("OrderSysIdC")
				path = beego.AppConfig.String("order_xml_path_c")
			} else if handBook.Type == handBookType2 {
				receiverId = beego.AppConfig.String("OrderReceiverIdE")
				sysId = beego.AppConfig.String("OrderSysIdE")
				path = beego.AppConfig.String("order_xml_path_e")
			} else {
				c.jsonResult(enums.JRCodeFailed, "错误手账册类型", nil)
			}
		}

		pathTemp := "./static/generate/order/" + strconv.FormatInt(id, 10) + "/temp/"
		// 报文名称
		mName := time.Now().Format(enums.BaseDateTimeSecondFormat) + "__" + m.ClientSeqNo
		fileName := mName + ".xml"

		signature.Object.Package.EnvelopInfo.FileName = fileName
		signature.Object.Package.EnvelopInfo.Version = beego.AppConfig.String("OrderVersion")
		signature.Object.Package.EnvelopInfo.BusinessId = beego.AppConfig.String("OrderBusinessId")
		signature.Object.Package.EnvelopInfo.MessageType = beego.AppConfig.String("OrderMessageType")
		signature.Object.Package.EnvelopInfo.SenderId = beego.AppConfig.String("OrderSenderId")
		signature.Object.Package.EnvelopInfo.ReceiverId = receiverId
		signature.Object.Package.EnvelopInfo.MessageId = m.ClientSeqNo
		signature.Object.Package.EnvelopInfo.SendTime = time.Now().Format(enums.RFC3339)
		signature.Object.Package.DataInfo.BussinessData.DelcareFlag = "0" // 0:暂存，1:申报
		signature.Object.Package.DataInfo.BussinessData.InvtMessage.SysId = sysId
		signature.Object.Package.DataInfo.BussinessData.InvtMessage.OperCusRegCode = beego.AppConfig.String("AgentCode")

		invtHeadType := xmlTemplate.InvtHeadType{
			SeqNo: m.SeqNo,
			//BondInvtNo:                   m.BondInvtNo,
			//ChgTmsCntstring:              m.ChgTmsCnt,
			//PutrecNostring:               m.PutrecNo,
			//InvtTypestring:               m.InvtType,
			//ClientSeqNostring:        m.ClientSeqNo,
			//BizopEtpsnostring:            m.BizopEtpsno,
			//BizopEtpsSccdstring:          m.BizopEtpsSccd,
			//BizopEtpsNmstring:            m.BizopEtpsNm,
			//RcvgdEtpsnostring:            m.RcvgdEtpsno,
			//RvsngdEtpsSccdstring:         m.RvsngdEtpsSccd,
			//RcvgdEtpsNmstring:            m.RcvgdEtpsNm,
			//DclEtpsnostring:              m.DclEtpsno,
			//DclEtpsSccdstring:            m.DclEtpsSccd,
			//DclEtpsNmstring:              m.DclEtpsNm,
			//InputCodestring:              m.InputCode,
			//InputCreditCodestring:        m.InputCreditCode,
			//InputNamestring:              m.InputName,
			//RltInvtNostring:              m.RltInvtNo,
			//RltPutrecNostring:            m.RltPutrecNo,
			//CorrEntryDclEtpsNostring:     m.CorrEntryDclEtpsNo,
			//CorrEntryDclEtpsSccdstring:   m.CorrEntryDclEtpsSccd,
			//CorrEntryDclEtpsNmstring:     m.CorrEntryDclEtpsNm,
			//RltEntryBizopEtpsnostring:    m.RltEntryBizopEtpsno,
			//RltEntryBizopEtpsSccdstring:  m.RltEntryBizopEtpsSccd,
			//RltEntryBizopEtpsNmstring:    m.RltEntryBizopEtpsNm,
			//RltEntryRcvgdEtpsnostring:    m.RltEntryRcvgdEtpsno,
			//RltEntryRvsngdEtpsSccdstring: m.RltEntryRvsngdEtpsSccd,
			//RltEntryRcvgdEtpsNmstring:    m.RltEntryRcvgdEtpsNm,
			//RltEntryDclEtpsnostring:      m.RltEntryDclEtpsno,
			//RltEntryDclEtpsSccdstring:    m.RltEntryDclEtpsSccd,
			//RltEntryDclEtpsNmstring:      m.RltEntryDclEtpsNm,
			//ImpexpPortcdstring:           m.ImpexpPortcd,
			//DclPlcCuscdstring:            m.DclPlcCuscd,
			//IEFlagstring:           m.IEFlag,
			//MtpckEndprdMarkcdstring:      m.MtpckEndprdMarkcd,
			//SupvModecdstring:             m.SupvModecd,
			//TrspModecdstring:             m.TrspModecd,
			//ApplyNostring:                m.ApplyNo,
			//ListTypestring:               m.ListType,
			//DclcusFlagstring:             m.DclcusFlag,
			//DclcusTypecdstring:           m.DclcusTypecd,
			//IcCardNostring:               beego.AppConfig.String("ICCode"),
			//DecTypestring:                m.DecType,
			//Rmkstring:                    m.Rmk,
			//StshipTrsarvNatcdstring:      m.StshipTrsarvNatcd,
			//EntryNostring:                m.EntryNo,
			//RltEntryNostring:             m.RltEntryNo,
			//DclTypecdstring:              m.DclTypecd,
			//GenDecFlagstring:             m.GenDecFlag,
		}

		invtListTypes := []xmlTemplate.InvtListType{}
		//for _, v := range m.OrderItems {
		//	invtListType := xmlTemplate.InvtListType{
		//		SeqNo:            v.SeqNo,
		//		GdsSeqno:         strconv.Itoa(v.GdsSeqno),
		//		PutrecSeqno:      strconv.Itoa(v.PutrecSeqno),
		//		GdsMtno:          v.GdsMtno,
		//		Gdecd:            v.Gdecd,
		//		GdsNm:            v.GdsNm,
		//		GdsSpcfModelDesc: v.GdsSpcfModelDesc,
		//		DclUnitcd:        v.DclUnitcd,
		//		LawfUnitcd:       v.LawfUnitcd,
		//		SecdLawfUnitcd:   v.SecdLawfUnitcd,
		//		Natcd:            v.Natcd,
		//		DclUprcAmt:       strconv.FormatFloat(v.DclUprcAmt, 'f', 4, 64),
		//		DclTotalAmt:      strconv.FormatFloat(v.DclTotalAmt, 'f', 2, 64),
		//		UsdStatTotalAmt:  strconv.FormatFloat(v.UsdStatTotalAmt, 'f', 5, 64),
		//		DclCurrcd:        v.DclCurrcd,
		//		LawfQty:          strconv.FormatFloat(v.LawfQty, 'f', 5, 64),
		//		SecdLawfQty:      strconv.FormatFloat(v.SecdLawfQty, 'f', 5, 64),
		//		WtSfVal:          strconv.FormatFloat(v.WtSfVal, 'f', 5, 64),
		//		FstSfVal:         strconv.FormatFloat(v.FstSfVal, 'f', 5, 64),
		//		SecdSfVal:        strconv.FormatFloat(v.SecdSfVal, 'f', 5, 64),
		//		DclQty:           strconv.FormatFloat(v.DclQty, 'f', 5, 64),
		//		GrossWt:          strconv.FormatFloat(v.GrossWt, 'f', 5, 64),
		//		NetWt:            strconv.FormatFloat(v.NetWt, 'f', 5, 64),
		//		UseCd:            v.UseCd,
		//		LvyrlfModecd:     v.LvyrlfModecd,
		//		UcnsVerno:        v.UcnsVerno,
		//		ClyMarkcd:        v.ClyMarkcd,
		//		EntryGdsSeqno:    strconv.Itoa(v.EntryGdsSeqno),
		//		ApplyTbSeqno:     strconv.Itoa(v.ApplyTbSeqno),
		//		DestinationNatcd: v.DestinationNatcd,
		//		ModfMarkcd:       v.ModfMarkcd,
		//		Rmk:              v.Rmk,
		//	}
		//
		//	invtListTypes = append(invtListTypes, invtListType)
		//}

		invtDecHeadType := xmlTemplate.InvtDecHeadType{}
		invtDecListType := []xmlTemplate.InvtDecListType{}

		signature.Object.Package.DataInfo.BussinessData.InvtMessage.InvtHeadType = invtHeadType
		signature.Object.Package.DataInfo.BussinessData.InvtMessage.InvtListType = invtListTypes
		signature.Object.Package.DataInfo.BussinessData.InvtMessage.InvtDecHeadType = invtDecHeadType
		signature.Object.Package.DataInfo.BussinessData.InvtMessage.InvtDecListType = invtDecListType

		output, err := xml.MarshalIndent(signature, "", "")
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

		err = file.WriteFile(pathTemp+fileName, output)
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

	sSting, err := enums.GetSectionWithInt(m.Status, "order_status")
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
		aStatus, err := enums.GetSectionWithString(v, "order_user_type")
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
	var orderCreatorName string // 制单人
	var orderList []*map[string]interface{}
	for _, v := range ms {
		orderItem := make(map[string]interface{})
		aStatus, err := enums.GetSectionWithInt(v.Status, "order_status")
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, "获取状态转中文出错", nil)
		}
		userType, err := enums.GetSectionWithString("制单人", "order_user_type")
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
	statusString, _ := enums.GetSectionWithInt(m.Status, "order_status")
	orderRecord := models.NewOrderRecord(0)
	orderRecord.Content = content
	orderRecord.BackendUser = &c.curUser
	orderRecord.Status = statusString
	orderRecord.Order = m
	return &orderRecord
}

// 是否能保存
func (c *BaseOrderController) getCanStore(m *models.Order, ieflag string) bool {
	if models.IsSuperAdmin(c.curUser.Id) {
		return true
	}
	if m == nil {
		return c.checkActionAuthor(c.controllerName, ieflag+"Audit")
	} else {
		if c.checkActionAuthor(c.controllerName, m.IEFlag+"Audit") {
			aStatus, _ := enums.GetSectionWithInt(m.Status, "order_status")
			if aStatus == "待审核" || aStatus == "审核中" {
				return true
			}
		} else if c.checkActionAuthor(c.controllerName, m.IEFlag+"Make") {
			aStatus, _ := enums.GetSectionWithInt(m.Status, "order_status")
			if aStatus == "待制单" || aStatus == "制单中" {
				return true
			}
		}
	}
	return false
}
