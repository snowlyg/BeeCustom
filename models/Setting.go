package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// SettingQueryParam 用于搜索的类
type SettingQueryParam struct {
	BaseQueryParam
}

// TableName 设置表名
func (a *Setting) TableName() string {
	return SettingTBName()
}

// Setting
type Setting struct {
	BaseModel
	Key   string `orm:"size(128)"`
	Value string `orm:"type(text)"`
	Rmk   string `orm:"size(128)"`
}

func NewSetting(id int64) Setting {
	return Setting{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

// 查询参数
func NewSettingQueryParam() SettingQueryParam {
	return SettingQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

// SettingOne 获取单条
func SettingOne(id int64) (*Setting, error) {
	m := NewSetting(id)
	o := orm.NewOrm()
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

// GetSettingByKey 获取单条
func GetSettingByKey(key string) (*Setting, error) {
	o := orm.NewOrm()
	m := NewSetting(0)

	query := o.QueryTable(SettingTBName())
	err := query.Filter("Key", key).One(&m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

// SettingTreeGrid 获取treegrid顺序的列表
func SettingTreeGrid(params *SettingQueryParam) ([]*Setting, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(SettingTBName())
	datas := make([]*Setting, 0)

	total, _ := query.Count()
	query = BaseListQuery(query, params.Sort, params.Order, params.Limit, params.Offset).RelatedSel()
	_, _ = query.All(&datas)

	return datas, total
}

// Save 添加、编辑页面 保存
func SettingSave(m *Setting) (*Setting, error) {
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

// 删除
func SettingDelete(id int64) (num int64, err error) {
	m := NewSetting(id)
	if num, err := BaseDelete(&m); err != nil {
		return num, err
	} else {
		return num, nil
	}
}
