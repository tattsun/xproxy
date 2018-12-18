package proxy

import (
	"net/http"
	"sync"
)

type loadBalancingProxy struct {
	mtx     *sync.Mutex
	i       int
	proxies []Proxy
}

func NewLoadBalancingProxy(proxies []Proxy) Proxy {
	return &loadBalancingProxy{
		mtx:     new(sync.Mutex),
		i:       0,
		proxies: proxies,
	}
}

func (p *loadBalancingProxy) incr() int {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	i := p.i
	i += 1
	if i >= len(p.proxies) {
		i = 0
	}
	p.i = i

	return i
}

func (p *loadBalancingProxy) HandleHTTP(w http.ResponseWriter, r *http.Request) error {
	proxy := p.proxies[p.incr()]
	return proxy.HandleHTTP(w, r)
}

func (p *loadBalancingProxy) HandleHTTPS(w http.ResponseWriter, r *http.Request) error {
	proxy := p.proxies[p.incr()]
	return proxy.HandleHTTPS(w, r)
}
