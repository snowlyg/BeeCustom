package models

import (
	"fmt"
	"time"

	"BeeCustom/utils"
	"github.com/astaxie/beego/orm"
)

// TableName 设置AnnotationFile表名
func (u *AnnotationFile) TableName() string {
	return AnnotationFileTBName()
}

// AnnotationFileQueryParam 用于查询的类
type AnnotationFileQueryParam struct {
	BaseQueryParam

	AnnotationId int64
}

// AnnotationFile 实体类
type AnnotationFile struct {
	BaseModel

	Type    string `orm:"column(type)" description:"附件类型"`
	Name    string `orm:"column(name);size(255)" description:"附件名称"`
	Url     string `orm:"column(url);size(255)" description:"附件url"`
	Creator string `orm:"column(creator);size(255)" description:"创建人"`
	Version string `orm:"column(version);size(255)" description:"版本"`

	Annotation   *Annotation `orm:"column(annotation_id);rel(fk)"`
	AnnotationId int64       `orm:"-" form:"AnnotationId"` // 关联管理会自动生成字段，此处不生成字段
}

func NewAnnotationFile(id int64) AnnotationFile {
	return AnnotationFile{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

// 查询参数
func NewAnnotationFileQueryParam() AnnotationFileQueryParam {
	return AnnotationFileQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

// AnnotationFilePageList 获取分页数据
func AnnotationFilePageList(params *AnnotationFileQueryParam) ([]*AnnotationFile, int64) {

	query := orm.NewOrm().QueryTable(AnnotationFileTBName())
	datas := make([]*AnnotationFile, 0)

	query = query.Filter("annotation_id", params.AnnotationId).RelatedSel()
	params.Sort = "Id"
	params.Order = "desc"

	total, _ := query.Count()
	query = BaseListQuery(query, params.Sort, params.Order, params.Limit, params.Offset)
	_, _ = query.All(&datas)

	return datas, total
}

// AnnotationFileOne 根据id获取单条
func AnnotationFileOneByTypeAndAnnotationId(m *AnnotationFile) error {
	m.Id = 0
	o := orm.NewOrm()
	if err := o.QueryTable(AnnotationFileTBName()).
		Filter("annotation_id", m.Annotation.Id).
		Filter("type", m.Type).
		One(m); err != nil {
		return err
	}

	return nil
}

// Save 添加、编辑页面 保存
func AnnotationFileSaveOrUpdate(m *AnnotationFile) error {
	o := orm.NewOrm()
	if err := AnnotationFileOneByTypeAndAnnotationId(m); err != nil && err.Error() != "<QuerySeter> no row found" {
		utils.LogDebug(fmt.Sprintf("AnnotationFileSave:%v", err))
		return err
	}

	if m.Id == 0 {
		if _, err := o.Insert(m); err != nil {
			utils.LogDebug(fmt.Sprintf("AnnotationFileSave:%v", err))
			return err
		}
	} else {
		if _, err := o.Update(m, "url"); err != nil {
			utils.LogDebug(fmt.Sprintf("AnnotationFileSave:%v", err))
			return err
		}
	}

	return nil
}
