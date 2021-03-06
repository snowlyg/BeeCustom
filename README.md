## BeeCustom

#### 更新日志
[更新日志](UPDATE.MD)

###### Require
- go 1.13.x

### 简介
BeeCustom 是基于 `Beego` 开发的关务管理系统。前端框架基于 `Layui` 进行开发，功能参考海关的单一窗口。

### 后端框架
1. 基于 `Beego` ，使用官方的 orm、cache、session、logs 等模块，感谢原作者提供了如此简单易用的框架;

### 前端框架
1. 基于 Layui

### 安装方法
本系统基于 `beego` 开发，默认使用 `mysql` 数据库，缓存 `redis` 

1. 安装 `golang` 环境

2. 安装 `BeeCustom`

```shell script
git clone https://e.coding.net/snowlyg/BeeCustom.git

#或者
git clone https://github.com/snowlyg/BeeCustom.git

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
 
- 运行
在 BeeCustom 目录使用 `beego` 官方提供的命令运行
改用 [gowatch](https://gitee.com/silenceper/gowatch) 解决修改前端文件不热加载问题， gowatch 打包效率更快

```shell script
bee run 

#or

gowatch 

```

### 增加了 gulp 打包前端资源
```shell script
npm install
gulp

```

### 使用组件
[组件列表](PLUGS.md)

### 问题
[问题汇总](ERRORS.md)

### 问题记录
[待处理BUG](BUG.md)

### 演示地址
[BeeCustom](https:bee.snowlyg.com)

> 账号/密码: admin/123456
>

### 打包
bee pack -be GOOS=linux

 

