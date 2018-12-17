# xproxy

WIP

## Examples

```
host: localhost
port: 1080
proxies:
  - name: internal
    type: noproxy
  - name: corp_proxy
    type: auth
    config:
      host: your-corp-proxy.example.com
      port: 8080
      username: hoge
      password: fuga
proxy_binds:
  - name: internal
    match:
      hosts: ["*.your-corp.com"]
      ips: ["192.168.1.0/24", "172.16.0.0/16"]
  - name: corp_proxy
    default: true
```