package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
)

// TableName 设置CompanyForeign表名
func (u *CompanyForeign) TableName() string {
	return CompanyForeignTBName()
}

// CompanyForeign 实体类
type CompanyForeign struct {
	BaseModel

	ForeignCompanyName    string   `orm:"column(foreign_company_name);size(200)" description:"外商公司名称"`
	ForeignCompanyPhone   string   `orm:"column(foreign_company_phone);size(200);null" description:"外商公司电话"`
	ForeignCompanyAddress string   `orm:"column(foreign_company_address);size(200);null" description:"外商公司地址"`
	ForeignCompanyChapter string   `orm:"column(foreign_company_chapter);size(255);null" description:"外商公司章"`
	ForeignType           string   `orm:"column(foreign_type);size(1)" description:"关联公司类型"`
	Company               *Company `orm:"column(company_id);rel(fk)"`
	CompanyId             int64    `orm:"-" form:"CompanyId"`
}

// CompanyForeignQueryParam 用于查询的类
type CompanyForeignQueryParam struct {
	BaseQueryParam

	CompanyId string
}

func NewCompanyForeign(id int64) CompanyForeign {
	return CompanyForeign{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

//查询参数
func NewCompanyForeignQueryParam() CompanyForeignQueryParam {
	return CompanyForeignQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

// CompanyPageList 获取分页数据
func CompanyForeignPageList(params *CompanyForeignQueryParam) ([]*CompanyForeign, int64) {
	query := orm.NewOrm().QueryTable(CompanyForeignTBName())
	datas := make([]*CompanyForeign, 0)

	query = query.Filter("company_id", params.CompanyId)

	total, _ := query.Count()
	query = BaseListQuery(query, params.Sort, params.Order, params.Limit, params.Offset)
	_, _ = query.All(&datas)

	return datas, total
}

// CompanyForeignOne 根据id获取单条
func CompanyForeignOne(id int64) (*CompanyForeign, error) {
	m := NewCompanyForeign(0)
	o := orm.NewOrm()
	if err := o.QueryTable(CompanyForeignTBName()).Filter("Id", id).RelatedSel().One(&m); err != nil {
		return nil, err
	}

	if &m == nil {
		return &m, errors.New("获取失败")
	}

	return &m, nil
}

//Save 添加、编辑页面 保存
func CompanyForeignSave(m *CompanyForeign) (*CompanyForeign, error) {
	o := orm.NewOrm()
	if m.Id == 0 {
		if err := getCompanyForeignBackendUser(m); err != nil {
			return nil, err
		}

		if _, err := o.Insert(m); err != nil {
			return nil, err
		}
	} else {
		if err := getCompanyForeignBackendUser(m); err != nil {
			return nil, err
		}

		if _, err := o.Update(m); err != nil {
			return nil, err
		}
	}

	return m, nil
}

//获取关联模型
func getCompanyForeignBackendUser(m *CompanyForeign) error {
	if bU, err := CompanyOne(m.CompanyId, false); err != nil {
		return err
	} else {
		m.Company = bU
	}
	return nil
}

//删除
func CompanyForeignDelete(id int64) (num int64, err error) {
	m := NewCompanyForeign(id)
	if num, err := BaseDelete(&m); err != nil {
		return num, err
	} else {
		return num, nil
	}
}
