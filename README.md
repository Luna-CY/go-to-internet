# go-to-internet

#### 介绍
基于TLS协议的代理服务器

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
