package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// TableName 设置BackendUser表名
func (a *BackendUser) TableName() string {
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
	Id                 int
	RealName           string    `orm:"size(32)"`
	UserName           string    `orm:"size(24)"`
	UserPwd            string    `json:"-"`
	Mobile             string    `orm:"size(16)"`
	Email              string    `orm:"size(256)"`
	Avatar             string    `orm:"size(256)"`
	ICCode             string    `orm:"column(i_c_code);size(255);null"`
	Chapter            string    `orm:"column(chapter);size(255);null" description:"签章"`
	EnterpriseId       string    `orm:"column(enterprise_id);size(255);null" description:"关联企业ID"`
	RoleId             int       `orm:"-" form:"RoleId"` //关联管理会自动生成 role_id 字段，此处不生成字段
	Role               *Role     `orm:"rel(fk)"`         // fk 的反向关系
	ResourceUrlForList []string  `orm:"-"`
	CreatedAt          time.Time `orm:"column(created_at);type(timestamp);null"`
	UpdatedAt          time.Time `orm:"column(updated_at);type(timestamp);null"`
	IsSuper            bool
	Status             int
}

// BackendUserPageList 获取分页数据
func BackendUserPageList(params *BackendUserQueryParam) ([]*BackendUser, int64) {
	query := orm.NewOrm().QueryTable(BackendUserTBName())

	datas := make([]*BackendUser, 0)

	//默认排序
	sortorder := "Id"
	if len(params.Sort) > 0 {
		sortorder = params.Sort
	}

	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}

	query = query.Filter("username__istartswith", params.UserNameLike)
	query = query.Filter("realname__istartswith", params.UserNameLike)

	if len(params.SearchStatus) > 0 {
		query = query.Filter("status", params.SearchStatus)
	}

	total, _ := query.Count()
	_, _ = query.OrderBy(sortorder).Limit(params.Limit, (params.Offset-1)*params.Limit).RelatedSel().All(&datas)

	return datas, total
}

// BackendUserOne 根据id获取单条
func BackendUserOne(id int) (*BackendUser, error) {
	o := orm.NewOrm()
	m := BackendUser{Id: id}

	err := o.QueryTable(BackendUserTBName()).RelatedSel().One(&m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

// BackendUserOneByUserName 根据用户名密码获取单条
func BackendUserOneByUserName(username, userpwd string) (*BackendUser, error) {
	m := BackendUser{}
	err := orm.NewOrm().QueryTable(BackendUserTBName()).Filter("username", username).Filter("userpwd", userpwd).One(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
