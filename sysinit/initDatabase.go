package sysinit

import (
	_ "BeeCustom/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	//_ "github.com/mattn/go-sqlite3"
	_ "github.com/go-sql-driver/mysql"
)

//初始化数据连接
func InitDatabase() {
	//读取配置文件，设置数据库参数
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
	case "mysql":
		dbCharset := beego.AppConfig.DefaultString(dbType+"::db_charset", "")
		_ = orm.RegisterDataBase(dbAlias, dbType, dbUser+":"+dbPwd+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?charset="+dbCharset, 30)
	}
	//如果是开发模式，则显示命令信息
	isDev := beego.AppConfig.DefaultString("runmode", "dev") == "dev"
	//自动建表
	_ = orm.RunSyncdb("default", false, isDev)
	if isDev {
		orm.Debug = isDev
	}
}
