package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"BeeCustom/utils"
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

	Name      string `orm:"size(32)" form:"Name" valid:"Required;MaxSize(32)"`
	Resources string `orm:"_"` // 设置多对多的反向关系
	//BackendUsers []*BackendUser `orm:"_"`                     //设置一对多关系
}

//初始化角色
func NewRole(id int64) Role {
	return Role{BaseModel: BaseModel{Id: id, CreatedAt: time.Now(), UpdatedAt: time.Now()}}
}

// RolePageList 获取分页数据
func RolePageList(params *RoleQueryParam) ([]*Role, int64) {
	query := orm.NewOrm().QueryTable(RoleTBName())
	data := make([]*Role, 0)
	if len(params.NameLike) > 0 {
		query = query.Filter("name__istartswith", params.NameLike)
	}

	query = BaseListQuery(query, params.Sort, params.Order, params.Limit, params.Offset)

	total, _ := query.Count()
	_, _ = query.All(&data)

	return data, total
}

//查询参数
func NewRoleQueryParam() RoleQueryParam {
	return RoleQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc", Offset: 0}}
}

// RoleDataList 获取角色列表
func RoleDataList(params *RoleQueryParam) []*Role {
	data, _ := RolePageList(params)
	return data
}

// RoleOne 获取单条
func RoleOne(id int64, hasResource bool) (*Role, error) {
	o := orm.NewOrm()
	m := NewRole(id)
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}

	if hasResource {
		perms := utils.E.GetPermissionsForUser(strconv.FormatInt(m.Id, 10))
		if len(perms) > 0 {
			m.Resources = strings.Join(perms[0], ",")
		}
	}

	return &m, nil
}

//Save 添加、编辑页面 保存
func RoleSave(m *Role, permIds []string) error {
	o := orm.NewOrm()
	if _, err := o.Insert(m); err != nil {
		return err
	}

	for _, permId := range permIds {
		_, err := utils.E.AddPermissionForUser(strconv.FormatInt(m.Id, 10), permId)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("AddPermissionForUser error:%v", err))
			return err
		}
	}

	return nil
}

//Save 添加、编辑页面 保存
func RoleUpdate(m *Role, permIds []string) (*Role, error) {
	o := orm.NewOrm()
	if _, err := o.Update(m, "Name", "UpdatedAt"); err != nil {
		return nil, err
	}
	for _, permId := range permIds {
		_, err := utils.E.AddPermissionForUser(strconv.FormatInt(m.Id, 10), permId)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("AddPermissionForUser error:%v", err))
			return nil, err
		}
	}

	return m, nil
}

//删除
func RoleDelete(id int64) (num int64, err error) {
	m := NewRole(id)
	if num, err := BaseDelete(&m); err != nil {
		utils.LogDebug(fmt.Sprintf("Delete Role:%v", err))
		return num, err
	} else {
		return num, nil
	}
}
