package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"BeeCustom/enums"
	"BeeCustom/models"
	"BeeCustom/utils"
)

type BaseAnnotationController struct {
	BaseController
}

// 列表数据
func (c *BaseAnnotationController) bDataGrid(impexpMarkcd string) {
	// 直接获取参数 getDataGridData()
	params := models.NewAnnotationQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	params.ImpexpMarkcd = impexpMarkcd

	// 获取数据列表和总数
	data, total, err := models.AnnotationPageList(&params)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据列表和总数失败", nil)
	}

	err = models.AnnotationGetRelations(data, "Company,BackendUsers")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "关联关系获取失败", nil)
	}

	// 格式化数据
	annotationList := c.TransformAnnotationList(data)
	c.ResponseList(annotationList, total)

	c.ServeJSON()
}

func (c *BaseAnnotationController) bIndex(impexpMarkcd, impexpMarkcdName string) {
	// 页面模板设置
	c.setTpl("annotation/index.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "annotation/index_footerjs.html"

	// 页面里按钮权限控制
	c.getActionData(impexpMarkcd, "Index", "Create", "Edit", "Make", "Audit", "Delete", "Distribute", "Recheck", "Push", "PushXml", "StoreError", "Change", "Restart", "Cancel")

	// 获取制单人
	backendUsers := models.GetCreateBackendUsers("AnnotationController.Make")
	c.Data["BackendUsers"] = backendUsers
	c.Data["ImpexpMarkcd"] = impexpMarkcd
	c.Data["ImpexpMarkcdName"] = impexpMarkcdName

	c.GetXSRFToken()
}

// 数据统计
func (c *BaseAnnotationController) bStatusCount(impexpMarkcd string) {
	// 直接获取参数 getDataGridData()
	params := models.NewAnnotationQueryParam()
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	params.ImpexpMarkcd = impexpMarkcd

	// 获取数据列表和总数
	data, err := models.AnnotationStatusCount(&params)
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
func (c *BaseAnnotationController) bCreate(impexpMarkcd string) {
	c.Data["ImpexpMarkcd"] = impexpMarkcd
	c.Data["canStore"] = c.getCanStore(nil, impexpMarkcd)
	c.getResponses(impexpMarkcd)
}

// Store 添加 新建 页面
func (c *BaseAnnotationController) bStore(impexpMarkcd string) {
	m := models.NewAnnotation(0)
	// 获取form里的值
	if err := c.ParseForm(&m); err != nil {
		utils.LogDebug(fmt.Sprintf("ParseForm:%v", err))
		c.jsonResult(enums.JRCodeFailed, "获取数据出错", m)
	}
	company, err := models.CompanyByManageCode(m.BizopEtpsno)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取客户出错", nil)
	}

	if err = c.updateAnnotationStatus(&m, "待审核"); err != nil {
		c.jsonResult(enums.JRCodeFailed, "添加失败", nil)
	}

	m.Company = company
	m.InputTime = time.Now()
	m.InvtDclTime = time.Now()
	m.EtpsInnerInvtNo = c.getEtpsInnerInvtNo(impexpMarkcd, m.DclPlcCuscd)

	c.validRequestData(m)

	// valid := validation.Validation{}
	// valid.Required(m.UserPwd, "密码")
	// valid.MinSize(m.UserPwd, 6, "密码")
	// valid.MaxSize(m.UserPwd, 18, "密码")
	//
	// if valid.HasErrors() {
	//	// 如果有错误信息，证明验证没通过
	//	// 打印错误信息
	//	for _, err := range valid.Errors {
	//		c.jsonResult(enums.JRCodeFailed, err.Key+err.Message, m)
	//	}
	// }

	if err := models.AnnotationSave(&m, ""); err != nil {
		c.jsonResult(enums.JRCodeFailed, "添加失败", m)
	} else {
		if err := c.setAnnotaionUserRelType(&m, nil, "创建人"); err != nil {
			c.jsonResult(enums.JRCodeFailed, "添加失败", m)
		}
		annotationRecord := c.newAnnotationRecord(&m, "创建订单")
		if err := models.AnnotationRecordSave(annotationRecord); err != nil {
			c.jsonResult(enums.JRCodeFailed, "添加失败", m)
		}

		c.jsonResult(enums.JRCodeSucc, "添加成功", m)
	}
}

// Edit 添加 编辑 页面
func (c *BaseAnnotationController) bEdit(id int64) {
	m, err := models.AnnotationOne(id, "")
	if err != nil {
		c.pageError("数据无效，请刷新后重试")
	}

	c.setStatusOnly(m, "审核中")

	// 获取制单人
	backendUsers := models.GetCreateBackendUsers("AnnotationController.Make")
	c.Data["BackendUsers"] = backendUsers
	c.Data["m"] = c.TransformAnnotation(m)
	c.Data["canStore"] = c.getCanStore(m, "")
	c.getResponses(m.ImpexpMarkcd)
}

// bMake 添加 编辑 页面
func (c *BaseAnnotationController) bMake(id int64) {
	m, err := models.AnnotationOne(id, "")
	if m != nil && id > 0 {
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}

	c.setStatusOnly(m, "制单中")
	c.Data["m"] = c.TransformAnnotation(m)
	c.Data["canStore"] = c.getCanStore(m, "")
	c.getResponses(m.ImpexpMarkcd)
}

// 编辑相关页面返回
func (c *BaseAnnotationController) getResponses(impexpMarkcd string) {
	var impexpMarkcdName string
	if impexpMarkcd == "I" {
		impexpMarkcdName = "进口"
	} else if impexpMarkcd == "E" {
		impexpMarkcdName = "出口"
	}
	// 页面里按钮权限控制
	c.getActionData(impexpMarkcd, "Audit", "Distribute", "ForRecheck")
	c.Data["ImpexpMarkcdName"] = impexpMarkcdName
	c.setTpl("annotation/change_create_edit_show.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "annotation/create_footerjs.html"
	c.GetXSRFToken()
}

// Cancel 取消订单
func (c *BaseAnnotationController) bCancel(id int64) {
	m, err := models.AnnotationOne(id, "")
	if m != nil && id > 0 {
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}

	if err = c.updateAnnotationStatus(m, "订单关闭"); err != nil {
		c.jsonResult(enums.JRCodeFailed, "取消失败", m)
	}

	if err := models.AnnotationSave(m, ""); err != nil {
		c.jsonResult(enums.JRCodeFailed, "取消失败", m)
	}

	annotationRecord := c.newAnnotationRecord(m, "取消订单")
	if err := models.AnnotationRecordSave(annotationRecord); err != nil {
		c.jsonResult(enums.JRCodeFailed, "取消失败", m)
	}

	c.jsonResult(enums.JRCodeSucc, "取消成功", m)

}

// Audit 审核通过订单
func (c *BaseAnnotationController) bAudit(id int64) {
	m, err := models.AnnotationOne(id, "")
	if m != nil && id > 0 {
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}

	if err := c.setAnnotaionUserRelType(m, nil, "审单人"); err != nil {
		c.jsonResult(enums.JRCodeFailed, "审核失败", m)
	}

	c.setStatusOnly(m, "审核通过")

	annotationRecord := c.newAnnotationRecord(m, "审核订单")
	if err := models.AnnotationRecordSave(annotationRecord); err != nil {
		c.jsonResult(enums.JRCodeFailed, "审核失败", m)
	}

	c.jsonResult(enums.JRCodeSucc, "审核通过", m)

}

// Distribute 分配
func (c *BaseAnnotationController) bDistribute(backendUserId, id int64) {

	bu, err := models.BackendUserOne(backendUserId)
	if bu != nil && backendUserId > 0 {
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}

	m, err := models.AnnotationOne(id, "")
	if m != nil && id > 0 {
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}

	}

	if err = c.setAnnotaionUserRelType(m, bu, "制单人,派单人"); err != nil {
		c.jsonResult(enums.JRCodeFailed, "派单失败", m)
	}

	if err = c.updateAnnotationStatus(m, "待制单"); err != nil {
		c.jsonResult(enums.JRCodeFailed, "派单失败", m)
	}

	if err := models.AnnotationSave(m, ""); err != nil {
		c.jsonResult(enums.JRCodeFailed, "派单失败", m)
	}

	annotationRecord := c.newAnnotationRecord(m, "派单："+bu.RealName)
	if err := models.AnnotationRecordSave(annotationRecord); err != nil {
		c.jsonResult(enums.JRCodeFailed, "派单失败", m)
	}

	c.jsonResult(enums.JRCodeSucc, "派单通过", m)

}

// Update 添加 编辑 页面
func (c *BaseAnnotationController) bUpdate(id int64) {

	m, err := models.AnnotationOne(id, "Company,AnnotationItems")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}

	// 获取form里的值
	if err := c.ParseForm(m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "ParseForm", m)
	}

	m.InputTime = time.Now()
	m.InvtDclTime = time.Now()

	c.validRequestData(m)

	// valid := validation.Validation{}
	// if len(m.UserPwd) > 0 {
	//	valid.MinSize(m.UserPwd, 6, "密码")
	//	valid.MaxSize(m.UserPwd, 18, "密码")
	// }
	//
	// if valid.HasErrors() {
	//	// 如果有错误信息，证明验证没通过
	//	// 打印错误信息
	//	for _, err := range valid.Errors {
	//		c.jsonResult(enums.JRCodeFailed, err.Key+err.Message, m)
	//	}
	// }

	if err := models.AnnotationSave(m, ""); err != nil {
		c.jsonResult(enums.JRCodeFailed, "编辑失败", m)
	}

	annotationRecord := c.newAnnotationRecord(m, "保存数据")
	if err := models.AnnotationRecordSave(annotationRecord); err != nil {
		c.jsonResult(enums.JRCodeFailed, "编辑失败", m)
	}

	c.jsonResult(enums.JRCodeSucc, "编辑成功", m)

}

