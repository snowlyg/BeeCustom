package models

import (
	"fmt"
	"time"

	"BeeCustom/utils"
	"github.com/astaxie/beego/orm"
)

// TableName 设置OrderItemLimitVin表名
func (u *OrderItemLimitVin) TableName() string {
	return OrderItemLimitVinTBName()
}

// OrderItemLimitVinQueryParam 用于查询的类
type OrderItemLimitVinQueryParam struct {
	BaseQueryParam

	OrderItemLimitId int64
}

// OrderItemLimitVin 实体类
type OrderItemLimitVin struct {
	BaseModel

	VinNo        string `orm:"column(vin_no);size(100);null" description:"VIN序号"`
	BillLadDate  string `orm:"column(bill_lad_date);size(19);null" description:"提/运单日期"`
	QualityQgp   string `orm:"column(quality_qgp);size(100);null" description:"质量保质期"`
	VinCode      string `orm:"column(vin_code);size(20);null" description:"车辆识别代码(VIN)"`
	MotorNo      string `orm:"column(motor_no);size(100);null" description:"发动机号或电机号"`
	InvoiceNo    string `orm:"column(invoice_no);size(30);null" description:"发票号"`
	InvoiceNum   string `orm:"column(invoice_num);size(14);null" description:"发票所列数量"`
	ProdCnnm     string `orm:"column(prod_cnnm);size(500);null" description:"品名（中文名称）"`
	ProdEnnm     string `orm:"column(prod_ennm);size(500);null" description:"品名（英文名称）"`
	ModelEn      string `orm:"column(model_en);size(500);null" description:"型号(英文)"`
	ChassisNo    string `orm:"column(chassis_no);size(20);null" description:"底盘(车架)号"`
	PricePerUnit string `orm:"column(price_per_unit);size(20);null" description:"单价"`

	OrderItemLimit   *OrderItemLimit `orm:"column(order_item_limit_id);rel(fk)"`
	OrderItemLimitId int64           `orm:"-" form:"OrderItemLimitId"` //关联管理会自动生成
}

func NewOrderItemLimitVin(id int64) OrderItemLimitVin {
	return OrderItemLimitVin{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

//Save 添加、编辑页面 保存
func OrderItemLimitVinSave(m *OrderItemLimitVin, files []string) error {
	o := orm.NewOrm()

	if m.Id == 0 {
		if _, err := o.Insert(m); err != nil {
			utils.LogDebug(fmt.Sprintf("OrderItemLimitVinSave:%v", err))
			return err
		}
	} else {
		if len(files) > 0 {
			if _, err := o.Update(m, files...); err != nil {
				utils.LogDebug(fmt.Sprintf("OrderItemLimitVinSave:%v", err))
				return err
			}
		} else {
			if _, err := o.Update(m); err != nil {
				utils.LogDebug(fmt.Sprintf("OrderItemLimitVinSave:%v", err))
				return err
			}
		}

	}

	return nil
}

//删除
func OrderItemLimitVinDelete(id int64) (num int64, err error) {
	m := NewOrderItemLimitVin(id)
	if num, err := BaseDelete(&m); err != nil {
		return num, err
	} else {
		return num, nil
	}
}
