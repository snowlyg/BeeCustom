package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// init 初始化
func init() {
	orm.RegisterModel(
		new(BackendUser),
		new(Resource),
		new(Role),
		new(Clearance),
		new(Ciq),
		new(HsCode),
		new(Company),
		new(CompanyContact),
		new(CompanyForeign),
		new(CompanySeal),
		new(HandBook),
		new(HandBookGood),
		new(HandBookUllage),
		new(Annotation),
		new(ClearanceUpdateTime),
		new(AnnotationUserRel),
		new(AnnotationItem),
		new(AnnotationRecord),
		new(AnnotationReturn),
	)
}

// TableName 下面是统一的表名管理
func TableName(name string) string {
	prefix := beego.AppConfig.String("db_dt_prefix")
	return prefix + name
}

// BackendUserTBName 获取 BackendUser 对应的表名称
func BackendUserTBName() string {
	return TableName("users")
}

// ResourceTBName 获取 Resource 对应的表名称
func ResourceTBName() string {
	return TableName("resource")
}

// RoleTBName 获取 Role 对应的表名称
func RoleTBName() string {
	return TableName("roles")
}

// ClearanceTBName 获取 Clearance 对应的表名称
func ClearanceTBName() string {
	return TableName("clearances")
}

// CiqTBName 获取 Ciq 对应的表名称
func CiqTBName() string {
	return TableName("ciqs")
}

// HsCodeTBName 获取 HsCode 对应的表名称
func HsCodeTBName() string {
	return TableName("hs_codes")
}

// CompanyTBName 获取 Company 对应的表名称
func CompanyTBName() string {
	return TableName("companies")
}

// CompanyContactTBName 获取 CompanyContact 对应的表名称
func CompanyContactTBName() string {
	return TableName("company_contacts")
}

// CompanyForeignTBName 获取 CompanyForeign 对应的表名称
func CompanyForeignTBName() string {
	return TableName("company_foreigns")
}

// CompanySealTBName 获取 CompanySeal 对应的表名称
func CompanySealTBName() string {
	return TableName("company_seals")
}

// HandBookTBName 获取 HandBook 对应的表名称
func HandBookTBName() string {
	return TableName("hand_books")
}

// HandBookGoodTBName 获取 HandBookGood 对应的表名称
func HandBookGoodTBName() string {
	return TableName("hand_book_goods")
}

// HandBookUllageTBName 获取 HandBookUllage 对应的表名称
func HandBookUllageTBName() string {
	return TableName("hand_book_ullages")
}

// AnnotationTBName 获取 Annotation 对应的表名称
func AnnotationTBName() string {
	return TableName("annotations")
}

// ClearanceUpdateTimeTBName 获取 ClearanceUpdateTime 对应的表名称
func ClearanceUpdateTimeTBName() string {
	return TableName("clearance_update_times")
}

// AnnotationUserRelTBName 获取 AnnotationUserRelTBName 对应的表名称
func AnnotationUserRelTBName() string {
	return "annotation_user_rel"
}

// AnnotationItemTBName 获取 AnnotationItemTBName 对应的表名称
func AnnotationItemTBName() string {
	return "annotation_items"
}

// AnnotationRecordTBName 获取 AnnotationRecordTBName 对应的表名称
func AnnotationRecordTBName() string {
	return "annotation_records"
}

// AnnotationReturnTBName 获取 AnnotationReturnTBName 对应的表名称
func AnnotationReturnTBName() string {
	return "annotation_returns"
}
