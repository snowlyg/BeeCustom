package transforms

import "BeeCustom/models"

type HsCode struct {
	Id          int64
	Code        string
	Name        string
	License     string
	GeneralRate float64
	OfferRate   float64
	ExportRate  float64
	TaxRate     float64
	ConsumeRate float64
	Unit1       string
	Unit2       string
	Declaration string
	Remark      string
	Unit1Name   string `gtf:"Func.GetUnitCode(Unit1,计量单位代码)"`
	Unit2Name   string `gtf:"Func.GetUnitCode(Unit2,计量单位代码)"`
	CreatedAt   string
	UpdatedAt   string
}

func (h *HsCode) GetUnitCode(v, t string) interface{} {
	cs := models.GetClearancesByTypes(t, false)
	for _, c := range cs {
		if c[0] == v {
			return c[1]
		}
	}
	return nil
}
