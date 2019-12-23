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

// 清单
func (c *BaseAnnotationController) bIndex(impexpMarkcd string) {
	// 页面模板设置
	c.setTpl("annotation/index.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "annotation/index_footerjs.html"

	// 页面里按钮权限控制
	c.getActionData(impexpMarkcd, "Index", "Create", "Edit", "Make", "ReMake", "Audit", "Delete", "Distribute", "Recheck", "Push", "PushXml", "StoreError", "Change", "Restart", "Cancel", "Copy")

	// 获取制单人
	backendUsers := models.GetCreateBackendUsers("AnnotationController.Make")
	c.Data["BackendUsers"] = backendUsers
	c.Data["IsDelete"] = false
	c.Data["ImpexpMarkcd"] = impexpMarkcd
	c.Data["ImpexpMarkcdName"] = enums.GetImpexpMarkcdCNName(impexpMarkcd)
	c.GetXSRFToken()
}

// 回收站
func (c *BaseAnnotationController) bRecycle(impexpMarkcd string) {
	// 页面模板设置
	c.setTpl("annotation/index.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "annotation/index_footerjs.html"

	// 页面里按钮权限控制
	c.getActionData(impexpMarkcd, "Restore", "ForceDelete")

	c.Data["IsDelete"] = true
	c.Data["ImpexpMarkcd"] = impexpMarkcd
	c.Data["ImpexpMarkcdName"] = enums.GetImpexpMarkcdCNName(impexpMarkcd)
	c.GetXSRFToken()
}

// 还原删除订单
func (c *BaseAnnotationController) bRestore(id int64) {
	m, err := models.AnnotationOne(id, "")
	if m != nil && id > 0 {
		m.DeletedAt = time.Time{}
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}

	c.setStatusOnly(m, "待审核", false)
	if err = models.AnnotationUpdate(m, []string{"Status", "DeletedAt", "StatusUpdatedAt"}); err != nil {
		c.jsonResult(enums.JRCodeFailed, "取消失败", m)
	}
	annotationRecord := c.newAnnotationRecord(m, "还原删除订单")
	if err := models.AnnotationRecordSave(annotationRecord); err != nil {
		c.jsonResult(enums.JRCodeFailed, "还原失败", m)
	}
	c.jsonResult(enums.JRCodeSucc, "还原成功", m)
}

// 彻底删除订单
func (c *BaseAnnotationController) bForceDelete(id int64) {
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
	if err = UpdateAnnotationStatus(&m, "待审核", false); err != nil {
		c.jsonResult(enums.JRCodeFailed, "添加失败", nil)
	}
	m.Company = company
	m.InputTime = time.Now()
	m.InvtDclTime = time.Now()
	m.EtpsInnerInvtNo = c.getEtpsInnerInvtNo(impexpMarkcd, m.DclPlcCuscd)
	c.validRequestData(m)
	if err := models.AnnotationUpdateOrSave(&m); err != nil {
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

// copy 复制
func (c *BaseAnnotationController) bCopy(id int64) {
	m, err := models.AnnotationOne(id, "Company,AnnotationItems")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据出错", nil)
	}
	if err := UpdateAnnotationStatus(m, "待审核", true); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", nil)
	}

	// 重置数据
	if m != nil {
		newAnnotation := models.NewAnnotation(0)
		newAnnotation = *m
		newAnnotation.Id = 0
		newAnnotation.InputTime = time.Now()
		newAnnotation.InvtDclTime = time.Now()
		newAnnotation.EtpsInnerInvtNo = c.getEtpsInnerInvtNo(m.ImpexpMarkcd, m.DclPlcCuscd)
		newAnnotation.SeqNo = ""
		newAnnotation.BondInvtNo = ""
		if err := models.AnnotationUpdateOrSave(&newAnnotation); err != nil {
			c.jsonResult(enums.JRCodeFailed, "操作失败", err)
		} else {
			// 复制表体
			for _, annotationItem := range m.AnnotationItems {
				newAnnotationItem := models.NewAnnotationItem(0)
				newAnnotationItem = *annotationItem
				newAnnotationItem.Id = 0
				newAnnotationItem.Annotation = m
				newAnnotationItem.AnnotationId = 0
				if err := models.AnnotationItemSave(&newAnnotationItem); err != nil {
					c.jsonResult(enums.JRCodeFailed, "操作失败", err)
				}
			}

			if err := c.setAnnotaionUserRelType(m, nil, "创建人"); err != nil {
				c.jsonResult(enums.JRCodeFailed, "操作失败", err)
			}
			annotationRecord := c.newAnnotationRecord(m, "创建订单")
			if err := models.AnnotationRecordSave(annotationRecord); err != nil {
				c.jsonResult(enums.JRCodeFailed, "操作失败", err)
			}
			c.jsonResult(enums.JRCodeSucc, "操作成功", m)
		}
	} else {
		c.jsonResult(enums.JRCodeFailed, "操作失败", nil)
	}

}

// Edit 添加 编辑 页面
func (c *BaseAnnotationController) bEdit(id int64) {
	m, err := models.AnnotationOne(id, "")
	if err != nil {
		c.pageError("数据无效，请刷新后重试")
	}
	c.setStatusOnly(m, "审核中", false)
	// 获取制单人
	backendUsers := models.GetCreateBackendUsers("AnnotationController.Make")
	c.Data["BackendUsers"] = backendUsers
	c.Data["m"] = models.TransformAnnotation(id, "AnnotationItems")
	c.Data["canStore"] = c.getCanStore(m, "")
	if m != nil {
		c.getResponses(m.ImpexpMarkcd)
	}
}

// bMake 制单
func (c *BaseAnnotationController) bMake(id int64) {
	m, err := models.AnnotationOne(id, "")
	if m != nil && id > 0 {
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}
	c.setStatusOnly(m, "制单中", false)
	c.Data["m"] = models.TransformAnnotation(id, "AnnotationItems")
	c.Data["canStore"] = c.getCanStore(m, "")
	if m != nil {
		c.getResponses(m.ImpexpMarkcd)
	}
}

// bReMake 驳回修改
func (c *BaseAnnotationController) bReMake(id int64) {
	m, err := models.AnnotationOne(id, "")
	if m != nil && id > 0 {
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}

	c.Data["m"] = models.TransformAnnotation(id, "AnnotationItems,AnnotationRecords")
	c.Data["canStore"] = c.getCanStore(m, "")
	if m != nil {
		c.getResponses(m.ImpexpMarkcd)
	}
}

// 编辑相关页面返回
func (c *BaseAnnotationController) getResponses(impexpMarkcd string) {
	// 页面里按钮权限控制
	c.getActionData(impexpMarkcd, "Audit", "Distribute", "ForRecheck", "Print", "ExtraRemark", "ReForRecheck")
	c.Data["ImpexpMarkcdName"] = enums.GetImpexpMarkcdCNName(impexpMarkcd)
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
	if err = UpdateAnnotationStatus(m, "订单关闭", false); err != nil {
		c.jsonResult(enums.JRCodeFailed, "取消失败", m)
	}
	if err := models.AnnotationUpdate(m, []string{"Status", "StatusUpdatedAt"}); err != nil {
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
	c.setStatusOnly(m, "审核通过", false)
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
	if err = UpdateAnnotationStatus(m, "待制单", false); err != nil {
		c.jsonResult(enums.JRCodeFailed, "派单失败", m)
	}
	if err := models.AnnotationUpdateOrSave(m); err != nil {
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
	if m != nil {
		m.InputTime = time.Now()
		m.InvtDclTime = time.Now()
	}
	c.validRequestData(m)
	if err := models.AnnotationUpdateOrSave(m); err != nil {
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
	if err = UpdateAnnotationStatus(m, "待复核", false); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}
	if err := models.AnnotationUpdate(m, []string{"Status", "StatusUpdatedAt"}); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}
	annotationRecord := c.newAnnotationRecord(m, "复核")
	if err := models.AnnotationRecordSave(annotationRecord); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}
	c.jsonResult(enums.JRCodeSucc, "操作成功", m)
}

// bRestart 重新开启
func (c *BaseAnnotationController) bRestart(id int64) {
	m, err := models.AnnotationOne(id, "Company,AnnotationItems")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}
	if err = UpdateAnnotationStatus(m, "待审核", true); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}
	if err := models.AnnotationUpdate(m, []string{"Status", "StatusUpdatedAt"}); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}
	annotationRecord := c.newAnnotationRecord(m, "重新开启订单")
	if err := models.AnnotationRecordSave(annotationRecord); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}
	c.jsonResult(enums.JRCodeSucc, "操作成功", m)
}

