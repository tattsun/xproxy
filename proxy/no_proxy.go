package proxy

import (
	"net"
	"net/http"
)

type noProxy struct {
	transport *http.Transport
}

func NewNoProxy() (Proxy, error) {
	return &noProxy{
		transport: &http.Transport{},
	}, nil
}

func (p *noProxy) HandleHTTP(w http.ResponseWriter, r *http.Request) error {
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

func (p *noProxy) HandleHTTPS(w http.ResponseWriter, r *http.Request) error {
	// connection will be closed in "transfer()"
	destConn, err := net.Dial("tcp", r.Host)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	hj, _ := w.(http.Hijacker)

	// connection will be closed in "transfer()"
	clientConn, _, err := hj.Hijack()
	if err != nil {
		destConn.Close()
		return err
	}

	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)
	return nil
}
