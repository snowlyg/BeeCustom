package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
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
	RealName    string        `orm:"size(32)" valid:"Required;MaxSize(32)"`
	UserName    string        `orm:"size(24)" valid:"Required;MaxSize(24)"`
	UserPwd     string        `orm:"size(256)"`
	Mobile      string        `orm:"size(16)" valid:"Required;Mobile"`
	Email       string        `orm:"size(256)" valid:"Required;Email"`
	Avatar      string        `orm:"size(256)"`
	ICCode      string        `orm:"column(i_c_code);size(255);null"`
	Chapter     string        `orm:"column(chapter);size(255);null" description:"签章"`
	IsSuper     bool          `valid:"Required"`
	Status      bool          `valid:"Required"`
	RoleId      int64         `orm:"-" form:"RoleId" valid:"Required"` //关联管理会自动生成 role_id 字段，此处不生成字段
	Role        *Role         `orm:"rel(fk)"`                          // fk 的反向关系
	Companies   []*Company    `orm:"reverse(many)"`                    //设置一对多关系
	Annotations []*Annotation `orm:"reverse(many)"`                    // 设置多对多的反向关系
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

	prefix := beego.AppConfig.DefaultString("db_dt_prefix", "bee_custom_")
	datas := make([]*BackendUser, 0)
	// 获取 QueryBuilder 对象. 需要指定数据库驱动参数。
	// 第二个返回值是错误对象，在这里略过
	qb1, _ := orm.NewQueryBuilder("mysql")

	// 构建查询对象
	qb1.Select(prefix + "roles.id").
		From(prefix + "roles").
		InnerJoin("role_resource").
		On(prefix + "roles.id = " + "role_resource." + prefix + "roles_id").
		InnerJoin(prefix + "resource").
		On(prefix + "resource.id = " + "role_resource." + prefix + "resource_id").
		And(prefix + "resource.url_for='" + roleResouceString + "'")
	subSql := qb1.String()

	qb, _ := orm.NewQueryBuilder("mysql")
	// 构建查询对象
	qb.Select(prefix+"users.real_name", prefix+"users.id").
		From(prefix + "users").
		LeftJoin(prefix + "roles").
		On(prefix + "users.id").
		In(subSql).
		//Where("age > ?").
		OrderBy("id").Desc()

	// 导出 SQL 语句
	sql := qb.String()

	//SELECT
	//DISTINCT
	//bee_custom_users.real_name,
	//	bee_custom_users.id
	//FROM
	//bee_custom_users
	//INNER JOIN bee_custom_roles
	//ON
	//bee_custom_users.role_id IN (
	//	SELECT
	//bee_custom_roles.id
	//FROM
	//bee_custom_roles
	//INNER JOIN role_resource ON bee_custom_roles.id = role_resource.bee_custom_roles_id
	//INNER JOIN bee_custom_resource ON bee_custom_resource.id = role_resource.bee_custom_resource_id
	//and bee_custom_resource.url_for = 'AnnotationController.Make'
	//
	//)
	//ORDER BY
	//id DESC DISTINCT
	replace := strings.Replace(sql, "SELECT", "SELECT DISTINCT", 1)

	o := orm.NewOrm()
	_, _ = o.Raw(replace).QueryRows(&datas)

	return datas
}

func BackendUserGetRelations(ms []*BackendUser, relations string) ([]*BackendUser, error) {
	o := orm.NewOrm()
	rs := strings.Split(relations, ",")
	for _, v := range ms {
		for _, rv := range rs {
			_, err := o.LoadRelated(v, rv)
			if err != nil {
				utils.LogDebug(fmt.Sprintf("LoadRelated:%v", err))
				return nil, err
			}
		}
	}

	return ms, nil
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

	if m.Role != nil {
		mr := m.Role
		// 获取关系字段，o.LoadRelated(v, "Roles") 这是关键
		// 查找该用户所属的角色
		if _, err := o.LoadRelated(mr, "Resources"); err != nil {
			return nil, err
		}

		m.Role = mr
	} else {
		return &m, errors.New("用户获取失败")
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
func BackendUserSave(m *BackendUser) (*BackendUser, error) {
	o := orm.NewOrm()
	if m.Id == 0 {
		//对密码进行加密
		m.UserPwd = utils.String2md5(m.UserPwd)

		if err := getBackendUserRole(m); err != nil {
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

		if err := getBackendUserRole(m); err != nil {
			return nil, err
		}

		if _, err := o.Update(m); err != nil {
			return nil, err
		}
	}

	return m, nil
}

//获取关联模型
func getBackendUserRole(m *BackendUser) error {
	if oR, err := RoleOne(m.RoleId); err != nil {
		return err
	} else {
		m.Role = oR
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
