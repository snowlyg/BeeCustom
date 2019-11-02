package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
)

// TableName 设置HandBook表名
func (u *HandBook) TableName() string {
	return HandBookTBName()
}

// HandBook 实体类
type HandBook struct {
	BaseModel

	CompanyNo               string    `orm:"column(company_no);size(20);null" description:"企业内部编号"`
	CompanyManageCode       string    `orm:"column(company_manage_code);size(50)" description:"经营单位代码"`
	CompanyManageName       string    `orm:"column(company_manage_name);size(100)" description:"经营单位名称"`
	ContractNumber          string    `orm:"column(contract_number);size(100)" description:"账册编号"`
	CompanyClientCode       string    `orm:"column(company_client_code);size(10)" description:"加工单位代码"`
	CompanyClientName       string    `orm:"column(company_client_name);size(100)" description:"加工单位"`
	PutrecNo                string    `orm:"column(putrec_no);size(20);null" description:"预录入号(预录入统一编号)"`
	ForeignTradeCompanyName string    `orm:"column(foreign_trade_company_name);size(100);null" description:"外商公司"`
	PermitNumber            string    `orm:"column(permit_number);size(20);null" description:"批准证编号"`
	Manualslx               string    `orm:"column(manualslx);size(20)" description:"账册类型"`
	SuperviseMode           string    `orm:"column(supervise_mode);size(50);null" description:"监管方式 (贸易方式)"`
	SuperviseModeCode       string    `orm:"column(supervise_mode_code);size(20);null" description:"监管方式 (贸易方式)代码"`
	Permitcn                string    `orm:"column(permitcn);size(20);null" description:"批文账册号"`
	TaxationxzCode          string    `orm:"column(taxationxz_code);size(255);null" description:"征免性质代码"`
	TaxationxzName          string    `orm:"column(taxationxz_name);size(200);null" description:"征免性质"`
	ProcessingMode          string    `orm:"column(processing_mode);size(20);null" description:"加工种类"`
	OriginalityPquantity    string    `orm:"column(originality_pquantity);size(10);null" description:"进口货物项数"`
	FinishedPquantity       string    `orm:"column(finished_pquantity);size(10);null" description:"出口货物项数"`
	BondedMode              string    `orm:"column(bonded_mode);size(20);null" description:"保税方式名称"`
	InAmount                float64   `orm:"column(in_amount);null;digits(17);decimals(4)" description:"进口总金额(实际进口总金额)"`
	ContractNo              string    `orm:"column(contract_no);size(10);null" description:"协议号"`
	OutAmount               float64   `orm:"column(out_amount);null;digits(17);decimals(4)" description:"出口总金额(实际出口总金额)"`
	InMoneyunitCode         string    `orm:"column(in_moneyunit_code);size(20);null" description:"进口币制代码"`
	InMoneyunit             string    `orm:"column(in_moneyunit);size(255);null" description:"进口币制"`
	InMoneyunitEn           string    `orm:"column(in_moneyunit_en);size(100);null" description:"进口币制英文"`
	OutMoneyunitCode        string    `orm:"column(out_moneyunit_code);size(20);null" description:"出口币制代码"`
	OutMoneyunit            string    `orm:"column(out_moneyunit);size(255);null" description:"出口币制"`
	OutMoneyunitEn          string    `orm:"column(out_moneyunit_en);size(100);null" description:"出口币制英文"`
	InContractNo            string    `orm:"column(in_contract_no);size(255);null" description:"进口合同号"`
	OutContractNo           string    `orm:"column(out_contract_no);size(255);null" description:"出口合同号"`
	WarehouseVolume         string    `orm:"column(warehouse_volume);size(10);null" description:"仓库体积"`
	WarehouseArea           string    `orm:"column(warehouse_area);size(255);null" description:"仓库面积"`
	FinishedRate            string    `orm:"column(finished_rate);size(255);null" description:"成本率"`
	UllageMode              string    `orm:"column(ullage_mode);size(255);null" description:"损耗模式率"`
	ThroughPut              float64   `orm:"column(through_put);null;digits(17);decimals(5)" description:"生产能力"`
	ThroughPutUnit          string    `orm:"column(through_put_unit);size(255);null" description:"生产能力单位"`
	MaxWorkingCapital       string    `orm:"column(max_working_capital);size(255);null" description:"最大周转资金"`
	MaxWorkingCapitalUnit   string    `orm:"column(max_working_capital_unit);size(255);null" description:"最大周转资金单位"`
	Remark                  string    `orm:"column(remark);size(1000);null" description:"备注"`
	AccountType             int8      `orm:"column(account_type)" description:"账册类别 1：普通账册；2：二期账册"`
	CompanyManageCreditCode string    `orm:"column(company_manage_credit_code);size(18);null" description:"经营单位社会信用代码"`
	CompanyClientCreditCode string    `orm:"column(company_client_credit_code);size(18);null" description:"加工单位社会信用代码"`
	AgentCode               string    `orm:"column(agent_code);size(50);null" description:"申报单位代码"`
	AgentCodeScc            string    `orm:"column(agent_code_scc);size(18);null" description:"申报单位社会信用代码"`
	AgentName               string    `orm:"column(agent_name);size(100);null" description:"申报单位名称"`
	CompanyClientAreaCode   string    `orm:"column(company_client_area_code);size(255);null" description:"加工企业地区代码"`
	CompanyClientAreaName   string    `orm:"column(company_client_area_name);size(100);null" description:"加工企业地区"`
	AplCompanyType          string    `orm:"column(apl_company_type);size(100);null" description:"申报企业类型"`
	AplType                 string    `orm:"column(apl_type);size(100);null" description:"申报类型"`
	CompanyDocNo            string    `orm:"column(company_doc_no);size(100);null" description:"企业档案库编号"`
	OriginCount             string    `orm:"column(origin_count);size(100);null" description:"料件项数"`
	FinishCount             string    `orm:"column(finish_count);size(100);null" description:"成品项数"`
	Ciq                     string    `orm:"column(ciq);size(100);null" description:"主管海关"`
	InputCompanyCode        string    `orm:"column(input_company_code);size(50);null" description:"录入单位代码"`
	InputCompanyCodeScc     string    `orm:"column(input_company_code_scc);size(18);null" description:"录入单位社会信用代码"`
	InputCompanyName        string    `orm:"column(input_company_name);size(100);null" description:"录入单位名称"`
	RecordAt                string    `orm:"column(record_at);size(255);null" description:"备案批准日期"`
	ChangeAt                string    `orm:"column(change_at);size(255);null" description:"更变批准日期"`
	LastComplateAt          string    `orm:"column(last_complate_at);size(255);null" description:"最近核销日期"`
	UllagePro               string    `orm:"column(ullage_pro);size(255);null" description:"单耗申报环节"`
	UllageProVersionFlag    string    `orm:"column(ullage_pro_version_flag);size(255);null" description:"单耗版本号控制标志"`
	ChangTimes              string    `orm:"column(chang_times);size(255);null" description:"账册变更次数(延期次数)"`
	BiggestInAmount         string    `orm:"column(biggest_in_amount);size(255);null" description:"最大进口金额(美元)"`
	ComplateDays            string    `orm:"column(complate_days);size(255);null" description:"核销周期"`
	ComplateType            string    `orm:"column(complate_type);size(255);null" description:"核销类型"`
	AccountStartFlag        string    `orm:"column(account_start_flag);size(255);null" description:"账册执行标志"`
	ComplateStatus          string    `orm:"column(complate_status);size(255);null" description:"核销方式"`
	AccountFunction         string    `orm:"column(account_function);size(255);null" description:"账册用途"`
	AplDate                 time.Time `orm:"column(apl_date);type(datetime);null" description:"申报日期"`
	UsefulLife              time.Time `orm:"column(useful_life);type(datetime);null" description:"结束有效期"`
	UsefulLifeDays          int       `orm:"-" `
	PreentryDate            time.Time `orm:"column(preentry_date);type(datetime)" description:"录入日期"`
	Company                 *Company  `orm:"column(company_id);rel(fk)"`
	CompanyId               int64     `orm:"-" form:"CompanyId"`
}

