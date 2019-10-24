package models

import (
	"fmt"
	"time"

	"BeeCustom/utils"

	"github.com/astaxie/beego/orm"
)

// ResourceQueryParam 用于搜索的类
type ResourceQueryParam struct {
	BaseQueryParam
	NameLike string
	IsParent bool
}

// TableName 设置表名
func (a *Resource) TableName() string {
	return ResourceTBName()
}

// Resource 权限控制资源表
type Resource struct {
	BaseModel

	Name         string `orm:"size(64)"`
	Rtype        int
	SonNum       int         `orm:"-"`
	Icon         string      `orm:"size(32)"`
	LinkUrl      string      `orm:"-"`
	UrlFor       string      `orm:"size(256)" Json:"-"`
	HtmlDisabled int         `orm:"-"`                 //在html里应用时是否可用
	Level        int         `orm:"-"`                 //第几级，从0开始
	ParentId     int64       `orm:"-" form:"ParentId"` //关联管理会自动生成 role_id 字段，此处不生成字段
	Parent       *Resource   `orm:"null;rel(fk)"`      // RelForeignKey relation
	Sons         []*Resource `orm:"reverse(many)"`     // fk 的反向关系
	Roles        []*Role     `orm:"reverse(many)"`     // 设置一对多的反向关系
}

func NewResource(id int64) Resource {
	return Resource{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

//查询参数
func NewResourceQueryParam() ResourceQueryParam {
	return ResourceQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

// ResourceOne 获取单条
func ResourceOne(id int64) (*Resource, error) {
	m := NewResource(id)
	o := orm.NewOrm()
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}

	if _, err := o.LoadRelated(&m, "Sons"); err != nil {
		utils.LogDebug(fmt.Sprintf("关联Sons权限出错：%s", err))
	}

	pr := NewResource(0)
	if m.Parent != nil && m.Parent.Id > 0 {
		pr.Id = m.Parent.Id
		if err := o.Read(&pr); err != nil {
			utils.LogDebug(fmt.Sprintf("关联Parent权限出错：%s", err))
		}
	}

	m.Parent = &pr
	return &m, nil
}

// ResourceTreeGrid 获取treegrid顺序的列表
func ResourceTreeGrid(params *ResourceQueryParam) ([]*Resource, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(ResourceTBName())
	datas := make([]*Resource, 0)

	if params.IsParent {
		query = query.Filter("parent_id__isnull", true)
	}

	if len(params.NameLike) > 0 {
		query = query.Filter("name__istartswith", params.NameLike)
	}

	total, _ := query.Count()
	query = BaseListQuery(query, params.Sort, params.Order, params.Limit, params.Offset)
	_, _ = query.All(&datas)

	//关联子权
	for _, v := range datas {
		if v.Parent == nil {
			if _, err := o.LoadRelated(v, "Sons"); err != nil {
				utils.LogDebug(fmt.Sprintf("关联子权限出错：%s", err))
			}
		}

	}

	return datas, total
}

// ResourceDataList 获取角色列表
func ResourceDataList(params *ResourceQueryParam) []*Resource {
	params.Limit = -1
	params.Sort = "Id"
	params.Order = "asc"
	data, _ := ResourceTreeGrid(params)

	return data
}

//Save 添加、编辑页面 保存
func ResourceSave(m *Resource) (*Resource, error) {
	o := orm.NewOrm()
	if m.ParentId != 0 {
		if pr, err := ResourceOne(m.ParentId); err != nil {
			return nil, err
		} else {
			m.Parent = pr
		}
	}

	if m.Id == 0 {
		if _, err := o.Insert(m); err != nil {
			return nil, err
		}
	} else {
		if _, err := o.Update(m, "Name", "Parent", "Rtype", "Sons", "Sons", "Icon", "UrlFor", "Roles", "UpdatedAt"); err != nil {
			return nil, err
		}
	}
	return m, nil
}

//删除
func ResourceDelete(id int64) (num int64, err error) {
	m := NewResource(id)
	if num, err := BaseDelete(&m); err != nil {
		return num, err
	} else {
		return num, nil
	}
}
