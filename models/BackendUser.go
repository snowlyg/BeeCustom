package models

import (
	"strings"

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
	RealName string `orm:"size(32)" valid:"Required;MaxSize(32)"`
	UserName string `orm:"size(24)" valid:"Required;MaxSize(24)"`
	UserPwd  string `orm:"size(256)" valid:"Required"`
	Mobile   string `orm:"size(16)" valid:"Required;Mobile"`
	Email    string `orm:"size(256)" valid:"Required;Email"`
	Avatar   string `orm:"size(256)"`
	ICCode   string `orm:"column(i_c_code);size(255);null"`
	Chapter  string `orm:"column(chapter);size(255);null" description:"签章"`
	RoleId   int64  `orm:"-" form:"RoleId" valid:"Required"` //关联管理会自动生成 role_id 字段，此处不生成字段
	Role     *Role  `orm:"rel(fk)"`                          // fk 的反向关系
	IsSuper  bool   `valid:"Required"`
	Status   bool   `valid:"Required"`
}

func NewBackendUser(id int64) BackendUser {
	return BackendUser{BaseModel: BaseModel{Id: id}}
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

// BackendUserOne 根据id获取单条
func BackendUserOne(id int64) (*BackendUser, error) {
	m := NewBackendUser(0)
	o := orm.NewOrm()
	if err := o.QueryTable(BackendUserTBName()).Filter("Id", id).RelatedSel().One(&m); err != nil {
		return nil, err
	}

	mr := m.Role
	// 获取关系字段，o.LoadRelated(v, "Roles") 这是关键
	// 查找该用户所属的角色
	if _, err := o.LoadRelated(mr, "Resources"); err != nil {
		return nil, err
	}

	m.Role = mr

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
func BackendUserSave(m *BackendUser) (*BackendUser, error) {
	o := orm.NewOrm()
	if m.Id == 0 {
		//对密码进行加密
		m.UserPwd = utils.String2md5(m.UserPwd)
		if oR, err := RoleOne(m.RoleId); err != nil {
			return nil, err
		} else {
			m.Role = oR
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

		if oR, err := RoleOne(m.RoleId); err != nil {
			return nil, err
		} else {
			m.Role = oR
		}

		if _, err := o.Update(m); err != nil {
			return nil, err
		}
	}

	return m, nil
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
