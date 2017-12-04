package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type Transformer interface {
	http.RoundTripper
	Director
	Modifier
}

type Director interface {
	Director(t Targeter) func(*http.Request)
}

type Modifier interface {
	ModifyResponse(*http.Response) error
}

type Forwarder interface {
	Forward(w http.ResponseWriter, r *http.Request, c Transformer)
}

type Proxy struct {
	target        Targeter
	FlushInterval time.Duration
	Logger        *log.Logger
}

func NewProxy(target Targeter) *Proxy {
	return &Proxy{target: target}
}

func (p *Proxy) Target() *url.URL {
	return p.target.Target()
}

func (p *Proxy) Forward(w http.ResponseWriter, r *http.Request, n Transformer) {
	proxy := &httputil.ReverseProxy{
		Director:       n.Director(p.target),
		Transport:      n,
		ModifyResponse: n.ModifyResponse,
		FlushInterval:  p.FlushInterval,
		ErrorLog:       p.Logger,
	}

	proxy.ServeHTTP(w, r)
}
