package routers

import (
	"BeeCustom/controllers"

	"github.com/astaxie/beego"
)

func init() {

	//后台用户路由
	beego.Router("/backenduser/index", &controllers.BackendUserController{}, "*:Index")
	beego.Router("/backenduser/create", &controllers.BackendUserController{}, "Get:Create")
	beego.Router("/backenduser/store", &controllers.BackendUserController{}, "Post:Store")
	beego.Router("/backenduser/datagrid", &controllers.BackendUserController{}, "Post:DataGrid")
	beego.Router("/backenduser/edit/?:id", &controllers.BackendUserController{}, "Get:Edit")
	beego.Router("/backenduser/update/?:id", &controllers.BackendUserController{}, "Patch:Update")
	beego.Router("/backenduser/delete/?:id", &controllers.BackendUserController{}, "Delete:Delete")

	//用户角色路由
	beego.Router("/role/index", &controllers.RoleController{}, "*:Index")
	beego.Router("/role/create", &controllers.RoleController{}, "Get:Create")
	beego.Router("/role/perm_lists/?:id", &controllers.RoleController{}, "Get:PermLists")
	beego.Router("/role/store", &controllers.RoleController{}, "Post:Store")
	beego.Router("/role/datagrid", &controllers.RoleController{}, "Post:DataGrid")
	beego.Router("/role/edit/?:id", &controllers.RoleController{}, "Get:Edit")
	beego.Router("/role/update/?:id", &controllers.RoleController{}, "Patch:Update")
	beego.Router("/role/delete/?:id", &controllers.RoleController{}, "Delete:Delete")
	beego.Router("/role/datalist", &controllers.RoleController{}, "Post:DataList")

	//资源路由
	beego.Router("/resource/index", &controllers.ResourceController{}, "*:Index")
	beego.Router("/resource/create", &controllers.ResourceController{}, "GET:Create")
	beego.Router("/resource/store", &controllers.ResourceController{}, "POST:Store")
	beego.Router("/resource/treegrid", &controllers.ResourceController{}, "POST:TreeGrid")
	beego.Router("/resource/edit/?:id", &controllers.ResourceController{}, "GET:Edit")
	beego.Router("/resource/update/?:id", &controllers.ResourceController{}, "PATCH:Update")
	beego.Router("/resource/delete/?:id", &controllers.ResourceController{}, "Delete:Delete")

	//通用选择面板
	//beego.Router("/resource/select", &controllers.ResourceController{}, "Get:Select")
	//用户有权管理的菜单列表（包括区域）
	//beego.Router("/resource/usermenutree", &controllers.ResourceController{}, "POST:UserMenuTree")
	//beego.Router("/resource/checkurlfor", &controllers.ResourceController{}, "POST:CheckUrlFor")

	beego.Router("/home/control", &controllers.HomeController{}, "*:Control")
	beego.Router("/home/login", &controllers.HomeController{}, "*:Login")
	beego.Router("/home/dologin", &controllers.HomeController{}, "Post:DoLogin")
	beego.Router("/home/logout", &controllers.HomeController{}, "*:Logout")
	beego.Router("/home/datareset", &controllers.HomeController{}, "Post:DataReset")

	beego.Router("/home/404", &controllers.HomeController{}, "*:Page404")
	beego.Router("/home/error/?:error", &controllers.HomeController{}, "*:Error")

	beego.Router("/", &controllers.HomeController{}, "*:Index")

}
