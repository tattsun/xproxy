package server

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
)

type ProxyServer struct {
	host string
	port string

	proxyHost     string
	proxyPort     string
	proxyUser     string
	proxyPass     string
	authorization string
}

func NewProxyServer(host string, port string, proxyHost string, proxyPort string, username string, password string) *ProxyServer {
	authInfo := []byte(username + ":" + password)
	authorization := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString(authInfo))
	return &ProxyServer{
		host:          host,
		port:          port,
		proxyHost:     proxyHost,
		proxyPort:     proxyPort,
		proxyUser:     username,
		proxyPass:     password,
		authorization: authorization,
	}
}

func (s *ProxyServer) Start() error {
	log.Println(s.host + ":" + s.port)
	server := &http.Server{
		Addr:    s.host + ":" + s.port,
		Handler: http.HandlerFunc(s.handleRequest),
	}
	return server.ListenAndServe()
}

func (s *ProxyServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	var err error
	log.Println(r.URL.String())
	if r.Method == "CONNECT" {
		err = s.handleHTTPS(w, r)
	} else {
		err = s.handleHTTP(w, r)
	}
	if err != nil {
		log.Println(err)
		// TODO: handle error
	}
}

func (s *ProxyServer) handleHTTP(w http.ResponseWriter, r *http.Request) error {
	hj, _ := w.(http.Hijacker)
	proxyUrl, _ := url.Parse(fmt.Sprintf("http://%s:%s@%s:%s", s.proxyUser, s.proxyPass, s.proxyHost, s.proxyPort))
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}

	r.RequestURI = ""

	if resp, err := client.Do(r); err != nil {
		return err
	} else if conn, _, err := hj.Hijack(); err != nil {
		return err
	} else {
		defer conn.Close()
		defer resp.Body.Close()
		resp.Write(conn)
	}

	return nil
}

func transfer(dest io.WriteCloser, src io.ReadCloser) {
	defer func() {
		if dest != nil {
			dest.Close()
		}
	}()
	defer func() {
		if src != nil {
			src.Close()
		}
	}()
	if dest != nil && src != nil {
		io.Copy(dest, src)
	}
}

func (s *ProxyServer) handleHTTPS(w http.ResponseWriter, r *http.Request) error {
	hj, _ := w.(http.Hijacker)

	if proxyConn, err := net.Dial("tcp", s.proxyHost+":"+s.proxyPort); err != nil {
		return err
	} else if clientConn, _, err := hj.Hijack(); err != nil {
		return err
	} else {
		r.Header.Set("Proxy-Authorization", s.authorization)
		r.Write(proxyConn)
		go transfer(proxyConn, clientConn)
		go transfer(clientConn, proxyConn)
		return nil
	}
}
