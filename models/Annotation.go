package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"BeeCustom/utils"
	"github.com/astaxie/beego/orm"
)

// TableName 设置Annotation表名
func (u *Annotation) TableName() string {
	return AnnotationTBName()
}

// AnnotationQueryParam 用于查询的类
type AnnotationQueryParam struct {
	BaseQueryParam
}

// Annotation 实体类
type Annotation struct {
	BaseModel

	Status                 int8         `orm:"column(status)" description:"核注清单状态"`
	BondInvtNo             string       `orm:"column(bond_invt_no);size(64);null" description:"清单编号 (返填)"`
	SeqNo                  string       `orm:"column(seq_no);size(18);null" description:"清单预录入统一编号 (返填)"`
	PutrecNo               string       `orm:"column(putrec_no);size(64)" description:"备案编号（手(账)册编号）"`
	EtpsInnerInvtNo        string       `orm:"column(etps_inner_invt_no);size(64);null" description:"企业内部清单编号 （企业内部编号） (企业自行编写)"`
	BizopEtpsSccd          string       `orm:"column(bizop_etps_sccd);size(18);null" description:"经营企业社会信用代码"`
	BizopEtpsno            string       `orm:"column(bizop_etpsno);size(10)" description:"经营企业编号"`
	BizopEtpsNm            string       `orm:"column(bizop_etps_nm);size(512)" description:"经营企业名称"`
	RcvgdEtpsno            string       `orm:"column(rcvgd_etpsno);size(10)" description:"收货企业编号（加工单位）"`
	RvsngdEtpsSccd         string       `orm:"column(rvsngd_etps_sccd);size(18);null" description:"收发货企业社会信用代码（加工单位）"`
	RcvgdEtpsNm            string       `orm:"column(rcvgd_etps_nm);size(512)" description:"收货企业名称（加工单位）"`
	DclEtpsSccd            string       `orm:"column(dcl_etps_sccd);size(18);null" description:"申报企业社会信用代码"`
	DclEtpsno              string       `orm:"column(dcl_etpsno);size(10)" description:"申报企业编号"`
	DclEtpsNm              string       `orm:"column(dcl_etps_nm);size(512)" description:"申报企业名称"`
	EntryNo                string       `orm:"column(entry_no);size(64);null" description:"对应报关单编号(返填)清单报关时使用。海关端报关单入库时，反填并反馈企业端"`
	ImpexpPortcd           string       `orm:"column(impexp_portcd);size(4);null" description:"进/出境关别"`
	ImpexpPortcdName       string       `orm:"column(impexp_portcd_name);size(100);null" description:"进/出境关别名称"`
	DclPlcCuscd            string       `orm:"column(dcl_plc_cuscd);size(4);null" description:"申报地关区代码（主管海关）"`
	DclPlcCuscdName        string       `orm:"column(dcl_plc_cuscd_name);size(100);null" description:"申报地关区代码名称"`
	ImpexpMarkcd           string       `orm:"column(impexp_markcd);size(4);null" description:"进出口标记代码 I：进口,E：出口"`
	MtpckEndprdMarkcd      string       `orm:"column(mtpck_endprd_markcd);size(4);null" description:"料件成品标记代码 I：料件，E：成品"`
	MtpckEndprdMarkcdName  string       `orm:"column(mtpck_endprd_markcd_name);size(100);null" description:"料件成品标记名称"`
	SupvModecd             string       `orm:"column(supv_modecd);size(6);null" description:"监管方式代码"`
	SupvModecdName         string       `orm:"column(supv_modecd_name);size(100);null" description:"监管方式名称"`
	TrspModecd             string       `orm:"column(trsp_modecd);size(6);null" description:"运输方式代码"`
	TrspModecdName         string       `orm:"column(trsp_modecd_name);size(100);null" description:"运输方式名称"`
	DclcusFlag             string       `orm:"column(dclcus_flag);size(1);null" description:"是否报关标志：1.报关2.非报关"`
	DclcusFlagName         string       `orm:"column(dclcus_flag_name);size(100);null" description:"是否报关标志名称"`
	DclcusTypecd           string       `orm:"column(dclcus_typecd);size(25);null" description:"报关类型代码(1.关联报关2.对应报关；当报关标志为“1.报关”时，企业可选择“关联报关单”/“对应报关单”；当报关标志填写为“2.非报关”时，报关标志填写为“2.非报关”该项不可填。)"`
	DclcusTypecdName       string       `orm:"column(dclcus_typecd_name);size(100);null" description:"报关类型名称)"`
	VrfdedMarkcd           string       `orm:"column(vrfded_markcd);size(4);null" description:"核扣标记代码（核扣标志） (返填)"`
	InvtIochkptStucd       string       `orm:"column(invt_iochkpt_stucd);size(4);null" description:"清单进出卡口状态代码 (返填)"`
	ApplyNo                string       `orm:"column(apply_no);size(64);null" description:"申请表编号(单一显示：申报表编号)"`
	ListType               string       `orm:"column(list_type);size(1);null" description:"流转类型 (非流转类不填写，流转类填写: A：加工贸易深加工结转、B：加工贸易余料结转、C：不作价设备结转)"`
	ListTypeName           string       `orm:"column(list_type_name);size(100);null" description:"流转类型名称"`
	InputCode              string       `orm:"column(input_code);size(10)" description:"录入企业编号"`
	InputCreditCode        string       `orm:"column(input_credit_code);size(18);null" description:"录入企业社会信用代码"`
	InputName              string       `orm:"column(input_name);size(255)" description:"录入单位名称"`
	ListStat               string       `orm:"column(list_stat);size(1);null" description:"清单状态 (返填)"`
	CorrEntryDclEtpsSccd   string       `orm:"column(corr_entry_dcl_etps_sccd);size(18);null" description:"对应报关单申报单位社会统一信用代码"`
	CorrEntryDclEtpsNo     string       `orm:"column(corr_entry_dcl_etps_no);size(10);null" description:"对应报关单申报单位代码(当报关类型DCLCUS_TYPECD字段为2时，该字段必填)"`
	CorrEntryDclEtpsNm     string       `orm:"column(corr_entry_dcl_etps_nm);size(512);null" description:"对应报关单申报单位名称(当报关类型DCLCUS_TYPECD字段为2时，该字段必填)"`
	DecType                string       `orm:"column(dec_type);size(1);null" description:"报关单类型"`
	DecTypeName            string       `orm:"column(dec_type_name);size(100);null" description:"报关单类型名称"`
	StshipTrsarvNatcd      string       `orm:"column(stship_trsarv_natcd);size(3);null" description:"起运/运抵国(地区）（启运国(地区)）"`
	StshipTrsarvNatcdName  string       `orm:"column(stship_trsarv_natcd_name);size(100);null" description:"起运/运抵国(地区）(启运国(地区))名称"`
	InvtType               string       `orm:"column(invt_type);size(1);null" description:"清单类型 (SAS项目新增)"`
	InvtTypeName           string       `orm:"column(invt_type_name);size(100);null" description:"清单类型名称"`
	EntryStucd             string       `orm:"column(entry_stucd);size(1);null" description:"报关状态 (SAS项目新增)"`
	PassportUsedTypeCd     string       `orm:"column(passport_used_type_cd);size(1);null" description:"核放单生成标志代码 (SAS项目新增:返填)"`
	Rmk                    string       `orm:"column(rmk);null" description:"备注"`
	DecRmk                 string       `orm:"column(dec_rmk);null" description:"报关单草稿(备注)"`
	DclTypecd              string       `orm:"column(dcl_typecd);size(1);null" description:"申报类型"`
	NeedEntryModified      string       `orm:"column(need_entry_modified);size(1);null" description:"报关单同步修改标志"`
	LevyBlAmt              string       `orm:"column(levy_bl_amt);size(25);null" description:"计征金额"`
	ChgTmsCnt              string       `orm:"column(chg_tms_cnt);size(64);null" description:"变更次数 (有变更时填写)"`
	RltInvtNo              string       `orm:"column(rlt_invt_no);size(64);null" description:"关联清单编号(结转类专用)"`
	RltPutrecNo            string       `orm:"column(rlt_putrec_no);size(64);null" description:"关联备案编号(结转类专用)"`
	RltEntryNo             string       `orm:"column(rlt_entry_no);size(64);null" description:"关联报关单编号(可录入或者系统自动生成报关单后返填二线取消报关的情况下使用，用于生成区外一般贸易报关单。暂未使用)"`
	RltEntryBizopEtpsSccd  string       `orm:"column(rlt_entry_bizop_etps_sccd);size(18);null" description:"关联报关单境内收发货人社会信用代码(二线取消报关的情况下使用，用于生成区外一般贸易报关单。暂未使用)"`
	RltEntryBizopEtpsno    string       `orm:"column(rlt_entry_bizop_etpsno);size(10);null" description:"关联报关单境内收发货人编号(当报关类型DCLCUS_TYPECD字段为1时，该字段必填报关类型为关联报关时必填。二线取消报关的情况下使用，用于生成区外一般贸易报关单。暂未使用)"`
	RltEntryBizopEtpsNm    string       `orm:"column(rlt_entry_bizop_etps_nm);size(512);null" description:"关联报关单境内收发货人名称(当报关类型DCLCUS_TYPECD字段为1时，该字段必填同上)"`
	RltEntryRvsngdEtpsSccd string       `orm:"column(rlt_entry_rvsngd_etps_sccd);size(18);null" description:"关联报关单收发货单位社会统一信用代码(二线取消报关的情况下使用，用于生成区外一般贸易报关单。暂未使用)"`
	RltEntryRcvgdEtpsno    string       `orm:"column(rlt_entry_rcvgd_etpsno);size(10);null" description:"关联报关单海关收发货单位编码(当报关类型DCLCUS_TYPECD字段为1时，该字段必填报关类型为关联报关时必填。二线取消报关的情况下使用，用于生成区外一般贸易报关单。)"`
	RltEntryRcvgdEtpsNm    string       `orm:"column(rlt_entry_rcvgd_etps_nm);size(512);null" description:"关联报关单收发货单位名称(当报关类型DCLCUS_TYPECD字段为1时，该字段必填)"`
	RltEntryDclEtpsSccd    string       `orm:"column(rlt_entry_dcl_etps_sccd);size(18);null" description:"关联报关单申报单位社会统一信用代码(二线取消报关的情况下使用，用于生成区外一般贸易报关单。暂未使用)"`
	RltEntryDclEtpsno      string       `orm:"column(rlt_entry_dcl_etpsno);size(10);null" description:"关联报关单海关申报单位编码(当报关类型DCLCUS_TYPECD字段为1时，该字段必填报关类型为关联报关时必填。二线取消报关的情况下使用，用于生成区外一般贸易报关单。)"`
	RltEntryDclEtpsNm      string       `orm:"column(rlt_entry_dcl_etps_nm);size(512);null" description:"关联报关单申报单位名称(当报关类型DCLCUS_TYPECD字段为1时，该字段必填)"`
	Param1                 string       `orm:"column(param1);size(19);null" description:"备用1"`
	Param2                 string       `orm:"column(param2);size(19);null" description:"备用2"`
	Param3                 string       `orm:"column(param3);size(19);null" description:"备用3"`
	SysId                  string       `orm:"column(sys_id);size(2)" description:"子系统ID 95 加工贸易账册系统;B1 加工贸易手册系统 ;B2 加工贸易担保管理系统;B3 保税货物流转系统二期 ;Z7 海关特殊监管区域管理系统;Z8 保税物流管理系统"`
	OperCusRegCode         string       `orm:"column(oper_cus_reg_code);size(10)" description:"操作卡的海关十位"`
	KeyName                string       `orm:"column(key_name);size(255);null" description:"签名所用的证书信息"`
	Version                string       `orm:"column(version);size(255);null" description:"版本编号"`
	BusinessId             string       `orm:"column(business_id);size(255);null" description:"业务单证号"`
	MessageId              string       `orm:"column(message_id);size(255);null" description:"报文唯一编号"`
	FileName               string       `orm:"column(file_name);size(255);null" description:"用户原始报文名，主要用于用户查询"`
	MessageType            string       `orm:"column(message_type);size(255);null" description:"报文类型"`
	SenderId               string       `orm:"column(sender_id);size(255);null" description:"发送方编号"`
	ReceiverId             string       `orm:"column(receiver_id);size(255);null" description:"接收方编号"`
	DelcareFlag            string       `orm:"column(delcare_flag);size(255)" description:"申报标志 0--暂存；1--申报"`
	ExtraRemark            string       `orm:"column(extra_remark);null" description:"附注"`
	Creator                string       `orm:"column(creator);size(50);null" description:"创建人"`
	GenDecFlag             string       `orm:"column(gen_dec_flag);size(2);null" description:"是否生成报关单 1 生成，2 不生成"`
	GenDecFlagName         string       `orm:"column(gen_dec_flag_name);size(100);null" description:"是否生成报关单名称"`
	InputTime              time.Time    `form:"-" orm:"column(input_time);type(datetime);null" description:"录入日期"`
	PrevdTime              time.Time    `form:"-" orm:"column(prevd_time);type(datetime);null" description:"预核扣时间"`
	FormalVrfdedTime       time.Time    `form:"-" orm:"column(formal_vrfded_time);type(datetime);null" description:"正式核扣时间 (返填)"`
	DeletedAt              time.Time    `form:"-" orm:"column(deleted_at);type(timestamp);null" `
	StatusUpdatedAt        time.Time    `form:"-" orm:"column(status_updated_at);type(datetime)" description:"状态更新时间"`
	InvtDclTime            time.Time    `form:"-" orm:"column(invt_dcl_time);type(datetime);null" description:"清单申报时间(清单申报日期)(返填)"`
	EntryDclTime           time.Time    `form:"-" orm:"column(entry_dcl_time);type(datetime);null" description:"报关单申报日期(返填)清单报关时使用。海关端报关单入库时，反填并反馈企业端"`
	BackendUser            *BackendUser `orm:"column(user_id);rel(fk)"`
	Company                *Company     `orm:"column(company_id);rel(fk)"`
	UserId                 int64        `orm:"-" form:"UserId"`    //关联管理会自动生成 UserId 字段，此处不生成字段
	CompanyId              int64        `orm:"-" form:"CompanyId"` //关联管理会自动生成 CompanyId 字段，此处不生成字段
	OrderId                int64        `orm:"-" form:"OrderId"`   //关联管理会自动生成 OrderId 字段，此处不生成字段
}

