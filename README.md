# go-to-internet

#### 介绍
基于TLS协议的透明代理服务器

特性：

- 简单：启动服务器简单、启动客户端简单，内置了对acme工具的支持，可快速申请https证书
- 安全：基于TLS加密信息传输，除了会话两端无法获取通信内容
- 快速：极简的代理握手协议，快速建立客户与目标的双向通道
- 支持多用户：在单域名以及端口下支持创建多用户
- 用户管理：每个用户支持对过期时间、传输速度、连接数量进行限制

#### Wiki

使用教程见[这里](https://blog.luna.xin/article/29.html)，或者查看[wiki](https://github.com/Luna-CY/go-to-internet/wiki)

#### 客户端GUI项目

- macos: [GoToNetUI-X](https://github.com/Luna-CY/GoToNetUI-X)
- android: [GoToNetUI-A (开发中)](https://github.com/Luna-CY/GoToNetUI-A)
- windows: GoToNetUI-W (计划中)

#### 本地开发指南
将`custom-root-ca/ca/certs/cacert.pem`根证书添加进系统的根证书库

一、使用`local.luna.xin`

- 绑定`local.luna.xin`域名到本地`127.0.0.1`
- 运行服务器`ser-go-to-net -H local.luna.xin -c cert/server.pem -k cert/server.key`启动

二、自定义域名
- 通过`custom-root-ca`签发自定义域名证书
- 运行服务器`ser-go-to-net -H 域名 -c 证书pem路径 -k 证书key路径`启动