// bForRecheck 申请复核
func (c *BaseAnnotationController) bForRecheck(id int64) {

	m, err := models.AnnotationOne(id, "Company,AnnotationItems")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}

	if err = c.updateAnnotationStatus(m, "待复核"); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}

	if err := models.AnnotationSave(m, "Status"); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}

	annotationRecord := c.newAnnotationRecord(m, "复核")
	if err := models.AnnotationRecordSave(annotationRecord); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}

	c.jsonResult(enums.JRCodeSucc, "操作成功", m)

}

// bRecheckPass 通过复核、驳回
func (c *BaseAnnotationController) bRecheckPassReject(id int64) {

	m, err := models.AnnotationOne(id, "Company,AnnotationItems")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}

	// 获取form里的值
	if err := c.ParseForm(m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "ParseForm", m)
	}

	c.validRequestData(m)

	// valid := validation.Validation{}
	// if len(m.UserPwd) > 0 {
	//	valid.MinSize(m.UserPwd, 6, "密码")
	//	valid.MaxSize(m.UserPwd, 18, "密码")
	// }
	//
	// if valid.HasErrors() {
	//	// 如果有错误信息，证明验证没通过
	//	// 打印错误信息
	//	for _, err := range valid.Errors {
	//		c.jsonResult(enums.JRCodeFailed, err.Key+err.Message, m)
	//	}
	// }

	if err := models.AnnotationSave(m, ""); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}

	annotationRecord := c.newAnnotationRecord(m, "通过复核")
	if err := models.AnnotationRecordSave(annotationRecord); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}

	if err := c.setAnnotaionUserRelType(m, nil, "复核人"); err != nil {
		c.jsonResult(enums.JRCodeFailed, "添加失败", m)
	}

	if err := enums.NewPDFGenerator(m.Id); err != nil {
		c.jsonResult(enums.JRCodeFailed, "添加失败", m)
	}

	c.jsonResult(enums.JRCodeSucc, "操作成功", m)

}

