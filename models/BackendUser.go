package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"BeeCustom/utils"
	"github.com/astaxie/beego/orm"
)

// TableName 设置BackendUser表名
func (u *BackendUser) TableName() string {
	return BackendUserTBName()
}

// BackendUserQueryParam 用于查询的类
type BackendUserQueryParam struct {
	BaseQueryParam

	UserNameLike string //模糊查询
	RealNameLike string //模糊查询
	Mobile       string //精确查询
	SearchStatus string //为空不查询，有值精确查询
}

// BackendUser 实体类
type BackendUser struct {
	BaseModel
	RealName          string              `orm:"size(32)" valid:"Required;MaxSize(32)"`
	UserName          string              `orm:"size(24)" valid:"Required;MaxSize(24)"`
	UserPwd           string              `orm:"size(256)"`
	Mobile            string              `orm:"size(16)" valid:"Required;Mobile"`
	Email             string              `orm:"size(256)" valid:"Required;Email"`
	Avatar            string              `orm:"size(256)"`
	ICCode            string              `orm:"column(i_c_code);size(255);null"`
	Chapter           string              `orm:"column(chapter);size(255);null" description:"签章"`
	IsSuper           bool                `valid:"Required"`
	Status            bool                `valid:"Required"`
	Companies         []*Company          `orm:"reverse(many)"` //设置一对多关系
	Annotations       []*Annotation       `orm:"reverse(many)"` // 设置多对多的反向关系
	AnnotationRecords []*AnnotationRecord `orm:"reverse(many)"` // 设置多对多的反向关系
	RoleIds           []interface{}       `orm:"-"`
	RoleNames         string              `orm:"-"`
}

func NewBackendUser(id int64) BackendUser {
	return BackendUser{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

//查询参数
func NewBackendUserQueryParam() BackendUserQueryParam {
	return BackendUserQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

// BackendUserPageList 获取分页数据
func BackendUserPageList(params *BackendUserQueryParam) ([]*BackendUser, int64) {
	query := orm.NewOrm().QueryTable(BackendUserTBName())
	datas := make([]*BackendUser, 0)
	query = query.Filter("username__istartswith", params.UserNameLike)

	if len(params.SearchStatus) > 0 {
		query = query.Filter("status", params.SearchStatus)
	}

	total, _ := query.Count()
	query = BaseListQuery(query, params.Sort, params.Order, params.Limit, params.Offset)
	_, _ = query.All(&datas)

	return datas, total
}

// GetCreateBackendUsers 制单人
func GetCreateBackendUsers(roleResouceString string) []*BackendUser {
	params := NewBackendUserQueryParam()
	//获取数据列表和总数
	datas, _ := BackendUserPageList(&params)
	for i, v := range datas {
		formatInt := strconv.FormatInt(v.Id, 10)
		hasRoleForUser, _ := utils.E.HasRoleForUser(formatInt, "1") //超级管理员
		if !utils.E.HasPermissionForUser(formatInt, roleResouceString) || hasRoleForUser {
			if i <= len(datas)-1 {
				datas = append(datas[:i], datas[i+1:]...) //删除
			}

		}
	}

	return datas
}

func BackendUsersGetRelations(ms []*BackendUser) ([]*BackendUser, error) {
	if len(ms) > 0 {
		for _, v := range ms {
			err := BackendUserGetRelations(v)
			if err != nil {
				return nil, err
			}
		}
	}

	return ms, nil
}

func BackendUserGetRelations(v *BackendUser) error {

	roleIdStrings, err := utils.E.GetRolesForUser(strconv.FormatInt(v.Id, 10))
	if err != nil {
		utils.LogDebug(fmt.Sprintf("GetRolesForUser error:%v", err))
	}

	var roleNames string
	if len(roleIdStrings) > 0 {
		for _, roleId := range roleIdStrings {
			id64, err := strconv.ParseInt(roleId, 10, 64)
			if err != nil {
				utils.LogDebug(fmt.Sprintf("ParseInt error:%v", err))
			}

			role, err := RoleOne(id64, false)
			if err != nil {
				utils.LogDebug(fmt.Sprintf("ParseInt error:%v", err))
			}

			roleNames += role.Name + ","
		}
	}

	v.RoleNames = roleNames

	return nil
}

// BackenduserDataList 获取用户列表
func BackenduserDataList(params *BackendUserQueryParam) []*BackendUser {
	data, _ := BackendUserPageList(params)
	return data
}

// BackendUserOne 根据id获取单条
func BackendUserOne(id int64) (*BackendUser, error) {
	m := NewBackendUser(0)
	o := orm.NewOrm()
	if err := o.QueryTable(BackendUserTBName()).Filter("Id", id).RelatedSel().One(&m); err != nil {
		return nil, err
	}

	if &m == nil {
		return &m, errors.New("用户获取失败")
	}

	return &m, nil
}

// BackendUserOneByUserName 根据用户名密码获取单条
func BackendUserOneByUserName(username, userpwd string) (*BackendUser, error) {
	m := NewBackendUser(0)
	err := orm.NewOrm().QueryTable(BackendUserTBName()).Filter("username", username).Filter("userpwd", userpwd).One(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

//Save 添加、编辑页面 保存
func BackendUserSave(m *BackendUser, roleIds []string) (*BackendUser, error) {
	o := orm.NewOrm()
	if m.Id == 0 {
		//对密码进行加密
		m.UserPwd = utils.String2md5(m.UserPwd)

		if err := setRoles(m, roleIds); err != nil {
			return nil, err
		}

		if _, err := o.Insert(m); err != nil {
			return nil, err
		}
	} else {
		if oM, err := BackendUserOne(m.Id); err != nil {
			return nil, err
		} else {
			m.UserPwd = strings.TrimSpace(m.UserPwd)
			m.CreatedAt = oM.CreatedAt

			if len(m.UserPwd) == 0 {
				//如果密码为空则不修改
				m.UserPwd = oM.UserPwd
			} else {
				m.UserPwd = utils.String2md5(m.UserPwd)
			}
			//本页面不修改头像和密码，直接将值附给新m
			m.Avatar = oM.Avatar
		}

		if err := setRoles(m, roleIds); err != nil {
			return nil, err
		}

		if _, err := o.Update(m); err != nil {
			return nil, err
		}
	}

	return m, nil
}

//设置角色
func setRoles(m *BackendUser, roleIds []string) error {
	if len(roleIds) > 0 && len(roleIds[0]) > 0 {
		if err := setBackendUserRole(m, roleIds); err != nil {
			return err
		}
	}

	return nil

}

//获取关联模型
func setBackendUserRole(m *BackendUser, roleIds []string) error {
	for _, roleId := range roleIds {
		_, err := utils.E.AddRoleForUser(strconv.FormatInt(m.Id, 10), roleId)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("AddRoleForUser error:%v", err))
			return err
		}
	}

	return nil
}

//Save 添加、编辑页面 保存
func BackendUserFreeze(m *BackendUser) (*BackendUser, error) {
	o := orm.NewOrm()
	if _, err := o.Update(m, "Status"); err != nil {
		return nil, err
	}

	return m, nil
}

//删除
func BackendUserDelete(id int64) (num int64, err error) {
	m := NewBackendUser(id)
	if num, err := BaseDelete(&m); err != nil {
		return num, err
	} else {
		return num, nil
	}
}
