# GhProxy-Go-0RTT

![pull](https://img.shields.io/docker/pulls/wjqserver/ghproxy-0rtt.svg)

使用Go实现的GHProxy,支持Git clone等文件拉取,支持Docker部署,支持速率限制

APACHE2.0 LICENSE VERSION

弃用[原实现](https://github.com/0-RTT/ghproxy-go)过于依赖的jiasu.in网页,更换前端界面,并增加Docker支持,并增加了高频请求限制,以避免CC

使用Caddy实现了符合[RFC 7234](https://httpwg.org/specs/rfc7234.html)的HTTP Cache

主程序 Fork & Sync From [0RTT/ghproxy-go](https://github.com/0-RTT/ghproxy-go)

页面及其余部分来自 [WJQSERVER-STUDIO/ghproxy-go](https://github.com/WJQSERVER-STUDIO/ghproxy-go)

# 使用示例

```
https://ghproxy0rtt.1888866.xyz/raw.githubusercontent.com/WJQSERVER-STUDIO/tools-stable/main/tools-stable-ghproxy.sh

git clone https://ghproxy0rtt.1888866.xyz/github.com/WJQSERVER-STUDIO/ghproxy-go.git
```

# Docker部署

- Docker-cli

```
docker run -p 8078:80 -v ./ghproxy/log/run:/data/ghproxy/log -v ./ghproxy/log/caddy:/data/caddy/log --restart always wjqserver/ghproxy-0rtt
```

- Docker-Compose

    参看[docker-compose.yml](https://github.com/WJQSERVER/ghproxy-go-0RTT/blob/main/docker-compose.yml)

# Caddy反代配置

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
