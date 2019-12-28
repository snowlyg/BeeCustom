package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"BeeCustom/enums"

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

	ImpexpMarkcd        string
	TrspModecd          string
	StatusString        string
	SearchTimeString    string
	EtpsInnerInvtNoLike string
	IsDelete            bool
}

// Annotation 实体类
type Annotation struct {
	BaseModel

	Status                   int8      `orm:"column(status)"  valid:"Required" description:"核注清单状态"`
	BondInvtNo               string    `orm:"column(bond_invt_no);size(64);null"  valid:"MaxSize(64)" description:"清单编号 (返填)"`
	SeqNo                    string    `orm:"column(seq_no);size(18);null" valid:"MaxSize(18)" description:"清单预录入统一编号 (返填)"`
	PutrecNo                 string    `orm:"column(putrec_no);size(64)"  valid:"Required;MaxSize(64)" description:"备案编号（手(账)册编号）"`
	EtpsInnerInvtNo          string    `orm:"column(etps_inner_invt_no);size(64);null"  valid:"Required;MaxSize(64)" description:"企业内部清单编号 （企业内部编号） (企业自行编写)"`
	BizopEtpsSccd            string    `orm:"column(bizop_etps_sccd);size(18);null"  valid:"Required;Length(18)" description:"经营企业社会信用代码"`
	BizopEtpsno              string    `orm:"column(bizop_etpsno);size(10)"  valid:"Required;Length(10)" description:"经营企业编号"`
	BizopEtpsNm              string    `orm:"column(bizop_etps_nm);size(512)" valid:"Required;MaxSize(512)" description:"经营企业名称"`
	RcvgdEtpsno              string    `orm:"column(rcvgd_etpsno);size(10)" valid:"Required;Length(10)" description:"收货企业编号（加工单位）"`
	RvsngdEtpsSccd           string    `orm:"column(rvsngd_etps_sccd);size(18);null" valid:"Required;Length(18)" description:"收发货企业社会信用代码（加工单位）"`
	RcvgdEtpsNm              string    `orm:"column(rcvgd_etps_nm);size(512)" valid:"Required;MaxSize(512)" description:"收货企业名称（加工单位）"`
	DclEtpsSccd              string    `orm:"column(dcl_etps_sccd);size(18);null" valid:"Required;Length(18)" description:"申报企业社会信用代码"`
	DclEtpsno                string    `orm:"column(dcl_etpsno);size(10)" valid:"Required;Length(10)" description:"申报企业编号"`
	DclEtpsNm                string    `orm:"column(dcl_etps_nm);size(512)" valid:"Required;MaxSize(512)" description:"申报企业名称"`
	EntryNo                  string    `orm:"column(entry_no);size(64);null" description:"对应报关单编号(返填)清单报关时使用。海关端报关单入库时，反填并反馈企业端"`
	ImpexpPortcd             string    `orm:"column(impexp_portcd);size(4);null" valid:"Required;Length(4)" description:"进/出境关别"`
	ImpexpPortcdName         string    `orm:"column(impexp_portcd_name);size(100);null" valid:"Required;MaxSize(100)" description:"进/出境关别名称"`
	DclPlcCuscd              string    `orm:"column(dcl_plc_cuscd);size(4);null" valid:"Required;Length(4)" description:"申报地关区代码（主管海关）"`
	DclPlcCuscdName          string    `orm:"column(dcl_plc_cuscd_name);size(100);null" valid:"Required;MaxSize(100)" description:"申报地关区代码名称"`
	ImpexpMarkcd             string    `orm:"column(impexp_markcd);size(4);null" valid:"Required;MaxSize(4)" description:"进出口标记代码 I：进口,E：出口"`
	MtpckEndprdMarkcd        string    `orm:"column(mtpck_endprd_markcd);size(4);null" valid:"Required;MaxSize(4)" description:"料件成品标记代码 I：料件，E：成品"`
	MtpckEndprdMarkcdName    string    `orm:"column(mtpck_endprd_markcd_name);size(100);null" valid:"Required;MaxSize(100)" description:"料件成品标记名称"`
	SupvModecd               string    `orm:"column(supv_modecd);size(6);null" valid:"Required;MaxSize(6)" description:"监管方式代码"`
	SupvModecdName           string    `orm:"column(supv_modecd_name);size(100);null" valid:"Required;MaxSize(100)" description:"监管方式名称"`
	TrspModecd               string    `orm:"column(trsp_modecd);size(6);null" valid:"Required;MaxSize(6)" description:"运输方式代码"`
	TrspModecdName           string    `orm:"column(trsp_modecd_name);size(100);null" valid:"Required;MaxSize(100)"  description:"运输方式名称"`
	DclcusFlag               string    `orm:"column(dclcus_flag);size(1);null" valid:"Required;Length(1)"  description:"是否报关标志：1.报关2.非报关"`
	DclcusFlagName           string    `orm:"column(dclcus_flag_name);size(100);null" valid:"Required;MaxSize(100)"  description:"是否报关标志名称"`
	DclcusTypecd             string    `orm:"column(dclcus_typecd);size(25);null"   description:"报关类型代码(1.关联报关2.对应报关；当报关标志为“1.报关”时，企业可选择“关联报关单”/“对应报关单”；当报关标志填写为“2.非报关”时，报关标志填写为“2.非报关”该项不可填。)"`
	DclcusTypecdName         string    `orm:"column(dclcus_typecd_name);size(100);null"   description:"报关类型名称)"`
	VrfdedMarkcd             string    `orm:"column(vrfded_markcd);size(4);null" description:"核扣标记代码（核扣标志） (返填)"`
	InvtIochkptStucd         string    `orm:"column(invt_iochkpt_stucd);size(4);null" description:"清单进出卡口状态代码 (返填)"`
	ApplyNo                  string    `orm:"column(apply_no);size(64);null" description:"申请表编号(单一显示：申报表编号)"`
	ListType                 string    `orm:"column(list_type);size(1);null" description:"流转类型 (非流转类不填写，流转类填写: A：加工贸易深加工结转、B：加工贸易余料结转、C：不作价设备结转)"`
	ListTypeName             string    `orm:"column(list_type_name);size(100);null" description:"流转类型名称"`
	InputCode                string    `orm:"column(input_code);size(10)" valid:"Required;Length(10)"  description:"录入企业编号"`
	InputCreditCode          string    `orm:"column(input_credit_code);size(18);null" valid:"Required;Length(18)"  description:"录入企业社会信用代码"`
	InputName                string    `orm:"column(input_name);size(255)" valid:"Required;MaxSize(255)"  description:"录入单位名称"`
	ListStat                 string    `orm:"column(list_stat);size(1);null" description:"清单状态 (返填)"`
	CorrEntryDclEtpsSccd     string    `orm:"column(corr_entry_dcl_etps_sccd);size(18);null" valid:"Required;Length(18)"  description:"对应报关单申报单位社会统一信用代码"`
	CorrEntryDclEtpsNo       string    `orm:"column(corr_entry_dcl_etps_no);size(10);null" valid:"Required;Length(10)"  description:"对应报关单申报单位代码(当报关类型DCLCUS_TYPECD字段为2时，该字段必填)"`
	CorrEntryDclEtpsNm       string    `orm:"column(corr_entry_dcl_etps_nm);size(512);null" valid:"Required;MaxSize(512)"  description:"对应报关单申报单位名称(当报关类型DCLCUS_TYPECD字段为2时，该字段必填)"`
	DecType                  string    `orm:"column(dec_type);size(1);null" description:"报关单类型"`
	DecTypeName              string    `orm:"column(dec_type_name);size(100);null" description:"报关单类型名称"`
	StshipTrsarvNatcd        string    `orm:"column(stship_trsarv_natcd);size(3);null" valid:"Required;Length(3)" description:"起运/运抵国(地区）（启运国(地区)）"`
	StshipTrsarvNatcdName    string    `orm:"column(stship_trsarv_natcd_name);size(100);null"  valid:"Required;MaxSize(100)" description:"起运/运抵国(地区）(启运国(地区))名称"`
	InvtType                 string    `orm:"column(invt_type);size(1);null" description:"清单类型 (SAS项目新增)"`
	InvtTypeName             string    `orm:"column(invt_type_name);size(100);null" description:"清单类型名称"`
	EntryStucd               string    `orm:"column(entry_stucd);size(1);null" description:"报关状态 (SAS项目新增)"`
	PassportUsedTypeCd       string    `orm:"column(passport_used_type_cd);size(1);null" description:"核放单生成标志代码 (SAS项目新增:返填)"`
	Rmk                      string    `orm:"column(rmk);null" description:"备注"`
	DecRmk                   string    `orm:"column(dec_rmk);null" description:"报关单草稿(备注)"`
	DclTypecd                string    `orm:"column(dcl_typecd);size(1);null" description:"申报类型"`
	NeedEntryModified        string    `orm:"column(need_entry_modified);size(1);null" description:"报关单同步修改标志"`
	LevyBlAmt                string    `orm:"column(levy_bl_amt);size(25);null" description:"计征金额"`
	ChgTmsCnt                string    `orm:"column(chg_tms_cnt);size(64);null" description:"变更次数 (有变更时填写)"`
	RltInvtNo                string    `orm:"column(rlt_invt_no);size(64);null" description:"关联清单编号(结转类专用)"`
	RltPutrecNo              string    `orm:"column(rlt_putrec_no);size(64);null" description:"关联备案编号(结转类专用)"`
	RltEntryNo               string    `orm:"column(rlt_entry_no);size(64);null" description:"关联报关单编号(可录入或者系统自动生成报关单后返填二线取消报关的情况下使用，用于生成区外一般贸易报关单。暂未使用)"`
	RltEntryBizopEtpsSccd    string    `orm:"column(rlt_entry_bizop_etps_sccd);size(18);null" description:"关联报关单境内收发货人社会信用代码(二线取消报关的情况下使用，用于生成区外一般贸易报关单。暂未使用)"`
	RltEntryBizopEtpsno      string    `orm:"column(rlt_entry_bizop_etpsno);size(10);null" description:"关联报关单境内收发货人编号(当报关类型DCLCUS_TYPECD字段为1时，该字段必填报关类型为关联报关时必填。二线取消报关的情况下使用，用于生成区外一般贸易报关单。暂未使用)"`
	RltEntryBizopEtpsNm      string    `orm:"column(rlt_entry_bizop_etps_nm);size(512);null" description:"关联报关单境内收发货人名称(当报关类型DCLCUS_TYPECD字段为1时，该字段必填同上)"`
	RltEntryRvsngdEtpsSccd   string    `orm:"column(rlt_entry_rvsngd_etps_sccd);size(18);null" description:"关联报关单收发货单位社会统一信用代码(二线取消报关的情况下使用，用于生成区外一般贸易报关单。暂未使用)"`
	RltEntryRcvgdEtpsno      string    `orm:"column(rlt_entry_rcvgd_etpsno);size(10);null" description:"关联报关单海关收发货单位编码(当报关类型DCLCUS_TYPECD字段为1时，该字段必填报关类型为关联报关时必填。二线取消报关的情况下使用，用于生成区外一般贸易报关单。)"`
	RltEntryRcvgdEtpsNm      string    `orm:"column(rlt_entry_rcvgd_etps_nm);size(512);null" description:"关联报关单收发货单位名称(当报关类型DCLCUS_TYPECD字段为1时，该字段必填)"`
	RltEntryDclEtpsSccd      string    `orm:"column(rlt_entry_dcl_etps_sccd);size(18);null" description:"关联报关单申报单位社会统一信用代码(二线取消报关的情况下使用，用于生成区外一般贸易报关单。暂未使用)"`
	RltEntryDclEtpsno        string    `orm:"column(rlt_entry_dcl_etpsno);size(10);null" description:"关联报关单海关申报单位编码(当报关类型DCLCUS_TYPECD字段为1时，该字段必填报关类型为关联报关时必填。二线取消报关的情况下使用，用于生成区外一般贸易报关单。)"`
	RltEntryDclEtpsNm        string    `orm:"column(rlt_entry_dcl_etps_nm);size(512);null" description:"关联报关单申报单位名称(当报关类型DCLCUS_TYPECD字段为1时，该字段必填)"`
	Param1                   string    `orm:"column(param1);size(19);null" description:"备用1"`
	Param2                   string    `orm:"column(param2);size(19);null" description:"备用2"`
	Param3                   string    `orm:"column(param3);size(19);null" description:"备用3"`
	ExtraRemark              string    `orm:"column(extra_remark);null" description:"附注"`
	GenDecFlag               string    `orm:"column(gen_dec_flag);size(2);null" valid:"Required;MaxSize(2)" description:"是否生成报关单 1 生成，2 不生成"`
	GenDecFlagName           string    `orm:"column(gen_dec_flag_name);size(100);null" valid:"Required;MaxSize(100)" description:"是否生成报关单名称"`
	RecheckErrorInputIds     string    `orm:"column(recheck_error_input_ids);type(text);null" description:"复核input id"`
	ItemRecheckErrorInputIds string    `orm:"column(item_recheck_error_input_ids);type(text);null" description:"复核input id"`
	InputTime                time.Time `form:"-" orm:"column(input_time);type(datetime);null" valid:"Required"  description:"录入日期"`
	PrevdTime                time.Time `form:"-" orm:"column(prevd_time);type(datetime);null" description:"预核扣时间"`
	FormalVrfdedTime         time.Time `form:"-" orm:"column(formal_vrfded_time);type(datetime);null" description:"正式核扣时间 (返填)"`
	StatusUpdatedAt          time.Time `form:"-" orm:"column(status_updated_at);type(datetime)" description:"状态更新时间"`
	InvtDclTime              time.Time `form:"-" orm:"column(invt_dcl_time);type(datetime);null" valid:"Required" description:"清单申报时间(清单申报日期)(返填)"`
	EntryDclTime             time.Time `form:"-" orm:"column(entry_dcl_time);type(datetime);null"  description:"报关单申报日期(返填)清单报关时使用。海关端报关单入库时，反填并反馈企业端"`
	DeletedAt                time.Time `form:"-" orm:"column(deleted_at);type(timestamp);null" `

	BackendUsers []*BackendUser `orm:"rel(m2m);rel_through(BeeCustom/models.AnnotationUserRel)"` // 设置一对多的反向关系
	Company      *Company       `orm:"column(company_id);rel(fk)"`
	CompanyId    int64          `orm:"-" form:"CompanyId"` // 关联管理会自动生成 CompanyId 字段，此处不生成字段
	HandBookId   int64          `orm:"column(hand_book_id)" form:"HandBookId"`
	Order        *Order         `orm:"null;rel(one);on_delete(set_null)"` //

	AnnotationItems   []*AnnotationItem   `orm:"reverse(many)"` // 设置一对多关系
	AnnotationRecords []*AnnotationRecord `orm:"reverse(many)"` // 设置一对多关系

	// SysId                  string    `orm:"column(sys_id);size(2)" description:"子系统ID 95 加工贸易账册系统;B1 加工贸易手册系统 ;B2 加工贸易担保管理系统;B3 保税货物流转系统二期 ;Z7 海关特殊监管区域管理系统;Z8 保税物流管理系统"`
	// OperCusRegCode         string    `orm:"column(oper_cus_reg_code);size(10)" description:"操作卡的海关十位"`
	// KeyName                string    `orm:"column(key_name);size(255);null" description:"签名所用的证书信息"`
	// Version                string    `orm:"column(version);size(255);null" description:"版本编号"`
	// BusinessId             string    `orm:"column(business_id);size(255);null" description:"业务单证号"`
	// MessageId              string    `orm:"column(message_id);size(255);null" description:"报文唯一编号"`
	// FileName               string    `orm:"column(file_name);size(255);null" description:"用户原始报文名，主要用于用户查询"`
	// MessageType            string    `orm:"column(message_type);size(255);null" description:"报文类型"`
	// SenderId               string    `orm:"column(sender_id);size(255);null" description:"发送方编号"`
	// ReceiverId             string    `orm:"column(receiver_id);size(255);null" description:"接收方编号"`
	// DelcareFlag            string    `orm:"column(delcare_flag);size(255)" description:"申报标志 0--暂存；1--申报"`

}