// bRecheck 复核
func (c *BaseAnnotationController) bRecheck(id int64) {
	m, err := models.AnnotationOne(id, "AnnotationItems")
	if err != nil {
		c.pageError("数据无效，请刷新后重试")
	}

	c.setStatusOnly(m, "复核中")

	annotation := c.TransformAnnotation(m)
	annotation["AnnotationItems"] = m.AnnotationItems
	c.Data["m"] = annotation
	c.setTpl("annotation/recheck.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "annotation/recheck_footerjs.html"

	c.GetXSRFToken()
}

// 删除
func (c *BaseAnnotationController) bDelete(id int64) {
	m, err := models.AnnotationOne(id, "")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "删除失败", err)
	}

	if _, err := models.AnnotationDelete(id); err != nil {
		c.jsonResult(enums.JRCodeFailed, "删除失败", err)
	}

	annotationRecord := c.newAnnotationRecord(m, "删除订单")
	if err := models.AnnotationRecordSave(annotationRecord); err != nil {
		c.jsonResult(enums.JRCodeFailed, "删除失败", m)
	}

	c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", 1), "")

}

// 清单订单号
func (c *BaseAnnotationController) getEtpsInnerInvtNo(iEFlag, customMasterName string) string {
	eiin := "QD" + iEFlag + customMasterName + time.Now().Format(enums.BaseDateTimeSecondFormat) + enums.CreateCaptcha()

	return eiin
}