// HandBookQueryParam 用于查询的类
type HandBookQueryParam struct {
	BaseQueryParam
	Type     string //模糊查询
	NameLike string //模糊查询
}

func NewHandBook(id int64) HandBook {
	return HandBook{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

//查询参数
func NewHandBookQueryParam() HandBookQueryParam {
	return HandBookQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

// HandBookOne 根据id获取单条
func HandBookOne(id int64) (*HandBook, error) {
	m := NewHandBook(0)
	o := orm.NewOrm()
	if err := o.QueryTable(HandBookTBName()).Filter("Id", id).RelatedSel().One(&m); err != nil {
		return nil, err
	}

	if &m == nil {
		return &m, errors.New("获取失败")
	}

	return &m, nil
}

//Save 添加、编辑页面 保存
func HandBookSave(m *HandBook) (*HandBook, error) {
	o := orm.NewOrm()
	if m.Id == 0 {
		if _, err := o.Insert(m); err != nil {
			return nil, err
		}
	} else {
		if _, err := o.Update(m); err != nil {
			return nil, err
		}
	}

	return m, nil
}

//删除
func HandBookDelete(id int64) (num int64, err error) {
	m := NewHandBook(id)
	if num, err := BaseDelete(&m); err != nil {
		return num, err
	} else {
		return num, nil
	}
}

//删除
func HandBookDeleteAll(clearanceType int8) (num int64, err error) {
	if num, err := BaseDeleteAll(clearanceType); err != nil {
		return num, err
	} else {
		return num, nil
	}
}

//批量插入
func InsertHandBookMulti(datas []*HandBook) (num int64, err error) {
	return BaseInsertMulti(len(datas), datas)
}
