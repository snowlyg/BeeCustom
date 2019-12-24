package models

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"BeeCustom/enums"

	"BeeCustom/utils"
	"github.com/astaxie/beego/orm"
)

// TableName 设置Order表名
func (u *Order) TableName() string {
	return OrderTBName()
}

// OrderQueryParam 用于查询的类
type OrderQueryParam struct {
	BaseQueryParam

	IEFlag           string
	StatusString     string
	SearchTimeString string
	TrafMode         string
	ClientSeqNoLike  string
}

// Order 实体类
type Order struct {
	BaseModel

	Status                    int8      `orm:"column(status)" description:"状态"`
	IEFlag                    string    `orm:"column(i_e_flag);size(1)" description:"进出口标志"`
	AplStatus                 string    `orm:"column(apl_status);size(1);null" description:"申报状态 （表单不需填写）"`
	SeqNo                     string    `orm:"column(seq_no);size(20);null" description:"统一编号 （表单不需填写）"`
	ClientSeqNo               string    `orm:"column(client_seq_no);size(100);null" description:"订单号 （表单不需填写）"`
	PreEntryId                string    `orm:"column(pre_entry_id);size(20);null" description:"预录入编号（表单不需填写）"`
	EntryId                   string    `orm:"column(entry_id);size(20);null" description:"海关编号/报关单号 （表单不需填写）"`
	DocumentCodeString        string    `orm:"column(document_code_string);size(255);null" description:"随附单证 （表单不需填写）"`
	ContainerCounts           int       `orm:"column(container_counts);null" description:"集装箱数 （表单不需填写）"`
	ForeignCompanyName        string    `orm:"column(foreign_company_name);size(100);null" description:"外商公司名称"`
	Remark                    string    `orm:"column(remark);size(1000);null" description:"唛头及备注"`
	CustomMaster              string    `orm:"column(custom_master);size(4);null" description:"申报地海关"`
	CustomMasterName          string    `orm:"column(custom_master_name);size(50);null" description:"申报地海关名称"`
	IEPort                    string    `orm:"column(i_e_port);size(4);null" description:"进境关别"`
	IEPortName                string    `orm:"column(i_e_port_name);size(50);null" description:"进境关别名称"`
	ManualNo                  string    `orm:"column(manual_no);size(12);null" description:"备案号，账册号，手册号"`
	ContrNo                   string    `orm:"column(contr_no);size(32);null" description:"合同协议号"`
	IEDate                    string    `orm:"column(i_e_date);size(8);null" description:"进/出口日期"`
	TradeCoScc                string    `orm:"column(trade_co_scc);size(18);null" description:"境内收发货人社会信用代码"`
	TradeCode                 string    `orm:"column(trade_code);size(10);null" description:"境内收发货人海关代码"`
	TradeCiqCode              string    `orm:"column(trade_ciq_code);size(10);null" description:"境内收发货人检验检疫编码"`
	TradeName                 string    `orm:"column(trade_name);size(70);null" description:"境内收发货人企业名称（中文）"`
	OverseasConsignorCode     string    `orm:"column(overseas_consignor_code);size(50);null" description:"境外发货人代码"`
	OverseasConsignorCname    string    `orm:"column(overseas_consignor_cname);size(150);null" description:"境外发货人名称"`
	OverseasConsignorEname    string    `orm:"column(overseas_consignor_ename);size(100);null" description:"境外发货人名称（外文）"`
	OverseasConsignorAddr     string    `orm:"column(overseas_consignor_addr);size(100);null" description:"境外发货人地址"`
	OverseasConsigneeCode     string    `orm:"column(overseas_consignee_code);size(50);null" description:"境外收货人编码"`
	OverseasConsigneeEname    string    `orm:"column(overseas_consignee_ename);size(400);null" description:"境外收货人名称(外文)"`
	DomesticConsigneeEname    string    `orm:"column(domestic_consignee_ename);size(400);null" description:"境内收货人名称（外文）"`
	OwnerCodeScc              string    `orm:"column(owner_code_scc);size(18);null" description:"消费使用单位/生产销售单位社会信用代"`
	OwnerCode                 string    `orm:"column(owner_code);size(10);null" description:"消费使用单位/生产销售单位海关代码"`
	OwnerCiqCode              string    `orm:"column(owner_ciq_code);size(10);null" description:"消费使用单位/生产销售单位检验检疫编"`
	OwnerName                 string    `orm:"column(owner_name);size(70);null" description:"消费使用单位/生产销售单位企业名称 (中文)"`
	AgentCodeScc              string    `orm:"column(agent_code_scc);size(18);null" description:"申报单位社会信用代码 "`
	AgentCode                 string    `orm:"column(agent_code);size(10);null" description:"申报单位海关代码"`
	DeclCiqCode               string    `orm:"column(decl_ciq_code);size(10);null" description:"申报单位检验检疫编码"`
	AgentName                 string    `orm:"column(agent_name);size(70);null" description:"申报单位名称企业名称（中文）"`
	AgentTel                  string    `orm:"column(agent_tel);size(70);null" description:"申报单位电话"`
	TrafMode                  string    `orm:"column(traf_mode);size(10);null" description:"运输方式代码"`
	TrafModeName              string    `orm:"column(traf_mode_name);size(50);null" description:"运输方式名称"`
	TrafName                  string    `orm:"column(traf_name);size(200);null" description:"运输工具名称"`
	TrafNameCode              string    `orm:"column(traf_name_code);size(200);null" description:"运输工具代码 (没使用)"`
	BillNo                    string    `orm:"column(bill_no);size(50);null" description:"提运单号"`
	TradeMode                 string    `orm:"column(trade_mode);size(4);null" description:"监管方式"`
	TradeModeName             string    `orm:"column(trade_mode_name);size(50);null" description:"监管方式名称"`
	CutMode                   string    `orm:"column(cut_mode);size(3);null" description:"征免性质"`
	CutModeName               string    `orm:"column(cut_mode_name);size(50);null" description:"征免性质名称"`
	LicenseNo                 string    `orm:"column(license_no);size(20);null" description:"许可证号"`
	TradeCountry              string    `orm:"column(trade_country);size(3);null" description:"启运国（地区）/运抵国(地区)"`
	TradeCountryName          string    `orm:"column(trade_country_name);size(50);null" description:"启运国（地区）/运抵国(地区)名称"`
	DistinatePort             string    `orm:"column(distinate_port);size(6);null" description:"经停港/指运港"`
	DistinatePortName         string    `orm:"column(distinate_port_name);size(50);null" description:"经停港/指运港名称"`
	TransMode                 string    `orm:"column(trans_mode);size(1);null" description:"成交方式"`
	TransModeName             string    `orm:"column(trans_mode_name);size(50);null" description:"成交方式名称"`
	FeeMark                   string    `orm:"column(fee_mark);size(1);null" description:"运费标记"`
	FeeMarkName               string    `orm:"column(fee_mark_name);size(50);null" description:"运费标记名称"`
	FeeCurr                   string    `orm:"column(fee_curr);size(3);null" description:"运费币制"`
	FeeCurrName               string    `orm:"column(fee_curr_name);size(50);null" description:"运费币制名称"`
	FeeRate                   float64   `orm:"column(fee_rate);null;digits(12);decimals(4)" description:"运费／率"`
	InsurMark                 string    `orm:"column(insur_mark);size(1);null" description:"保险费标记"`
	InsurMarkName             string    `orm:"column(insur_mark_name);size(50);null" description:"保险费标记名称"`
	InsurCurr                 string    `orm:"column(insur_curr);size(3);null" description:"保险费币制"`
	InsurCurrName             string    `orm:"column(insur_curr_name);size(50);null" description:"保险费币制名称"`
	InsurRate                 float64   `orm:"column(insur_rate);null;digits(12);decimals(4)" description:"保险费／率"`
	OtherMark                 string    `orm:"column(other_mark);size(1);null" description:"杂费标志"`
	OtherMarkName             string    `orm:"column(other_mark_name);size(50);null" description:"杂费标志名称"`
	OtherCurr                 string    `orm:"column(other_curr);size(3);null" description:"杂费币制"`
	OtherCurrName             string    `orm:"column(other_curr_name);size(50);null" description:"杂费币制名称"`
	OtherRate                 float64   `orm:"column(other_rate);null;digits(12);decimals(4)" description:"杂费／率"`
	PackNo                    int       `orm:"column(pack_no);null" description:"件数"`
	WrapType                  string    `orm:"column(wrap_type);size(2);null" description:"包装种类"`
	WrapTypeName              string    `orm:"column(wrap_type_name);size(50);null" description:"包装种类名称"`
	GrossWet                  float64   `orm:"column(gross_wet);null;digits(17);decimals(5)" description:"毛重（KG）"`
	NetWt                     float64   `orm:"column(net_wt);null;digits(17);decimals(5)" description:"净重（KG）"`
	TradeAreaCode             string    `orm:"column(trade_area_code);size(3);null" description:"贸易国（地区）"`
	TradeAreaName             string    `orm:"column(trade_area_name);size(50);null" description:"贸易国（地区）名称"`
	EntyPortCode              string    `orm:"column(enty_port_code);size(6);null" description:"入境口岸/离境口岸"`
	EntyPortName              string    `orm:"column(enty_port_name);size(50);null" description:"入境口岸/离境口岸代码"`
	GoodsPlace                string    `orm:"column(goods_place);size(100);null" description:"货物存放地点（海关监管作业场所、分拨仓库、定点加工厂、隔离检疫场、企业自有仓库）"`
	DespPortCode              string    `orm:"column(desp_port_code);size(6);null" description:"启运港"`
	DespPortName              string    `orm:"column(desp_port_name);size(100);null" description:"启运港名称"`
	EntryType                 string    `orm:"column(entry_type);size(1);null" description:"报关单类型"`
	EntryTypeName             string    `orm:"column(entry_type_name);size(50);null" description:"报关单类型名称"`
	EdiId                     string    `orm:"column(edi_id);size(1);null" description:"报关标志"`
	Type                      string    `orm:"column(type);size(2);null" description:"单据类型(业务事项)"`
	NoteS                     string    `orm:"column(note_s);size(500);null" description:"备注"`
	PromiseItmes              string    `orm:"column(promise_itmes);size(50);null" description:"业务选项(其他事项确认) array [特殊关系确认，价格影响确认，支付特权使用费确认]"`
	MarkNo                    string    `orm:"column(mark_no);size(400);null" description:"标记唛码标记唛码"`
	BillType                  string    `orm:"column(bill_type);size(1);null" description:"备案清单类型,1:普通备案清单 2:先进区后报关 3:分送集报备案清单,4:分送集报报关单"`
	ChkSurety                 int8      `orm:"column(chk_surety);null" description:"担保验放标志,0:否；1:是"`
	CheckFlow                 int8      `orm:"column(check_flow);null" description:"查验分流,0:表示不是查验分流；1:表示是查验分流"`
	TaxAaminMark              int8      `orm:"column(tax_aamin_mark);null" description:"税收征管标记,0:无； 1:有"`
	OrgCode                   string    `orm:"column(org_code);size(10);null" description:"检验检疫受理机关"`
	OrgCodeName               string    `orm:"column(org_code_name);size(50);null" description:"检验检疫受理机关名称"`
	VsaOrgCode                string    `orm:"column(vsa_org_code);size(10);null" description:"领证机关"`
	VsaOrgCodeName            string    `orm:"column(vsa_org_code_name);size(50);null" description:"领证机关名称"`
	InspOrgCode               string    `orm:"column(insp_org_code);size(10);null" description:"口岸商检机关"`
	InspOrgName               string    `orm:"column(insp_org_name);size(50);null" description:"口岸商检机关名称"`
	PurpOrgCode               string    `orm:"column(purp_org_code);size(10);null" description:"目的地检验检疫机关"`
	PurpOrgName               string    `orm:"column(purp_org_name);size(50);null" description:"目的地检验检疫机关名称"`
	CorrelationNo             string    `orm:"column(correlation_no);size(500);null" description:"关联号码"`
	CorrelationReasonFlag     string    `orm:"column(correlation_reason_flag);size(2);null" description:"关联理由"`
	CorrelationReasonFlagName string    `orm:"column(correlation_reason_flag_name);size(50);null" description:"关联理由名称"`
	DecUsers                  string    `orm:"column(dec_users);type(text);null" description:"使用单位联系人 array [[use_org_person_code:使用单位联系人,use_org_person_tel:使用单位联系电话]]"`
	DecRequestCerts           string    `orm:"column(dec_request_certs);type(text);null" description:"报关单申请单证信息 （检验检疫申报要素） array [[app_cert_code:代码,app_cert_name:代码,appl_ori:正本数量 appl_copy_quan ：副本数量]]"`
	DecOtherPacks             string    `orm:"column(dec_other_packs);type(text);null" description:"报关单其他包装信息 array [[pack_qty=>包装件数(默认0,留空),pack_type=>包装材料种类]]"`
	SpecDeclFlag              string    `orm:"column(spec_decl_flag);type(text);null" description:"特殊业务标识 （用 ， 号分割的字符串）  0:未勾选；1:勾选]"`
	EntQualifTypes            string    `orm:"column(ent_qualif_types);type(text);null" description:"企业资质"`
	DeclaratioMaterialCode    string    `orm:"column(declaratio_material_code);size(10);null" description:"企业承诺信息 (证明/声明材料代码) :进口填写101040，出口填写102053"`
	RelId                     string    `orm:"column(rel_id);size(18);null" description:"关联报关单"`
	RelManNo                  string    `orm:"column(rel_man_no);size(12);null" description:"关联备案"`
	BonNo                     string    `orm:"column(bon_no);size(32);null" description:"保税/监管场地 (监管仓号)"`
	CusFie                    string    `orm:"column(cus_fie);size(8);null" description:"场地代码(货场代码）"`
	DecNo                     string    `orm:"column(dec_no);size(13);null" description:"报关员号"`
	DecBpNo                   string    `orm:"column(dec_bp_no);size(32);null" description:"报关员联系方式"`
	VoyNo                     string    `orm:"column(voy_no);size(32);null" description:"航次号"`
	DespDate                  string    `orm:"column(desp_date);size(8);null" description:"启运日期"`
	CmplDschrgDt              string    `orm:"column(cmpl_dschrg_dt);size(8);null" description:"卸毕日期"`
	BLNo                      string    `orm:"column(b_l_no);size(50);null" description:"BL/号"`
	OrigBoxFlag               string    `orm:"column(orig_box_flag);size(255);null" description:"原箱运送"`
	OperType                  string    `orm:"column(oper_type);size(1);null" description:"操作类型 G：报关单暂存（转关提前报关单暂存）"`
	ContactSafe               string    `orm:"column(contact_safe);size(255);null" description:"合同 保险"`
	CusFieName                string    `orm:"column(cus_fie_name);size(200);null" description:"货场名称"`
	IsOther                   int8      `orm:"column(is_other)" description:"是否异地报关"`
	IsSync                    int8      `orm:"column(is_sync)" description:"是否同步关务通"`
	RecheckErrorInputIds      string    `orm:"column(recheck_error_input_ids);type(text);null" description:"复核input id"`
	ItemRecheckErrorInputIds  string    `orm:"column(item_recheck_error_input_ids);type(text);null" description:"复核input id"`
	StatusUpdatedAt           time.Time `form:"-" orm:"column(status_updated_at);type(datetime)" description:"状态更新时间"`
	AplDate                   time.Time `form:"-" orm:"column(apl_date);type(datetime);null" description:"申报日期 （表单不需填写）"`
	ContactSignDate           time.Time `form:"-" orm:"column(contact_sign_date);type(datetime);null" description:"合同签约日期（进出口日期前一个月）"`
	DeletedAt                 time.Time `form:"-" orm:"column(deleted_at);type(timestamp);null"`

	BackendUsers    []*BackendUser    `orm:"rel(m2m);rel_through(BeeCustom/models.OrderUserRel)"` // 设置一对多的反向关系
	Company         *Company          `orm:"column(company_id);rel(fk)"`
	CompanyId       int64             `orm:"-" form:"CompanyId"` // 关联管理会自动生成 CompanyId 字段，此处不生成字段
	Annotation      *Annotation       `orm:"reverse(one)"`
	HandBookId      int64             `orm:"column(hand_book_id)" form:"HandBookId"`
	OrderItems      []*OrderItem      `orm:"reverse(many)"` // 设置一对多关系
	OrderContainers []*OrderContainer `orm:"reverse(many)"` // 设置一对多关系
	OrderDocuments  []*OrderDocument  `orm:"reverse(many)"` // 设置一对多关系
	OrderFiles      []*OrderFile      `orm:"reverse(many)"` // 设置一对多关系
	OrderRecords    []*OrderRecord    `orm:"reverse(many)"` // 设置一对多关系

}