// 更新状态和状态更新时间
func (c *BaseAnnotationController) updateAnnotationStatus(m *models.Annotation, StatusString string) error {
	aStatus, err := enums.GetSectionWithString(StatusString, "annotation_status")
	if err != nil {
		utils.LogDebug(fmt.Sprintf("转换清单状态出错:%v", err))
		return err
	}

	// 禁止状态回退
	if m.Status < aStatus {
		m.Status = aStatus
		m.StatusUpdatedAt = time.Now()
	}

	return nil
}

// 更新状态和状态更新时间
func (c *BaseAnnotationController) setAnnotaionUserRelType(m *models.Annotation, bu *models.BackendUser, userTypes string) error {
	rs := strings.Split(userTypes, ",")
	for _, v := range rs {
		aur := models.NewAnnotationUserRel(0)
		aStatus, err := enums.GetSectionWithString(v, "annotation_user_type")
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
		aur.Annotation = m
		aur.UserType = aStatus

		if err = models.AnnotationUserRelSave(&aur); err != nil {
			return err
		}
	}
	return nil
}

// TransformAnnotationList 格式化列表数据
func (c *BaseAnnotationController) TransformAnnotationList(ms []*models.Annotation) []*map[string]interface{} {
	var annotationCreatorName string // 制单人
	var annotationList []*map[string]interface{}
	for _, v := range ms {
		annotationItem := make(map[string]interface{})
		aStatus, err := enums.GetSectionWithInt(v.Status, "annotation_status")
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, "获取状态转中文出错", nil)
		}

		userType, err := enums.GetSectionWithString("制单人", "annotation_user_type")
		if err != nil {
			utils.LogDebug(fmt.Sprintf("转换制单人类型出错:%v", err))
		}
		if len(v.BackendUsers) > 0 {
			for _, bu := range v.BackendUsers {
				abur, err := models.AnnotationUserRelByUserIdAndAnnotationId(bu.Id, v.Id, userType)
				if err != nil {
					c.jsonResult(enums.JRCodeFailed, "获取制单人出错", nil)
				}
				if abur.Id != 0 {
					annotationCreatorName = bu.RealName
				}

			}
		}

		annotationItem["Id"] = strconv.FormatInt(v.Id, 10)
		annotationItem["StatusString"] = aStatus
		annotationItem["PutrecNo"] = v.PutrecNo
		annotationItem["ImpexpPortcd"] = v.ImpexpPortcd
		annotationItem["ImpexpPortcdName"] = v.ImpexpPortcdName
		annotationItem["BondInvtNo"] = v.BondInvtNo
		annotationItem["EntryNo"] = v.EntryNo
		annotationItem["SupvModecdName"] = v.SupvModecdName
		annotationItem["TrspModecdName"] = v.TrspModecdName
		annotationItem["InvtDclTime"] = v.InvtDclTime.Format(enums.BaseDateTimeFormat)
		annotationItem["EtpsInnerInvtNo"] = v.EtpsInnerInvtNo
		annotationItem["CompanyName"] = v.Company.Name
		annotationItem["DeclareName"] = annotationCreatorName

		annotationList = append(annotationList, &annotationItem)
	}

	return annotationList
}

