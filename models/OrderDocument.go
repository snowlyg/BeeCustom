package models

import (
	"fmt"
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

// Save 添加、编辑页面 保存
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

// 删除
func OrderDocumentDelete(id int64) (num int64, err error) {
	m := NewOrderDocument(id)
	if num, err := BaseDelete(&m); err != nil {
		return num, err
	} else {
		return num, nil
	}
}
