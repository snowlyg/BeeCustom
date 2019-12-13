package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"BeeCustom/utils"
	"github.com/astaxie/beego/orm"
)

// TableName 设置OrderItem表名
func (u *OrderItem) TableName() string {
	return OrderItemTBName()
}

// OrderItemQueryParam 用于查询的类
type OrderItemQueryParam struct {
	BaseQueryParam

	OrderId int64
}

// OrderItem 实体类
type OrderItem struct {
	BaseModel

	GNo                    int     `orm:"column(g_no)" description:"项号(序号)"`
	ContrItem              string  `orm:"column(contr_item);size(19);null" description:"备案号"`
	CodeTS                 string  `orm:"column(code_t_s);size(10)" description:"商品编码"`
	GName                  string  `orm:"column(g_name);size(255)" description:"商品名称"`
	GModel                 string  `orm:"column(g_model);size(255);null" description:"商品规格、型号"`
	GQty                   float64 `orm:"column(g_qty);null;digits(14);decimals(5)" description:"成交数量"`
	GUnit                  string  `orm:"column(g_unit);size(3);null" description:"成交计量单位"`
	GUnitName              string  `orm:"column(g_unit_name);size(50);null" description:"成交计量单位名称"`
	DeclPrice              float64 `orm:"column(decl_price);null;digits(14);decimals(4)" description:"单价"`
	DeclTotal              float64 `orm:"column(decl_total);null;digits(15);decimals(2)" description:"总价"`
	TradeCurr              string  `orm:"column(trade_curr);size(3);null" description:"币制"`
	TradeCurrName          string  `orm:"column(trade_curr_name);size(50);null" description:"币制名称"`
	FirstQty               float64 `orm:"column(first_qty);null;digits(19);decimals(5)" description:"法定第一数量"`
	FirstUnit              string  `orm:"column(first_unit);size(3);null" description:"法定第一单位"`
	FirstUnitName          string  `orm:"column(first_unit_name);size(3);null" description:"法定第一单位名称"`
	SecondQty              float64 `orm:"column(second_qty);null;digits(19);decimals(5)" description:"法定第二数量"`
	SecondUnit             string  `orm:"column(second_unit);size(3);null" description:"法定第二单位"`
	SecondUnitName         string  `orm:"column(second_unit_name);size(3);null" description:"法定第二单位名称"`
	ExgVersion             string  `orm:"column(exg_version);size(8);null" description:"加工成品单耗版本号"`
	ExgNo                  string  `orm:"column(exg_no);size(30);null" description:"货号"`
	DestinationCountry     string  `orm:"column(destination_country);size(3);null" description:"最终目的国"`
	DestinationCountryName string  `orm:"column(destination_country_name);size(50);null" description:"最终目的国名称"`
	OriginCountry          string  `orm:"column(origin_country);size(3);null" description:"原产国（地区）"`
	OriginCountryName      string  `orm:"column(origin_country_name);size(50);null" description:"原产国（地区）名称"`
	OrigPlaceCode          string  `orm:"column(orig_place_code);size(6);null" description:"原产地区代码"`
	OrigPlaceCodeName      string  `orm:"column(orig_place_code_name);size(100);null" description:"原产地区代码名称"`
	DistrictCode           string  `orm:"column(district_code);size(5);null" description:"境内目的地/境内货源地"`
	DistrictCodeName       string  `orm:"column(district_code_name);size(200);null" description:"境内目的地/境内货源地名称"`
	DestCode               string  `orm:"column(dest_code);size(6);null" description:"目的地代码"`
	DestCodeName           string  `orm:"column(dest_code_name);size(6);null" description:"目的地代码名称"`
	DutyMode               string  `orm:"column(duty_mode);size(1);null" description:"征免类型 (方式)"`
	DutyModeName           string  `orm:"column(duty_mode_name);size(50);null" description:"征免类型名称"`
	ClassMark              string  `orm:"column(class_mark);size(1);null" description:"归类标志"`
	CiqCode                string  `orm:"column(ciq_code);size(20);null" description:"检验检疫3位编码"`
	CiqName                string  `orm:"column(ciq_name);size(400);null" description:"检验检疫名称"`
	DeclGoodsEname         string  `orm:"column(decl_goods_ename);size(100);null" description:"商品英文名称"`
	Stuff                  string  `orm:"column(stuff);size(400);null" description:"成份/原料/组份"`
	ProdValidDt            string  `orm:"column(prod_valid_dt);size(20);null" description:"产品有效期"`
	ProdQgp                string  `orm:"column(prod_qgp);size(20);null" description:"产品保质期(天)"`
	EngManEntCnm           string  `orm:"column(eng_man_ent_cnm);size(100);null" description:"境外生产企业名称"`
	GoodsSpec              string  `orm:"column(goods_spec);null" description:"检验检疫货物规格"`
	GoodsModel             string  `orm:"column(goods_model);null" description:"检验检疫货物型号"`
	GoodsBrand             string  `orm:"column(goods_brand);null" description:"货物品牌"`
	ProduceDate            string  `orm:"column(produce_date);null" description:"生产日期"`
	ProdBatchNo            string  `orm:"column(prod_batch_no);null" description:"生产批号:货物的生产批号"`
	GoodsAttr              string  `orm:"column(goods_attr);type(text);null" description:"货物属性代码 array"`
	GoodsAttrName          string  `orm:"column(goods_attr_name);type(text);null" description:"货物属性name"`
	Purpose                string  `orm:"column(purpose);size(2);null" description:"用途"`
	PurposeName            string  `orm:"column(purpose_name);size(100);null" description:"用途名称"`
	NoDangFlag             string  `orm:"column(no_dang_flag);size(1);null" description:"非危险化学品"`
	NoDangFlagName         string  `orm:"column(no_dang_flag_name);size(1);null" description:"非危险化学品"`
	DangName               string  `orm:"column(dang_name);size(80);null" description:"危险货物名称"`
	UnCode                 string  `orm:"column(un_code);size(20);null" description:"UN编码"`
	DangPackType           string  `orm:"column(dang_pack_type);size(4);null" description:"危包类别"`
	DangPackTypeName       string  `orm:"column(dang_pack_type_name);size(200);null" description:"危包类别名称"`
	DangPackSpec           string  `orm:"column(dang_pack_spec);size(24);null" description:"危包规格"`
	DangPackSpecName       string  `orm:"column(dang_pack_spec_name);size(100);null" description:"危包规格"`

	DeletedAt time.Time `orm:"column(deleted_at);type(timestamp);null"`

	OrderItemLimits []*OrderItemLimit `orm:"reverse(many)"` // 设置一对多关系

	Order   *Order `orm:"column(order_id);rel(fk)"`
	OrderId int64  `orm:"-" form:"OrderId"` //关联管理会自动生成 CompanyId 字段，此处不生成字段
}