// TransformAnnotation 格式化列表数据
func (c *BaseAnnotationController) TransformAnnotation(v *models.Annotation) map[string]interface{} {

	annotationItem := make(map[string]interface{})
	aStatus, err := enums.GetSectionWithInt(v.Status, "annotation_status")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取状态转中文出错", nil)
	}

	annotationItem["Id"] = strconv.FormatInt(v.Id, 10)
	annotationItem["StatusString"] = aStatus
	annotationItem["PutrecNo"] = v.PutrecNo
	annotationItem["ImpexpPortcd"] = v.ImpexpPortcd
	annotationItem["ImpexpPortcdName"] = v.ImpexpPortcdName
	annotationItem["BondInvtNo"] = v.BondInvtNo
	annotationItem["EntryNo"] = v.EntryNo
	annotationItem["EtpsInnerInvtNo"] = v.EtpsInnerInvtNo
	annotationItem["CompanyName"] = v.Company.Name
	annotationItem["SeqNo"] = v.SeqNo
	annotationItem["BizopEtpsSccd"] = v.BizopEtpsSccd
	annotationItem["BizopEtpsno"] = v.BizopEtpsno
	annotationItem["BizopEtpsNm"] = v.BizopEtpsNm
	annotationItem["RcvgdEtpsno"] = v.RcvgdEtpsno
	annotationItem["RvsngdEtpsSccd"] = v.RvsngdEtpsSccd
	annotationItem["RcvgdEtpsNm"] = v.RcvgdEtpsNm
	annotationItem["DclEtpsSccd"] = v.DclEtpsSccd
	annotationItem["DclEtpsno"] = v.DclEtpsno
	annotationItem["DclEtpsNm"] = v.DclEtpsNm
	annotationItem["DclPlcCuscd"] = v.DclPlcCuscd
	annotationItem["DclPlcCuscdName"] = v.DclPlcCuscdName
	annotationItem["ImpexpMarkcd"] = v.ImpexpMarkcd
	annotationItem["MtpckEndprdMarkcd"] = v.MtpckEndprdMarkcd
	annotationItem["MtpckEndprdMarkcdName"] = v.MtpckEndprdMarkcdName
	annotationItem["SupvModecd"] = v.SupvModecd
	annotationItem["SupvModecdName"] = v.SupvModecdName
	annotationItem["TrspModecd"] = v.TrspModecd
	annotationItem["TrspModecdName"] = v.TrspModecdName
	annotationItem["DclcusFlag"] = v.DclcusFlag
	annotationItem["DclcusFlagName"] = v.DclcusFlagName
	annotationItem["DclcusTypecd"] = v.DclcusTypecd
	annotationItem["DclcusTypecdName"] = v.DclcusTypecdName
	annotationItem["VrfdedMarkcd"] = v.VrfdedMarkcd
	annotationItem["InvtIochkptStucd"] = v.InvtIochkptStucd
	annotationItem["ApplyNo"] = v.ApplyNo
	annotationItem["ListType"] = v.ListType
	annotationItem["ListTypeName"] = v.ListTypeName
	annotationItem["InputCode"] = v.InputCode
	annotationItem["InputCreditCode"] = v.InputCreditCode
	annotationItem["InputName"] = v.InputName
	annotationItem["ListStat"] = v.ListStat
	annotationItem["CorrEntryDclEtpsSccd"] = v.CorrEntryDclEtpsSccd
	annotationItem["CorrEntryDclEtpsNo"] = v.CorrEntryDclEtpsNo
	annotationItem["CorrEntryDclEtpsNm"] = v.CorrEntryDclEtpsNm
	annotationItem["DecType"] = v.DecType
	annotationItem["DecTypeName"] = v.DecTypeName
	annotationItem["StshipTrsarvNatcd"] = v.StshipTrsarvNatcd
	annotationItem["StshipTrsarvNatcdName"] = v.StshipTrsarvNatcdName
	annotationItem["InvtType"] = v.InvtType
	annotationItem["InvtTypeName"] = v.InvtTypeName
	annotationItem["EntryStucd"] = v.EntryStucd
	annotationItem["PassportUsedTypeCd"] = v.PassportUsedTypeCd
	annotationItem["Rmk"] = v.Rmk
	annotationItem["DecRmk"] = v.DecRmk
	annotationItem["DclTypecd"] = v.DclTypecd
	annotationItem["NeedEntryModified"] = v.NeedEntryModified
	annotationItem["LevyBlAmt"] = v.LevyBlAmt
	annotationItem["ChgTmsCnt"] = v.ChgTmsCnt
	annotationItem["RltInvtNo"] = v.RltInvtNo
	annotationItem["RltPutrecNo"] = v.RltPutrecNo
	annotationItem["RltEntryNo"] = v.RltEntryNo
	annotationItem["RltEntryBizopEtpsSccd"] = v.RltEntryBizopEtpsSccd
	annotationItem["RltEntryBizopEtpsno"] = v.RltEntryBizopEtpsno
	annotationItem["RltEntryBizopEtpsNm"] = v.RltEntryBizopEtpsNm
	annotationItem["RltEntryRvsngdEtpsSccd"] = v.RltEntryRvsngdEtpsSccd
	annotationItem["RltEntryRcvgdEtpsno"] = v.RltEntryRcvgdEtpsno
	annotationItem["RltEntryRcvgdEtpsNm"] = v.RltEntryRcvgdEtpsNm
	annotationItem["RltEntryDclEtpsSccd"] = v.RltEntryDclEtpsSccd
	annotationItem["RltEntryDclEtpsno"] = v.RltEntryDclEtpsno
	annotationItem["RltEntryDclEtpsNm"] = v.RltEntryDclEtpsNm
	annotationItem["Param1"] = v.Param1
	annotationItem["Param2"] = v.Param2
	annotationItem["Param3"] = v.Param3
	annotationItem["ExtraRemark"] = v.ExtraRemark
	annotationItem["GenDecFlag"] = v.GenDecFlag
	annotationItem["GenDecFlagName"] = v.GenDecFlagName
	annotationItem["HandBookId"] = strconv.FormatInt(v.HandBookId, 10)
	annotationItem["InputTime"] = enums.GetDateTimeString(&v.InputTime, enums.BaseDateTimeFormat)
	annotationItem["PrevdTime"] = enums.GetDateTimeString(&v.PrevdTime, enums.BaseDateTimeFormat)
	annotationItem["FormalVrfdedTime"] = enums.GetDateTimeString(&v.FormalVrfdedTime, enums.BaseDateTimeFormat)
	annotationItem["EntryDclTime"] = enums.GetDateTimeString(&v.EntryDclTime, enums.BaseDateTimeFormat)
	annotationItem["InvtDclTime"] = enums.GetDateTimeString(&v.InvtDclTime, enums.BaseDateTimeFormat)

	return annotationItem
}

