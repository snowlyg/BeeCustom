package models

import (
	"fmt"
	"time"

	"BeeCustom/utils"
	"github.com/astaxie/beego/orm"
)

// TableName 设置OrderReturn表名
func (u *OrderReturn) TableName() string {
	return OrderReturnTBName()
}

// OrderReturnQueryParam 用于查询的类
type OrderReturnQueryParam struct {
	BaseQueryParam

	OrderId int64
}

// OrderReturn 实体类
type OrderReturn struct {
	BaseModel

	Channel    string    `orm:"column(channel);size(255);null" description:"代码：回执代码"`
	CusCiqNo   string    `orm:"column(cus_ciq_no);size(50);null" description:"统一编号"`
	EntryId    string    `orm:"column(entry_id);size(50);null" description:"海关编号(报关单号)"`
	Note       string    `orm:"column(note);size(1000);null" description:"回执内容"`
	NoticeDate time.Time `orm:"column(notice_date);type(datetime);null" description:"回执时间"`
	Remark     string    `orm:"column(remark);size(500);null" description:"备注"`

	DeletedAt time.Time `form:"-" orm:"column(deleted_at);type(timestamp);null" `

	Order   *Order `orm:"column(annotation_id);rel(fk)"`
	OrderId int64  `orm:"-" form:"OrderId"` // 关联管理会自动生成字段，此处不生成字段
}

func NewOrderReturn(id int64) OrderReturn {
	return OrderReturn{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

// 查询参数
func NewOrderReturnQueryParam() OrderReturnQueryParam {
	return OrderReturnQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

// OrderReturnPageList 获取分页数据
func OrderReturnPageList(params *OrderReturnQueryParam) ([]*OrderReturn, int64) {

	query := orm.NewOrm().QueryTable(OrderReturnTBName())
	datas := make([]*OrderReturn, 0)

	query = query.Filter("annotation_id", params.OrderId).RelatedSel()
	params.Sort = "Id"
	params.Order = "desc"

	total, _ := query.Count()
	query = BaseListQuery(query, params.Sort, params.Order, params.Limit, params.Offset)
	_, _ = query.All(&datas)

	return datas, total
}

// OrderReturnOne 根据id获取单条
func OrderReturnOneByStatusAndOrderId(m *OrderReturn) error {
	m.Id = 0
	o := orm.NewOrm()
	if err := o.QueryTable(OrderReturnTBName()).
		Filter("annotation_id", m.Order.Id).
		RelatedSel().One(m); err != nil {
		return err
	}

	return nil
}

// Save 添加、编辑页面 保存
func OrderReturnSave(m *OrderReturn) error {
	o := orm.NewOrm()

	if _, err := o.Insert(m); err != nil {
		utils.LogDebug(fmt.Sprintf("OrderReturnSave:%v", err))
		return err
	}

	return nil
}
