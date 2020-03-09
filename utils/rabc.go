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
	dbType := beego.AppConfig.String("db_type")
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
		dns = fmt.Sprintf("%s.db", dbName)
		break
	case "mysql":
		dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPwd, dbHost, dbPort, dbName)
		break
	default:
		LogCritical(fmt.Sprintf("Database driver is not allowed:%v", dbType))
	}

	a := beegoormadapter.NewAdapter(dbType, dns)
	E = casbin.NewEnforcer("conf/rbac_model.conf", a)
	if err != nil {
		LogDebug(fmt.Sprintf("NewEnforcer error:%v", err))
	}

	err = E.LoadPolicy()

	if err != nil {
		LogDebug(fmt.Sprintf("LoadPolicy error:%v", err))
	}

}