func NewAnnotation(id int64) Annotation {
	return Annotation{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

// 查询参数
func NewAnnotationQueryParam() AnnotationQueryParam {
	return AnnotationQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

// AnnotationPageList 获取分页数据
func AnnotationStatusCount(params *AnnotationQueryParam) (orm.Params, error) {

	var maps []orm.Params
	rows := orm.Params{
		"审核通过":  0,
		"待制单":   0,
		"待复核":   0,
		"单一处理中": 0,
		"已完成":   0,
	}
	o := orm.NewOrm()

	sql := "SELECT "
	sql += "count( CASE WHEN STATUS = 3 THEN 1 END ) AS '审核通过',"
	sql += "count( CASE WHEN STATUS = 5 THEN 1 END ) AS '待制单',"
	sql += "count( CASE WHEN STATUS = 7 THEN 1 END ) AS '待复核',"
	sql += "count( CASE WHEN STATUS = 12 THEN 1 END ) AS '单一处理中',"
	sql += "count( CASE WHEN STATUS = 13 THEN 1 END ) AS '已完成' "
	sql = GetCommonListSql(sql, params)

	_, err := o.Raw(sql).Values(&maps)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("Raw:%v", err))
		return nil, err
	}

	if len(maps) > 0 {
		rows = maps[0]
	}

	return rows, nil

}

// AnnotationPageList 获取分页数据
func AnnotationPageList(params *AnnotationQueryParam) ([]*Annotation, int64, error) {

	datas := make([]*Annotation, 0)

	sql := "SELECT * "
	sql = GetCommonListSql(sql, params)
	if len(params.StatusString) > 0 && params.StatusString != "全部订单" {
		aStatusS, _ := GetSettingRValueByKey("annotationStatus", false)
		aStatus, _ := enums.TransformCnToInt(aStatusS, params.StatusString)
		sql += " AND status = " + strconv.Itoa(int(aStatus))
	}

	// 默认排序
	sortorder := "Id"
	if len(params.Sort) > 0 {
		sortorder = params.Sort
	}

	sql += " ORDER BY " + sortorder
	if params.Order == "desc" {
		sql += " DESC "
	} else {
		sql += " ASC "
	}

	o := orm.NewOrm()
	// 总数量
	total, err := o.Raw(sql).QueryRows(&datas)
	if err != nil {
		return nil, 0, err
	}

	if params.Limit != -1 {
		limit := strconv.Itoa(int(params.Limit))
		offset := strconv.Itoa(int((params.Offset - 1) * params.Limit))
		sql += " LIMIT " + offset + "," + limit
	}

	// 分页数据
	_, err = o.Raw(sql).QueryRows(&datas)
	if err != nil {
		return nil, 0, err
	}

	return datas, total, nil
}

func AnnotationGetRelations(ms []*Annotation, relations string) error {
	if len(relations) > 0 {
		o := orm.NewOrm()
		rs := strings.Split(relations, ",")
		for _, v := range ms {
			for _, rv := range rs {
				_, err := o.LoadRelated(v, rv)
				if err != nil {
					utils.LogDebug(fmt.Sprintf("LoadRelated:%v", err))
					return err
				}
			}
		}
	}
	return nil
}

// AnnotationOne 根据id获取单条
func AnnotationOne(id int64, relations string) (*Annotation, error) {
	m := NewAnnotation(0)
	o := orm.NewOrm()
	if err := o.QueryTable(AnnotationTBName()).Filter("Id", id).RelatedSel().One(&m); err != nil {
		utils.LogDebug(fmt.Sprintf("find annotation:%v", err))
		return nil, err
	}

	if len(relations) > 0 {
		rs := strings.Split(relations, ",")
		for _, rv := range rs {
			_, err := o.LoadRelated(&m, rv)
			if err != nil {
				utils.LogDebug(fmt.Sprintf("LoadRelated:%v", err))
				return nil, err
			}
		}

	}

	return &m, nil
}

// Annotations
func GetAnnotations(handBookId int64) ([]*Annotation, int64, error) {
	var ms []*Annotation
	o := orm.NewOrm()
	tatol, err := o.QueryTable(AnnotationTBName()).Filter("hand_book_id", handBookId).RelatedSel().All(&ms)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("find annotation:%v", err))
		return nil, 0, err
	}

	return ms, tatol, nil
}

