package models

import (
	"github.com/astaxie/beego/orm"
)

// TableName 设置Article表名
func (u *Article) TableName() string {
	return ArticleTBName()
}

// ArticleQueryParam 用于查询的类
type ArticleQueryParam struct {
	BaseQueryParam

	Type int8 //模糊查询
}

// Article 实体类
type Article struct {
	BaseModel

	Type     int8   `orm:"column(type);size(255)" `
	Title    string `orm:"column(title);size(255)"`
	Content  string `orm:"column(content);type(text);null" `
	Overview string `orm:"column(overview);size(255);null"`
	Origin   string `orm:"column(origin);size(255);null" `
	NewTime  string `orm:"column(newtime);size(255);null" `
}

func NewArticle(id int64) Article {
	return Article{BaseModel: BaseModel{Id: id}}
}

//查询参数
func NewArticleQueryParam() ArticleQueryParam {
	return ArticleQueryParam{BaseQueryParam: BaseQueryParam{Limit: -1, Sort: "Id", Order: "asc"}}
}

// ArticlePageList 获取分页数据
func ArticlePageList(params *ArticleQueryParam) ([]*Article, int64) {
	query := orm.NewOrm().QueryTable(ArticleTBName())
	datas := make([]*Article, 0)

	query = query.Filter("type", params.Type)

	total, _ := query.Count()
	query = BaseListQuery(query, params.Sort, params.Order, params.Limit, params.Offset)
	_, _ = query.All(&datas)

	return datas, total
}

// GetArticleByTitle 根据清单号获取单条
func GetArticleByTitle(title string) (*Article, error) {
	m := NewArticle(0)
	o := orm.NewOrm()
	if err := o.QueryTable(ArticleTBName()).Filter("title", title).One(&m); err != nil {
		return nil, err
	}

	return &m, nil
}

// 批量删除
func ArticleDeleteAll() (num int64, err error) {
	o := orm.NewOrm()
	if num, err := o.QueryTable(ArticleTBName()).Filter("code__isnull", false).Delete(); err != nil {
		return num, err
	} else {
		return num, nil
	}

}

// 批量插入
func InsertArticleMulti(datas []*Article) (num int64, err error) {
	return BaseInsertMulti(datas)
}
