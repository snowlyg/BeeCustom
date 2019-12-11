package controllers

import (
	"BeeCustom/enums"
	"BeeCustom/models"
)

type OrderController struct {
	BaseOrderController
}

func (c *OrderController) Prepare() {
	// 先执行
	c.BaseController.Prepare()
	// 如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	perms := []string{}
	c.checkAuthor(perms)

	// 如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	// 权限控制里会进行登录验证，因此这里不用再作登录验证
	// c.checkLogin()
}

func (c *OrderController) IIndex() {
	c.bIndex("I")
}

func (c *OrderController) EIndex() {
	c.bIndex("E")
}

// 列表数据
func (c *OrderController) IDataGrid() {
	c.bDataGrid("I")
}

// 列表数据
func (c *OrderController) EDataGrid() {
	c.bDataGrid("E")
}

// 数据统计
func (c *OrderController) IStatusCount() {
	c.bStatusCount("I")
}

// 数据统计
func (c *OrderController) EStatusCount() {
	c.bStatusCount("E")
	c.bStatusCount("E")
}

//  Create 添加 新建 页面
func (c *OrderController) ICreate() {
	c.bCreate("I")
}

//  Create 添加 新建 页面
func (c *OrderController) ECreate() {
	c.bCreate("E")
}

//  Store 添加 新建 页面
func (c *OrderController) IStore() {
	c.bStore("I")
}

//  Store 添加 新建 页面
func (c *OrderController) EStore() {
	c.bStore("E")
}

//  Edit 添加 编辑 页面
func (c *OrderController) IEdit() {
	Id, _ := c.GetInt64(":id", 0)
	c.bEdit(Id)
}

//  Edit 添加 编辑 页面
func (c *OrderController) EEdit() {
	Id, _ := c.GetInt64(":id", 0)
	c.bEdit(Id)
}

//  IMake 制单
func (c *OrderController) IMake() {
	Id, _ := c.GetInt64(":id", 0)
	c.bMake(Id)
}

//  EMake 制单
func (c *OrderController) EMake() {
	Id, _ := c.GetInt64(":id", 0)
	c.bMake(Id)
}

//  IReMake 驳回修改
func (c *OrderController) IReMake() {
	Id, _ := c.GetInt64(":id", 0)
	c.bReMake(Id)
}

//  EReMake 驳回修改
func (c *OrderController) EReMake() {
	Id, _ := c.GetInt64(":id", 0)
	c.bReMake(Id)
}

//  Cancel 取消订单
func (c *OrderController) ICancel() {
	Id, _ := c.GetInt64(":id", 0)
	c.bCancel(Id)
}

//  Cancel 取消订单
func (c *OrderController) ECancel() {
	Id, _ := c.GetInt64(":id", 0)

	c.bCancel(Id)

}

//  Audit 审核通过订单
func (c *OrderController) IAudit() {
	Id, _ := c.GetInt64(":id", 0)
	c.bAudit(Id)
}

//  Audit 审核通过订单
func (c *OrderController) EAudit() {
	Id, _ := c.GetInt64(":id", 0)
	c.bAudit(Id)
}

//  Distribute 分配
func (c *OrderController) IDistribute() {
	BackendUserId, _ := c.GetInt64("BackendUserId")
	Id, _ := c.GetInt64(":id", 0)
	c.bDistribute(BackendUserId, Id)
}

//  Distribute 分配
func (c *OrderController) EDistribute() {
	BackendUserId, _ := c.GetInt64("BackendUserId")
	Id, _ := c.GetInt64(":id", 0)
	c.bDistribute(BackendUserId, Id)
}

//  Update 添加 编辑 页面
func (c *OrderController) IUpdate() {
	Id, _ := c.GetInt64(":id", 0)
	c.bUpdate(Id)
}

//  Update 添加 编辑 页面
func (c *OrderController) EUpdate() {
	Id, _ := c.GetInt64(":id", 0)
	c.bUpdate(Id)
}