func NewAnnotation(id int64) Annotation {
	return Annotation{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

//查询参数
func NewAnnotationQueryParam() AnnotationQueryParam {
	return AnnotationQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

// AnnotationPageList 获取分页数据
func AnnotationPageList(params *AnnotationQueryParam) ([]*Annotation, int64) {
	query := orm.NewOrm().QueryTable(AnnotationTBName())
	datas := make([]*Annotation, 0)

	total, _ := query.Count()
	query = BaseListQuery(query, params.Sort, params.Order, params.Limit, params.Offset)
	_, _ = query.All(&datas)

	return datas, total
}

func AnnotationGetRelations(ms []*Annotation, relations string) ([]*Annotation, error) {
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

// AnnotationOne 根据id获取单条
func AnnotationOne(id int64) (*Annotation, error) {
	m := NewAnnotation(0)
	o := orm.NewOrm()
	if err := o.QueryTable(AnnotationTBName()).Filter("Id", id).RelatedSel().One(&m); err != nil {
		return nil, err
	}

	if &m == nil {
		return &m, errors.New("用户获取失败")
	}

	return &m, nil
}

// AnnotationOneByUserName 根据用户名密码获取单条
func AnnotationOneByUserName(username, userpwd string) (*Annotation, error) {
	m := NewAnnotation(0)
	err := orm.NewOrm().QueryTable(AnnotationTBName()).Filter("username", username).Filter("userpwd", userpwd).One(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

//Save 添加、编辑页面 保存
func AnnotationSave(m *Annotation) (*Annotation, error) {
	o := orm.NewOrm()
	if m.Id == 0 {
		if _, err := o.Insert(m); err != nil {
			utils.LogDebug(fmt.Sprintf("AnnotationSave:%v", err))
			return nil, err
		}
	} else {

		if _, err := o.Update(m); err != nil {
			utils.LogDebug(fmt.Sprintf("AnnotationSave:%v", err))
			return nil, err
		}
	}

	return m, nil
}

//Save 添加、编辑页面 保存
func AnnotationFreeze(m *Annotation) (*Annotation, error) {
	o := orm.NewOrm()
	if _, err := o.Update(m, "Status"); err != nil {
		return nil, err
	}

	return m, nil
}

//删除
func AnnotationDelete(id int64) (num int64, err error) {
	m := NewAnnotation(id)
	if num, err := BaseDelete(&m); err != nil {
		return num, err
	} else {
		return num, nil
	}
}
