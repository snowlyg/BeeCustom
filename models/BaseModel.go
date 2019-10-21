package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type BaseModel struct {
	Id        int64
	CreatedAt time.Time `orm:"column(created_at);type(timestamp);null"`
	UpdatedAt time.Time `orm:"column(updated_at);type(timestamp);null"`
}

func BaseListQuery(query orm.QuerySeter, sort, order string, limit, offset int64) orm.QuerySeter {

	//默认排序
	sortorder := "Id"
	if len(sort) > 0 {
		sortorder = sort
	}

	if order == "desc" {
		sortorder = "-" + sortorder
	}

	query.OrderBy(sortorder)

	if limit == -1 {
		query = query.Limit(limit, (offset-1)*limit).RelatedSel()
	}

	return query
}

//删除
func BaseDelete(m interface{}) (num int64, err error) {

	o := orm.NewOrm()
	if num, err := o.Delete(m); err != nil {
		return num, err
	} else {
		return num, nil
	}

}
