package models

import (
	"BeeCustom/utils"
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

// TableName 设置BackendUser表名
//func (u *AnnotationUserRel) TableName() string {
//	return AnnotationUserRelTBName()
//}

// AnnotationUserRelQueryParam 用于查询的类
type AnnotationUserRelQueryParam struct {
	BaseQueryParam

	UserType      string
	BackendUserId string
	AnnotationId  string
}

// AnnotationUserRel 实体类
type AnnotationUserRel struct {
	BaseModel

	Annotation  *Annotation  `orm:"rel(fk)"` //设置一对多关系
	BackendUser *BackendUser `orm:"rel(fk)"` //设置一对多关系
	UserType    int8         `orm:"column(type)"`
}

func NewAnnotationUserRel(id int64) AnnotationUserRel {
	return AnnotationUserRel{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

//查询参数
func NewAnnotationUserRelQueryParam() AnnotationUserRelQueryParam {
	return AnnotationUserRelQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

// BackendUserOne 根据id获取单条
func AnnotationUserRelByUserIdAndAnnotationId(userId, annotationId int64, userType int8) (*AnnotationUserRel, error) {
	m := NewAnnotationUserRel(0)
	o := orm.NewOrm()
	if err := o.QueryTable(AnnotationUserRelTBName()).
		Filter("backend_user_id", userId).
		Filter("annotation_id", annotationId).
		Filter("type", userType).
		One(&m); err != nil && err.Error() != "<QuerySeter> no row found" {
		utils.LogDebug(fmt.Sprintf("AnnotationUserRelByUserIdAndAnnotationId:%v", err))
		return nil, err
	}

	if &m == nil {
		return &m, errors.New("用户获取失败")
	}

	return &m, nil
}

//Save 添加、编辑页面 保存
func AnnotationUserRelSave(m *AnnotationUserRel) error {
	o := orm.NewOrm()

	if err := o.Read(m); err != nil {
		utils.LogDebug(fmt.Sprintf("AnnotationUserRelRead:%v", err))
		return err
	}

	if m.Id != 0 {
		if _, err := o.Update(m); err != nil {
			utils.LogDebug(fmt.Sprintf("AnnotationUserRelInsert:%v", err))
			return err
		}

	} else {
		if _, err := o.Insert(m); err != nil {
			utils.LogDebug(fmt.Sprintf("AnnotationUserRelInsert:%v", err))
			return err
		}

	}

	return nil
}
