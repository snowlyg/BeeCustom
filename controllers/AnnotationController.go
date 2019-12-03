package controllers

import (
	"BeeCustom/enums"
	"BeeCustom/models"
)

type AnnotationController struct {
	BaseAnnotationController
}

func (c *AnnotationController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	perms := []string{}
	c.checkAuthor(perms)

	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()
}

func (c *AnnotationController) IIndex() {
	c.bIndex("I")
}

func (c *AnnotationController) EIndex() {
	c.bIndex("E")
}

//列表数据
func (c *AnnotationController) IDataGrid() {
	c.bDataGrid("I")
}

//列表数据
func (c *AnnotationController) EDataGrid() {
	c.bDataGrid("E")
}

//数据统计
func (c *AnnotationController) IStatusCount() {
	c.bStatusCount("I")
}

//数据统计
func (c *AnnotationController) EStatusCount() {
	c.bStatusCount("E")
	c.bStatusCount("E")
}

// Create 添加 新建 页面
func (c *AnnotationController) ICreate() {
	c.bCreate("I")
}

// Create 添加 新建 页面
func (c *AnnotationController) ECreate() {
	c.bCreate("E")
}

// Store 添加 新建 页面
func (c *AnnotationController) IStore() {
	c.bStore("I")
}

// Store 添加 新建 页面
func (c *AnnotationController) EStore() {
	c.bStore("E")
}

// Edit 添加 编辑 页面
func (c *AnnotationController) IEdit() {
	Id, _ := c.GetInt64(":id", 0)
	c.bEdit(Id)
}

// Edit 添加 编辑 页面
func (c *AnnotationController) EEdit() {
	Id, _ := c.GetInt64(":id", 0)
	c.bEdit(Id)
}

// IMake 制单
func (c *AnnotationController) IMake() {
	Id, _ := c.GetInt64(":id", 0)
	c.bMake(Id)
}

// EMake 制单
func (c *AnnotationController) EMake() {
	Id, _ := c.GetInt64(":id", 0)
	c.bMake(Id)
}

// IReMake 驳回修改
func (c *AnnotationController) IReMake() {
	Id, _ := c.GetInt64(":id", 0)
	c.bReMake(Id)
}

// EReMake 驳回修改
func (c *AnnotationController) EReMake() {
	Id, _ := c.GetInt64(":id", 0)
	c.bReMake(Id)
}

// Cancel 取消订单
func (c *AnnotationController) ICancel() {
	Id, _ := c.GetInt64(":id", 0)
	c.bCancel(Id)
}

// Cancel 取消订单
func (c *AnnotationController) ECancel() {
	Id, _ := c.GetInt64(":id", 0)

	c.bCancel(Id)

}

// Audit 审核通过订单
func (c *AnnotationController) IAudit() {
	Id, _ := c.GetInt64(":id", 0)
	c.bAudit(Id)
}

// Audit 审核通过订单
func (c *AnnotationController) EAudit() {
	Id, _ := c.GetInt64(":id", 0)
	c.bAudit(Id)
}

// Distribute 分配
func (c *AnnotationController) IDistribute() {
	BackendUserId, _ := c.GetInt64("BackendUserId")
	Id, _ := c.GetInt64(":id", 0)
	c.bDistribute(BackendUserId, Id)
}

// Distribute 分配
func (c *AnnotationController) EDistribute() {
	BackendUserId, _ := c.GetInt64("BackendUserId")
	Id, _ := c.GetInt64(":id", 0)
	c.bDistribute(BackendUserId, Id)
}

// Update 添加 编辑 页面
func (c *AnnotationController) IUpdate() {
	Id, _ := c.GetInt64(":id", 0)
	c.bUpdate(Id)
}

// Update 添加 编辑 页面
func (c *AnnotationController) EUpdate() {
	Id, _ := c.GetInt64(":id", 0)
	c.bUpdate(Id)
}

// Copy 添加 编辑 页面
func (c *AnnotationController) ICopy() {
	Id, _ := c.GetInt64(":id", 0)
	c.bCopy(Id)
}

// Copy 添加 编辑 页面
func (c *AnnotationController) ECopy() {
	Id, _ := c.GetInt64(":id", 0)
	c.bCopy(Id)
}

