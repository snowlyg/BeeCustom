package models

import (
	"fmt"
	"time"

	"BeeCustom/utils"
	"github.com/astaxie/beego/orm"
)

// TableName 设置OrderFile表名
func (u *OrderFile) TableName() string {
	return OrderFileTBName()
}

// OrderFileQueryParam 用于查询的类
type OrderFileQueryParam struct {
	BaseQueryParam

	OrderId int64
}

// OrderFile 实体类
type OrderFile struct {
	BaseModel

	EdocID string `orm:"column(edoc_i_d);size(64)" description:"文件名、随附单据编号 命名规则是：申报口岸+随附单据类别代码+SW+18位流
水号"`
	EdocCode      string  `orm:"column(edoc_code);size(8);null" description:"随附单据文件类别"`
	EdocCodeName  string  `orm:"column(edoc_code_name);size(50);null" description:"随附单据文件类别名称"`
	EdocFomatType string  `orm:"column(edoc_fomat_type);size(2);null" description:"随附单据格式类型,S:结构化；US:非结构化"`
	OpNote        string  `orm:"column(op_note);size(255);null" description:"操作说明（重传原因）可选"`
	EdocCopId     string  `orm:"column(edoc_cop_id);size(64);null" description:"随附单据文件名"`
	EdocOwnerCode string  `orm:"column(edoc_owner_code);size(10);null" description:"所属单位海关编号 (申报单位)"`
	SignUnit      string  `orm:"column(sign_unit);size(10);null" description:"签名单位代码 (申报单位)"`
	SignTime      string  `orm:"column(sign_time);size(17);null" description:"签名时间 yyyyMMdd hh:ss:mm 申报时间"`
	EdocOwnerName string  `orm:"column(edoc_owner_name);size(100);null" description:"所属单位名称 (申报单位)"`
	EdocSize      string  `orm:"column(edoc_size);size(12);null" description:"随附单据附件文件大小 （KB） 可选"`
	EdocCopUrl    string  `orm:"column(edoc_cop_url);size(300);null" description:"随附单据文件地址"`
	Creator       string  `orm:"column(creator);size(50);null" description:"操作人"`
	Version       float64 `orm:"column(version)" description:"版本号"`

	Order   *Order `orm:"column(order_id);rel(fk)"`
	OrderId int64  `orm:"-" form:"OrderId"` // 关联管理会自动生成字段，此处不生成字段
}

func NewOrderFile(id int64) OrderFile {
	return OrderFile{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

// 查询参数
func NewOrderFileQueryParam() OrderFileQueryParam {
	return OrderFileQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

// OrderFilePageList 获取分页数据
func OrderFilePageList(params *OrderFileQueryParam) ([]*OrderFile, int64) {

	query := orm.NewOrm().QueryTable(OrderFileTBName())
	datas := make([]*OrderFile, 0)

	query = query.Filter("order_id", params.OrderId).RelatedSel()
	params.Sort = "Id"
	params.Order = "desc"

	total, _ := query.Count()
	query = BaseListQuery(query, params.Sort, params.Order, params.Limit, params.Offset)
	_, _ = query.All(&datas)

	return datas, total
}

// OrderFileOne 根据id获取单条
func OrderFileOneByTypeAndOrderId(m *OrderFile) error {
	m.Id = 0
	o := orm.NewOrm()
	if err := o.QueryTable(OrderFileTBName()).
		Filter("order_id", m.Order.Id).
		Filter("edoc_code", m.EdocCode).
		One(m); err != nil {
		return err
	}

	return nil
}

// Save 添加、编辑页面 保存
func OrderFileSaveOrUpdate(m *OrderFile) error {
	o := orm.NewOrm()
	if err := OrderFileOneByTypeAndOrderId(m); err != nil && err.Error() != "<QuerySeter> no row found" {
		utils.LogDebug(fmt.Sprintf("OrderFileSave:%v", err))
		return err
	}

	if m.Id == 0 {
		if _, err := o.Insert(m); err != nil {
			utils.LogDebug(fmt.Sprintf("OrderFileSave:%v", err))
			return err
		}
	} else {
		if _, err := o.Update(m, "EdocCopUrl"); err != nil {
			utils.LogDebug(fmt.Sprintf("OrderFileSave:%v", err))
			return err
		}
	}

	return nil
}
