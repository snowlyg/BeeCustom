package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

// TableName 设置表名
func (a *Role) TableName() string {
	return RoleTBName()
}

// RoleQueryParam 用于搜索的类
type RoleQueryParam struct {
	BaseQueryParam
	NameLike string
}

// Role 用户角色 实体类
type Role struct {
	Id              int    `form:"Id"`
	Name            string `form:"Name"`
	Seq             int
	RoleResourceRel []*RoleResourceRel `orm:"reverse(many)" json:"-"` // 设置一对多的反向关系
	BackendUsers    []*BackendUser     `orm:"reverse(many)"`          //设置一对多关系
	CreatedAt       time.Time          `orm:"column(created_at);type(timestamp);null"`
	UpdatedAt       time.Time          `orm:"column(updated_at);type(timestamp);null"`
}

// RolePageList 获取分页数据
func RolePageList(params *RoleQueryParam) ([]*Role, int64) {

	query := orm.NewOrm().QueryTable(RoleTBName())
	data := make([]*Role, 0)

	//默认排序
	sortorder := "Id"
	if len(params.Sort) > 0 {
		sortorder = params.Sort
	}

	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}

	query = query.Filter("name__istartswith", params.NameLike)

	total, _ := query.Count()
	_, _ = query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit).All(&data)

	return data, total
}

// RoleDataList 获取角色列表
func RoleDataList(params *RoleQueryParam) []*Role {
	params.Limit = -1
	params.Sort = "Seq"
	params.Order = "asc"
	data, _ := RolePageList(params)
	return data
}

// RoleOne 获取单条
func RoleOne(id int) (*Role, error) {
	o := orm.NewOrm()
	m := Role{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
