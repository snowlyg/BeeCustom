## BeeCustom

# 简介
BeeCustom 是基于 `Beego` 开发的关务管理系统。前端框架基于 `Layui` 进行开发，功能已经贴近海关的单一窗口。

# 后端框架
1. 基于 `Beego` ，使用官方的 orm、cache、session、logs 等模块，感谢原作者提供了如此简单易用的框架;

# 前端框架
1. 基于 Layui；


# 安装方法

本系统基于 `beego` 开发，默认使用 `mysql` 数据库，缓存 `redis` 

1. 安装 `golang` 环境

2. 安装本系统
```
git clone https://git.dev.tencent.com/Dreamfish/BeeCustom.git

```

4. 修改配置文件 conf/app.conf

 需要配置 `mysql` 和 `redis` 的参数
 
5. 运行
在 BeeCustom 目录使用 `beego` 官方提供的命令运行
```
bee run 
```

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

# 参考项目
 - [lhtzbj12/sdrms](https://gitee.com/lhtzbj12/sdrms/tree/master)