// 仅仅更新状态
func (c *BaseAnnotationController) setStatusOnly(m *models.Annotation, statusString string) {
	if err := c.updateAnnotationStatus(m, statusString); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", nil)
	}
	if err := models.AnnotationUpdateStatus(m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}
}

// 操作记录
func (c *BaseAnnotationController) newAnnotationRecord(m *models.Annotation, content string) *models.AnnotationRecord {

	statusString, _ := enums.GetSectionWithInt(m.Status, "annotation_status")
	annotationRecord := models.NewAnnotationRecord(0)
	annotationRecord.Content = content
	annotationRecord.BackendUser = &c.curUser
	annotationRecord.Status = statusString
	annotationRecord.Annotation = m

	return &annotationRecord
}

// 是否能保存
func (c *BaseAnnotationController) getCanStore(m *models.Annotation, impexpMarkcd string) bool {

	if models.IsSuperAdmin(c.curUser.Id) {
		return true
	}

	if m == nil {
		return c.checkActionAuthor(c.controllerName, impexpMarkcd+"Audit")
	} else {
		if c.checkActionAuthor(c.controllerName, m.ImpexpMarkcd+"Audit") {
			aStatus, _ := enums.GetSectionWithInt(m.Status, "annotation_status")
			if aStatus == "待审核" || aStatus == "审核中" {
				return true
			}

		} else if c.checkActionAuthor(c.controllerName, m.ImpexpMarkcd+"Make") {
			if m == nil {
				return false
			}

			aStatus, _ := enums.GetSectionWithInt(m.Status, "annotation_status")
			if aStatus == "待制单" || aStatus == "制单中" {
				return true
			}
		}
	}

	return false
}