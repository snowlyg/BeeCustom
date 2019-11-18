package utils

import (
	"fmt"

	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/cache/redis"
	beegoormadapter "github.com/casbin/beego-orm-adapter"
	"github.com/casbin/casbin"
)

var E *casbin.Enforcer

//初始化数据连接
func InitRabc() {

	var dns string

	//数据库类别
	dbType := beego.AppConfig.DefaultString("db_type", "mysql")
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

		dns = fmt.Sprintf("%s.db", dbName)
		// 注册casbin
		a := beegoormadapter.NewAdapter("sqlite3", dns, true)
		E, _ = casbin.NewEnforcer("conf/rbac_model.conf", a)
		break
	case "mysql":

		dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPwd, dbHost, dbPort, dbName)

		a := beegoormadapter.NewAdapter("mysql", dns)
		E, _ = casbin.NewEnforcer("conf/rbac_model.conf", a)

	default:
		LogCritical(fmt.Sprintf("Database driver is not allowed:%v", dbType))
	}

	_ = E.LoadPolicy()

}