//  Copy 添加 编辑 页面
func (c *OrderController) ICopy() {
	Id, _ := c.GetInt64(":id", 0)
	c.bCopy(Id)
}

//  Copy 添加 编辑 页面
func (c *OrderController) ECopy() {
	Id, _ := c.GetInt64(":id", 0)
	c.bCopy(Id)
}

//  ForRecheck 申请复核
func (c *OrderController) IForRecheck() {
	Id, _ := c.GetInt64(":id", 0)
	c.bForRecheck(Id)
}

//  ForRecheck 申请复核
func (c *OrderController) EForRecheck() {
	Id, _ := c.GetInt64(":id", 0)
	c.bForRecheck(Id)
}

//  ReForRecheck 重新申请复核
func (c *OrderController) IReForRecheck() {
	Id, _ := c.GetInt64(":id", 0)
	c.bReForRecheck(Id)
}

//  ReForRecheck 申请复核
func (c *OrderController) EReForRecheck() {
	Id, _ := c.GetInt64(":id", 0)
	c.bReForRecheck(Id)
}

//  IRecheck 复核
func (c *OrderController) IRecheck() {
	Id, _ := c.GetInt64(":id", 0)
	c.bRecheck(Id)
}

//  ERecheck 复核
func (c *OrderController) ERecheck() {
	Id, _ := c.GetInt64(":id", 0)
	c.bRecheck(Id)
}

//  RecheckPass 复核通过
func (c *OrderController) IRecheckPass() {
	c.bRecheckPassReject("复核通过", "pass", "复核通过凭证")
}

//  RecheckPass 复核通过
func (c *OrderController) ERecheckPass() {
	c.bRecheckPassReject("复核通过", "pass", "复核通过凭证")
}

//  RecheckReject 复核不通过
func (c *OrderController) IRecheckReject() {
	c.bRecheckPassReject("复核不通过", "reject", "复核驳回凭证")
}

//  RecheckReject 复核不通过
func (c *OrderController) ERecheckReject() {
	c.bRecheckPassReject("复核不通过", "reject", "复核驳回凭证")
}

//  PushXml 已提交单一
func (c *OrderController) IPushXml() {
	id, _ := c.GetInt64(":id")
	c.bPushXml(id)
}

//  PushXml 已提交单一
func (c *OrderController) EPushXml() {
	id, _ := c.GetInt64(":id")
	c.bPushXml(id)
}

//  Print 打印
func (c *OrderController) IPrint() {
	id, _ := c.GetInt64(":id")
	c.bPrint(id)
}

//  Print 打印
func (c *OrderController) EPrint() {
	id, _ := c.GetInt64(":id")
	c.bPrint(id)
}

//  Remark 附注
func (c *OrderController) IRemark() {
	id, _ := c.GetInt64(":id")
	remark := c.GetString("Remark")
	c.bRemark(id, remark)
}

//  Remark 附注
func (c *OrderController) ERemark() {
	id, _ := c.GetInt64(":id")
	remark := c.GetString("Remark")
	c.bRemark(id, remark)
}

//  AuditFirstRejectLog 附注
func (c *OrderController) IAuditFirstRejectLog() {
	id, _ := c.GetInt64(":id")
	c.bAuditFirstRejectLog(id)
}

//  AuditFirstRejectLog 附注
func (c *OrderController) EAuditFirstRejectLog() {
	id, _ := c.GetInt64(":id")
	c.bAuditFirstRejectLog(id)
}

//  Restart 重启
func (c *OrderController) IRestart() {
	id, _ := c.GetInt64(":id")
	c.bRestart(id)
}

//  Restart 附注
func (c *OrderController) ERestart() {
	id, _ := c.GetInt64(":id")
	c.bRestart(id)
}

// 删除
func (c *OrderController) IDelete() {
	id, _ := c.GetInt64(":id")
	c.bDelete(id)
}

// 删除
func (c *OrderController) EDelete() {
	id, _ := c.GetInt64(":id")
	c.bDelete(id)
}

// 客户管理联系人信息
func (c *OrderController) CompanyAdminUser() {
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
