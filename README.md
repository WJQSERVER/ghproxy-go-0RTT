# ghproxy-go
安装go
```
sudo apt update
sudo apt upgrade
wget https://golang.org/dl/go1.22.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz
sudo nano /etc/profile
export PATH=$PATH:/usr/local/go/bin
source /etc/profile
go version
```
配置ghproxy.service
```
sudo bash -c 'cat << EOF > /etc/systemd/system/ghproxy.service
[Unit]
Description=ghproxy
After=network.target

[Service]
ExecStart=/usr/local/go/bin/go run /main.go所在路径
Restart=always
User=root
Group=root
WorkingDirectory=/main.go所在路径

[Install]
WantedBy=multi-user.target
EOF'
```
示例：
```
[Unit]
Description=ghproxy
After=network.target

[Service]
ExecStart=/usr/local/go/bin/go run /www/wwwroot/gh.jiasu.in/main.go
Restart=always
User=root
Group=root
WorkingDirectory=/www/wwwroot/gh.jiasu.in

[Install]
WantedBy=multi-user.target
```
设置开机自启：```systemctl enable ghproxy.service```

启动：```systemctl start ghproxy.service```

重启：```systemctl restart ghproxy.service```

查询运行状态：```systemctl status ghproxy.service```

配置nginx反代
```
    
    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
  ```

