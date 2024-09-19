# GhProxy-Go-0RTT

![pull](https://img.shields.io/docker/pulls/wjqserver/ghproxy-0rtt.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/WJQSERVER/ghproxy-go-0RTT)](https://goreportcard.com/report/github.com/WJQSERVER/ghproxy-go-0RTT)

使用Go实现的GHProxy,支持Git clone等文件拉取,支持Docker部署,支持速率限制

## 项目说明

## 起源与项目特点

**本项目自v0.2.0版本开始,对底层核心实现代码进行重构,彻底脱离原有实现**

弃用[原实现](https://github.com/0-RTT/ghproxy-go)过于依赖的jiasu.in网页,更换前端界面,并增加Docker支持,并增加了高频请求限制,以避免CC

使用Caddy实现了符合[RFC 7234](https://httpwg.org/specs/rfc7234.html)的HTTP Cache

0.1.3版本前,主程序 Fork & Sync From [0RTT/ghproxy-go](https://github.com/0-RTT/ghproxy-go) 

从0.1.3版本开始,自行进行改进与维护,并自行实现功能,如0.1.4的外部配置文件和日志模块

页面及其余部分资源,项目内容类似,故对其进行修改并复用 [WJQSERVER-STUDIO/ghproxy-go](https://github.com/WJQSERVER-STUDIO/ghproxy-go)

### LICENSE

本项目继承于[WJQSERVER-STUDIO/ghproxy-go](https://github.com/WJQSERVER-STUDIO/ghproxy-go)的APACHE2.0 LICENSE VERSION

## 使用示例

```
https://ghproxy0rtt.1888866.xyz/raw.githubusercontent.com/WJQSERVER-STUDIO/tools-stable/main/tools-stable-ghproxy.sh

git clone https://ghproxy0rtt.1888866.xyz/github.com/WJQSERVER-STUDIO/ghproxy-go.git
```

## 部署说明

### Docker部署

- Docker-cli

```
docker run -p 8078:80 -v ./ghproxy/log/run:/data/ghproxy/log -v ./ghproxy/log/caddy:/data/caddy/log --restart always wjqserver/ghproxy-0rtt
```

- Docker-Compose

    参看[docker-compose.yml](https://github.com/WJQSERVER/ghproxy-go-0RTT/blob/main/docker-compose.yml)

### 外部配置文件

本项目采用config.yaml作为外部配置,默认配置如下
使用Docker部署时,慎重修改config.yaml,以免造成不必要的麻烦

```
port: 8080 # 监听端口
host: "127.0.0.1" # 监听地址
sizelimit: 131072000 # 125MB (文件大小默认限制)
logfilepath: "/data/ghproxy/log/ghproxy-0rtt.log" # 日志存储目录
CorsAllowOrigins: true # 是否允许跨域请求
```

### Caddy反代配置

```
example.com {
    reverse_proxy {
        to 127.0.0.1:8078
        header_up X-Real-IP {remote_host}	    
        header_up X-Real-IP {http.request.header.CF-Connecting-IP}
        header_up X-Forwarded-For {http.request.header.CF-Connecting-IP}
        header_up X-Forwarded-Proto {http.request.header.CF-Visitor}
    }    
}
```

## TODO & BETA

### TODO

- [x] 允许更多参数通过config结构传入
- [x] 改进程序效率

### BETA

- [x] Docker Pull 代理 (DEV版本内实现)
