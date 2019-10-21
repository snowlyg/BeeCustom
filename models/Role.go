package models

import (
	"github.com/astaxie/beego/orm"
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
	BaseModel

	Name         string         `form:"Name"`
	Resources    []*Resource    `orm:"rel(m2m)"`      // 设置一对多的反向关系
	BackendUsers []*BackendUser `orm:"reverse(many)"` //设置一对多关系
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
	params.Sort = "Id"
	params.Order = "asc"
	data, _ := RolePageList(params)
	return data
}

// RoleOne 获取单条
func RoleOne(id int64) (*Role, error) {
	m := Role{BaseModel: BaseModel{Id: id}}

	o := orm.NewOrm()
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}

	// 获取关系字段，o.LoadRelated(v, "Roles") 这是关键
	// 查找该用户所属的角色
	if _, err := o.LoadRelated(&m, "Resources"); err != nil {
		return nil, err
	}

	return &m, nil
}

//Save 添加、编辑页面 保存
func RoleSave(m *Role, perm_ids string) (*Role, error) {

	o := orm.NewOrm()
	if m.Id == 0 {
		if _, err := o.Insert(m); err != nil {
			return nil, err
		}

	} else {
		if _, err := o.Update(m, "Name", "UpdatedAt"); err != nil {
			return nil, err
		}
	}

	m2m := o.QueryM2M(&m, "Resources")
	if _, err := m2m.Clear(); err != nil {
		return nil, err
	}

	if _, err := m2m.Add(); err != nil {
		return nil, err
	}

	return m, nil

}

//删除
func RoleDelete(id int64) (num int64, err error) {

	if num, err := BaseDelete(&Role{BaseModel: BaseModel{Id: id}}); err != nil {
		return num, err
	} else {
		return num, nil
	}

}