// bReForRecheck 重新申请复核
func (c *BaseAnnotationController) bReForRecheck(id int64) {
	m, err := models.AnnotationOne(id, "Company,AnnotationItems")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}
	if err = UpdateAnnotationStatus(m, "待复核", true); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}

	if err := models.AnnotationUpdate(m, []string{"Status", "StatusUpdatedAt"}); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}
	annotationRecord := c.newAnnotationRecord(m, "复核")
	if err := models.AnnotationRecordSave(annotationRecord); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}
	c.jsonResult(enums.JRCodeSucc, "操作成功", m)
}

// bRecheckPass 通过复核、驳回
func (c *BaseAnnotationController) bRecheckPassReject(statusString, action, actionName string) {
	Id, _ := c.GetInt64(":id", 0)
	m, err := models.AnnotationOne(Id, "Company,AnnotationItems")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m)
	}
	c.validRequestData(m)
	if err = UpdateAnnotationStatus(m, statusString, false); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}
	recheckErrorInputIds := c.GetString("RecheckErrorInputIds")
	itemRecheckErrorInputIds := c.GetString("ItemRecheckErrorInputIds")
	m.RecheckErrorInputIds = recheckErrorInputIds
	m.ItemRecheckErrorInputIds = itemRecheckErrorInputIds
	if err := models.AnnotationUpdate(m, []string{"Status", "StatusUpdatedAt", "RecheckErrorInputIds", "ItemRecheckErrorInputIds"}); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}
	content := c.GetString("Content")
	remark := c.GetString("Remark")
	if len(content) > 0 {
		statusString += ":" + content
	}
	annotationRecord := c.newAnnotationRecord(m, statusString)
	annotationRecord.Remark = remark
	if err := models.AnnotationRecordSave(annotationRecord); err != nil {
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
		m.EtpsInnerInvtNo,
		"annotation_recheck_pdf",
		action,
		"annotation",
		"",
		username,
		password,
		10,
	}

	if ffp, err := enums.NewPDFGenerator(&pdfData); err != nil {
		c.jsonResult(enums.JRCodeFailed, "添加失败", m)
	} else {
		aFile := models.NewAnnotationFile(0)
		aFile.Url = strings.Replace(ffp, ".", "", 1)
		aFile.Type = actionName
		aFile.Name = actionName
		aFile.Creator = c.curUser.RealName
		aFile.Annotation = m
		aFile.Version = "1.0"
		err = models.AnnotationFileSaveOrUpdate(&aFile)
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, "添加失败", err)
		}
	}

	c.jsonResult(enums.JRCodeSucc, "操作成功", m)
}

