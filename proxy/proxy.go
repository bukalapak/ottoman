package proxy

import (
	"log"
	"net/http"

	"github.com/subosito/httpx"
)

type Director interface {
	Director(t Targeter) func(*http.Request)
}

type Forwarder interface {
	Forward(w http.ResponseWriter, r *http.Request, c Director)
}

type Proxy struct {
	target    Targeter
	Transport http.RoundTripper
	Logger    *log.Logger
}

func NewProxy(target Targeter) *Proxy {
	return &Proxy{target: target}
}

func (p *Proxy) Forward(w http.ResponseWriter, r *http.Request, n Director) {
	proxy := &httpx.ReverseProxy{
		Director:  n.Director(p.target),
		Transport: p.Transport,
		ErrorLog:  p.Logger,
	}

	proxy.ServeHTTP(w, r)
}
