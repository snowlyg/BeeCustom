package models

import (
	"errors"
	"html"
	"strings"
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
	Key    string `orm:"size(128)"`
	Value  string `orm:"type(text)"`
	RValue string `orm:"type(text)"`
	Rmk    string `orm:"size(128)"`
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

	m.Value = html.UnescapeString(m.Value)

	return &m, nil
}

// GetSettingRValueByKey 获取单条
func GetSettingRValueByKey(key string, isSinge bool) (map[string]string, error) {
	o := orm.NewOrm()
	m := NewSetting(0)

	query := o.QueryTable(SettingTBName())
	err := query.Filter("Key", key).One(&m)
	if err != nil {
		return nil, err
	}

	rValue := map[string]string{}
	for _, v := range strings.Split(m.RValue, ",") {
		if isSinge {
			rValue["-1"] = v
			return rValue, nil
		}

		iv := strings.Split(v, ":")
		if len(iv) > 1 && len(iv[0]) > 0 && len(iv[1]) > 0 {
			rValue[iv[0]] = iv[1]
		} else {
			rValue["-1"] = v
		}

	}

	return rValue, nil
}

// GetSettingValueByKey 获取单条
func GetSettingValueByKey(key string) (string, error) {
	value, err := GetSettingRValueByKey(key, true)
	if err != nil {
		return "", err
	}

	if len(value) != 1 {
		return "", errors.New("获取数据格式错误")
	}

	return value["-1"], nil
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
	//html.EscapeString(m.Value)
	//html.UnescapeString(content)
	m.Value = html.EscapeString(m.Value)
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