// bRecheck 复核
func (c *BaseAnnotationController) bRecheck(id int64) {
	m, err := models.AnnotationOne(id, "AnnotationItems")
	if err != nil {
		c.pageError("数据无效，请刷新后重试")
	}
	if m != nil {
		c.getActionData(m.ImpexpMarkcd, "RecheckPass", "RecheckReject")
	}
	c.setStatusOnly(m, "复核中", false)
	annotation := models.TransformAnnotation(id, "AnnotationItems")
	c.Data["m"] = annotation
	c.setTpl("annotation/recheck.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "annotation/recheck_footerjs.html"
	// 页面里按钮权限控制
	c.GetXSRFToken()
}

// bPrint 打印
func (c *BaseAnnotationController) bPrint(id int64) {
	m, err := models.AnnotationOne(id, "AnnotationItems")
	if err != nil {
		c.pageError("数据无效，请刷新后重试")
	}
	if m != nil {
		// 生成 pdf 凭证
		// basic auth 认证用户名和密码
		username, _ := models.GetSettingValueByKey("pdf_username")
		password, _ := models.GetSettingValueByKey("pdf_password")
		pdfData := enums.PdfData{
			m.Id,
			m.EtpsInnerInvtNo,
			"annotation_pdf",
			"report",
			"annotation",
			"",
			username,
			password,
			10,
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
func (c *BaseAnnotationController) bPushXml(id int64) {
	m, err := models.AnnotationOne(id, "AnnotationItems")
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
		handBookTypeS, _ := models.GetSettingRValueByKey("handBookType", false)
		handBookType1, _, _ := enums.TransformCnToInt(handBookTypeS, "手册")
		handBookType2, _, _ := enums.TransformCnToInt(handBookTypeS, "账册")
		if handBook == nil {
			c.jsonResult(enums.JRCodeFailed, "错误手账册类型", nil)
		} else {
			if handBook.Type == handBookType1 {
				receiverId, _ = models.GetSettingValueByKey("AnnotationReceiverIdC")
				sysId, _ = models.GetSettingValueByKey("AnnotationSysIdC")
				path, _ = models.GetSettingValueByKey("annotation_xml_path_c")
			} else if handBook.Type == handBookType2 {
				receiverId, _ = models.GetSettingValueByKey("AnnotationReceiverIdE")
				sysId, _ = models.GetSettingValueByKey("AnnotationSysIdE")
				path, _ = models.GetSettingValueByKey("annotation_xml_path_e")
			} else {
				c.jsonResult(enums.JRCodeFailed, "错误手账册类型", nil)
			}
		}

		pathTemp := "./static/generate/annotation/" + strconv.FormatInt(id, 10) + "/temp/"
		// 报文名称
		mName := time.Now().Format(enums.BaseDateTimeSecondFormat) + "_" + m.EtpsInnerInvtNo
		fileName := mName + ".xml"

		signature.Object.Package.EnvelopInfo.FileName = fileName
		signature.Object.Package.EnvelopInfo.Version, _ = models.GetSettingValueByKey("AnnotationVersion")
		signature.Object.Package.EnvelopInfo.BusinessId, _ = models.GetSettingValueByKey("AnnotationBusinessId")
		signature.Object.Package.EnvelopInfo.MessageType, _ = models.GetSettingValueByKey("AnnotationMessageType")
		signature.Object.Package.EnvelopInfo.SenderId, _ = models.GetSettingValueByKey("AnnotationSenderId")
		signature.Object.Package.EnvelopInfo.ReceiverId = receiverId
		signature.Object.Package.EnvelopInfo.MessageId = m.EtpsInnerInvtNo
		signature.Object.Package.EnvelopInfo.SendTime = time.Now().Format(enums.RFC3339)

		signature.Object.Package.DataInfo.BussinessData.DelcareFlag = "0" // 0:暂存，1:申报
		signature.Object.Package.DataInfo.BussinessData.InvtMessage.SysId = sysId
		signature.Object.Package.DataInfo.BussinessData.InvtMessage.OperCusRegCode, _ = models.GetSettingValueByKey("AgentCode")

		invtHeadType := xmlTemplate.InvtHeadType{}
		enums.SetObjValueFromObj(&invtHeadType, m) // 设置数据到 xml 结构体

		iCCode, _ := models.GetSettingValueByKey("ICCode")
		invtHeadType.IcCardNo = iCCode

		rmk := xmlTemplate.Cdata{Value: m.Rmk}
		invtHeadType.Rmk = rmk

		var invtListTypes []xmlTemplate.InvtListType
		for _, v := range m.AnnotationItems {
			invtListType := xmlTemplate.InvtListType{}
			enums.SetObjValueFromObj(&invtListType, v) // 设置数据到 xml 结构体

			gdsNm := xmlTemplate.Cdata{Value: v.GdsNm}
			invtListType.GdsNm = gdsNm

			rmk := xmlTemplate.Cdata{Value: v.Rmk}
			invtListType.Rmk = rmk

			invtListTypes = append(invtListTypes, invtListType)
		}

		invtDecHeadType := xmlTemplate.InvtDecHeadType{}
		var invtDecListType []xmlTemplate.InvtDecListType

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

// bExtraRemark 附注
func (c *BaseAnnotationController) bExtraRemark(id int64, extraRemark string) {
	m, err := models.AnnotationOne(id, "")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", err)
	}
	m.ExtraRemark = extraRemark
	if err = models.AnnotationUpdate(m, []string{"ExtraRemark"}); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", err)
	}
	c.jsonResult(enums.JRCodeSucc, "操作成功", err)
}

// bAuditFirstRejectLog 驳回原因
func (c *BaseAnnotationController) bAuditFirstRejectLog(id int64) {
	m, err := models.AnnotationOne(id, "")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", err)
	}

	aRecord := models.NewAnnotationRecord(0)
	aRecord.Annotation = m

	aStatusS, err := c.getAnnotationStatus("orderStatus")
	sSting, err, _ := enums.TransformIntToCn(aStatusS, m.Status)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", err)
	}

	if err = models.AnnotationRecordOneByStatusAndAnnotationId(&aRecord, sSting); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", err)
	}

	c.jsonResult(enums.JRCodeSucc, "操作成功", aRecord)
}

