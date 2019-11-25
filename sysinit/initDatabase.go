package sysinit

import (
	"fmt"

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
	dbType := beego.AppConfig.DefaultString("db_type", "mysql")
	//连接名称
	dbAlias := beego.AppConfig.DefaultString(dbType+"::db_alias", "default")
	//数据库名称
	dbName := beego.AppConfig.DefaultString(dbType+"::db_name", "bee_custom")
	//数据库连接用户名
	dbUser := beego.AppConfig.DefaultString(dbType+"::db_user", "root")
	//数据库连接用户名
	dbPwd := beego.AppConfig.DefaultString(dbType+"::db_pwd", "")
	//数据库IP（域名）
	dbHost := beego.AppConfig.DefaultString(dbType+"::db_host", "127.0.0.1")
	//数据库端口
	dbPort := beego.AppConfig.DefaultString(dbType+"::db_port", "3306")
	switch dbType {
	case "sqlite3":
		_ = orm.RegisterDataBase(dbAlias, dbType, dbName)
		break
	case "mysql":
		dbCharset := beego.AppConfig.DefaultString(dbType+"::db_charset", "")
		dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPwd, dbHost, dbPort, dbName)
		_ = orm.RegisterDataBase(dbAlias, dbType, dns+"?charset="+dbCharset, 30)

	default:
		utils.LogCritical(fmt.Sprintf("Database driver is not allowed:%v", dbType))
	}

	//如果是开发模式，则显示命令信息
	isDev := beego.AppConfig.DefaultString("runmode", "dev") == "dev"
	//自动建表
	_ = orm.RunSyncdb("default", false, isDev)
	if isDev {
		orm.Debug = isDev
	}

	username := beego.AppConfig.DefaultString("pdf_username", "bee_custom_pdf")
	password := beego.AppConfig.DefaultString("pdf_password", "nvWQ8qE6kUtSURHhSQvWa2BZ3ct0eDOo")
	beego.InsertFilter("/pdf/*", beego.BeforeRouter, auth.Basic(username, password))
}
