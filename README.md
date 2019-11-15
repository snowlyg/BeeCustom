## BeeCustom

# 简介
BeeCustom 是基于 `Beego` 开发的关务管理系统。前端框架基于 `Layui` 进行开发，功能已经贴近海关的单一窗口。

# 后端框架
1. 基于 `Beego` ，使用官方的 orm、cache、session、logs 等模块，感谢原作者提供了如此简单易用的框架;

# 前端框架
1. 基于 Layui


# 安装方法

本系统基于 `beego` 开发，默认使用 `mysql` 数据库，缓存 `redis` 

1. 安装 `golang` 环境

2. 安装 `BeeCustom`

```
git clone https://git.dev.tencent.com/Dreamfish/BeeCustom.git

```
 
3.加载依赖管理包 使用 gopm 管理包

 ``` 

  go get -v -u github.com/gpmgo/gopm
  
  // 拉取依赖到缓存目录
  gopm get 

  // 拉取依赖到缓存目录
  gopm build

  //运行
  ./BeeCustom
  
```

4. 修改配置文件 conf/app.conf

 需要配置 `mysql` 和 `redis` 的参数
 
5. session 使用 redis 管理,新建 session 表
```
 CREATE TABLE `session` (
        `session_key` char(64) NOT NULL,
        `session_data` blob,
        `session_expiry` int(11) unsigned NOT NULL,
        PRIMARY KEY (`session_key`)
    ) ENGINE=MyISAM DEFAULT CHARSET=utf8;
```
 
 运行
在 BeeCustom 目录使用 `beego` 官方提供的命令运行
```
bee run 

改用 gowatch 解决修改前端文件不热加载问题， gowatch 打包效率更快
```

# 增加了 gulp 打包前端资源
```
npm install

gulp

```

# 使用组件
 1. 数据格式化 [https://github.com/IamBusy/amoeba](https://github.com/IamBusy/amoeba)

在浏览器里打开 http://localhost:8080 进行访问


# 问题
 `beego` 升级到v1.10.1后，启动本项目时报错
 ```
 cannot find package "github.com/gomodule/redigo/redis"
 ```
 解决方法很简单，只需要在终端运行下面命令，下载需要的包即可
 
 ```
 go get github.com/gomodule/redigo/redis
 ```

启动报错 `beego runtime error: invalid memory address or nil pointer dereference`
没有开启 session 
```
beego.BConfig.WebConfig.Session.SessionOn = true
```
或者在 `app.conf` 添加
```
sessionon = true
```
# 参考项目
 - [lhtzbj12/sdrms](https://gitee.com/lhtzbj12/sdrms/tree/master)

