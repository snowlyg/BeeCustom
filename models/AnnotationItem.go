package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"BeeCustom/utils"
	"github.com/astaxie/beego/orm"
)

// TableName 设置AnnotationItem表名
func (u *AnnotationItem) TableName() string {
	return AnnotationItemTBName()
}

// AnnotationItemQueryParam 用于查询的类
type AnnotationItemQueryParam struct {
	BaseQueryParam

	AnnotationId int64
}

// AnnotationItem 实体类
type AnnotationItem struct {
	BaseModel

	SeqNo                string      `orm:"column(seq_no);size(18);null" description:"中心统一编号 (首次导入时自动生成并返填，非首次导入须填写)"`
	GdsSeqno             int         `orm:"column(gds_seqno)"  valid:"Required" description:"商品序号"`
	PutrecSeqno          int         `orm:"column(putrec_seqno)"  description:"备案序号(对应底账序号）"`
	GdsMtno              string      `orm:"column(gds_mtno);size(32)" valid:"Required;MaxSize(32)"   description:"商品料号"`
	Gdecd                string      `orm:"column(gdecd);size(10)" valid:"Required;MaxSize(10)"   description:"商品编码"`
	GdsNm                string      `orm:"column(gds_nm);size(512)" valid:"Required;MaxSize(512)"   description:"商品名称"`
	GdsSpcfModelDesc     string      `orm:"column(gds_spcf_model_desc);size(512)" valid:"Required;MaxSize(512)"  description:"商品规格型号"`
	DclUnitcd            string      `orm:"column(dcl_unitcd);size(3);null" valid:"Required;MaxSize(3)" description:"申报计量单位"`
	DclUnitcdName        string      `orm:"column(dcl_unitcd_name);size(100);null" valid:"Required;MaxSize(100)" description:"申报计量单位名称"`
	LawfUnitcd           string      `orm:"column(lawf_unitcd);size(3);null" valid:"Required;MaxSize(3)" description:"法定计量单位"`
	LawfUnitcdName       string      `orm:"column(lawf_unitcd_name);size(100);null" valid:"Required;MaxSize(100)" description:"法定计量单位名称"`
	SecdLawfUnitcd       string      `orm:"column(secd_lawf_unitcd);size(3);null" description:"法定第二计量单位"`
	SecdLawfUnitcdName   string      `orm:"column(secd_lawf_unitcd_name);size(100);null" description:"法定第二计量单位名称"`
	Natcd                string      `orm:"column(natcd);size(3)" valid:"Required;MaxSize(3)" description:"原产国(地区)"`
	NatcdName            string      `orm:"column(natcd_name);size(100)" valid:"Required;MaxSize(100)" description:"原产国(地区)名称"`
	DclUprcAmt           float64     `orm:"column(dcl_uprc_amt);digits(19);decimals(4)" valid:"Required" description:"企业申报单价"`
	DclTotalAmt          float64     `orm:"column(dcl_total_amt);digits(19);decimals(2)" valid:"Required" description:"企业申报总价"`
	UsdStatTotalAmt      float64     `orm:"column(usd_stat_total_amt);null;digits(25);decimals(5)" description:"美元统计总金额"`
	DclCurrcd            string      `orm:"column(dcl_currcd);size(3)" valid:"Required;MaxSize(3)" description:"币制"`
	DclCurrcdName        string      `orm:"column(dcl_currcd_name);size(100)" valid:"Required;MaxSize(100)" description:"币制名称"`
	LawfQty              float64     `orm:"column(lawf_qty);digits(19);decimals(5)"  valid:"Required" description:"法定数量"`
	SecdLawfQty          float64     `orm:"column(secd_lawf_qty);null;digits(19);decimals(5)" description:"第二法定数量 (当法定第二计量单位为空时，该项为非必填)"`
	WtSfVal              float64     `orm:"column(wt_sf_val);null;digits(19);decimals(5)" description:"重量比例因子"`
	FstSfVal             float64     `orm:"column(fst_sf_val);null;digits(19);decimals(5)" description:"第一比例因子"`
	SecdSfVal            float64     `orm:"column(secd_sf_val);null;digits(19);decimals(5)" description:"第二比例因子"`
	DclQty               float64     `orm:"column(dcl_qty);digits(19);decimals(5)" valid:"Required" description:"申报数量"`
	GrossWt              float64     `orm:"column(gross_wt);null;digits(19);decimals(5)" description:"毛重"`
	NetWt                float64     `orm:"column(net_wt);null;digits(19);decimals(5)" description:"净重"`
	UseCd                string      `orm:"column(use_cd);size(4);null" description:"用途代码 (取消该字段使用；不需要填写)"`
	LvyrlfModecd         string      `orm:"column(lvyrlf_modecd);size(6)" valid:"Required;MaxSize(6)" description:"征免方式"`
	LvyrlfModecdName     string      `orm:"column(lvyrlf_modecd_name);size(100)" valid:"Required;MaxSize(100)" description:"征免方式名称"`
	UcnsVerno            string      `orm:"column(ucns_verno);size(8);null" description:"单耗版本号(账册由开关控制是否必填。需看单耗该字段如何定义)"`
	EntryGdsSeqno        int         `orm:"column(entry_gds_seqno);null" description:"报关单商品序号(企业可录入，如果企业不录入，系统自动返填)(填写后自动归并报关单商品)"`
	ClyMarkcd            string      `orm:"column(cly_markcd);size(4);null" description:"归类标志"`
	ApplyTbSeqno         int         `orm:"column(apply_tb_seqno);null" description:"流转申报表序号(流转类专用。用于建立清单商品与流转申请表商品之间的关系)"`
	DestinationNatcd     string      `orm:"column(destination_natcd);size(3)" valid:"Required;MaxSize(3)" description:"最终目的国"`
	DestinationNatcdName string      `orm:"column(destination_natcd_name);size(100)" valid:"Required;MaxSize(100)" description:"最终目的国名称"`
	ModfMarkcd           string      `orm:"column(modf_markcd);size(1)" valid:"Required;MaxSize(1)" description:"修改标志"`
	Rmk                  string      `orm:"column(rmk);null" description:"备注"`
	Param3               string      `orm:"column(param3);size(19);null" description:"（单一显示：自动备案序号）备用3"`
	EntrySeqNo           string      `orm:"column(entry_seq_no);size(20);null" description:"报关单预录入号"`
	ModfMarkcdName       string      `orm:"column(modf_markcd_name);size(100);null" valid:"Required;MaxSize(100)" description:"修改标志 0-未修改 1-修改 2-删除 3-增加"`
	DeletedAt            time.Time   `orm:"column(deleted_at);type(timestamp);null"`
	Annotation           *Annotation `orm:"column(annotation_id);rel(fk)"`
	AnnotationId         int64       `orm:"-" form:"AnnotationId"` //关联管理会自动生成 CompanyId 字段，此处不生成字段
}

