package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// init 初始化
func init() {
	orm.RegisterModel(new(BackendUser), new(Resource), new(Role), new(Clearance), new(Ciq), new(HsCode), new(Company))
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
