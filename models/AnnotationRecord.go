package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"BeeCustom/utils"
	"github.com/astaxie/beego/orm"
)

// TableName 设置AnnotationRecord表名
func (u *AnnotationRecord) TableName() string {
	return AnnotationRecordTBName()
}

// AnnotationRecordQueryParam 用于查询的类
type AnnotationRecordQueryParam struct {
	BaseQueryParam

	AnnotationId int64
}

// AnnotationRecord 实体类
type AnnotationRecord struct {
	BaseModel

	Content       string       `orm:"column(content)" description:"办理记录内容"`
	Status        string       `orm:"column(status);size(255);null" description:"办理记录内容时状态"`
	Remark        string       `orm:"column(remark);size(255);null" description:"备注"`
	DeletedAt     time.Time    `orm:"column(deleted_at);type(timestamp);null"`
	BackendUser   *BackendUser `orm:"column(user_id);rel(fk)"`
	BackendUserId int64        `orm:"-" form:"BackendUserId"` //关联管理会自动生成字段，此处不生成字段
	Annotation    *Annotation  `orm:"column(annotation_id);rel(fk)"`
	AnnotationId  int64        `orm:"-" form:"AnnotationId"` //关联管理会自动生成字段，此处不生成字段
}

func NewAnnotationRecord(id int64) AnnotationRecord {
	return AnnotationRecord{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

//查询参数
func NewAnnotationRecordQueryParam() AnnotationRecordQueryParam {
	return AnnotationRecordQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

// AnnotationRecordPageList 获取分页数据
func AnnotationRecordPageList(params *AnnotationRecordQueryParam) ([]*AnnotationRecord, int64) {

	query := orm.NewOrm().QueryTable(AnnotationRecordTBName())
	datas := make([]*AnnotationRecord, 0)

	query = query.Filter("annotation_id", params.AnnotationId).RelatedSel()
	params.Sort = "Id"
	params.Order = "desc"

	total, _ := query.Count()
	query = BaseListQuery(query, params.Sort, params.Order, params.Limit, params.Offset)
	_, _ = query.All(&datas)

	return datas, total
}

func AnnotationRecordGetRelations(ms []*AnnotationRecord, relations string) ([]*AnnotationRecord, error) {
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

// AnnotationRecordOne 根据id获取单条
func AnnotationRecordOne(id int64) (*AnnotationRecord, error) {
	m := NewAnnotationRecord(0)
	o := orm.NewOrm()
	if err := o.QueryTable(AnnotationRecordTBName()).Filter("Id", id).RelatedSel().One(&m); err != nil {
		return nil, err
	}

	if &m == nil {
		return &m, errors.New("清单获取失败")
	}

	return &m, nil
}

//Save 添加、编辑页面 保存
func AnnotationRecordSave(m *AnnotationRecord) error {
	o := orm.NewOrm()
	if m.Id == 0 {
		if _, err := o.Insert(m); err != nil {
			utils.LogDebug(fmt.Sprintf("AnnotationRecordSave:%v", err))
			return err
		}
	} else {
		if _, err := o.Update(m); err != nil {
			utils.LogDebug(fmt.Sprintf("AnnotationRecordSave:%v", err))
			return err
		}
	}

	return nil
}

//删除
func AnnotationRecordDelete(id int64) (num int64, err error) {
	m := NewAnnotationRecord(id)
	if num, err := BaseDelete(&m); err != nil {
		return num, err
	} else {
		return num, nil
	}
}