func NewAnnotationItem(id int64) AnnotationItem {
	return AnnotationItem{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

//查询参数
func NewAnnotationItemQueryParam() AnnotationItemQueryParam {
	return AnnotationItemQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

// AnnotationItemPageList 获取分页数据
func AnnotationItemPageList(params *AnnotationItemQueryParam) ([]*AnnotationItem, int64) {

	query := orm.NewOrm().QueryTable(AnnotationItemTBName())
	datas := make([]*AnnotationItem, 0)

	query = query.Filter("annotation_id", params.AnnotationId)

	total, _ := query.Count()
	query = BaseListQuery(query, params.Sort, params.Order, params.Limit, params.Offset)
	_, _ = query.All(&datas)

	return datas, total
}

// AnnotationItemPageList 获取分页数据
func AnnotationItemsByAnnotationId(aId int64) ([]*AnnotationItem, error) {

	datas := make([]*AnnotationItem, 0)
	_, err := orm.NewOrm().QueryTable(AnnotationItemTBName()).Filter("annotation_id", aId).All(&datas)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("AnnotationItemsByAnnotationId error :%v", err))
		return nil, err
	}

	return datas, nil
}

func AnnotationItemGetRelations(ms []*AnnotationItem, relations string) ([]*AnnotationItem, error) {
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

// AnnotationItemOne 根据id获取单条
func AnnotationItemOne(id int64) (*AnnotationItem, error) {
	m := NewAnnotationItem(0)
	o := orm.NewOrm()
	if err := o.QueryTable(AnnotationItemTBName()).Filter("Id", id).RelatedSel().One(&m); err != nil {
		return nil, err
	}

	if &m == nil {
		return &m, errors.New("清单获取失败")
	}

	return &m, nil
}

//Save 添加、编辑页面 保存
func AnnotationItemSave(m *AnnotationItem) error {
	o := orm.NewOrm()

	//进出口原产国和目的国是相反的数据
	if m.Annotation.ImpexpMarkcd == "E" {
		natcd := m.Natcd
		natcdName := m.NatcdName
		destinationNatcd := m.DestinationNatcd
		destinationNatcdName := m.DestinationNatcdName

		m.Natcd = destinationNatcd
		m.NatcdName = destinationNatcdName
		m.DestinationNatcd = natcd
		m.DestinationNatcdName = natcdName
	}

	if m.Id == 0 {
		if _, err := o.Insert(m); err != nil {
			utils.LogDebug(fmt.Sprintf("AnnotationItemSave:%v", err))
			return err
		}
	} else {

		if _, err := o.Update(m); err != nil {
			utils.LogDebug(fmt.Sprintf("AnnotationItemSave:%v", err))
			return err
		}
	}

	return nil
}

//AnnotationItemUpdateAll 添加、编辑页面 保存
func AnnotationItemUpdateAll(aid int64, m *AnnotationItem) error {
	o := orm.NewOrm()
	qs := o.QueryTable(AnnotationItemTBName()).Filter("annotation_id", aid)

	var params orm.Params
	if len(m.Natcd) > 0 {
		params = orm.Params{
			"dcl_currcd":      m.DclCurrcd,
			"dcl_currcd_name": m.DclCurrcdName,
			"natcd":           m.Natcd,
			"natcd_name":      m.NatcdName,
		}
	} else if len(m.DestinationNatcd) > 0 {
		params = orm.Params{
			"dcl_currcd":      m.DclCurrcd,
			"dcl_currcd_name": m.DclCurrcdName,
			"natcd":           m.DestinationNatcd,
			"natcd_name":      m.DestinationNatcdName,
		}
	}

	if params != nil {
		_, err := qs.Update(params)
		if err != nil {
			return err
		}
	} else {
		return errors.New("未更新")
	}

	return nil
}

//删除
func AnnotationItemDelete(id int64) (num int64, err error) {
	m := NewAnnotationItem(id)
	if num, err := BaseDelete(&m); err != nil {
		return num, err
	} else {
		return num, nil
	}
}
