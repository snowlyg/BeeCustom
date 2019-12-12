package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"BeeCustom/utils"
	"github.com/astaxie/beego/orm"
)

// TableName 设置OrderItemLimit表名
func (u *OrderItemLimit) TableName() string {
	return OrderItemLimitTBName()
}

// OrderItemLimitQueryParam 用于查询的类
type OrderItemLimitQueryParam struct {
	BaseQueryParam

	OrderId int64
}

// OrderItemLimit 实体类
type OrderItemLimit struct {
	BaseModel

	GoodsNo             string `orm:"column(goods_no);size(9)" description:"序号"`
	LicTypeCode         string `orm:"column(lic_type_code);size(5);null" description:"许可证类别代码"`
	LicTypeName         string `orm:"column(lic_type_name);size(100);null" description:"许可证类别名称"`
	LicenceNo           string `orm:"column(licence_no);size(40);null" description:"许可证编码"`
	LicWrtofDetailNo    string `orm:"column(lic_wrtof_detail_no);size(4);null" description:"核销货物序号"`
	LicWrtofQty         int    `orm:"column(lic_wrtof_qty);null" description:"核销数量"`
	LicWrtofQtyUnit     string `orm:"column(lic_wrtof_qty_unit);size(3);null" description:"核销数量单位"`
	LicWrtofQtyUnitName string `orm:"column(lic_wrtof_qty_unit_name);size(50);null" description:"核销数量单位名称"`

	DeletedAt time.Time `orm:"column(deleted_at);type(timestamp);null"`

	OrderItem   *OrderItem `orm:"column(order_item_id);rel(fk)"`
	OrderItemId int64      `orm:"-" form:"OrderItemId"` //关联管理会自动生成 CompanyId 字段，此处不生成字段
}

func NewOrderItemLimit(id int64) OrderItemLimit {
	return OrderItemLimit{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

//查询参数
func NewOrderItemLimitQueryParam() OrderItemLimitQueryParam {
	return OrderItemLimitQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

// OrderItemLimitPageList 获取分页数据
func OrderItemLimitPageList(params *OrderItemLimitQueryParam) ([]*OrderItemLimit, int64) {

	query := orm.NewOrm().QueryTable(OrderItemLimitTBName())
	datas := make([]*OrderItemLimit, 0)

	query = query.Filter("order_id", params.OrderId)

	total, _ := query.Count()
	query = BaseListQuery(query, params.Sort, params.Order, params.Limit, params.Offset)
	_, _ = query.All(&datas)

	return datas, total
}

// OrderItemLimitPageList 获取分页数据
func OrderItemLimitsByOrderId(aId int64) ([]*OrderItemLimit, error) {

	datas := make([]*OrderItemLimit, 0)
	_, err := orm.NewOrm().QueryTable(OrderItemLimitTBName()).Filter("order_id", aId).All(&datas)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("OrderItemLimitsByOrderId error :%v", err))
		return nil, err
	}

	return datas, nil
}

func OrderItemLimitGetRelations(ms []*OrderItemLimit, relations string) ([]*OrderItemLimit, error) {
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

// OrderItemLimitOne 根据id获取单条
func OrderItemLimitOne(id int64) (*OrderItemLimit, error) {
	m := NewOrderItemLimit(0)
	o := orm.NewOrm()
	if err := o.QueryTable(OrderItemLimitTBName()).Filter("Id", id).RelatedSel().One(&m); err != nil {
		return nil, err
	}

	if &m == nil {
		return &m, errors.New("数据获取失败")
	}

	return &m, nil
}

//Save 添加、编辑页面 保存
func OrderItemLimitSave(m *OrderItemLimit) error {
	o := orm.NewOrm()

	//进出口原产国和目的国是相反的数据
	//if m.Order.ImpexpMarkcd == "E" {
	//	natcd := m.Natcd
	//	natcdName := m.NatcdName
	//	destinationNatcd := m.DestinationNatcd
	//	destinationNatcdName := m.DestinationNatcdName
	//
	//	m.Natcd = destinationNatcd
	//	m.NatcdName = destinationNatcdName
	//	m.DestinationNatcd = natcd
	//	m.DestinationNatcdName = natcdName
	//}

	if m.Id == 0 {
		if _, err := o.Insert(m); err != nil {
			utils.LogDebug(fmt.Sprintf("OrderItemLimitSave:%v", err))
			return err
		}
	} else {

		if _, err := o.Update(m); err != nil {
			utils.LogDebug(fmt.Sprintf("OrderItemLimitSave:%v", err))
			return err
		}
	}

	return nil
}

//OrderItemLimitUpdateAll 添加、编辑页面 保存
func OrderItemLimitUpdateAll(aid int64, m *OrderItemLimit) error {
	o := orm.NewOrm()
	qs := o.QueryTable(OrderItemLimitTBName()).Filter("order_id", aid)

	var params orm.Params
	//if len(m.Natcd) > 0 {
	//	params = orm.Params{
	//		"dcl_currcd":      m.DclCurrcd,
	//		"dcl_currcd_name": m.DclCurrcdName,
	//		"natcd":           m.Natcd,
	//		"natcd_name":      m.NatcdName,
	//	}
	//} else if len(m.DestinationNatcd) > 0 {
	//	params = orm.Params{
	//		"dcl_currcd":      m.DclCurrcd,
	//		"dcl_currcd_name": m.DclCurrcdName,
	//		"natcd":           m.DestinationNatcd,
	//		"natcd_name":      m.DestinationNatcdName,
	//	}
	//}

	if params != nil {
		_, err := qs.Update(params)
		if err != nil {
			return err
		}
	} else {
		return errors.New("未更新")
	}

	return nil
}

//删除
func OrderItemLimitDelete(id int64) (num int64, err error) {
	m := NewOrderItemLimit(id)
	if num, err := BaseDelete(&m); err != nil {
		return num, err
	} else {
		return num, nil
	}
}
