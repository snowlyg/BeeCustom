package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
)

// TableName 设置Clearance表名
func (u *Clearance) TableName() string {
	return ClearanceTBName()
}

// Clearance 实体类
type Clearance struct {
	BaseModel

	Type                int8   `orm:"column(type)" description:"参数类别"`
	CustomsCode         string `orm:"column(customs_code);size(255)" description:"海关编码"  valid:"Required;MaxSize(255)"`
	Name                string `orm:"column(name);size(255)" description:"名称"  valid:"Required;MaxSize(255)"`
	ShortName           string `orm:"column(short_name);size(255);null" description:"简称"`
	EnName              string `orm:"column(en_name);size(255);null" description:"英文名称"`
	InspectionCode      string `orm:"column(inspection_code);size(255);null" description:"商检编码"`
	ShortEnName         string `orm:"column(short_en_name);size(255);null" description:"英文简称"`
	MandatoryLevel      string `orm:"column(mandatory_level);size(255);null" description:"强制级别(企业产品许可类别)"`
	CertificateType     string `orm:"column(certificate_type);size(255);null" description:"证书类别(企业产品许可类别)"`
	StatisticalUnitCode string `orm:"column(statistical_unit_code);size(255);null" description:"对应统计计量单位代码(计量单位代码表)"`
	ConversionRate      string `orm:"column(conversion_rate);size(255);null" description:"换算率(计量单位代码表)"`
	NatureMark          string `orm:"column(nature_mark);size(255);null" description:"国内地区性质标记(国内地区代码)"`
	Iso2                string `orm:"column(iso2);size(255);null" description:"iso2(原产地区代码表)"`
	Iso3                string `orm:"column(iso3);size(255);null" description:"iso3(原产地区代码表)"`
	TypeCode            string `orm:"column(type_code);size(255);null" description:"分类代码(原产地区代码表)"`
	OldCustomCode       string `orm:"column(old_custom_code);size(255);null" description:"原报关代码"`
	OldCustomName       string `orm:"column(old_custom_name);size(255);null" description:"原报关名称"`
	OldCiqCode          string `orm:"column(old_ciq_code);size(255);null" description:"原报检代码"`
	OldCiqName          string `orm:"column(old_ciq_name);size(255);null" description:"原报检名称"`
	Remark              string `orm:"column(remark);size(1000);null" description:"备注"`
}

// ClearanceQueryParam 用于查询的类
type ClearanceQueryParam struct {
	BaseQueryParam
	Type     string //模糊查询
	NameLike string //模糊查询
}

func NewClearance(id int64) Clearance {
	return Clearance{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

//查询参数
func NewClearanceQueryParam() ClearanceQueryParam {
	return ClearanceQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

// ClearancePageList 获取分页数据
func ClearancePageList(params *ClearanceQueryParam) ([]*Clearance, int64) {
	query := orm.NewOrm().QueryTable(ClearanceTBName())
	datas := make([]*Clearance, 0)

	clearanceType := "0"
	if len(params.Type) > 0 {
		clearanceType = params.Type
	}

	query = query.Filter("type", clearanceType)

	if len(params.NameLike) > 0 {
		cond := orm.NewCondition()
		cond1 := cond.And("customs_code__istartswith", params.NameLike).
			Or("name__istartswith", params.NameLike).
			Or("short_name__istartswith", params.NameLike).
			Or("en_name__istartswith", params.NameLike)
		query = query.SetCond(cond1)
	}

	total, _ := query.Count()
	query = BaseListQuery(query, params.Sort, params.Order, params.Limit, params.Offset)
	_, _ = query.All(&datas)

	return datas, total
}

// ClearanceOne 根据id获取单条
func ClearanceOne(id int64) (*Clearance, error) {
	m := NewClearance(0)
	o := orm.NewOrm()
	if err := o.QueryTable(ClearanceTBName()).Filter("Id", id).RelatedSel().One(&m); err != nil {
		return nil, err
	}

	if &m == nil {
		return &m, errors.New("获取失败")
	}

	return &m, nil
}

//Save 添加、编辑页面 保存
func ClearanceSave(m *Clearance) (*Clearance, error) {
	o := orm.NewOrm()
	if m.Id == 0 {
		if _, err := o.Insert(m); err != nil {
			return nil, err
		}
	} else {
		if _, err := o.Update(m); err != nil {
			return nil, err
		}
	}

	return m, nil
}

//删除
func ClearanceDelete(id int64) (num int64, err error) {
	m := NewClearance(id)
	if num, err := BaseDelete(&m); err != nil {
		return num, err
	} else {
		return num, nil
	}
}

//删除
func ClearanceDeleteAll(clearanceType int8) (num int64, err error) {
	if num, err := BaseDeleteAll(clearanceType); err != nil {
		return num, err
	} else {
		return num, nil
	}
}

//批量插入
func InsertClearanceMulti(datas []*Clearance) (num int64, err error) {
	return BaseInsertMulti(len(datas), datas)
}
