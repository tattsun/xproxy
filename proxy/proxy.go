package proxy

import "net/http"

type Proxy interface {
	HandleHTTP(w http.ResponseWriter, r *http.Request) error
	HandleHTTPS(w http.ResponseWriter, r *http.Request) error
}
