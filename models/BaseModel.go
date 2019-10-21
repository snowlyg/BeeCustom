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

//删除
func BaseDelete(m interface{}) (num int64, err error) {

	o := orm.NewOrm()
	if num, err := o.Delete(m); err != nil {
		return num, err
	} else {
		return num, nil
	}

}
