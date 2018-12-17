package proxy

import (
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/url"
)

type ParentProxyConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

func (c *ParentProxyConfig) URL() (*url.URL, error) {
	rawURL := fmt.Sprintf("http://%s:%s@%s:%s", c.Username, c.Password, c.Host, c.Port)
	return url.Parse(rawURL)
}

type AuthProxyConfig struct {
	Parent ParentProxyConfig
}

type authProxy struct {
	transport     *http.Transport
	parentAddress string
	authorization string
}

func NewAuthProxy(config *AuthProxyConfig) (Proxy, error) {
	proxyURL, err := config.Parent.URL()
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	basicAuthInfo := base64.StdEncoding.EncodeToString([]byte(config.Parent.Username + ":" + config.Parent.Password))
	authorization := fmt.Sprintf("Basic %s", basicAuthInfo)

	return &authProxy{
		transport:     transport,
		parentAddress: config.Parent.Host + ":" + config.Parent.Port,
		authorization: authorization,
	}, nil
}

func (p *authProxy) HandleHTTP(w http.ResponseWriter, r *http.Request) error {
	hj, _ := w.(http.Hijacker)
	cli := http.Client{
		Transport: p.transport,
	}
	r.RequestURI = ""

	resp, err := cli.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	conn, _, err := hj.Hijack()
	if err != nil {
		return err
	}
	defer conn.Close()

	return resp.Write(conn)
}

func (p *authProxy) HandleHTTPS(w http.ResponseWriter, r *http.Request) error {
	hj, _ := w.(http.Hijacker)

	// connection will be closed in "transfer()"
	parentConn, err := net.Dial("tcp", p.parentAddress)
	if err != nil {
		return err
	}

	// connection will be closed in "transfer()"
	clientConn, _, err := hj.Hijack()
	if err != nil {
		parentConn.Close()
		return err
	}

	r.Header.Set("Proxy-Authorization", p.authorization)
	if err := r.Write(parentConn); err != nil {
		return err
	}
	go transfer(parentConn, clientConn)
	go transfer(clientConn, parentConn)
	return nil
}
