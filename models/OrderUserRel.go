package models

import (
	"BeeCustom/utils"
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

// OrderUserRelQueryParam 用于查询的类
type OrderUserRelQueryParam struct {
	BaseQueryParam

	UserType      string
	BackendUserId string
	OrderId       string
}

// OrderUserRel 实体类
type OrderUserRel struct {
	BaseModel

	Order       *Order       `orm:"rel(fk)"` //设置一对多关系
	BackendUser *BackendUser `orm:"rel(fk)"` //设置一对多关系
	UserType    int8         `orm:"column(type)"`
}

func NewOrderUserRel(id int64) OrderUserRel {
	return OrderUserRel{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

//查询参数
func NewOrderUserRelQueryParam() OrderUserRelQueryParam {
	return OrderUserRelQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

// BackendUserOne 根据id获取单条
func OrderUserRelByUserIdAndOrderId(userId, orderId int64, userType int8) (*OrderUserRel, error) {
	m := NewOrderUserRel(0)
	o := orm.NewOrm()
	if err := o.QueryTable(OrderUserRelTBName()).
		Filter("backend_user_id", userId).
		Filter("order_id", orderId).
		Filter("type", userType).
		One(&m); err != nil && err.Error() != "<QuerySeter> no row found" {
		utils.LogDebug(fmt.Sprintf("OrderUserRelByUserIdAndOrderId:%v", err))
		return nil, err
	}

	if &m == nil {
		return &m, errors.New("用户获取失败")
	}

	return &m, nil
}

//Save 添加、编辑页面 保存
func OrderUserRelSave(m *OrderUserRel) error {
	o := orm.NewOrm()

	if err := o.Read(m); err != nil && err.Error() != "<QuerySeter> no row found" {
		utils.LogDebug(fmt.Sprintf("OrderUserRelRead:%v", err))
		return err
	}

	if m.Id != 0 {
		if _, err := o.Update(m); err != nil {
			utils.LogDebug(fmt.Sprintf("OrderUserRelInsert:%v", err))
			return err
		}

	} else {
		if _, err := o.Insert(m); err != nil {
			utils.LogDebug(fmt.Sprintf("OrderUserRelInsert:%v", err))
			return err
		}

	}

	return nil
}