func NewOrderItem(id int64) OrderItem {
	return OrderItem{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

func OrderItemGetRelations(ms []*OrderItem, relations string) ([]*OrderItem, error) {
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

// OrderItemOne 根据id获取单条
func OrderItemOne(id int64) (*OrderItem, error) {
	m := NewOrderItem(0)
	o := orm.NewOrm()
	if err := o.QueryTable(OrderItemTBName()).Filter("Id", id).RelatedSel().One(&m); err != nil {
		return nil, err
	}

	if &m == nil {
		return &m, errors.New("数据获取失败")
	}

	return &m, nil
}

//Save 添加、编辑页面 保存
func OrderItemSave(m *OrderItem, fields []string) error {
	o := orm.NewOrm()
	if m.Id == 0 {
		if _, err := o.Insert(m); err != nil {
			utils.LogDebug(fmt.Sprintf("OrderItemSave:%v", err))
			return err
		}
	} else {

		if len(fields) > 0 {
			if _, err := o.Update(m, fields...); err != nil {
				utils.LogDebug(fmt.Sprintf("OrderItemSave:%v", err))
				return err
			}
		} else {
			if _, err := o.Update(m); err != nil {
				utils.LogDebug(fmt.Sprintf("OrderItemSave:%v", err))
				return err
			}
		}
	}

	return nil
}

// 删除
func OrderItemDelete(id int64) (num int64, err error) {
	m := NewOrderItem(id)
	if num, err := BaseDelete(&m); err != nil {
		return num, err
	} else {
		return num, nil
	}
}
