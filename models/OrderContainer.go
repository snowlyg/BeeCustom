package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"BeeCustom/utils"
	"github.com/astaxie/beego/orm"
)

// TableName 设置OrderContainer表名
func (u *OrderContainer) TableName() string {
	return OrderContainerTBName()
}

// OrderContainerFieldNames 设置OrderItemLimitVin填充名称
func OrderContainerFieldNames() []string {
	return []string{
		"ContainerId",
		"ContainerMd",
		"ContainerMdName",
		"ContainerWt",
		"LclFlag",
		"LclFlagName",
		"GoodsNo",
		"GoodsContaWt",
	}
}

// OrderContainerQueryParam 用于查询的类
type OrderContainerQueryParam struct {
	BaseQueryParam

	OrderId int64
}

// OrderContainer 实体类
type OrderContainer struct {
	BaseModel
	ContainerId     string    `orm:"column(container_id);size(11)" description:"集装箱号"`
	ContainerMd     string    `orm:"column(container_md);size(2);null" description:"集装箱规格"`
	ContainerMdName string    `orm:"column(container_md_name);size(200);null" description:"集装箱规格名称"`
	ContainerWt     float64   `orm:"column(container_wt);null;digits(13);decimals(5)" description:"自重（KG）"`
	LclFlag         int8      `orm:"column(lcl_flag);null" description:"拼箱规格（拼箱标识）拼箱标识可选项：0:否；1:是"`
	LclFlagName     string    `orm:"column(lcl_flag_name)size(100);null" description:"拼箱规格名称"`
	GoodsNo         string    `orm:"column(goods_no);size(255);null" description:"商品项号关系:商品项号用半角逗号分隔，如“1,3”，该节点长度为255"`
	GoodsContaWt    float64   `orm:"column(goods_conta_wt);null;digits(17);decimals(5)" description:"箱货重量:集装箱箱体自重（千克）+ 装载货物重量（千克）,不计算该值，报文也不发送"`
	DeletedAt       time.Time `orm:"column(deleted_at);type(timestamp);null"`
	Order           *Order    `orm:"column(order_id);rel(fk)"`
	OrderId         int64     `orm:"-" form:"OrderId"` // 关联管理会自动生成 CompanyId 字段，此处不生成字段
}

func NewOrderContainer(id int64) OrderContainer {
	return OrderContainer{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

//查询参数
func NewOrderContainerQueryParam() OrderContainerQueryParam {
	return OrderContainerQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

// OrderContainerPageList 获取分页数据
func OrderContainerPageList(params *OrderContainerQueryParam) ([]*OrderContainer, int64) {

	query := orm.NewOrm().QueryTable(OrderContainerTBName())
	datas := make([]*OrderContainer, 0)

	query = query.Filter("order_id", params.OrderId)

	total, _ := query.Count()
	query = BaseListQuery(query, params.Sort, params.Order, params.Limit, params.Offset)
	_, _ = query.All(&datas)

	return datas, total
}

// OrderContainerPageList 获取分页数据
func OrderContainersByOrderId(aId int64) ([]*OrderContainer, error) {

	datas := make([]*OrderContainer, 0)
	_, err := orm.NewOrm().QueryTable(OrderContainerTBName()).Filter("order_id", aId).All(&datas)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("OrderContainersByOrderId error :%v", err))
		return nil, err
	}

	return datas, nil
}

func OrderContainerGetRelations(ms []*OrderContainer, relations string) ([]*OrderContainer, error) {
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

// OrderContainerOne 根据id获取单条
func OrderContainerOne(id int64) (*OrderContainer, error) {
	m := NewOrderContainer(0)
	o := orm.NewOrm()
	if err := o.QueryTable(OrderContainerTBName()).Filter("Id", id).RelatedSel().One(&m); err != nil {
		return nil, err
	}

	if &m == nil {
		return &m, errors.New("数据获取失败")
	}

	return &m, nil
}

// Save 添加、编辑页面 保存
func OrderContainerSave(m *OrderContainer, fields []string) error {
	o := orm.NewOrm()
	if m.Id == 0 {
		if _, err := o.Insert(m); err != nil {
			utils.LogDebug(fmt.Sprintf("OrderContainerSave:%v", err))
			return err
		}
	} else {
		if len(fields) > 0 {
			if _, err := o.Update(m, fields...); err != nil {
				utils.LogDebug(fmt.Sprintf("OrderContainerSave:%v", err))
				return err
			}
		} else {
			if _, err := o.Update(m); err != nil {
				utils.LogDebug(fmt.Sprintf("OrderContainerSave:%v", err))
				return err
			}
		}
	}

	return nil
}

// 删除
func OrderContainerDelete(id int64) (num int64, err error) {
	m := NewOrderContainer(id)
	if num, err := BaseDelete(&m); err != nil {
		return num, err
	} else {
		return num, nil
	}
}
