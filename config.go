package main

type Proxy struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

type Match struct {
	Hosts []string `yaml:"hosts"`
	IPs   []string `yaml:"ips"`
}

type ProxyBind struct {
	Match Match `yaml:"match"`
}

type Config struct {
	Proxies    []Proxy     `yaml:"proxies"`
	ProxyBinds []ProxyBind `yaml:"proxy_binds"`
}
