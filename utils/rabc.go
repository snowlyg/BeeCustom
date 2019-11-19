package utils

import (
	"fmt"

	"BeeCustom/middleware/beego-orm-adapter"
	"github.com/astaxie/beego"
	"github.com/casbin/casbin"
)

var E *casbin.Enforcer
var dns string
var err error

//初始化数据连接
func InitRabc() {
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
		break
	case "mysql":
		dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPwd, dbHost, dbPort, dbName)
		break
	default:
		LogCritical(fmt.Sprintf("Database driver is not allowed:%v", dbType))
	}

	a := beegoormadapter.NewAdapter(dbType, dns)
	E, err = casbin.NewEnforcer("middleware/beego-orm-adapter/examples/rbac_model.conf", a)
	if err != nil {
		LogDebug(fmt.Sprintf("NewEnforcer error:%v", err))
	}

	err = E.LoadPolicy()

	if err != nil {
		LogDebug(fmt.Sprintf("LoadPolicy error:%v", err))
	}

}