// ForRecheck 申请复核
func (c *AnnotationController) IForRecheck() {
	Id, _ := c.GetInt64(":id", 0)
	c.bForRecheck(Id)
}

// ForRecheck 申请复核
func (c *AnnotationController) EForRecheck() {
	Id, _ := c.GetInt64(":id", 0)
	c.bForRecheck(Id)
}

// ReForRecheck 重新申请复核
func (c *AnnotationController) IReForRecheck() {
	Id, _ := c.GetInt64(":id", 0)
	c.bReForRecheck(Id)
}

// ReForRecheck 申请复核
func (c *AnnotationController) EReForRecheck() {
	Id, _ := c.GetInt64(":id", 0)
	c.bReForRecheck(Id)
}

// IRecheck 复核
func (c *AnnotationController) IRecheck() {
	Id, _ := c.GetInt64(":id", 0)
	c.bRecheck(Id)
}

// ERecheck 复核
func (c *AnnotationController) ERecheck() {
	Id, _ := c.GetInt64(":id", 0)
	c.bRecheck(Id)
}

// RecheckPass 复核通过
func (c *AnnotationController) IRecheckPass() {
	c.bRecheckPassReject("复核通过", "pass", "复核通过凭证")
}

// RecheckPass 复核通过
func (c *AnnotationController) ERecheckPass() {
	c.bRecheckPassReject("复核通过", "pass", "复核通过凭证")
}

// RecheckReject 复核不通过
func (c *AnnotationController) IRecheckReject() {
	c.bRecheckPassReject("复核不通过", "reject", "复核驳回凭证")
}

// RecheckReject 复核不通过
func (c *AnnotationController) ERecheckReject() {
	c.bRecheckPassReject("复核不通过", "reject", "复核驳回凭证")
}

// PushXml 已提交单一
func (c *AnnotationController) IPushXml() {
	id, _ := c.GetInt64(":id")
	c.bPushXml(id)
}

// PushXml 已提交单一
func (c *AnnotationController) EPushXml() {
	id, _ := c.GetInt64(":id")
	c.bPushXml(id)
}

// Print 打印
func (c *AnnotationController) IPrint() {
	id, _ := c.GetInt64(":id")
	c.bPrint(id)
}

// Print 打印
func (c *AnnotationController) EPrint() {
	id, _ := c.GetInt64(":id")
	c.bPrint(id)
}

// ExtraRemark 附注
func (c *AnnotationController) IExtraRemark() {
	id, _ := c.GetInt64(":id")
	extraRemark := c.GetString("ExtraRemark")
	c.bExtraRemark(id, extraRemark)
}

// ExtraRemark 附注
func (c *AnnotationController) EExtraRemark() {
	id, _ := c.GetInt64(":id")
	extraRemark := c.GetString("ExtraRemark")
	c.bExtraRemark(id, extraRemark)
}

// AuditFirstRejectLog 附注
func (c *AnnotationController) IAuditFirstRejectLog() {
	id, _ := c.GetInt64(":id")
	c.bAuditFirstRejectLog(id)
}

// AuditFirstRejectLog 附注
func (c *AnnotationController) EAuditFirstRejectLog() {
	id, _ := c.GetInt64(":id")
	c.bAuditFirstRejectLog(id)
}

// Restart 重启
func (c *AnnotationController) IRestart() {
	id, _ := c.GetInt64(":id")
	c.bRestart(id)
}

// Restart 附注
func (c *AnnotationController) ERestart() {
	id, _ := c.GetInt64(":id")
	c.bRestart(id)
}

//删除
func (c *AnnotationController) IDelete() {
	id, _ := c.GetInt64(":id")
	c.bDelete(id)
}

//删除
func (c *AnnotationController) EDelete() {
	id, _ := c.GetInt64(":id")
	c.bDelete(id)
}

//客户管理联系人信息
func (c *AnnotationController) CompanyAdminUser() {
	id, _ := c.GetInt64(":id")
	annotation, err := models.AnnotationOne(id, "Company")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "数据查询失败", err)
	}

	companyId := annotation.Company.Id
	adminCompanyContact, err := models.GetAdminCompanyContactByCompanyId(companyId)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "数据查询失败", err)
	}

	c.Data["json"] = adminCompanyContact
	c.ServeJSON()
}
