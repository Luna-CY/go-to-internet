# go-to-internet

#### 介绍
基于TLS协议的透明代理服务器

特性：

- 简单：启动服务器简单、启动客户端简单，内置了对acme工具的支持，可快速申请https证书
- 安全：基于TLS加密信息传输，除了会话两端无法获取通信内容
- 快速：极简的代理握手协议，快速建立客户与目标的双向通道
- 支持多用户：在单域名以及端口下支持创建多用户
- 用户管理：每个用户支持对过期时间、传输速度、连接数量进行限制

#### 安装说明
一、下载二进制包

根据需要下载服务器、客户端的二进制包，服务器的二进制包内包含`manager-go-to-net`工具

二、手动安装

拉取git仓库，分别build服务器与客户端

服务器：`go build src/main/ser-go-to-net`

管理工具：`go build src/main/manager-go-to-net`

客户端：`go build src/main/cli-go-to-net`

#### 使用教程

使用教程见[这里](https://blog.luna.xin/article/29.html)

#### 使用说明: 服务端

一、 `ser-go-to-net`服务默认使用`/etc/go-to-net/users.josn`来管理用户配置，无需手动编辑该文件，`manager-go-to-net`工具对用户管理提供了支持

通过`manager-go-to-net`的子命令`user`来管理用户信息

- `manager-go-to-net user -list`查看现有的用户信息，该命令显示用户的相关信息，格式：`用户名称 : 过期时间 : 速度限制`
- `manager-go-to-net user -add -u USERNAME -p PASSWORD`添加新的用户，有关用户的其他参数可以通过`manager-go-to-net user -help`来获取
- `manager-go-to-net user -upd -u USERNAME -p NEW_PASSWORD`编辑一个现有用户，`-upd`参数的使用方法可以通过`manager-go-to-net user -help`来获取
- `manager-go-to-net user -del -u USERNAME`删除一个现有用户

二、 `ser-go-to-net`服务命令必须提供一个域名及相关TLS证书，域名需要自行到云服务商处购买并解析到服务器的ip，`manager-go-tonet`的子命令`acme`对申请TLS证书提供了支持

- `manager-go-to-net acme -install`命令支持安装`acme.sh`证书申请工具，如果已经安装该工具可以忽略该命令
- `manager-go-to-net acme -issue -nginx|-standalone -hostname YOUR_HOST`支持通过`acme.sh`工具来申请证书，证书相关的具体信息请参阅`acme.sh`的官方文档

三、 默认情况下`ser-go-to-net`命令能够根据域名查找`acme.sh`工具申请的证书，其默认路径一般在`/root/.acme.sh/YOUR_HOST`目录下，如果不使用默认规则，可以通过运行参数来指定证书位置，使用方式可以通过`ser-go-to-net -help`来获取

- `ser-go-to-net -H YOUR_HOST`命令在启动时默认查找`/root/.acme.sh/YOUR_HOST/fullchain.cer`以及`/root/.acme.sh/YOUR_HOST/YOUR_HOST.key`两个文件
- `ser-go-to-net -H YOUR_HOST -c /path/to/fullchain.cer -k /path/to/YOUR_HOST.key`命令在启动时使用指定的证书

#### 使用说明: 客户端

`cli-go-to-net`允许指定服务器域名、服务器端口、本地监听地址、本地监听端口、用户名、用户密码等

- `cli-go-to-net -sh YOUR_HOST -u USERNAME -p PASSWORD`命令将在连接到服务器的`443`端口，并监听本地`127.0.0.1`的`1280`端口
- `cli-go-to-net -sh YOUR_HOST -sp 4433 -la 0.0.0.0 -lp 1234 -u USERNAME -p PASSWORD`命令将连接到服务器的`4433`端口，并监听本地`0.0.0.0`的`1234`端口

#### 客户端GUI项目

- macos: [GoToNetUI-X](https://gitee.com/Luna-CY/GoToNetUI-X)
- android: [GoToNetUI-A (开发中)](https://gitee.com/Luna-CY/GoToNetUI-A)
- windows: GoToNetUI-W (计划中)

#### 本地开发指南
将`custom-root-ca/ca/certs/cacert.pem`根证书添加进系统的根证书库

一、使用`local.luna.xin`

- 绑定`local.luna.xin`域名到本地`127.0.0.1`
- 运行服务器`ser-go-to-net -H local.luna.xin -c cert/server.pem -k cert/server.key`启动

二、自定义域名
- 通过`custom-root-ca`签发自定义域名证书
- 运行服务器`ser-go-to-net -H 域名 -c 证书pem路径 -k 证书key路径`启动
