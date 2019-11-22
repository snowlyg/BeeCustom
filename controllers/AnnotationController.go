package controllers

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
	c.bIndex("I", "进口")
}

func (c *AnnotationController) EIndex() {
	c.bIndex("E", "出口")
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
	c.bStatusCount("I")
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

// Edit 添加 编辑 页面
func (c *AnnotationController) IMake() {
	Id, _ := c.GetInt64(":id", 0)
	c.bMake(Id)
}

// Edit 添加 编辑 页面
func (c *AnnotationController) EMake() {
	Id, _ := c.GetInt64(":id", 0)
	c.bMake(Id)
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

// ForRecheck 复核
func (c *AnnotationController) IRecheck() {
	Id, _ := c.GetInt64(":id", 0)
	c.bRecheck(Id)
}

// ForRecheck 复核
func (c *AnnotationController) ERecheck() {
	Id, _ := c.GetInt64(":id", 0)
	c.bRecheck(Id)
}

// RecheckPass 复核通过
func (c *AnnotationController) IRecheckPass() {
	Id, _ := c.GetInt64(":id", 0)
	c.bRecheckPassReject(Id)
}

// RecheckPass 复核通过
func (c *AnnotationController) ERecheckPass() {
	Id, _ := c.GetInt64(":id", 0)
	c.bRecheckPassReject(Id)
}

// RecheckReject 复核通过
func (c *AnnotationController) IRecheckReject() {
	Id, _ := c.GetInt64(":id", 0)
	c.bRecheckPassReject(Id)
}

// RecheckReject 复核通过
func (c *AnnotationController) ERecheckReject() {
	Id, _ := c.GetInt64(":id", 0)
	c.bRecheckPassReject(Id)
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
