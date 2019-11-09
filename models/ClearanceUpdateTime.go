package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
)

// TableName 设置ClearanceUpdateTime表名
func (u *ClearanceUpdateTime) TableName() string {
	return ClearanceUpdateTimeTBName()
}

// Clearance 实体类
type ClearanceUpdateTime struct {
	BaseModel

	Type          int8      `orm:"column(type)" description:"参数类别"`
	LastUpdatedAt time.Time `orm:"column(last_updated_at);type(timestamp);null"`
}

// ClearanceQueryParam 用于查询的类
type ClearanceUpdateTimeQueryParam struct {
	BaseQueryParam
	Type     string //模糊查询
	NameLike string //模糊查询
}

func NewClearanceUpdateTime(id int64) ClearanceUpdateTime {
	return ClearanceUpdateTime{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

//查询参数
func NewClearanceUpdateTimeQueryParam() ClearanceUpdateTimeQueryParam {
	return ClearanceUpdateTimeQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

// ClearanceUpdateTimePageList 获取分页数据
func ClearanceUpdateTimePageList(params *ClearanceUpdateTimeQueryParam) ([]*ClearanceUpdateTime, int64) {
	query := orm.NewOrm().QueryTable(ClearanceUpdateTimeTBName())
	datas := make([]*ClearanceUpdateTime, 0)

	total, _ := query.Count()
	query = BaseListQuery(query, params.Sort, params.Order, params.Limit, params.Offset)
	_, _ = query.All(&datas)

	return datas, total
}

// ClearanceUpdateTimeOne 根据id获取单条
func ClearanceUpdateTimeOne(id int64) (*ClearanceUpdateTime, error) {
	m := NewClearanceUpdateTime(0)
	o := orm.NewOrm()
	if err := o.QueryTable(ClearanceUpdateTimeTBName()).Filter("Id", id).RelatedSel().One(&m); err != nil {
		return nil, err
	}

	if &m == nil {
		return &m, errors.New("获取失败")
	}

	return &m, nil
}

// ClearanceUpdateTimeOne 根据id获取单条
func GetLastUpdteTimeByClearanceType(cType int8) (*ClearanceUpdateTime, error) {
	m := NewClearanceUpdateTime(0)
	o := orm.NewOrm()
	if err := o.QueryTable(ClearanceUpdateTimeTBName()).Filter("Type", cType).RelatedSel().One(&m); err != nil {
		return nil, err
	}

	if &m == nil {
		return &m, errors.New("获取失败")
	}

	return &m, nil
}

//Save 添加、编辑页面 保存
func ClearanceUpdateTimeSave(m *ClearanceUpdateTime) (*ClearanceUpdateTime, error) {
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
