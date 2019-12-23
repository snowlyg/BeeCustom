## 问题汇总 ERRORS

 1.`beego` 升级到v1.10.1后，启动本项目时报错
 ```
 cannot find package "github.com/gomodule/redigo/redis"
 ```
 解决方法很简单，只需要在终端运行下面命令，下载需要的包即可
 
 ```
 go get github.com/gomodule/redigo/redis
 ```

2.启动报错 `beego runtime error: invalid memory address or nil pointer dereference`
没有开启 session 
```
beego.BConfig.WebConfig.Session.SessionOn = true
```
或者在 `app.conf` 添加
```
sessionon = true
```

3. webhook 自动部署代码

`webhook` 自动部署代码，修改 `*.conf` 文件会失败， 需要到服务器手动部署。