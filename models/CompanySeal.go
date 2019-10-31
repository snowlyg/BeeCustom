package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
)

// TableName 设置CompanySeal表名
func (u *CompanySeal) TableName() string {
	return CompanySealTBName()
}

// CompanySeal 实体类
type CompanySeal struct {
	BaseModel

	Url       string   `orm:"column(url);size(255);null" description:"签章地址"`
	SealName  string   `orm:"column(seal_name);" description:"签章类型 1：公章，2：合同章"`
	Company   *Company `orm:"column(company_id);rel(fk)"`
	CompanyId int64    `orm:"-" form:"CompanyId"`
}

// CompanySealQueryParam 用于查询的类
type CompanySealQueryParam struct {
	BaseQueryParam

	CompanyId string
	SealName  string
}

func NewCompanySeal(id int64) CompanySeal {
	return CompanySeal{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

//查询参数
func NewCompanySealQueryParam() CompanySealQueryParam {
	return CompanySealQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

// CompanyPageList 获取分页数据
func CompanySealPageList(params *CompanySealQueryParam) ([]*CompanySeal, int64) {
	query := orm.NewOrm().QueryTable(CompanySealTBName())
	datas := make([]*CompanySeal, 0)

	query = query.Filter("company_id", params.CompanyId)
	if len(params.SealName) > 0 {
		query = query.Filter("seal_name", params.SealName)
	}

	total, _ := query.Count()
	query = BaseListQuery(query, params.Sort, params.Order, params.Limit, params.Offset)
	_, _ = query.All(&datas)

	return datas, total
}

// CompanySealOne 根据id获取单条
func CompanySealOne(id int64) (*CompanySeal, error) {
	m := NewCompanySeal(0)
	o := orm.NewOrm()
	if err := o.QueryTable(CompanySealTBName()).Filter("Id", id).RelatedSel().One(&m); err != nil {
		return nil, err
	}

	if &m == nil {
		return &m, errors.New("获取失败")
	}

	return &m, nil
}

//Save 添加、编辑页面 保存
func CompanySealSave(m *CompanySeal) (*CompanySeal, error) {
	o := orm.NewOrm()

	if m.Id == 0 {
		if err := getCompanySealBackendUser(m); err != nil {
			return nil, err
		}

		if _, err := o.Insert(m); err != nil {
			return nil, err
		}
	} else {
		if err := getCompanySealBackendUser(m); err != nil {
			return nil, err
		}

		if _, err := o.Update(m); err != nil {
			return nil, err
		}
	}

	return m, nil
}

//获取关联模型
func getCompanySealBackendUser(m *CompanySeal) error {
	if bU, err := CompanyOne(m.CompanyId, false); err != nil {
		return err
	} else {
		m.Company = bU
	}
	return nil
}

//删除
func CompanySealDelete(id int64) (num int64, err error) {
	m := NewCompanySeal(id)
	if num, err := BaseDelete(&m); err != nil {
		return num, err
	} else {
		return num, nil
	}
}