// GetAnnotationByEtpsInnerInvtNo 根据清单号获取单条
func GetAnnotationByEtpsInnerInvtNo(etpsInnerInvtNo string) (*Annotation, error) {
	m := NewAnnotation(0)
	o := orm.NewOrm()
	if err := o.QueryTable(AnnotationTBName()).Filter("etps_inner_invt_no", etpsInnerInvtNo).One(&m); err != nil {
		// utils.LogDebug(fmt.Sprintf("find annotation:%v", err))
		return nil, err
	}

	return &m, nil
}

// GetAnnotationBySeqNo 根据清单预录入编号获取单条
func GetAnnotationBySeqNo(seqNo string) (*Annotation, error) {
	m := NewAnnotation(0)
	o := orm.NewOrm()
	if err := o.QueryTable(AnnotationTBName()).Filter("seq_no", seqNo).One(&m); err != nil {
		// utils.LogDebug(fmt.Sprintf("find annotation:%v", err))
		return nil, err
	}

	return &m, nil
}

// Save 添加、编辑页面 保存
func AnnotationUpdateOrSave(m *Annotation) error {
	var err error
	o := orm.NewOrm()
	if m.Id == 0 {
		if _, err = o.Insert(m); err != nil {
			utils.LogDebug(fmt.Sprintf("AnnotationSave:%v", err))
			return err
		}
	} else {
		_, err = o.Update(m)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("AnnotationSave:%v", err))
			return err
		}
	}

	return nil
}