func NewOrder(id int64) Order {
	return Order{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

// 查询参数
func NewOrderQueryParam() OrderQueryParam {
	return OrderQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

// OrderPageList 获取分页数据
func OrderStatusCount(params *OrderQueryParam) (orm.Params, error) {

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
	sql = GetOrderCommonListSql(sql, params)

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

// OrderPageList 获取分页数据
func OrderPageList(params *OrderQueryParam) ([]*Order, int64, error) {
	orders := make([]*Order, 0)
	sql := "SELECT * "
	sql = GetOrderCommonListSql(sql, params)
	if len(params.StatusString) > 0 && params.StatusString != "全部订单" {
		aStatusS, _ := GetSettingRValueByKey("orderStatus", false)
		aStatus, _, _ := enums.TransformCnToInt(aStatusS, params.StatusString)
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
	total, err := o.Raw(sql).QueryRows(&orders)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("QueryRows%v", err))
		return nil, 0, err
	}

	if params.Limit != -1 {
		limit := strconv.Itoa(int(params.Limit))
		offset := strconv.Itoa(int((params.Offset - 1) * params.Limit))
		sql += " LIMIT " + offset + "," + limit
	}

	// 分页数据
	_, err = o.Raw(sql).QueryRows(&orders)
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func OrderGetRelations(ms []*Order, relations string) error {
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

// OrderOne 根据id获取单条
func OrderOne(id int64, relations string) (*Order, error) {
	m := NewOrder(0)
	o := orm.NewOrm()
	if err := o.QueryTable(OrderTBName()).Filter("Id", id).RelatedSel().One(&m); err != nil {
		utils.LogDebug(fmt.Sprintf("find order:%v", err))
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

// GetOrderByClientSeqNo 根据清单号获取单条
func GetOrderByClientSeqNo(clientSeqNo string) (*Order, error) {
	m := NewOrder(0)
	o := orm.NewOrm()
	if err := o.QueryTable(OrderTBName()).Filter("client_seq_no", clientSeqNo).One(&m); err != nil {
		// utils.LogDebug(fmt.Sprintf("find order:%v", err))
		return nil, err
	}

	return &m, nil
}

// GetOrderBySeqNo 根据清单预录入编号获取单条
func GetOrderBySeqNo(seqNo string) (*Order, error) {
	m := NewOrder(0)
	o := orm.NewOrm()
	if err := o.QueryTable(OrderTBName()).Filter("seq_no", seqNo).One(&m); err != nil {
		// utils.LogDebug(fmt.Sprintf("find order:%v", err))
		return nil, err
	}

	return &m, nil
}

// Save 添加、编辑页面 保存
func OrderUpdateOrSave(m *Order) error {
	var err error
	o := orm.NewOrm()
	if m.Id == 0 {
		if _, err = o.Insert(m); err != nil {
			utils.LogDebug(fmt.Sprintf("OrderSave:%v", err))
			return err
		}
	} else {
		_, err = o.Update(m)

		if err != nil {
			utils.LogDebug(fmt.Sprintf("OrderSave:%v", err))
			return err
		}
	}

	return nil
}

// 保存附注
func OrderUpdate(m *Order, arg []string) error {
	var err error
	o := orm.NewOrm()

	_, err = o.Update(m, arg...)

	if err != nil {
		utils.LogDebug(fmt.Sprintf("OrderSave:%v", err))
		return err
	}

	return nil
}

// 删除
func OrderDelete(id int64) (num int64, err error) {
	m := NewOrder(id)
	if num, err := BaseDelete(&m); err != nil {
		return num, err
	} else {
		return num, nil
	}
}

// 列表公用sql
func GetOrderCommonListSql(sql string, params *OrderQueryParam) string {
	sql += " FROM " + OrderTBName()
	sql += enums.GetOrderAnnotationDateTime(params.SearchTimeString, "apl_date")
	sql += " AND i_e_flag = '" + params.IEFlag + "'"
	if len(params.ClientSeqNoLike) > 0 {
		sql += " AND client_seq_no LIKE '%" + params.ClientSeqNoLike + "%'"
		sql += " OR entry_id LIKE '%" + params.ClientSeqNoLike + "%'"
	}

	return sql
}

// TransformOrder 格式化列表数据
func TransformOrder(id int64, relation string, isRechek bool) map[string]interface{} {
	v, _ := OrderOne(id, relation)
	orderItem := make(map[string]interface{})
	aStatusS, err := GetSettingRValueByKey("orderStatus", false)
	aStatus, err, _ := enums.TransformIntToCn(aStatusS, v.Status)
	if err != nil {
		return nil
	}

	// 转换表头复核标记
	// recheckErrorInputIds := strings.Replace(strings.Replace(strings.Replace(v.RecheckErrorInputIds, `id":"`, "", -1), `[{"`, "", -1), `"}]`, "", -1)
	// recheckErrorInputIdsSlice := strings.Split(recheckErrorInputIds, `"},{"`)
	recheckErrorInputIdsSlice := enums.TramsformStringToSlice(v.RecheckErrorInputIds, `[{"`, `"},{"`, `id":"`, `"}]`)

	// 转换表体复核标记
	// itemRecheckErrorInputIds := strings.Replace(strings.Replace(strings.Replace(v.ItemRecheckErrorInputIds, `index":`, "", -1), `[{"`, "", -1), `"]}]`, "", -1)
	// itemRecheckErrorInputIdsSlice := strings.Split(itemRecheckErrorInputIds, `"]},{"`)
	itemRecheckErrorInputIdsSlice := enums.TramsformStringToSlice(v.ItemRecheckErrorInputIds, `[{"`, `"]},{"`, `index":`, `"]}]`)
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

	orderItem["Id"] = strconv.FormatInt(v.Id, 10)
	orderItem["StatusString"] = aStatus
	orderItem["IEFlag"] = v.IEFlag
	orderItem["IEFlagName"] = enums.GetImpexpMarkcdCNName(v.IEFlag)
	orderItem["RecheckErrorInputIds"] = recheckErrorInputIdsSlice
	orderItem["ItemRecheckErrorInputIds"] = itemRecheckErrorInputIdsSlices
	orderItem["HandBookId"] = strconv.FormatInt(v.HandBookId, 10)
	orderItem["SeqNo"] = v.SeqNo                                         // "seq_no);size(20);null" description:"统一编号 （表单不需填写）"`
	orderItem["ClientSeqNo"] = v.ClientSeqNo                             // "client_seq_no);size(100);null" description:"订单号 （表单不需填写）"`
	orderItem["PreEntryId"] = v.PreEntryId                               // "pre_entry_id);size(20);null" description:"预录入编号（表单不需填写）"`
	orderItem["EntryId"] = v.EntryId                                     // "entry_id);size(20);null" description:"海关编号/报关单号 （表单不需填写）"`
	orderItem["DocumentCodeString"] = v.DocumentCodeString               // "document_code_string);size(255);null" description:"随附单证 （表单不需填写）"`
	orderItem["ContainerCounts"] = enums.IsIZore(v.ContainerCounts)      // container_counts);null" description:"集装箱数 （表单不需填写）"`
	orderItem["ForeignCompanyName"] = v.ForeignCompanyName               // "foreign_company_name);size(100);null" description:"外商公司名称"`
	orderItem["Remark"] = v.Remark                                       // "remark);size(1000);null" description:"唛头及备注"`
	orderItem["CustomMaster"] = v.CustomMaster                           // "custom_master);size(4);null" description:"申报地海关"`
	orderItem["CustomMasterName"] = v.CustomMasterName                   // "custom_master_name);size(50);null" description:"申报地海关名称"`
	orderItem["IEPort"] = v.IEPort                                       // "i_e_port);size(4);null" description:"进境关别"`
	orderItem["IEPortName"] = v.IEPortName                               // "i_e_port_name);size(50);null" description:"进境关别名称"`
	orderItem["ManualNo"] = v.ManualNo                                   // "manual_no);size(12);null" description:"备案号，账册号，手册号"`
	orderItem["ContrNo"] = v.ContrNo                                     // "contr_no);size(32);null" description:"合同协议号"`
	orderItem["IEDate"] = v.IEDate                                       // "i_e_date);size(8);null" description:"进/出口日期"`
	orderItem["TradeCoScc"] = v.TradeCoScc                               // "trade_co_scc);size(18);null" description:"境内收发货人社会信用代码"`
	orderItem["TradeCode"] = v.TradeCode                                 // "trade_code);size(10);null" description:"境内收发货人海关代码"`
	orderItem["TradeCiqCode"] = v.TradeCiqCode                           // "trade_ciq_code);size(10);null" description:"境内收发货人检验检疫编码"`
	orderItem["TradeName"] = v.TradeName                                 // "trade_name);size(70);null" description:"境内收发货人企业名称（中文）"`
	orderItem["OverseasConsignorCode"] = v.OverseasConsignorCode         // "overseas_consignor_code);size(50);null" description:"境外发货人代码"`
	orderItem["OverseasConsignorCname"] = v.OverseasConsignorCname       // "overseas_consignor_cname);size(150);null" description:"境外发货人名称"`
	orderItem["OverseasConsignorEname"] = v.OverseasConsignorEname       // "overseas_consignor_ename);size(100);null" description:"境外发货人名称（外文）"`
	orderItem["OverseasConsignorAddr"] = v.OverseasConsignorAddr         // "overseas_consignor_addr);size(100);null" description:"境外发货人地址"`
	orderItem["OverseasConsigneeCode"] = v.OverseasConsigneeCode         // "overseas_consignee_code);size(50);null" description:"境外收货人编码"`
	orderItem["OverseasConsigneeEname"] = v.OverseasConsigneeEname       // "overseas_consignee_ename);size(400);null" description:"境外收货人名称(外文)"`
	orderItem["DomesticConsigneeEname"] = v.DomesticConsigneeEname       // "domestic_consignee_ename);size(400);null" description:"境内收货人名称（外文）"`
	orderItem["OwnerCodeScc"] = v.OwnerCodeScc                           // "owner_code_scc);size(18);null" description:"消费使用单位/生产销售单位社会信用代"`
	orderItem["OwnerCode"] = v.OwnerCode                                 // "owner_code);size(10);null" description:"消费使用单位/生产销售单位海关代码"`
	orderItem["OwnerCiqCode"] = v.OwnerCiqCode                           // "owner_ciq_code);size(10);null" description:"消费使用单位/生产销售单位检验检疫编"`
	orderItem["OwnerName"] = v.OwnerName                                 // "owner_name);size(70);null" description:"消费使用单位/生产销售单位企业名称 (中文)"`
	orderItem["AgentCodeSccAgentCodeScc"] = v.AgentCodeScc               // "agent_code_scc);size(18);null" description:"申报单位社会信用代码 "`
	orderItem["AgentCode"] = v.AgentCode                                 // "agent_code);size(10);null" description:"申报单位海关代码"`
	orderItem["DeclCiqCode"] = v.DeclCiqCode                             // "decl_ciq_code);size(10);null" description:"申报单位检验检疫编码"`
	orderItem["AgentName"] = v.AgentName                                 // "agent_name);size(70);null" description:"申报单位名称企业名称（中文）"`
	orderItem["AgentTel"] = v.AgentTel                                   // "agent_tel);size(70);null" description:"申报单位电话"`
	orderItem["TrafMode"] = v.TrafMode                                   // "traf_mode);size(10);null" description:"运输方式代码"`
	orderItem["TrafModeName"] = v.TrafModeName                           // "traf_mode_name);size(50);null" description:"运输方式名称"`
	orderItem["TrafName"] = v.TrafName                                   // "traf_name);size(200);null" description:"运输工具名称"`
	orderItem["TrafNameCode"] = v.TrafNameCode                           // "traf_name_code);size(200);null" description:"运输工具代码 (没使用)"`
	orderItem["BillNo"] = v.BillNo                                       // "bill_no);size(50);null" description:"提运单号"`
	orderItem["TradeMode"] = v.TradeMode                                 // "trade_mode);size(4);null" description:"监管方式"`
	orderItem["TradeModeName"] = v.TradeModeName                         // "trade_mode_name);size(50);null" description:"监管方式名称"`
	orderItem["CutMode"] = v.CutMode                                     // "cut_mode);size(3);null" description:"征免性质"`
	orderItem["CutModeName"] = v.CutModeName                             // "cut_mode_name);size(50);null" description:"征免性质名称"`
	orderItem["LicenseNo"] = v.LicenseNo                                 // "license_no);size(20);null" description:"许可证号"`
	orderItem["TradeCountry"] = v.TradeCountry                           // "trade_country);size(3);null" description:"启运国（地区）/运抵国(地区)"`
	orderItem["TradeCountryName"] = v.TradeCountryName                   // "trade_country_name);size(50);null" description:"启运国（地区）/运抵国(地区)名称"`
	orderItem["DistinatePort "] = v.DistinatePort                        // "distinate_port);size(6);null" description:"经停港/指运港"`
	orderItem["DistinatePortName"] = v.DistinatePortName                 // "distinate_port_name);size(50);null" description:"经停港/指运港名称"`
	orderItem["TransMode"] = v.TransMode                                 // "trans_mode);size(1);null" description:"成交方式"`
	orderItem["TransModeName"] = v.TransModeName                         // "trans_mode_name);size(50);null" description:"成交方式名称"`
	orderItem["FeeMark"] = v.FeeMark                                     // "fee_mark);size(1);null" description:"运费标记"`
	orderItem["FeeMarkName"] = v.FeeMarkName                             // "fee_mark_name);size(50);null" description:"运费标记名称"`
	orderItem["FeeCurr"] = v.FeeCurr                                     // "fee_curr);size(3);null" description:"运费币制"`
	orderItem["FeeCurrName"] = v.FeeCurrName                             // "fee_curr_name);size(50);null" description:"运费币制名称"`
	orderItem["FeeRate"] = enums.IsFloatZore(v.FeeRate)                  // fee_rate);null;digits(12);decimals(4)" description:"运费／率"`
	orderItem["InsurMark"] = v.InsurMark                                 // "insur_mark);size(1);null" description:"保险费标记"`
	orderItem["InsurMarkName"] = v.InsurMarkName                         // "insur_mark_name);size(50);null" description:"保险费标记名称"`
	orderItem["InsurCurr"] = v.InsurCurr                                 // "insur_curr);size(3);null" description:"保险费币制"`
	orderItem["InsurCurrName"] = v.InsurCurrName                         // "insur_curr_name);size(50);null" description:"保险费币制名称"`
	orderItem["InsurRate"] = enums.IsFloatZore(v.InsurRate)              // insur_rate);null;digits(12);decimals(4)" description:"保险费／率"`
	orderItem["OtherMark"] = v.OtherMark                                 // "other_mark);size(1);null" description:"杂费标志"`
	orderItem["OtherMarkName"] = v.OtherMarkName                         // "other_mark_name);size(50);null" description:"杂费标志名称"`
	orderItem["OtherCurr"] = v.OtherCurr                                 // "other_curr);size(3);null" description:"杂费币制"`
	orderItem["OtherCurrName"] = v.OtherCurrName                         // "other_curr_name);size(50);null" description:"杂费币制名称"`
	orderItem["OtherRate"] = enums.IsFloatZore(v.OtherRate)              // `orm:"column(other_rate);null;digits(12);decimals(4)" description:"杂费／率"`
	orderItem["PackNo"] = enums.IsIZore(v.PackNo)                        // pack_no);null" description:"件数"`
	orderItem["WrapType"] = v.WrapType                                   // "wrap_type);size(2);null" description:"包装种类"`
	orderItem["WrapTypeName"] = v.WrapTypeName                           // "wrap_type_name);size(50);null" description:"包装种类名称"`
	orderItem["GrossWet"] = enums.IsFloatZore(v.GrossWet)                // gross_wet);null;digits(17);decimals(5)" description:"毛重（KG）"`
	orderItem["NetWt"] = enums.IsFloatZore(v.NetWt)                      // net_wt);null;digits(17);decimals(5)" description:"净重（KG）"`
	orderItem["TradeAreaCode"] = v.TradeAreaCode                         // "trade_area_code);size(3);null" description:"贸易国（地区）"`
	orderItem["TradeAreaName"] = v.TradeAreaName                         // "trade_area_name);size(50);null" description:"贸易国（地区）名称"`
	orderItem["EntyPortCode"] = v.EntyPortCode                           // "enty_port_code);size(6);null" description:"入境口岸/离境口岸"`
	orderItem["EntyPortName"] = v.EntyPortName                           // "enty_port_name);size(50);null" description:"入境口岸/离境口岸代码"`
	orderItem["GoodsPlace"] = v.GoodsPlace                               // "goods_place);size(100);null" description:"货物存放地点（海关监管作业场所、分拨仓库、定点加工厂、隔离检疫场、企业自有仓库）"`
	orderItem["DespPortCode"] = v.DespPortCode                           // "desp_port_code);size(6);null" description:"启运港"`
	orderItem["DespPortName"] = v.DespPortName                           // "desp_port_name);size(100);null" description:"启运港名称"`
	orderItem["EntryType "] = v.EntryType                                // "entry_type);size(1);null" description:"报关单类型"`
	orderItem["EntryTypeName"] = v.EntryTypeName                         // "entry_type_name);size(50);null" description:"报关单类型名称"`
	orderItem["EdiId"] = v.EdiId                                         // "edi_id);size(1);null" description:"报关标志"`
	orderItem["Type"] = v.Type                                           // "type);(2);null" description:"单据类型(业务事项)"`
	orderItem["NoteS"] = v.NoteS                                         // "note_s);size(500);null" description:"备注"`
	orderItem["PromiseItmes"] = v.PromiseItmes                           // "promise_itmes);size(50);null" description:"业务选项(其他事项确认) array [特殊关系确认，价格影响确认，支付特权使用费确认]"`
	orderItem["MarkNo"] = v.MarkNo                                       // "mark_no);size(400);null" description:"标记唛码标记唛码"`
	orderItem["BillType"] = v.BillType                                   // "bill_type);size(1);null" description:"备案清单类型,1:普通备案清单 2:先进区后报关 3:分送集报备案清单,4:分送集报报关单"`
	orderItem["ChkSurety"] = v.ChkSurety                                 // chk_surety);null" description:"担保验放标志,0:否；1:是"`
	orderItem["CheckFlow"] = v.CheckFlow                                 // check_flow);null" description:"查验分流,0:表示不是查验分流；1:表示是查验分流"`
	orderItem["TaxAaminMark"] = v.TaxAaminMark                           // tax_aamin_mark);null" description:"税收征管标记,0:无； 1:有"`
	orderItem["OrgCode"] = v.OrgCode                                     // "org_code);size(10);null" description:"检验检疫受理机关"`
	orderItem["OrgCodeName"] = v.OrgCodeName                             // "org_code_name);size(50);null" description:"检验检疫受理机关名称"`
	orderItem["VsaOrgCode"] = v.VsaOrgCode                               // "vsa_org_code);size(10);null" description:"领证机关"`
	orderItem["VsaOrgCodeName"] = v.VsaOrgCodeName                       // "vsa_org_code_name);size(50);null" description:"领证机关名称"`
	orderItem["InspOrgCode"] = v.InspOrgCode                             // "insp_org_code);size(10);null" description:"口岸商检机关"`
	orderItem["InspOrgName"] = v.InspOrgName                             // "insp_org_name);size(50);null" description:"口岸商检机关名称"`
	orderItem["PurpOrgCode"] = v.PurpOrgCode                             // "purp_org_code);size(10);null" description:"目的地检验检疫机关"`
	orderItem["PurpOrgName"] = v.PurpOrgName                             // "purp_org_name);size(50);null" description:"目的地检验检疫机关名称"`
	orderItem["CorrelationNo "] = v.CorrelationNo                        // "correlation_no);size(500);null" description:"关联号码"`
	orderItem["CorrelationReasonFlag"] = v.CorrelationReasonFlag         // "correlation_reason_flag);size(2);null" description:"关联理由"`
	orderItem["CorrelationReasonFlagName"] = v.CorrelationReasonFlagName // "correlation_reason_flag_name);size(50);null" description:"关联理由名称"`
	orderItem["DecUsers"] = v.DecUsers                                   // "dec_users);null" description:"使用单位联系人 array [[use_org_person_code:使用单位联系人,use_org_person_tel:使用单位联系电话]]"`
	orderItem["EntQualifTypes"] = v.EntQualifTypes                       // "dec_request_certs);null" description:"报关单申请单证信息 （检验检疫申报要素） array [[app_cert_code:代码,app_cert_name:代码,appl_ori:正本数量 appl_copy_quan ：副本数量]]"`
	orderItem["DecRequestCerts"] = v.DecRequestCerts                     // "dec_request_certs);null" description:"报关单申请单证信息 （检验检疫申报要素） array [[app_cert_code:代码,app_cert_name:代码,appl_ori:正本数量 appl_copy_quan ：副本数量]]"`
	orderItem["DecOtherPacks"] = v.DecOtherPacks                         // "dec_other_packs);null" description:"报关单其他包装信息 array [[pack_qty=>包装件数(默认0,留空),pack_type=>包装材料种类]]"`
	orderItem["SpecDeclFlag"] = v.SpecDeclFlag                           // "spec_decl_flag);size(100);null" description:"特殊业务标识 （用 ， 号分割的字符串）  0:未勾选；1:勾选]"`
	orderItem["DeclaratioMaterialCode"] = v.DeclaratioMaterialCode       // "declaratio_material_code);size(10);null" description:"企业承诺信息 (证明/声明材料代码) :进口填写101040，出口填写102053"`
	orderItem["RelId"] = v.RelId                                         // "rel_id);size(18);null" description:"关联报关单"`
	orderItem["RelManNo"] = v.RelManNo                                   // "rel_man_no);size(12);null" description:"关联备案"`
	orderItem["BonNo"] = v.BonNo                                         // "bon_no);size(32);null" description:"保税/监管场地 (监管仓号)"`
	orderItem["CusFie"] = v.CusFie                                       // "cus_fie);size(8);null" description:"场地代码(货场代码）"`
	orderItem["DecNo"] = v.DecNo                                         // "dec_no);size(13);null" description:"报关员号"`
	orderItem["DecBpNo"] = v.DecBpNo                                     // "dec_bp_no);size(32);null" description:"报关员联系方式"`
	orderItem["VoyNo"] = v.VoyNo                                         // "voy_no);size(32);null" description:"航次号"`
	orderItem["DespDate"] = v.DespDate                                   // "desp_date);size(8);null" description:"启运日期"`
	orderItem["CmplDschrgDt"] = v.CmplDschrgDt                           // "cmpl_dschrg_dt);size(8);null" description:"卸毕日期"`
	orderItem["BLNo"] = v.BLNo                                           // "b_l_no);size(50);null" description:"BL/号"`
	orderItem["OrigBoxFlag"] = v.OrigBoxFlag                             // "orig_box_flag);size(255);null" description:"原箱运送"`
	orderItem["OperType"] = v.OperType                                   // "oper_type);size(1);null" description:"操作类型 G：报关单暂存（转关提前报关单暂存）"`
	orderItem["ContactSafe"] = v.ContactSafe                             // "contact_safe);size(255);null" description:"合同 保险"`
	orderItem["CusFieName"] = v.CusFieName                               // "cus_fie_name);size(200);null" description:"货场名称"`
	orderItem["IsOther"] = v.IsOther                                     // is_other)" description:"是否异地报关"`
	orderItem["IsSync"] = v.IsSync

	if isRechek {
		var newItem []OrderItem
		var otherOrderItem [][]OrderItem
		var orderContatinersData string
		var orderDocumentsData string
		items, _ := OrderItemGetRelations(v.OrderItems, "OrderItemLimits")
		for ii, v := range items {
			orderItemLimits, _ := OrderItemLimitGetRelations(v.OrderItemLimits, "OrderItemLimitVins")
			v.OrderItemLimits = orderItemLimits

			if ii < 6 {
				newItem = append(newItem, *v)
			}

		}
		itemLn := (float64(len(items)) - float64(6)) / float64(14)
		i2 := int(math.Ceil(itemLn))
		for i := 0; i < i2; i++ {
			var newOtherOrderItem []OrderItem
			for ii, v := range items {
				orderItemLimits, _ := OrderItemLimitGetRelations(v.OrderItemLimits, "OrderItemLimitVins")
				v.OrderItemLimits = orderItemLimits

				if ii >= i*14+6 && ii < (i+1)*14+6 {
					newOtherOrderItem = append(newOtherOrderItem, *v)
				}
			}
			otherOrderItem = append(otherOrderItem, newOtherOrderItem)
		}

		orderItem["OrderItems"] = newItem
		orderItem["OtherOrderItems"] = otherOrderItem

		if len(v.OrderDocuments) == 0 {
			orderDocumentsData = "随附单证0:代理报关委托协议（电子）;"
		} else {
			orderDocumentsData = "随附单证" + strconv.Itoa(len(v.OrderDocuments)) + ":代理报关委托协议（电子）;"
			for _, v := range v.OrderDocuments {
				orderDocumentsData += v.DocuCodeName + ";" + v.CertCode + ";"
			}
		}
		orderItem["OrderDocumentsData"] = orderDocumentsData

		if len(v.OrderContainers) == 0 {
			orderContatinersData = "0;无"
		} else {
			orderContatinersData = strconv.Itoa(len(v.OrderContainers)) + ";"
			for _, v := range v.OrderContainers {
				orderContatinersData += v.ContainerId + ";"
			}
		}
		orderItem["OrderContatinersData"] = orderContatinersData

	} else {
		orderItem["OrderContainers"] = v.OrderContainers
		orderItem["OrderDocuments"] = v.OrderDocuments
		items, _ := OrderItemGetRelations(v.OrderItems, "OrderItemLimits")
		for _, v := range items {
			orderItemLimits, _ := OrderItemLimitGetRelations(v.OrderItemLimits, "OrderItemLimitVins")
			v.OrderItemLimits = orderItemLimits
		}
		orderItem["OrderItems"] = items
	}
	orderItem["ContactSignDate"] = enums.GetDateTimeString(&v.ContactSignDate, enums.BaseDateFormat) // contact_sign_date);type(datetime);null" description:"合同签约日期（进出口日期前一个月）"`
	orderItem["AplDate"] = enums.GetDateTimeString(&v.AplDate, enums.BaseDateFormat)
	return orderItem
}
