package models

import (
	"github.com/astaxie/beego/orm"
)

// TableName 设置HsCode表名
func (u *HsCode) TableName() string {
	return HsCodeTBName()
}

// HsCodeQueryParam 用于查询的类
type HsCodeQueryParam struct {
	BaseQueryParam

	NameLike string //模糊查询
}

// HsCode 实体类
type HsCode struct {
	BaseModel

	Code        string  `orm:"column(code);size(255)" description:"编码"`
	Name        string  `orm:"column(name);size(255)" description:"名称"`
	License     string  `orm:"column(license);size(255);null" description:"许可证代码"`
	GeneralRate float64 `orm:"column(general_rate);null;digits(17);decimals(4)" description:"普通税率"`
	OfferRate   float64 `orm:"column(offer_rate);null;digits(17);decimals(4)" description:"优惠税率"`
	ExportRate  float64 `orm:"column(export_rate);null;digits(17);decimals(4)" description:"出口税率"`
	TaxRate     float64 `orm:"column(tax_rate);null;digits(17);decimals(4)" description:"增值税率"`
	ConsumeRate float64 `orm:"column(consume_rate);null;digits(17);decimals(4)" description:"消费税率"`
	Unit1       string  `orm:"column(unit1);size(255);null" description:"第一法定单位"`
	Unit2       string  `orm:"column(unit2);size(255);null" description:"第二法定单位"`
	Declaration string  `orm:"column(declaration);size(255);null" description:"申报要素"`
	Remark      string  `orm:"column(remark);size(255);null" description:"备注"`
	Unit1Name   string  `orm:"column(unit1_name);size(255);null" description:"第一单位名称"`
	Unit2Name   string  `orm:"column(unit2_name);size(255);null" description:"第二单位名称"`
}

func NewHsCode(id int64) HsCode {
	return HsCode{BaseModel: BaseModel{Id: id}}
}

//查询参数
func NewHsCodeQueryParam() HsCodeQueryParam {
	return HsCodeQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

// HsCodePageList 获取分页数据
func HsCodePageList(params *HsCodeQueryParam) ([]*HsCode, int64) {
	query := orm.NewOrm().QueryTable(HsCodeTBName())
	datas := make([]*HsCode, 0)

	if len(params.NameLike) > 0 {
		cond := orm.NewCondition()
		cond1 := cond.And("name__istartswith", params.NameLike).
			Or("code__istartswith", params.NameLike)
		query = query.SetCond(cond1)
	}

	total, _ := query.Count()
	query = BaseListQuery(query, params.Sort, params.Order, params.Limit, params.Offset)
	_, _ = query.All(&datas)

	return datas, total
}

// GetHsCodeByCode
func GetHsCodeByCode(hsCodeS string) ([]*HsCode, error) {
	datas := make([]*HsCode, 0)
	query := orm.NewOrm().QueryTable(HsCodeTBName())
	query = query.Distinct().Filter("code", hsCodeS)
	if _, err := query.All(&datas); err != nil {
		return nil, err
	}

	return datas, nil
}

// 批量删除
func HsCodeDeleteAll() (num int64, err error) {
	o := orm.NewOrm()
	if num, err := o.QueryTable(HsCodeTBName()).Filter("code__isnull", false).Delete(); err != nil {
		return num, err
	} else {
		return num, nil
	}

}

// 批量插入
func InsertHsCodeMulti(datas []*HsCode) (num int64, err error) {
	return BaseInsertMulti(datas)
}
