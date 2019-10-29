package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
)

// TableName 设置Ciq表名
func (u *Ciq) TableName() string {
	return CiqTBName()
}

// CiqQueryParam 用于查询的类
type CiqQueryParam struct {
	BaseQueryParam

	NameLike string //模糊查询
}

// Ciq 实体类
type Ciq struct {
	BaseModel

	Hs          string    `orm:"column(hs);size(12)" description:"HS编码"`
	Name        string    `orm:"column(name);size(250)" description:"商品中文名称"`
	CiqCode     string    `orm:"column(ciq_code);size(13)" description:"CIQ代码(类别码)"`
	CiqName     string    `orm:"column(ciq_name);size(512);null" description:"CIQ代码中文代码(与HS货物名称差别)"`
	Version     string    `orm:"column(version);size(255)" description:"版本"`
	VersionDate time.Time `orm:"column(version_date);type(datetime)" description:"版本日期"`
	Mark        int8      `orm:"column(mark)" description:"有效标识"`
	Status      string    `orm:"column(status);size(255);null" description:"状态"`
}

func NewCiq(id int64) Ciq {
	return Ciq{BaseModel: BaseModel{Id: id}}
}

//查询参数
func NewCiqQueryParam() CiqQueryParam {
	return CiqQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

// CiqPageList 获取分页数据
func CiqPageList(params *CiqQueryParam) ([]*Ciq, int64) {
	query := orm.NewOrm().QueryTable(CiqTBName())
	datas := make([]*Ciq, 0)

	if len(params.NameLike) > 0 {
		cond := orm.NewCondition()
		cond1 := cond.And("name__istartswith", params.NameLike).
			Or("hs__istartswith", params.NameLike).
			Or("name__istartswith", params.NameLike)
		query = query.SetCond(cond1)
	}

	total, _ := query.Count()
	query = BaseListQuery(query, params.Sort, params.Order, params.Limit, params.Offset)
	_, _ = query.All(&datas)

	return datas, total
}

// CiqOne 根据id获取单条
func CiqOne(id int64) (*Ciq, error) {
	m := NewCiq(0)
	o := orm.NewOrm()
	if err := o.QueryTable(CiqTBName()).Filter("Id", id).RelatedSel().One(&m); err != nil {
		return nil, err
	}

	if &m == nil {
		return &m, errors.New("获取失败")
	}

	return &m, nil
}
