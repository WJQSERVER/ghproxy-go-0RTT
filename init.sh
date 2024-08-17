#!/bin/bash

if [ ! -f /data/caddy/config/Caddyfile ]; then
    cp /data/caddy/Caddyfile /data/caddy/config/Caddyfile
fi

if [ ! -f /data/ghproxy/config/config.yaml ]; then
    cp /data/ghproxy/config.yaml /data/ghproxy/config/config.yaml
fi

/data/caddy/caddy run --config /data/caddy/config/Caddyfile > /data/ghproxy/log/caddy.log 2>&1 &

/data/ghproxy/ghproxy > /data/ghproxy/log/ghproxy.log 2>&1 &

while [[ true ]]; do
    sleep 1
done    

