package controllers

//
type RABCController struct {
	BaseController
}

//
//type Userrole struct {
//	Id         int64
//	Rolename   string `json:"name"`
//	Rolenumber string
//	Status     string `json:"role"`
//	Level      string
//}
//
//type Tree struct {
//	Id    int64  `json:"id"`
//	Nodes []Tree `json:"nodes"`
//}
//
//func (c *RABCController) Prepare() {
//	//先执行
//	c.BaseController.Prepare()
//	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
//	perms := []string{
//		"IIndex",
//		"ICreate",
//		"IEdit",
//		"IMake",
//		"IAduit",
//		"IDelete",
//		"EIndex",
//		"ECreate",
//		"EEdit",
//		"EMake",
//		"EAduit",
//		"EDelete",
//	}
//	c.checkAuthor(perms)
//
//	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
//	//权限控制里会进行登录验证，因此这里不用再作登录验证
//	//c.checkLogin()
//}
//
//
//func (c *RABCController) Test() {
//	// Check the permission.
//	//请求的资源v1/v2/aaa.jpg
//	//得到资源的扩展名Suffix-jpg，输入enforce中间那个(?i:pdf)不分大小写
//	alice, _ := sysinit.E.Enforce("alice", "/v1/v2/aaa.jpg", "write", "jpg")
//	utils.LogInfo(alice)
//	bobP, _ := sysinit.E.Enforce("bob", "/v1/v2/aaa.PDF", "delete", "PDF")
//	utils.LogInfo(bobP)
//	bobJ, _ := sysinit.E.Enforce("bob", "/v1/v2/aaa.jpg", "write", "jpg")
//	utils.LogInfo(bobJ)
//	bobT, _ := sysinit.E.Enforce("bob", "/v1/v2/aaa.ttt", "read", "ttt")
//	utils.LogInfo(bobT)
//}
//
//func (c *RABCController) Get() {
//	id := c.Ctx.Input.Param(":id")
//
//	c.Data["Id"] = id
//	c.Data["Ip"] = c.Ctx.Input.IP()
//
//	params := models.NewRoleQueryParam()
//	//获取数据列表和总数
//	roles := models.RoleDataList(&params)
//
//	if id != "" { //如果选中了用户，则显示用户所具有的角色
//		userroles, _ := sysinit.E.GetRolesForUser(id)
//		userrole := make([]Userrole, 0)
//		var level string
//		level = "2"
//		for _, v1 := range roles {
//			for _, v2 := range userroles {
//				ridNum, err := strconv.ParseInt(strings.Replace(v2, "role_", "", -1), 10, 64)
//				if err != nil {
//					beego.Error(err)
//				}
//				if ridNum == v1.Id {
//					level = "1" //if (row.Level === "1") checked: true
//				}
//			}
//			aa := make([]Userrole, 1)
//			aa[0].Id = v1.Id
//			aa[0].Rolename = v1.Name
//			aa[0].Level = level
//			userrole = append(userrole, aa...)
//			aa = make([]Userrole, 0)
//			level = "2"
//		}
//		c.Data["json"] = userrole //用户所具有的角色，勾选
//		c.ServeJSON()
//	} else {
//		c.Data["json"] = roles //角色列表
//		c.ServeJSON()
//	}
//}
//
////添加角色
//func (c *RABCController) Post() {
//
//	var role models.Role
//	role.Name = c.Input().Get("rolename")
//
//	permIds := c.GetString("perm_ids")
//	err := models.RoleSave(&role, permIds)
//	if err == nil && role.Id > 0 {
//		c.Data["json"] = "ok"
//		c.ServeJSON()
//	} else {
//		// c.Rsp(false, err.Error())
//		beego.Error(err)
//		c.Data["json"] = "wrong"
//		c.ServeJSON()
//		// return
//	}
//}
//
////AddPolicy(sec string, ptype string, rule []string)
////添加用户角色
////先删除用户所有角色
//func (c *RABCController) UserRole() {
//
//	uid := c.GetString("uid") //secofficeid
//	_, _ = sysinit.E.DeleteRolesForUser(uid) //数据库没有删掉！
//	ids := c.GetString("ids") //roleid
//	if ids != "" {
//		array := strings.Split(ids, ",")
//		for _, v1 := range array {
//			_, _ = sysinit.E.AddGroupingPolicy(uid, "role_"+v1) //management_api.go
//		}
//
//	}
//	c.Data["json"] = "ok"
//	c.ServeJSON()
//}
//
////给角色赋项目目录的权限
////先删除角色对于这个项目的所有权限
//func (c *RABCController) RolePermission() {
//	var success bool
//	var nodeidint int
//	var projurl, action, suf1, suf string
//	var err error
//	roleids := c.GetString("roleids")
//	rolearray := strings.Split(roleids, ",")
//
//	permissionids := c.GetString("permissionids")
//	permissionarray := strings.Split(permissionids, ",")
//	switch permissionarray[0] {
//	case "添加成果":
//		action = "POST"
//	case "编辑成果":
//		action = "PUT"
//	case "删除成果":
//		action = "DELETE"
//	case "读取成果":
//		action = "GET"
//	}
//
//	sufids := c.GetString("sufids")
//	sufarray := strings.Split(sufids, ",")
//	switch sufids {
//	case "任意":
//		suf = ".*"
//	case "":
//		suf = "(?i:PDF)"
//	case "PDF":
//		suf = "(?i:PDF)"
//	}
//	treeids := c.GetString("treeids") //项目目录id，25001,25002
//	treearray := strings.Split(treeids, ",")
//
//	treenodeids := c.GetString("treenodeids") //项目目录的nodeid 0.0.0-0.0.1-0.1.0-0.1.0
//	treenodearray := strings.Split(treenodeids, ",")
//	projectid := c.GetString("projid")
//
//
//	//取出项目目录的顶级
//	var nodesid, nodesids []string
//	if len(treenodearray) > 1 {
//		nodesids, err = highest(treenodearray, nodesid, 0)
//		if err != nil {
//			beego.Error(err)
//		}
//	} else {
//		nodesids = []string{"0"} //append(nodesids, "0")
//	}
//
//
//	//删除这些角色、项目id、权限的全部权限
//	for _, v1 := range rolearray {
//		// var paths []beegoormadapter.CasbinRule
//		o := orm.NewOrm()
//		qs := o.QueryTable("casbin_rule")
//		if action == "GET" {
//			_, err := qs.Filter("PType", "p").Filter("v0", "role_"+v1).Filter("v1__contains", "/"+projectid+"/").Filter("v2", action).Filter("v3", suf).Delete()
//			if err != nil {
//				beego.Error(err)
//			}
//		} else {
//			_, err := qs.Filter("PType", "p").Filter("v0", "role_"+v1).Filter("v1__contains", "/"+projectid+"/").Filter("v2", action).Delete()
//			if err != nil {
//				beego.Error(err)
//			}
//		}
//	}
//
//	_ = sysinit.E.LoadPolicy() //重载权限
//	_, _ = sysinit.E.RemoveFilteredPolicy(1, "/onlyoffice/"+strconv.FormatInt(attachments[0].Id, 10))
//	if treeids != "" {
//		for _, v1 := range rolearray {
//			for _, v2 := range permissionarray {
//				//定义读取、添加、修改、删除
//				switch v2 {
//				case "添加成果":
//					action = "POST"
//					suf = ".*"
//				case "编辑成果":
//					action = "PUT"
//					suf = ".*"
//				case "删除成果":
//					action = "DELETE"
//					suf = ".*"
//				case "读取成果":
//					action = "GET"
//					for i, v4 := range sufarray {
//						if v4 == "任意" {
//							suf = ".*"
//							break
//						} else if v4 == "" { //用户没展开则读取不到table4的select
//							suf = "(?i:PDF)"
//							break
//						} else {
//							suf1 = "(?i:" + v4 + ")"
//							if i == 0 {
//								suf = suf1
//							} else {
//								suf = suf + "," + suf1
//							}
//						}
//					}
//				}
//
//				for _, v3 := range nodesids {
//					nodeidint, err = strconv.Atoi(v3)
//					if err != nil {
//						beego.Error(err)
//					}
//					//id转成64位
//					pidNum, err := strconv.ParseInt(treearray[nodeidint], 10, 64)
//					if err != nil {
//						beego.Error(err)
//					}
//
//					//根据projid取出路径
//					proj, err := models.GetProj(pidNum)
//					if err != nil {
//						beego.Error(err)
//					}
//					if proj.ParentIdPath == "" || proj.ParentIdPath == "$#" {
//						projurl = "/" + strconv.FormatInt(proj.Id, 10) + "/*"
//					} else {
//						// projurl = "/" + strings.Replace(proj.ParentIdPath, "-", "/", -1) + "/" + treearray[nodeidint] + "/*"
//						projurl = "/" + strings.Replace(strings.Replace(proj.ParentIdPath, "#", "/", -1), "$", "", -1) + strconv.FormatInt(proj.Id, 10) + "/*"
//					}
//
//					sufarray := strings.Split(suf, ",")
//					for _, v5 := range sufarray {
//						success, _ = sysinit.E.AddPolicy("role_"+v1, projurl, action, v5) //来自casbin\management_api.go
//						//这里应该用AddPermissionForUser()，来自casbin\rbac_api.go
//					}
//				}
//			}
//		}
//	} else {
//		success = true
//	}
//
//
//	if success == true {
//		c.Data["json"] = "ok"
//	} else {
//		c.Data["json"] = "wrong"
//	}
//	c.ServeJSON()
//}
//
////迭代查出最高级的树状目录
////nodesid是数组的序号
////nodeid是节点号：0.0.1   0.0.1.0
//func highest(nodeid, nodesid []string, i int) (nodesid1 []string, err error) {
//	if i == 0 {
//		nodesid = append(nodesid, "0")
//	}
//	var i1 int
//	for i1 = i; i1 < len(nodeid)-1; i1++ {
//		matched, err := regexp.MatchString("(?i:"+nodeid[i]+")", nodeid[i1+1])
//		// fmt.Println(matched)
//		if err != nil {
//			beego.Error(err)
//		}
//		if !matched {
//			i = i1 + 1
//			nodesid = append(nodesid, strconv.Itoa(i1+1))
//			break
//		} else {
//			if i == len(nodeid)-2 {
//				return nodesid, err
//			}
//		}
//	}
//	if i1 < len(nodeid)-1 {
//		nodesid, err = highest(nodeid, nodesid, i)
//	}
//	return nodesid, err
//}
//
////查询角色所具有的权限对应的项目目录
//func (c *RABCController) GetRolePermission() {
//
//	roleid := c.GetString("roleid") //角色id
//	action := c.GetString("action")
//	projectid := c.GetString("projectid")
//	sufids := c.GetString("sufids") //扩展名
//
//	var suf string
//	switch sufids {
//	case "任意":
//		suf = ".*"
//	case "":
//		suf = "(?i:PDF)"
//	case "PDF":
//		suf = "(?i:PDF)"
//	}
//
//	var paths []beegoormadapter.CasbinRule
//
//	o := orm.NewOrm()
//	qs := o.QueryTable("casbin_rule")
//	if action == "GET" || action == "" {
//		_, err := qs.Filter("PType", "p").Filter("v0", "role_"+roleid).Filter("v1__contains", "/"+projectid+"/").Filter("v2", "GET").Filter("v3", suf).All(&paths)
//		if err != nil {
//			beego.Error(err)
//		}
//
//	} else {
//		_, err := qs.Filter("PType", "p").Filter("v0", "role_"+roleid).Filter("v1__contains", "/"+projectid+"/").Filter("v2", action).All(&paths)
//		if err != nil {
//			beego.Error(err)
//		}
//	}
//
//	var projids []string
//	for _, v1 := range paths {
//		projid := strings.Replace(v1.V1, "/*", "", -1)
//		projids = append(projids, path.Base(projid))
//	}
//
//	c.Data["json"] = projids
//	c.ServeJSON()
//}
//
//func (c *RABCController) Delete() {
//	roleids := c.GetString("ids")
//	rolearray := strings.Split(roleids, ",")
//	for _, v1 := range rolearray {
//
//		idNum, err := strconv.ParseInt(v1, 10, 64)
//		if err != nil {
//			beego.Error(err)
//		}
//		_, err = models.RoleDelete(idNum)
//		if err != nil {
//			utils.LogDebug(fmt.Sprintf("Delete Role:%v",err))
//		} else {
//			c.Data["json"] = "ok"
//			c.ServeJSON()
//		}
//	}
//}
//
//
//func (c *RABCController) Update() {
//	var role models.Role
//	roleid := c.Input().Get("roleid")
//	idNum, err := strconv.ParseInt(roleid, 10, 64)
//	if err != nil {
//		utils.LogDebug(fmt.Sprintf("ParseInt:%v",err))
//	}
//	role.Id = idNum
//	role.Name = c.Input().Get("rolename")
//
//	permIds := c.GetString("perm_ids")
//	err = models.RoleSave(&role,permIds)
//	if err == nil {
//		c.Data["json"] = "ok"
//		c.ServeJSON()
//	} else {
//		utils.LogDebug(fmt.Sprintf("Update Role:%v",err))
//	}
//}
