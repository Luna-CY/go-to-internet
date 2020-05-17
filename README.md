# go-to-internet

#### 介绍
基于TLS协议的代理服务器

特性：

- 简单：启动服务器简单、启动客户端简单，目标就是干什么都简单
- 安全：基于TLS加密信息传输，除了会话两端无法获取通信内容
- 快速：极简的代理握手协议，快速建立客户与目标的双向通道

#### 安装说明
一、下载二进制包

根据需要下载服务器、客户端的二进制包，服务器的二进制包内包含`manager-go-to-net`工具

二、手动安装

拉取git仓库，分别build服务器与客户端

服务器：`go build src/main/ser-go-to-net`

管理工具：`go build src/main/manager-go-to-net`

客户端：`go build src/main/cli-go-to-net`

#### 使用说明
- `ser-go-to-net` 
代理服务器，服务器命令仅支持在*unix环境下使用，不支持windows环境，客户端请根据系统环境进行下载
运行服务器需要域名与证书，自行购买域名并解析到服务器，可以通过`manager-go-to-net`管理工具安装acme并申请证书，请使用命令`manager-go-to-net --help`查看帮助信息

- `manager-go-to-net`
管理工具，包含用户管理以及Acme辅助工具

- `cli-go-to-net`
代理客户端，根据系统环境选择

#### 本地开发指南
将`custom-root-ca/ca/certs/cacert.pem`根证书添加进系统的根证书库

一、使用`local.luna.xin`

- 绑定`local.luna.xin`域名到本地`127.0.0.1`
- 运行服务器`ser-go-to-net -H local.luna.xin -c cert/server.pem -k cert/server.key`启动

二、自定义域名
- 通过`custom-root-ca`签发自定义域名证书
- 运行服务器`ser-go-to-net -H 域名 -c 证书pem路径 -k 证书key路径`启动
