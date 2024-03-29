package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"BeeCustom/utils"

	"github.com/astaxie/beego/orm"
)

// TableName 设置HandBookGood表名
func (u *HandBookGood) TableName() string {
	return HandBookGoodTBName()
}

// HandBookGood 实体类
type HandBookGood struct {
	BaseModel

	Type                int8    `orm:"column(type)" description:"货物类型：1-成品，2-料件"`
	Serial              string  `orm:"column(serial);size(255)" description:"项号"`
	RecordNo            string  `orm:"column(record_no);size(20);null" description:"货号"`
	HsCode              string  `orm:"column(hs_code);size(10)" description:"商品编码"`
	Name                string  `orm:"column(name);size(50);null" description:"商品名称"`
	ClassificationMark  string  `orm:"column(classification_mark);size(10);null" description:"分类标志"`
	Special             string  `orm:"column(special);size(255)" description:"规格型号"`
	UnitOne             string  `orm:"column(unit_one);size(100)" description:"第一单位名称"`
	UnitTwo             string  `orm:"column(unit_two);size(100);null" description:"第二单位名称"`
	UnitThree           string  `orm:"column(unit_three);size(100);null" description:"第三单位名称"`
	Price               float64 `orm:"column(price);null;digits(17);decimals(4)" description:"单价"`
	Moneyunit           string  `orm:"column(moneyunit);size(200);null" description:"币制名称"`
	Quantity            float64 `orm:"column(quantity);null;digits(19);decimals(5)" description:"申报数量"`
	MaxAllowance        float64 `orm:"column(max_allowance);null;digits(19);decimals(5)" description:"最大余量"`
	InitialQuantity     float64 `orm:"column(initial_quantity);null;digits(19);decimals(5)" description:"初始数量"`
	UnitTwoProportion   float64 `orm:"column(unit_two_proportion);null;digits(19);decimals(5)" description:"第二单位比例"`
	UnitThreeProportion float64 `orm:"column(unit_three_proportion);null;digits(19);decimals(5)" description:"第三单位比例"`
	WeightProportion    float64 `orm:"column(weight_proportion);null;digits(19);decimals(5)" description:"重量比例因子"`
	Taxationlx          string  `orm:"column(taxationlx);size(200);null" description:"征免类型"`
	DeclareMode         string  `orm:"column(declare_mode);size(255);null" description:"申报类别"`
	Remark              string  `orm:"column(remark);size(255);null" description:"备注"`
	HandleMark          string  `orm:"column(handle_mark);size(10);null" description:"处理标志(修改标志)"`
	CompanyActionFlag   string  `orm:"column(company_action_flag);size(255);null" description:"企业执行标志"`
	CustomActionFlag    string  `orm:"column(custom_action_flag);size(255);null" description:"海关执行标志"`
	StartCount          string  `orm:"column(start_count);size(255);null" description:"期初数量"`
	CountControlFlag    string  `orm:"column(count_control_flag);size(255);null" description:"数量控制标志"`
	BigCount            string  `orm:"column(big_count);size(255);null" description:"批准最大余数量"`
	UllageFlag          string  `orm:"column(ullage_flag);size(255);null" description:"单耗质疑标志"`
	ConsultMark         string  `orm:"column(consult_mark);size(255);null" description:"磋商标志"`
	MainMark            string  `orm:"column(main_mark);size(255);null" description:"主料标志"`
	Amount              float64 `orm:"column(amount);digits(17);decimals(4)" description:"总价（申报总价）"`
	Manuplace           string  `orm:"column(manuplace);size(50)" description:"产销国(地区)"`
	GoodAttr            string  `orm:"column(good_attr);size(255);null" description:"商品属性"`
	SeqNo               string  `orm:"column(seq_no);size(18);null" description:"预录入统一编号 (返填)"`

	HandBook   *HandBook `orm:"column(hand_book_id);rel(fk)"`
	HandBookId int64     `orm:"-" form:"HandBookId"`
}

// HandBookGoodQueryParam 用于查询的类
type HandBookGoodQueryParam struct {
	BaseQueryParam

	Type       int8
	HandBookId int64
	Serial     string
}

func NewHandBookGood(id int64) HandBookGood {
	return HandBookGood{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

// 查询参数
func NewHandBookGoodQueryParam() HandBookGoodQueryParam {
	return HandBookGoodQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

// HandBookGoodPageList 获取分页数据
func HandBookGoodPageList(params *HandBookGoodQueryParam) ([]*HandBookGood, int64) {
	query := orm.NewOrm().QueryTable(HandBookGoodTBName())
	data := make([]*HandBookGood, 0)

	query = query.Distinct().Filter("hand_book_id", params.HandBookId).Filter("Type", params.Type)
	if len(params.Serial) > 0 {
		query = query.Filter("Serial", params.Serial)
	}

	total, _ := query.Count()
	query = BaseListQuery(query, params.Sort, params.Order, params.Limit, params.Offset)

	_, _ = query.All(&data)

	return data, total
}

func HandBookGoodGetRelations(v *HandBookGood, relations string) (*HandBookGood, error) {
	o := orm.NewOrm()
	rs := strings.Split(relations, ",")
	for _, rv := range rs {
		_, err := o.LoadRelated(v, rv)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("LoadRelated:%v", err))
			return nil, err
		}

	}

	return v, nil
}

// GetHandBookGoodById 根据id获取单条
func GetHandBookGoodById(params *HandBookGoodQueryParam) (*HandBookGood, error) {
	m := NewHandBookGood(0)

	query := orm.NewOrm().QueryTable(HandBookGoodTBName())
	query = query.Distinct().Filter("hand_book_id", params.HandBookId).Filter("Type", params.Type)
	if len(params.Serial) > 0 {
		query = query.Filter("Serial", params.Serial)
	}

	if err := query.One(&m); err != nil {
		return nil, err
	}

	if &m == nil {
		return nil, errors.New("获取失败")
	}

	return &m, nil
}

// GetHandBookGoodBySerial 根据Serial获取单条 成品
func GetHandBookGoodBySerial(serial string) (*HandBookGood, error) {
	m := NewHandBookGood(0)
	o := orm.NewOrm()
	if err := o.QueryTable(HandBookGoodTBName()).Filter("Serial", serial).Filter("Type", 1).One(&m); err != nil && err.Error() != "<QuerySeter> no row found" {
		return nil, err
	}

	if &m == nil {
		return &m, errors.New("获取失败")
	}

	return &m, nil
}

// 批量插入
func InsertHandBookGoodMulti(datas []*HandBookGood) (num int64, err error) {
	return BaseInsertMulti(datas)
}
