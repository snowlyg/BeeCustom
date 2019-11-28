package models

import (
	"fmt"
	"time"

	"BeeCustom/utils"
	"github.com/astaxie/beego/orm"
)

// TableName 设置AnnotationReturn表名
func (u *AnnotationReturn) TableName() string {
	return AnnotationReturnTBName()
}

// AnnotationReturnQueryParam 用于查询的类
type AnnotationReturnQueryParam struct {
	BaseQueryParam

	AnnotationId int64
}

// AnnotationReturn 实体类
type AnnotationReturn struct {
	BaseModel

	CheckInfo    string `orm:"column(check_info);null;type(text)" description:"回执信息"`
	DealFlag     string `orm:"column(deal_flag);size(10);null" description:"回执代码"`
	EtpsPreentNo string `orm:"column(etps_preent_no);size(64);null" description:"企业预录入编号"`

	BusinessId   string    `orm:"column(business_id);size(64);null" description:"业务编号"`
	ManageResult string    `orm:"column(manage_result);size(1);null" description:"处理结果"`
	Reason       string    `orm:"column(reason);size(255);null" description:"检查信息"`
	CreateDate   time.Time `orm:"column(create_date);type(datetime);null" description:"处理日期"`
	Rmk          string    `orm:"column(rmk);size(255);null" description:"备注"`

	Annotation   *Annotation `orm:"column(annotation_id);rel(fk)"`
	AnnotationId int64       `orm:"-" form:"AnnotationId"` // 关联管理会自动生成字段，此处不生成字段
}

func NewAnnotationReturn(id int64) AnnotationReturn {
	return AnnotationReturn{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

// 查询参数
func NewAnnotationReturnQueryParam() AnnotationReturnQueryParam {
	return AnnotationReturnQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

// AnnotationReturnPageList 获取分页数据
func AnnotationReturnPageList(params *AnnotationReturnQueryParam) ([]*AnnotationReturn, int64) {

	query := orm.NewOrm().QueryTable(AnnotationReturnTBName())
	datas := make([]*AnnotationReturn, 0)

	query = query.Filter("annotation_id", params.AnnotationId).RelatedSel()
	params.Sort = "Id"
	params.Order = "desc"

	total, _ := query.Count()
	query = BaseListQuery(query, params.Sort, params.Order, params.Limit, params.Offset)
	_, _ = query.All(&datas)

	return datas, total
}

// AnnotationReturnOne 根据id获取单条
func AnnotationReturnOneByStatusAndAnnotationId(m *AnnotationReturn) error {
	m.Id = 0
	o := orm.NewOrm()
	if err := o.QueryTable(AnnotationReturnTBName()).
		Filter("annotation_id", m.Annotation.Id).
		RelatedSel().One(m); err != nil {
		return err
	}

	return nil
}

// Save 添加、编辑页面 保存
func AnnotationReturnSave(m *AnnotationReturn) error {
	o := orm.NewOrm()

	if _, err := o.Insert(m); err != nil {
		utils.LogDebug(fmt.Sprintf("AnnotationReturnSave:%v", err))
		return err
	}

	return nil
}
