package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"BeeCustom/utils"
	"BeeCustom/xlsx"
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
	Permitcn                string    `orm:"column(permitcn);size(20);null" description:"批文账册号"`
	TaxationxzName          string    `orm:"column(taxationxz_name);size(200);null" description:"征免性质"`
	ProcessingMode          string    `orm:"column(processing_mode);size(20);null" description:"加工种类"`
	OriginalityPquantity    string    `orm:"column(originality_pquantity);size(10);null" description:"进口货物项数"`
	FinishedPquantity       string    `orm:"column(finished_pquantity);size(10);null" description:"出口货物项数"`
	BondedMode              string    `orm:"column(bonded_mode);size(20);null" description:"保税方式名称"`
	InAmount                float64   `orm:"column(in_amount);null;digits(17);decimals(4)" description:"进口总金额(实际进口总金额)"`
	ContractNo              string    `orm:"column(contract_no);size(10);null" description:"协议号"`
	OutAmount               float64   `orm:"column(out_amount);null;digits(17);decimals(4)" description:"出口总金额(实际出口总金额)"`
	InMoneyunit             string    `orm:"column(in_moneyunit);size(255);null" description:"进口币制"`
	OutMoneyunit            string    `orm:"column(out_moneyunit);size(255);null" description:"出口币制"`
	InContractNo            string    `orm:"column(in_contract_no);size(255);null" description:"进口合同号"`
	OutContractNo           string    `orm:"column(out_contract_no);size(255);null" description:"出口合同号"`
	WarehouseVolume         string    `orm:"column(warehouse_volume);size(10);null" description:"仓库体积"`
	WarehouseArea           string    `orm:"column(warehouse_area);size(255);null" description:"仓库面积"`
	FinishedRate            string    `orm:"column(finished_rate);size(255);null" description:"成本率"`
	UllageMode              string    `orm:"column(ullage_mode);size(255);null" description:"损耗模式率"`
	ThroughPut              float64   `orm:"column(through_put);null;digits(17);decimals(5)" description:"生产能力"`
	ThroughPutUnit          string    `orm:"column(through_put_unit);size(255);null" description:"生产能力单位"`
	MaxWorkingCapital       string    `orm:"column(max_working_capital);size(255);null" description:"最大周转资金"`
	Remark                  string    `orm:"column(remark);size(1000);null" description:"备注"`
	Type                    int8      `orm:"column(type)" description:"账册类别 1：普通账册；2：二期账册"`
	CompanyManageCreditCode string    `orm:"column(company_manage_credit_code);size(18);null" description:"经营单位社会信用代码"`
	CompanyClientCreditCode string    `orm:"column(company_client_credit_code);size(18);null" description:"加工单位社会信用代码"`
	AgentCode               string    `orm:"column(agent_code);size(50);null" description:"申报单位代码"`
	AgentCodeScc            string    `orm:"column(agent_code_scc);size(18);null" description:"申报单位社会信用代码"`
	AgentName               string    `orm:"column(agent_name);size(100);null" description:"申报单位名称"`
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
	Bank                    string    `orm:"column(bank);size(255);null" description:"台账银行"`
	StopInOutFlag           string    `orm:"column(stop_in_out_flag);size(100);null" description:"暂停进出口标记"`
	FirstExportAt           string    `orm:"column(first_export_at);size(255);null" description:"首次出口日期"`
	CompanyContact          string    `orm:"column(company_contact);size(255);null" description:"企业联系人"`
	CompanyContactPhone     string    `orm:"column(company_contact_phone);size(255);null" description:"联系人手机号"`
	StopChangeFlag          string    `orm:"column(stop_change_flag);size(100);null" description:"暂停变更标记"`
	SelfAuditFlag           string    `orm:"column(self_audit_flag);size(255);null" description:"自核资格标记"`
	ManualChangTimes        string    `orm:"column(manual_chang_times);size(255);null" description:"手册变更次数"`
	ManualType              string    `orm:"column(manual_type);size(255);null" description:"手册类型"`

	UsefulLifeDays int             `orm:"-" `
	PreentryDate   time.Time       `orm:"column(preentry_date);type(datetime)" description:"录入日期"`
	Company        *Company        `orm:"column(company_id);rel(fk)"`
	CompanyId      int64           `orm:"-" form:"CompanyId"`
	HandBookGoods  []*HandBookGood `orm:"reverse(many)"` //设置一对多关系
}

// HandBookQueryParam 用于查询的类
type HandBookQueryParam struct {
	BaseQueryParam
	Type     string //模糊查询
	NameLike string //模糊查询
}

// HandBookImportParam 用于导入的类
type HandBookImportParam struct {
	xlsx.BaseImportParam

	HandBookGoodType int8

	HandBook        HandBook
	HandBookGoods   []*HandBookGood
	HandBookUllages []*HandBookUllage
}

// HandBookGoodImportParam 用于导入的类
type HandBookGoodImportParam struct {
	ExcelNameString    string
	ExcelTitleString   string
	HandBookTypeString string
}

func NewHandBook(id int64) HandBook {
	return HandBook{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

func HandBookGetRelations(v *HandBook, relations string) (*HandBook, error) {
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

// HandBookOne 根据id获取单条
func HandBookOne(id int64, relations string) (*HandBook, error) {
	m := NewHandBook(0)
	o := orm.NewOrm()
	if err := o.QueryTable(HandBookTBName()).Filter("Id", id).RelatedSel().One(&m); err != nil {
		return nil, err
	}

	if len(relations) > 0 {
		_, err := HandBookGetRelations(&m, relations)
		if err != nil {
			return nil, err
		}
	}

	if &m == nil {
		return &m, errors.New("获取失败")
	}

	return &m, nil
}

// GetHandBookByContractNumber 根据contractNumber获取单条
func GetHandBookByContractNumber(contractNumber string) (*HandBook, error) {
	m := NewHandBook(0)
	o := orm.NewOrm()
	if err := o.QueryTable(HandBookTBName()).Filter("ContractNumber", contractNumber).One(&m); err != nil {
		utils.LogDebug(fmt.Sprintf("GetHandBookByContractNumber:%v", err))
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
