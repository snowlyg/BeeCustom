package models

import (
	"BeeCustom/utils"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

// TableName 设置Company表名
func (u *Company) TableName() string {
	return CompanyTBName()
}

// Company 实体类
type Company struct {
	BaseModel

	Number              string    `orm:"column(number);size(10)" description:"海关编号" valid:"Required;Length(10)"`
	Name                string    `orm:"column(name);size(200)" description:"全称" valid:"Required;MaxSize(255)"`
	Short               string    `orm:"column(short);size(255);null" description:"简称"`
	Registration        string    `orm:"column(registration);size(10);null" description:"商检注册号" valid:"Required;Length(10)"`
	Address             string    `orm:"column(address);size(200);null" description:"地址"`
	DeclareType         int8      `orm:"column(declare_type);null" description:"申报方式:代理，自理"`
	RegistrationCode    string    `orm:"column(registration_code);size(100);null" description:"产地证注册号"`
	Phone               string    `orm:"column(phone);size(20);null" description:"电话"`
	Fax                 string    `orm:"column(fax);size(20);null" description:"传真"`
	CreditCode          string    `orm:"column(credit_code);size(18)" description:"信用代码" valid:"Required;Length(18)"`
	BusinessName        string    `orm:"column(business_name);size(100);null" description:"经营单位名称" valid:"Length(10)"`
	BusinessCode        string    `orm:"column(business_code);size(10);null" description:"经营单位代码"`
	Bank                string    `orm:"column(bank);size(50);null" description:"银行账号" valid:"MaxSize(20)"`
	CustomCertification int8      `orm:"column(custom_certification);null" description:"海关认证：高级认证，一般认证，一般信用，失信企业"`
	CompanyType         int8      `orm:"column(company_type);null" description:"企业类别：往来客户，保税仓，供应商，代理报关公司，代理报检公司，物流公司"`
	CompanyKind         int8      `orm:"column(company_kind);null" description:"企业性质：国有，合作，合资， 独资， 集体，私营"`
	ControlRating       int8      `orm:"column(control_rating);null" description:"风控评级"`
	Remark              string    `orm:"column(remark);size(1000);null" description:"备注"`
	IsOpenSubEmail      int8      `orm:"column(is_open_sub_email);null" `
	IsOpenSubPhone      int8      `orm:"column(is_open_sub_phone);null" `
	SubPhone            string    `orm:"column(sub_phone);null" description:"订阅手机 多个用，隔开"`
	SubEmail            string    `orm:"column(sub_email);null" description:"订阅邮箱 多个用，隔开"`
	SubContentCheck     int8      `orm:"column(sub_content_check);null" description:"订阅内容 审核通过"`
	SubContentSubmit    int8      `orm:"column(sub_content_submit);null" description:"订阅内容 已提交海关处理"`
	SubContentReject    int8      `orm:"column(sub_content_reject);null" description:"订阅内容 驳回信息"`
	SubContentPass      int8      `orm:"column(sub_content_pass);null" description:"订阅内容 机关放行"`
	StatementDate       int8      `orm:"column(statement_date)" description:"生成账单日期"`
	IsTrade             int8      `orm:"column(is_trade)" description:"境内收发货单位 是否开启"`
	IsOwner             int8      `orm:"column(is_owner)" description:"生产销售单位 是否开启"`
	Business            string    `orm:"column(business);size(255);null" description:"营业执照 file_path"`
	BusinessAuditStatus int8      `orm:"column(business_audit_status);null" description:"营业执照审核状态"`
	BusinessAuditAt     time.Time `orm:"column(business_audit_at);type(datetime);null" description:"营业执照审核时间"`
	Tax                 int8      `orm:"column(tax)" description:"税率"`

	BackendUserId int64        `orm:"-" form:"BackendUserId"`
	BackendUser   *BackendUser `orm:"column(user_id);rel(fk);null"`

	CompanyContacts []*CompanyContact `orm:"reverse(many)"` //设置一对多关系
	CompanyForeigns []*CompanyForeign `orm:"reverse(many)"` //设置一对多关系
	CompanySeals    []*CompanySeal    `orm:"reverse(many)"` //设置一对多关系
	HandBooks       []*HandBook       `orm:"reverse(many)"` //设置一对多关系
	Annotations     []*Annotation     `orm:"reverse(many)"` //设置一对多关系
}

// CompanyQueryParam 用于查询的类
type CompanyQueryParam struct {
	BaseQueryParam

	NameLike   string //模糊查询
	SearchWord string //模糊查询

}

func NewCompany(id int64) Company {
	return Company{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

//查询参数
func NewCompanyQueryParam() CompanyQueryParam {
	return CompanyQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

// CompanyPageList 获取分页数据
func CompanyPageList(params *CompanyQueryParam) ([]*Company, int64) {
	query := orm.NewOrm().QueryTable(CompanyTBName())
	data := make([]*Company, 0)

	if len(params.NameLike) > 0 {
		cond := orm.NewCondition()
		cond1 := cond.AndCond(cond.And("number", params.NameLike)).
			OrCond(cond.And("name__istartswith", params.NameLike)).
			OrCond(cond.And("credit_code", params.NameLike)).
			OrCond(cond.And("business_code", params.NameLike))
		query = query.SetCond(cond1)
	}

	if len(params.SearchWord) > 0 {
		query = query.Distinct().Filter("HandBooks__contract_number__iexact", params.SearchWord)
	}

	total, _ := query.Count()
	query = BaseListQuery(query, params.Sort, params.Order, params.Limit, params.Offset)

	_, _ = query.All(&data)

	return data, total
}

func CompaniesGetRelations(rs []*Company, relations string) ([]*Company, error) {
	for _, rv := range rs {
		err := CompanyGetRelations(rv, relations)
		if err != nil {
			return nil, err
		}
	}

	return rs, nil
}

func CompanyGetRelations(v *Company, relations string) error {
	o := orm.NewOrm()
	rs := strings.Split(relations, ",")
	for _, rv := range rs {
		_, err := o.LoadRelated(v, rv)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("LoadRelated:%v", err))
			return err
		}

		for _, hv := range v.HandBooks {
			hv.UsefulLifeDays = int(math.Round(hv.UsefulLife.Sub(time.Now()).Hours() / 24))
		}

	}

	return nil
}

// CompanyByManageCode 根据海关编码 获取单条
func CompanyByManageCode(manageCode string) (*Company, error) {
	m := NewCompany(0)
	o := orm.NewOrm()

	if err := o.QueryTable(CompanyTBName()).Filter("Number", manageCode).One(&m); err != nil {
		utils.LogDebug(fmt.Sprintf("CompanyByManageCode:%v", err))
		return nil, err
	}

	if &m == nil {
		return &m, errors.New("获取失败")
	}

	return &m, nil
}

// CompanyOne 根据id获取单条
func CompanyOne(id int64, relations string) (*Company, error) {
	m := NewCompany(0)
	o := orm.NewOrm()

	if err := o.QueryTable(CompanyTBName()).Filter("Id", id).One(&m); err != nil {
		return nil, err
	}

	if len(relations) > 0 {
		err := CompanyGetRelations(&m, relations)
		if err != nil {
			return nil, err
		}
	}

	if &m == nil {
		return &m, errors.New("获取失败")
	}

	return &m, nil
}

//Save 添加、编辑页面 保存
func CompanySave(m *Company) (*Company, error) {
	o := orm.NewOrm()
	if m.Id == 0 {
		if err := getCompanyBackendUser(m); err != nil {
			return nil, err
		}

		if _, err := o.Insert(m); err != nil {
			return nil, err
		}
	} else {
		if err := getCompanyBackendUser(m); err != nil {
			return nil, err
		}

		if _, err := o.Update(m); err != nil {
			return nil, err
		}
	}

	return m, nil
}

//获取关联模型
func getCompanyBackendUser(m *Company) error {
	if bU, err := BackendUserOne(m.BackendUserId); err != nil {
		return err
	} else {
		m.BackendUser = bU
	}
	return nil
}

//删除
func CompanyDelete(id int64) (num int64, err error) {
	m := NewCompany(id)
	if num, err := BaseDelete(&m); err != nil {
		return num, err
	} else {
		return num, nil
	}
}