// 删除
func (c *BaseAnnotationController) bDelete(id int64) {
	m, err := models.AnnotationOne(id, "")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "删除失败", err)
	}

	m.DeletedAt = time.Now()
	if err := models.AnnotationUpdate(m, []string{"DeletedAt"}); err != nil {
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
func (c *BaseAnnotationController) setAnnotaionUserRelType(m *models.Annotation, bu *models.BackendUser, userTypes string) error {
	rs := strings.Split(userTypes, ",")
	for _, v := range rs {
		aur := models.NewAnnotationUserRel(0)
		aStatusS, err := c.getAnnotationStatus("annotationUserType")
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
		aur.Annotation = m
		aur.UserType = aStatus
		if err = models.AnnotationUserRelSave(&aur); err != nil {
			return err
		}
	}
	return nil
}

func (c *BaseAnnotationController) getAnnotationStatus(status string) (map[string]string, error) {
	return models.GetSettingRValueByKey(status, false)
}

// TransformAnnotationList 格式化列表数据
func (c *BaseAnnotationController) TransformAnnotationList(ms []*models.Annotation) []*map[string]interface{} {

	var annotationList []*map[string]interface{}
	for _, v := range ms {
		annotationCreatorName := "" // 制单人
		annotationItem := make(map[string]interface{})
		aStatusS, err := c.getAnnotationStatus("orderStatus")
		aStatus, err, _ := enums.TransformIntToCn(aStatusS, v.Status)

		if err != nil {
			c.jsonResult(enums.JRCodeFailed, "获取状态转中文出错", nil)
		}
		userTypeS, err := c.getAnnotationStatus("annotationUserType")
		userType, err, _ := enums.TransformCnToInt(userTypeS, "制单人")
		if err != nil {
			utils.LogDebug(fmt.Sprintf("转换制单人类型出错:%v", err))
		}
		if len(v.BackendUsers) > 0 {
			for _, bu := range v.BackendUsers {
				abur, err := models.AnnotationUserRelByUserIdAndAnnotationId(bu.Id, v.Id, userType)
				if err != nil {
					c.jsonResult(enums.JRCodeFailed, "获取制单人出错", nil)
				}
				if abur != nil && abur.Id != 0 {
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

// 仅仅更新状态
func (c *BaseAnnotationController) setStatusOnly(m *models.Annotation, statusString string, isRestart bool) {
	if err := UpdateAnnotationStatus(m, statusString, isRestart); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", nil)
	}
	if err := models.AnnotationUpdate(m, []string{"Status", "StatusUpdatedAt"}); err != nil {
		c.jsonResult(enums.JRCodeFailed, "操作失败", m)
	}
}

// 操作记录
func (c *BaseAnnotationController) newAnnotationRecord(m *models.Annotation, content string) *models.AnnotationRecord {
	aStatusS, _ := c.getAnnotationStatus("orderStatus")
	statusString, _, _ := enums.TransformIntToCn(aStatusS, m.Status)
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
		aStatusS, _ := c.getAnnotationStatus("orderStatus")
		aStatus, _, _ := enums.TransformIntToCn(aStatusS, m.Status)
		if c.checkActionAuthor(c.controllerName, m.ImpexpMarkcd+"Audit") {
			if aStatus == "待审核" || aStatus == "审核中" {
				return true
			}
		} else if c.checkActionAuthor(c.controllerName, m.ImpexpMarkcd+"Make") {
			if aStatus == "待制单" || aStatus == "制单中" {
				return true
			}
		}
	}
	return false
}
