package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/pkg/errors"
	"github.com/tattsun/xproxy/proxy"
)

type Server struct {
	host    string
	port    string
	handler http.HandlerFunc
}

type Binding struct {
	proxy   proxy.Proxy
	matcher *proxy.Matcher
}

func NewServer(host string, port string, config *Config) (*Server, error) {
	proxies := make(map[string]proxy.Proxy)
	for _, p := range config.Proxies {
		if p.Name == "" {
			return nil, errors.New("proxy name is empty")
		}
		switch p.Type {
		case "auth":
			host, ok := p.Config["host"]
			if !ok {
				return nil, errors.Errorf(`proxy "%s" has not config.host`, p.Type)
			}
			port, ok := p.Config["port"]
			if !ok {
				return nil, errors.Errorf(`proxy "%s" has not config.port`, p.Type)
			}
			username, ok := p.Config["username"]
			if !ok {
				return nil, errors.Errorf(`proxy "%s" has not config.username`, p.Type)
			}
			password, ok := p.Config["password"]
			if !ok {
				return nil, errors.Errorf(`proxy "%s" has not config.password`, p.Type)
			}
			conf := proxy.AuthProxyConfig{
				Parent: proxy.ParentProxyConfig{
					Host:     host,
					Port:     port,
					Username: username,
					Password: password,
				},
			}
			proxy, err := proxy.NewAuthProxy(&conf)
			if err != nil {
				return nil, errors.Wrapf(err, `failed to create proxy "%s"`, p.Name)
			}
			proxies[p.Name] = proxy
			break
		case "noproxy":
			proxies[p.Name] = proxy.NewNoProxy()
			break
		default:
			return nil, errors.Errorf(`invalid proxy type "%s"`, p.Type)
		}
	}

	bindings := make([]*Binding, len(config.ProxyBinds))
	for i, bind := range config.ProxyBinds {
		hosts := bind.Match.Hosts
		if hosts == nil {
			hosts = make([]string, 0)
		}
		if bind.Default {
			hosts = append(hosts, "*")
		}

		ips := bind.Match.IPs
		if ips == nil {
			ips = make([]string, 0)
		}

		p, ok := proxies[bind.Name]
		if !ok {
			return nil, errors.Errorf(`proxy "%s" is not defined`, bind.Name)
		}

		matcher, err := proxy.NewMatcher(hosts, ips)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create matcher")
		}

		bindings[i] = &Binding{
			matcher: matcher,
			proxy:   p,
		}
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		log.Print(r.URL.String())
		for _, binding := range bindings {
			if binding.matcher.Host(r.Host) {
				if r.Method == "CONNECT" {
					binding.proxy.HandleHTTPS(w, r)
				} else {
					binding.proxy.HandleHTTP(w, r)
				}
				return
			}

			addr, err := net.ResolveIPAddr("ip", r.Host)
			if err != nil {
				continue
			}
			if binding.matcher.IP(addr.IP) {
				if r.Method == "CONNECT" {
					binding.proxy.HandleHTTPS(w, r)
				} else {
					binding.proxy.HandleHTTP(w, r)
				}
				return
			}
		}

		w.WriteHeader(403)
		fmt.Fprint(w, "binding not found")
	}

	return &Server{
		host:    host,
		port:    port,
		handler: handler,
	}, nil
}

func (s *Server) Start() error {
	server := http.Server{
		Addr:    s.host + ":" + s.port,
		Handler: s.handler,
	}
	return server.ListenAndServe()
}
