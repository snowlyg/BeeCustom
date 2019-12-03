package models

import (
	"fmt"
	"time"

	"BeeCustom/utils"
	"github.com/astaxie/beego/orm"
)

// TableName 设置OrderRecord表名
func (u *OrderRecord) TableName() string {
	return OrderRecordTBName()
}

// OrderRecordQueryParam 用于查询的类
type OrderRecordQueryParam struct {
	BaseQueryParam

	OrderId int64
}

// OrderRecord 实体类
type OrderRecord struct {
	BaseModel

	Content       string       `orm:"column(content)" description:"办理记录内容"`
	Status        string       `orm:"column(status);size(255);null" description:"办理记录内容时状态"`
	Remark        string       `orm:"column(remark);size(255);null" description:"备注"`
	BackendUser   *BackendUser `orm:"column(user_id);rel(fk)"`
	BackendUserId int64        `orm:"-" form:"BackendUserId"` // 关联管理会自动生成字段，此处不生成字段
	Order         *Order       `orm:"column(order_id);rel(fk)"`
	OrderId       int64        `orm:"-" form:"OrderId"` // 关联管理会自动生成字段，此处不生成字段
}

func NewOrderRecord(id int64) OrderRecord {
	return OrderRecord{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

// 查询参数
func NewOrderRecordQueryParam() OrderRecordQueryParam {
	return OrderRecordQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

// OrderRecordPageList 获取分页数据
func OrderRecordPageList(params *OrderRecordQueryParam) ([]*OrderRecord, int64) {
	query := orm.NewOrm().QueryTable(OrderRecordTBName())
	datas := make([]*OrderRecord, 0)
	query = query.Filter("order_id", params.OrderId).RelatedSel()
	params.Sort = "Id"
	params.Order = "desc"

	total, _ := query.Count()
	query = BaseListQuery(query, params.Sort, params.Order, params.Limit, params.Offset)
	_, _ = query.All(&datas)

	return datas, total
}

// OrderRecordOne 根据id获取单条
func OrderRecordOneByStatusAndOrderId(m *OrderRecord, status string) error {
	m.Id = 0
	o := orm.NewOrm()
	if err := o.QueryTable(OrderRecordTBName()).
		Filter("order_id", m.Order.Id).
		Filter("status", status).
		RelatedSel().One(m); err != nil {
		return err
	}

	return nil
}

// Save 添加、编辑页面 保存
func OrderRecordSave(m *OrderRecord) error {
	o := orm.NewOrm()
	err := OrderRecordOneByStatusAndOrderId(m, m.Status)
	if err != nil && err.Error() != "<QuerySeter> no row found" {
		utils.LogDebug(fmt.Sprintf("OrderRecordOneByStatusAndOrderId:%v", err))
		return err
	}

	if m.Id == 0 {
		if _, err := o.Insert(m); err != nil {
			utils.LogDebug(fmt.Sprintf("OrderRecordSave:%v", err))
			return err
		}
	} else {
		if _, err := o.Update(m); err != nil {
			utils.LogDebug(fmt.Sprintf("OrderRecordUpdate:%v", err))
			return err
		}
	}

	return nil
}
