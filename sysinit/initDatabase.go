package sysinit

import (
	"fmt"

	"BeeCustom/models"
	"BeeCustom/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/plugins/auth"

	//_ "github.com/mattn/go-sqlite3"
	_ "github.com/go-sql-driver/mysql"
)

//初始化数据连接
func InitDatabase() {

	var dns string

	//数据库类别
	dbType := beego.AppConfig.String("db_type")
	//连接名称
	dbAlias := beego.AppConfig.String(dbType + "::db_alias")
	//数据库名称
	dbName := beego.AppConfig.String(dbType + "::db_name")
	//数据库连接用户名
	dbUser := beego.AppConfig.String(dbType + "::db_user")
	//数据库连接用户名
	dbPwd := beego.AppConfig.String(dbType + "::db_pwd")
	//数据库IP（域名）
	dbHost := beego.AppConfig.String(dbType + "::db_host")
	//数据库端口
	dbPort := beego.AppConfig.String(dbType + "::db_port")
	switch dbType {
	case "sqlite3":
		_ = orm.RegisterDataBase(dbAlias, dbType, dbName)
		break
	case "mysql":
		dbCharset := beego.AppConfig.String(dbType + "::db_charset")
		dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPwd, dbHost, dbPort, dbName)
		_ = orm.RegisterDataBase(dbAlias, dbType, dns+"?charset="+dbCharset, 30)

	default:
		utils.LogCritical(fmt.Sprintf("Database driver is not allowed:%v", dbType))
	}

	//如果是开发模式，则显示命令信息
	isDev := beego.AppConfig.String("runmode") == "dev"
	//自动建表
	_ = orm.RunSyncdb("default", false, isDev)
	if isDev {
		orm.Debug = isDev
	}

	//basicAuth 认证
	username, _ := models.GetSettingValueByKey("pdf_username")
	password, _ := models.GetSettingValueByKey("pdf_password")
	beego.InsertFilter("/pdf/*", beego.BeforeRouter, auth.Basic(username, password))
}
