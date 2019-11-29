package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
)

// TableName 设置CompanyContact表名
func (u *CompanyContact) TableName() string {
	return CompanyContactTBName()
}

// CompanyContact 实体类
type CompanyContact struct {
	BaseModel

	Name          string   `orm:"column(name);size(50)" description:"姓名" valid:"Required;MaxSize(255)"`
	Email         string   `orm:"column(email);size(100)" description:"email"valid:"Required;Email"`
	Phone         string   `orm:"column(phone);size(13)" description:"手机" valid:"Required;Mobile"`
	Password      string   `orm:"column(password);size(100)" description:"密码"`
	Offer         string   `orm:"column(offer);size(20)" description:"职位:法定代表，报关员，财务，其它" valid:"Required"`
	ICCode        string   `orm:"column(i_c_code);size(20);null" description:"ic 卡号"`
	Remark        string   `orm:"column(remark);size(500);null" description:"备注"`
	IsAdmin       int8     `orm:"column(is_admin)" description:"是否是企业管理员"`
	Frozen        int8     `orm:"column(frozen)" description:"是否禁用" valid:"Required"`
	RememberToken string   `orm:"column(remember_token);size(100);null"`
	Chapter       string   `orm:"column(chapter);size(255);null" description:"签章"`
	Company       *Company `orm:"column(company_id);rel(fk)"`
	CompanyId     int64    `orm:"-" form:"CompanyId"`
}

// CompanyContactQueryParam 用于查询的类
type CompanyContactQueryParam struct {
	BaseQueryParam

	CompanyId string
	IsAdmin   bool
}

func NewCompanyContact(id int64) CompanyContact {
	return CompanyContact{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

//查询参数
func NewCompanyContactQueryParam() CompanyContactQueryParam {
	return CompanyContactQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

// CompanyContactPageList 获取分页数据
func CompanyContactPageList(params *CompanyContactQueryParam) ([]*CompanyContact, int64) {
	query := orm.NewOrm().QueryTable(CompanyContactTBName())
	datas := make([]*CompanyContact, 0)

	if params.IsAdmin {
		query = query.Filter("is_admin", params.IsAdmin)
	}

	query = query.Filter("company_id", params.CompanyId)

	total, _ := query.Count()
	query = BaseListQuery(query, params.Sort, params.Order, params.Limit, params.Offset)
	_, _ = query.All(&datas)

	return datas, total
}

// CompanyContactOne 根据id获取单条
func CompanyContactOne(id int64) (*CompanyContact, error) {
	m := NewCompanyContact(0)
	o := orm.NewOrm()
	if err := o.QueryTable(CompanyContactTBName()).Filter("Id", id).RelatedSel().One(&m); err != nil {
		return nil, err
	}

	if &m == nil {
		return &m, errors.New("获取失败")
	}

	return &m, nil
}

// GetAdminCompanyContactByCompanyId 根据id获取单条
func GetAdminCompanyContactByCompanyId(id int64) (*CompanyContact, error) {
	m := NewCompanyContact(0)
	o := orm.NewOrm()
	if err := o.QueryTable(CompanyContactTBName()).Filter("company_id", id).Filter("is_admin", true).One(&m); err != nil {
		return nil, err
	}

	if &m == nil {
		return &m, errors.New("获取失败")
	}

	return &m, nil
}

//Save 添加、编辑页面 保存
func CompanyContactSave(m *CompanyContact) (*CompanyContact, error) {
	o := orm.NewOrm()
	if m.Id == 0 {
		if err := getCompanyContactCompany(m); err != nil {
			return nil, err
		}

		if _, err := o.Insert(m); err != nil {
			return nil, err
		}
	} else {
		if err := getCompanyContactCompany(m); err != nil {
			return nil, err
		}

		if _, err := o.Update(m); err != nil {
			return nil, err
		}
	}

	return m, nil
}

//获取关联模型
func getCompanyContactCompany(m *CompanyContact) error {
	if bU, err := CompanyOne(m.CompanyId, ""); err != nil {
		return err
	} else {
		m.Company = bU
	}

	return nil
}

//删除
func CompanyContactDelete(id int64) (num int64, err error) {
	m := NewCompanyContact(id)
	if num, err := BaseDelete(&m); err != nil {
		return num, err
	} else {
		return num, nil
	}
}
