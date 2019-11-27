package sysinit

import (
	"encoding/gob"

	"BeeCustom/models"
	"BeeCustom/tasks"
	"BeeCustom/utils"

	_ "github.com/astaxie/beego/session/redis"
)

func init() {

	//
	gob.Register(models.BackendUser{})

	//初始化日志
	utils.InitLogs()
	//初始化缓存
	utils.InitCache()
	//初始化权限
	utils.InitRabc()
	//自定义末模板方法
	utils.InitFunc()
	//定时任务
	tasks.InitTask()
	//ws
	utils.InitWs()
	//初始化数据库
	InitDatabase()

}