// 保存
func AnnotationUpdate(m *Annotation, arg []string) error {
	var err error
	o := orm.NewOrm()

	_, err = o.Update(m, arg...)

	if err != nil {
		utils.LogDebug(fmt.Sprintf("AnnotationSave:%v", err))
		return err
	}

	return nil
}

// 删除
func AnnotationDelete(id int64) (num int64, err error) {
	m := NewAnnotation(id)
	if num, err := BaseDelete(&m); err != nil {
		return num, err
	} else {
		return num, nil
	}
}

// 列表公用sql
func GetCommonListSql(sql string, params *AnnotationQueryParam) string {
	sql += " FROM " + AnnotationTBName()
	sql += enums.GetOrderAnnotationDateTime(params.SearchTimeString, "invt_dcl_time")
	sql += " AND impexp_markcd = '" + params.ImpexpMarkcd + "'"
	if len(params.EtpsInnerInvtNoLike) > 0 {
		sql += " AND (etps_inner_invt_no LIKE '%" + params.EtpsInnerInvtNoLike + "%'" + " OR bond_invt_no LIKE '%" + params.EtpsInnerInvtNoLike + "%')"
	}
	if len(params.TrspModecd) > 0 {
		sql += " AND trsp_modecd = '" + params.TrspModecd + "'"
	}

	// 是否删除
	if params.IsDelete {
		sql += " AND  deleted_at  IS NOT NULL"
	} else {
		sql += " AND deleted_at IS NULL "
	}

	return sql
}

