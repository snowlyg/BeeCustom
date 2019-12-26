package transforms

import (
	"BeeCustom/models"
)

type HandBookGood struct {
	Id                  int64
	Type                int8
	Serial              string
	RecordNo            string
	HsCode              string
	Name                string
	ClassificationMark  string
	Special             string
	UnitOne             string
	UnitTwo             string
	UnitThree           string
	UnitOneCode         string `gtf:"Func.GetUnitCode(UnitOne,计量单位代码)"`
	UnitTwoCode         string `gtf:"Func.GetUnitCode(UnitTwo,计量单位代码)"`
	UnitThreeCode       string `gtf:"Func.GetUnitCode(UnitThreeCode,计量单位代码)"`
	Price               float64
	Moneyunit           string
	MoneyunitCode       string `gtf:"Func.GetMoneyunitCode(Moneyunit,货币代码)"`
	Quantity            float64
	MaxAllowance        float64
	InitialQuantity     float64
	UnitTwoProportion   float64
	UnitThreeProportion float64
	WeightProportion    float64
	Taxationlx          string
	DeclareMode         string
	Remark              string
	HandleMark          string
	CompanyActionFlag   string
	CustomActionFlag    string
	StartCount          string
	CountControlFlag    string
	BigCount            string
	UllageFlag          string
	ConsultMark         string
	MainMark            string
	Amount              float64
	Manuplace           string
	GoodAttr            string
	SeqNo               string
	CreatedAt           string
	UpdatedAt           string
}

func (h *HandBookGood) GetMoneyunitCode(v, t string) interface{} {
	cs := models.GetClearancesByTypes(t, true)
	for _, c := range cs {
		if c[0] == v {
			return c[1]
		}
	}
	return nil
}

func (h *HandBookGood) GetUnitCode(v, t string) interface{} {
	cs := models.GetClearancesByTypes(t, false)
	for _, c := range cs {
		if c[0] == v {
			return c[1]
		}
	}
	return nil
}
