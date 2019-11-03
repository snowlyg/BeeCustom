package models

import (
	"BeeCustom/utils"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

// TableName 设置HandBookUllage表名
func (u *HandBookUllage) TableName() string {
	return HandBookUllageTBName()
}

// HandBookUllage 实体类
type HandBookUllage struct {
	BaseModel

	OriginalityProNo      int8          `orm:"column(originality_pro_no)" description:"料件项号"`
	OriginalityProName    string        `orm:"column(originality_pro_name);size(50)" description:"料件名称"`
	OriginalityProSpecial string        `orm:"column(originality_pro_special);size(1000);null" description:"料件规格"`
	OriginalityProUCode   string        `orm:"column(originality_pro_u_code);size(255);null" description:"料件单位代码"`
	OriginalityProU       string        `orm:"column(originality_pro_u);size(255);null" description:"料件单位"`
	OnlyUllage            string        `orm:"column(only_ullage);size(255);null" description:"单耗"`
	Ullage                string        `orm:"column(ullage);size(255);null" description:"损耗"`
	Gedition              string        `orm:"column(gedition);size(100);null" description:"成品版本号"`
	Serial                string        `orm:"column(serial);size(255);null" description:"序号"`
	FinishProNo           string        `orm:"column(finish_pro_no);size(255);null" description:"成品序号"`
	FinishRecordNo        string        `orm:"column(finish_record_no);size(255);null" description:"成品料号"`
	FinishHsCode          string        `orm:"column(finish_hs_code);size(255);null" description:"成品商品编码"`
	FinishName            string        `orm:"column(finish_name);size(50);null" description:"成品名称"`
	FinishSpecial         string        `orm:"column(finish_special);size(1000);null" description:"成品规格"`
	FinishSpecialUCode    string        `orm:"column(finish_special_u_code);size(255);null" description:"成品计量单位代码"`
	FinishSpecialU        string        `orm:"column(finish_special_u);size(200);null" description:"成品计量单位"`
	OriginalityRecordNo   string        `orm:"column(originality_record_no);size(255);null" description:"料件料号"`
	OriginalityHsCode     string        `orm:"column(originality_hs_code);size(50);null" description:"料件商品编码"`
	OnlyUllageVersion     string        `orm:"column(only_ullage_version);size(255);null" description:"单耗版本号"`
	OneUllage             string        `orm:"column(one_ullage);size(255);null" description:"净耗"`
	NoUllage              string        `orm:"column(no_ullage);size(255);null" description:"无形损耗率"`
	OnlyUllageStatus      string        `orm:"column(only_ullage_status);size(255);null" description:"单耗申报状态"`
	ChangeMark            string        `orm:"column(change_mark);size(255);null" description:"处理标志(修改标志)"`
	BondedRate            string        `orm:"column(bonded_rate);size(255);null" description:"保税料件比例%"`
	CompanyExecuteFlag    string        `orm:"column(company_execute_flag);size(255);null" description:"企业执行标志"`
	OnlyUllageAt          time.Time     `orm:"column(only_ullage_at);type(datetime);null" description:"单耗有效期"`
	UllageFlag            time.Time     `orm:"column(ullage_flag);type(datetime);null" description:"单耗质疑标志"`
	TalkFlag              time.Time     `orm:"column(talk_flag);type(datetime);null" description:"磋商标志"`
	Remark                string        `orm:"column(remark);size(255);null" description:"备注"`
	HandBookGood          *HandBookGood `orm:"column(hand_book_good_id);rel(fk)"`
	HandBookGoodId        int64         `orm:"-" form:"HandBookGoodId"`
}

// HandBookUllageQueryParam 用于查询的类
type HandBookUllageQueryParam struct {
	BaseQueryParam
	Type     string //模糊查询
	NameLike string //模糊查询
}

func NewHandBookUllage(id int64) HandBookUllage {
	return HandBookUllage{BaseModel: BaseModel{id, time.Now(), time.Now()}}
}

//查询参数
func NewHandBookUllageQueryParam() HandBookUllageQueryParam {
	return HandBookUllageQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

func HandBookUllageGetRelations(v *HandBookUllage, relations string) (*HandBookUllage, error) {
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

// HandBookUllageOne 根据id获取单条
func HandBookUllageOne(id int64, relations string) (*HandBookUllage, error) {
	m := NewHandBookUllage(0)
	o := orm.NewOrm()
	if err := o.QueryTable(HandBookUllageTBName()).Filter("Id", id).RelatedSel().One(&m); err != nil {
		return nil, err
	}

	if len(relations) > 0 {
		_, err := HandBookUllageGetRelations(&m, relations)
		if err != nil {
			return nil, err
		}
	}

	if &m == nil {
		return &m, errors.New("获取失败")
	}

	return &m, nil
}

//删除
func HandBookUllageDelete(id int64) (num int64, err error) {
	m := NewHandBookUllage(id)
	if num, err := BaseDelete(&m); err != nil {
		return num, err
	} else {
		return num, nil
	}
}

//删除
func HandBookUllageDeleteAll(clearanceType int8) (num int64, err error) {
	if num, err := BaseDeleteAll(clearanceType); err != nil {
		return num, err
	} else {
		return num, nil
	}
}

//批量插入
func InsertHandBookUllageMulti(datas []*HandBookUllage) (num int64, err error) {
	return BaseInsertMulti(len(datas), datas)
}