// TransformAnnotation 格式化数据
func TransformAnnotation(id int64, relation string) map[string]interface{} {

	v, _ := AnnotationOne(id, relation)
	annotationItem := make(map[string]interface{})
	aStatusS, err := GetSettingRValueByKey("orderStatus", false)
	aStatus, err, _ := enums.TransformIntToCn(aStatusS, v.Status)
	if err != nil {
		return nil
	}
	//转换表头复核标记
	recheckErrorInputIds := strings.Replace(strings.Replace(strings.Replace(v.RecheckErrorInputIds, `id":"`, "", -1), `[{"`, "", -1), `"}]`, "", -1)
	recheckErrorInputIdsSlice := strings.Split(recheckErrorInputIds, `"},{"`)

	//转换表体复核标记
	itemRecheckErrorInputIds := strings.Replace(strings.Replace(strings.Replace(v.ItemRecheckErrorInputIds, `index":`, "", -1), `[{"`, "", -1), `"]}]`, "", -1)
	itemRecheckErrorInputIdsSlice := strings.Split(itemRecheckErrorInputIds, `"]},{"`)
	var itemRecheckErrorInputIdsSlices []map[int][]string
	for _, v := range itemRecheckErrorInputIdsSlice {
		itemRecheckErrorInputIdsSlices1 := map[int][]string{}
		itemRecheckErrorInputIdsSlice1 := strings.Split(v, `,"id":["`)
		if len(itemRecheckErrorInputIdsSlice1) > 1 && len(itemRecheckErrorInputIdsSlice1[1]) > 0 {
			itemRecheckErrorInputIdsSlice1[1] = strings.Replace(itemRecheckErrorInputIdsSlice1[1], `"`, "", -1)
			itemRecheckErrorInputIdsSlice2 := strings.Split(itemRecheckErrorInputIdsSlice1[1], `,`)
			i, _ := strconv.ParseInt(itemRecheckErrorInputIdsSlice1[0], 10, 64)
			itemRecheckErrorInputIdsSlices1[int(i)] = itemRecheckErrorInputIdsSlice2
			itemRecheckErrorInputIdsSlices = append(itemRecheckErrorInputIdsSlices, itemRecheckErrorInputIdsSlices1)
		}
	}

	annotationItem["Id"] = strconv.FormatInt(v.Id, 10)
	annotationItem["StatusString"] = aStatus
	annotationItem["PutrecNo"] = v.PutrecNo
	annotationItem["ImpexpPortcd"] = v.ImpexpPortcd
	annotationItem["ImpexpPortcdName"] = v.ImpexpPortcdName
	annotationItem["BondInvtNo"] = v.BondInvtNo
	annotationItem["EntryNo"] = v.EntryNo
	annotationItem["EtpsInnerInvtNo"] = v.EtpsInnerInvtNo
	annotationItem["CompanyName"] = v.Company.Name
	annotationItem["SeqNo"] = v.SeqNo
	annotationItem["BizopEtpsSccd"] = v.BizopEtpsSccd
	annotationItem["BizopEtpsno"] = v.BizopEtpsno
	annotationItem["BizopEtpsNm"] = v.BizopEtpsNm
	annotationItem["RcvgdEtpsno"] = v.RcvgdEtpsno
	annotationItem["RvsngdEtpsSccd"] = v.RvsngdEtpsSccd
	annotationItem["RcvgdEtpsNm"] = v.RcvgdEtpsNm
	annotationItem["DclEtpsSccd"] = v.DclEtpsSccd
	annotationItem["DclEtpsno"] = v.DclEtpsno
	annotationItem["DclEtpsNm"] = v.DclEtpsNm
	annotationItem["DclPlcCuscd"] = v.DclPlcCuscd
	annotationItem["DclPlcCuscdName"] = v.DclPlcCuscdName
	annotationItem["ImpexpMarkcd"] = v.ImpexpMarkcd
	annotationItem["ImpexpMarkcdName"] = enums.GetImpexpMarkcdCNName(v.ImpexpMarkcd)
	annotationItem["RecheckErrorInputIds"] = recheckErrorInputIdsSlice
	annotationItem["ItemRecheckErrorInputIds"] = itemRecheckErrorInputIdsSlices
	annotationItem["MtpckEndprdMarkcd"] = v.MtpckEndprdMarkcd
	annotationItem["MtpckEndprdMarkcdName"] = v.MtpckEndprdMarkcdName
	annotationItem["SupvModecd"] = v.SupvModecd
	annotationItem["SupvModecdName"] = v.SupvModecdName
	annotationItem["TrspModecd"] = v.TrspModecd
	annotationItem["TrspModecdName"] = v.TrspModecdName
	annotationItem["DclcusFlag"] = v.DclcusFlag
	annotationItem["DclcusFlagName"] = v.DclcusFlagName
	annotationItem["DclcusTypecd"] = v.DclcusTypecd
	annotationItem["DclcusTypecdName"] = v.DclcusTypecdName
	annotationItem["VrfdedMarkcd"] = v.VrfdedMarkcd
	annotationItem["InvtIochkptStucd"] = v.InvtIochkptStucd
	annotationItem["ApplyNo"] = v.ApplyNo
	annotationItem["ListType"] = v.ListType
	annotationItem["ListTypeName"] = v.ListTypeName
	annotationItem["InputCode"] = v.InputCode
	annotationItem["InputCreditCode"] = v.InputCreditCode
	annotationItem["InputName"] = v.InputName
	annotationItem["ListStat"] = v.ListStat
	annotationItem["CorrEntryDclEtpsSccd"] = v.CorrEntryDclEtpsSccd
	annotationItem["CorrEntryDclEtpsNo"] = v.CorrEntryDclEtpsNo
	annotationItem["CorrEntryDclEtpsNm"] = v.CorrEntryDclEtpsNm
	annotationItem["DecType"] = v.DecType
	annotationItem["DecTypeName"] = v.DecTypeName
	annotationItem["StshipTrsarvNatcd"] = v.StshipTrsarvNatcd
	annotationItem["StshipTrsarvNatcdName"] = v.StshipTrsarvNatcdName
	annotationItem["InvtType"] = v.InvtType
	annotationItem["InvtTypeName"] = v.InvtTypeName
	annotationItem["EntryStucd"] = v.EntryStucd
	annotationItem["PassportUsedTypeCd"] = v.PassportUsedTypeCd
	annotationItem["Rmk"] = v.Rmk
	annotationItem["DecRmk"] = v.DecRmk
	annotationItem["DclTypecd"] = v.DclTypecd
	annotationItem["NeedEntryModified"] = v.NeedEntryModified
	annotationItem["LevyBlAmt"] = v.LevyBlAmt
	annotationItem["ChgTmsCnt"] = v.ChgTmsCnt
	annotationItem["RltInvtNo"] = v.RltInvtNo
	annotationItem["RltPutrecNo"] = v.RltPutrecNo
	annotationItem["RltEntryNo"] = v.RltEntryNo
	annotationItem["RltEntryBizopEtpsSccd"] = v.RltEntryBizopEtpsSccd
	annotationItem["RltEntryBizopEtpsno"] = v.RltEntryBizopEtpsno
	annotationItem["RltEntryBizopEtpsNm"] = v.RltEntryBizopEtpsNm
	annotationItem["RltEntryRvsngdEtpsSccd"] = v.RltEntryRvsngdEtpsSccd
	annotationItem["RltEntryRcvgdEtpsno"] = v.RltEntryRcvgdEtpsno
	annotationItem["RltEntryRcvgdEtpsNm"] = v.RltEntryRcvgdEtpsNm
	annotationItem["RltEntryDclEtpsSccd"] = v.RltEntryDclEtpsSccd
	annotationItem["RltEntryDclEtpsno"] = v.RltEntryDclEtpsno
	annotationItem["RltEntryDclEtpsNm"] = v.RltEntryDclEtpsNm
	annotationItem["Param1"] = v.Param1
	annotationItem["Param2"] = v.Param2
	annotationItem["Param3"] = v.Param3
	annotationItem["ExtraRemark"] = v.ExtraRemark
	annotationItem["GenDecFlag"] = v.GenDecFlag
	annotationItem["GenDecFlagName"] = v.GenDecFlagName
	annotationItem["HandBookId"] = strconv.FormatInt(v.HandBookId, 10)
	annotationItem["InputTime"] = enums.GetDateTimeString(&v.InputTime, enums.BaseDateFormat)
	annotationItem["PrevdTime"] = enums.GetDateTimeString(&v.PrevdTime, enums.BaseDateFormat)
	annotationItem["FormalVrfdedTime"] = enums.GetDateTimeString(&v.FormalVrfdedTime, enums.BaseDateFormat)
	annotationItem["EntryDclTime"] = enums.GetDateTimeString(&v.EntryDclTime, enums.BaseDateFormat)
	annotationItem["InvtDclTime"] = enums.GetDateTimeString(&v.InvtDclTime, enums.BaseDateFormat)
	annotationItem["AnnotationItems"] = v.AnnotationItems
	annotationItem["AnnotationRecords"] = v.AnnotationRecords

	return annotationItem
}
