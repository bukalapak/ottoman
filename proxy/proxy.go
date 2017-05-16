package proxy

import (
	"log"
	"net/http"

	"github.com/subosito/httpx"
)

type Transformer interface {
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
	target    Targeter
	Transport http.RoundTripper
	Logger    *log.Logger
}

func NewProxy(target Targeter) *Proxy {
	return &Proxy{target: target}
}

func (p *Proxy) Forward(w http.ResponseWriter, r *http.Request, n Transformer) {
	proxy := &httpx.ReverseProxy{
		Director:       n.Director(p.target),
		Transport:      p.Transport,
		ModifyResponse: n.ModifyResponse,
		ErrorLog:       p.Logger,
	}

	proxy.ServeHTTP(w, r)
}
