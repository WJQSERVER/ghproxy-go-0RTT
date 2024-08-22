# GhProxy-Go-0RTT

APACHE2.0 LICENSE VERSION

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
        to 172.20.20.221:80
        header_up X-Real-IP {remote_host}	    
        header_up X-Real-IP {http.request.header.CF-Connecting-IP}
        header_up X-Forwarded-For {http.request.header.CF-Connecting-IP}
        header_up X-Forwarded-Proto {http.request.header.CF-Visitor}
    }    
}
```
