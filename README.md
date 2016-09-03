# agent
Agent is a service application which is part of new generation of DevOps framework or specification focuses on building channels.inspired by docker.

## 框架特点
- 支持绝大部分操作系统
- 架构离散化,无需负载均衡,容易排查问题,自由分组.
- 服务器安全性独立.
- 安全,基于tls双向验证,参数签名.
- 可控性强
- 提供RESTful接口,方便内管系统开发.
- 提供强大的命令行功能,网络瘫痪时运维人员也可控
- 插件功能,
- 基于证书扩展字段的功能授权,路由等,方便不同团队隔离
- 支持自动反向代理
- *可支持动态配置文件的分发管理*
- *可支持程序版本的细粒度分发管理*
- *可支持统计数据/日志的收发*


## NOTICE

I DO NOT IMPLEMENT FUNCTIONS LIKE HOLDING OLD WOMAN ACROSS STREET (ﾒ ﾟ皿ﾟ)ﾒ,I CREATE TRAFFIC LIGHT AND ZEBRA CROSSING!

目前仅在MacOS上检查通过。

## REQUIREMENT
Go Verion 1.7+

## USAGE
```shell
// start a server with debug,console, no-root privilege,certificate's extension authenticate mode.
./agent daemon start -c -r=false -D --cert-exten-auth
```
>INFO[0000:645] Daemon has completed initialization
>INFO[0000:645] Daemon start                                  agentid=56F0F0B9-10E6-450E-96B6-AA212C4909AD buildtime=20160903122212 version=0.1.0
>INFO[0000:645] API listen on 127.0.0.1:3567

```shell
// launch a version request with specifying client certificate path
./agent version --tlscert=./cert/client-cert.pem --tlskey=./cert/client-key.pem
```
>Client:
> API version:  0.1.0
> Go version:   go1.7
> Built:        20160903122212
> OS/Arch:      darwin/amd64
> Kernel:       15.6.0 
>
>Server:
> API version:  0.1.0
> Go version:   go1.7
> Built:        20160903122212
> OS/Arch:      darwin/amd64
> Kernel:       15.6.0


## DEVELOP
New functionality should be implement by two steps at least: declare interface at route point and implement in daemon package.it will be better to implement client's cli command. 