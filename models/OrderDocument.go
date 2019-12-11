package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"BeeCustom/utils"
	"github.com/astaxie/beego/orm"
)

// TableName 设置OrderDocument表名
func (u *OrderDocument) TableName() string {
	return OrderDocumentTBName()
}

// OrderDocumentQueryParam 用于查询的类
type OrderDocumentQueryParam struct {
	BaseQueryParam

	OrderId int64
}

// OrderDocument 实体类
type OrderDocument struct {
	BaseModel
	DocuCode     string    `orm:"column(docu_code);size(3);null" description:"随附单证代码"`
	DocuCodeName string    `orm:"column(docu_code_name);size(100);null" description:"随附单证代码名称"`
	CertCode     string    `orm:"column(cert_code);size(32);null" description:"随附单证编码"`
	DeletedAt    time.Time `orm:"column(deleted_at);type(timestamp);null"`
	Order        *Order    `orm:"column(order_id);rel(fk)"`
	OrderId      int64     `orm:"-" form:"OrderId"` // 关联管理会自动生成 CompanyId 字段，此处不生成字段
}

func NewOrderDocument(id int64) OrderDocument {
	return OrderDocument{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

//查询参数
func NewOrderDocumentQueryParam() OrderDocumentQueryParam {
	return OrderDocumentQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

// OrderDocumentPageList 获取分页数据
func OrderDocumentPageList(params *OrderDocumentQueryParam) ([]*OrderDocument, int64) {

	query := orm.NewOrm().QueryTable(OrderDocumentTBName())
	datas := make([]*OrderDocument, 0)

	query = query.Filter("order_id", params.OrderId)

	total, _ := query.Count()
	query = BaseListQuery(query, params.Sort, params.Order, params.Limit, params.Offset)
	_, _ = query.All(&datas)

	return datas, total
}

// OrderDocumentPageList 获取分页数据
func OrderDocumentsByOrderId(aId int64) ([]*OrderDocument, error) {

	datas := make([]*OrderDocument, 0)
	_, err := orm.NewOrm().QueryTable(OrderDocumentTBName()).Filter("order_id", aId).All(&datas)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("OrderDocumentsByOrderId error :%v", err))
		return nil, err
	}

	return datas, nil
}

func OrderDocumentGetRelations(ms []*OrderDocument, relations string) ([]*OrderDocument, error) {
	if len(relations) > 0 {
		o := orm.NewOrm()
		rs := strings.Split(relations, ",")
		for _, v := range ms {
			for _, rv := range rs {
				_, err := o.LoadRelated(v, rv)
				if err != nil {
					utils.LogDebug(fmt.Sprintf("LoadRelated:%v", err))
					return nil, err
				}
			}
		}
	}
	return ms, nil
}

// OrderDocumentOne 根据id获取单条
func OrderDocumentOne(id int64) (*OrderDocument, error) {
	m := NewOrderDocument(0)
	o := orm.NewOrm()
	if err := o.QueryTable(OrderDocumentTBName()).Filter("Id", id).RelatedSel().One(&m); err != nil {
		return nil, err
	}

	if &m == nil {
		return &m, errors.New("数据获取失败")
	}

	return &m, nil
}

//Save 添加、编辑页面 保存
func OrderDocumentSave(m *OrderDocument) error {
	o := orm.NewOrm()
	if m.Id == 0 {
		if _, err := o.Insert(m); err != nil {
			utils.LogDebug(fmt.Sprintf("OrderDocumentSave:%v", err))
			return err
		}
	} else {

		if _, err := o.Update(m); err != nil {
			utils.LogDebug(fmt.Sprintf("OrderDocumentSave:%v", err))
			return err
		}
	}

	return nil
}

//删除
func OrderDocumentDelete(id int64) (num int64, err error) {
	m := NewOrderDocument(id)
	if num, err := BaseDelete(&m); err != nil {
		return num, err
	} else {
		return num, nil
	}
}
