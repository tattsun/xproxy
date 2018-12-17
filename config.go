package main

import (
	"io"

	"gopkg.in/yaml.v2"
)

type Proxy struct {
	Name   string            `yaml:"name"`
	Type   string            `yaml:"type"`
	Config map[string]string `yaml:"config"`
}

type Match struct {
	Hosts []string `yaml:"hosts"`
	IPs   []string `yaml:"ips"`
}

type ProxyBind struct {
	Proxy   string `yaml:"proxy"`
	Match   Match  `yaml:"match"`
	Default bool   `yaml:"default"`
}

type Config struct {
	Host       string      `yaml:"host"`
	Port       string      `yaml:"port"`
	Proxies    []Proxy     `yaml:"proxies"`
	ProxyBinds []ProxyBind `yaml:"proxy_binds"`
}

func ParseConfig(r io.Reader) (*Config, error) {
	var config Config
	if err := yaml.NewDecoder(r).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
