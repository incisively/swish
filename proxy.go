package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Proxy struct {
	Listen  string
	Target  string
	Handler *httputil.ReverseProxy
}

func NewProxy(params *ProxyParams) *Proxy {
	p := &Proxy{
		Listen:	params.Listen,
		Target:	params.Target,
	}
	targetURL := p.TargetURL()
	p.Handler = httputil.NewSingleHostReverseProxy(&targetURL)

	return p
}

func (proxy *Proxy) SetTarget(target string) {
	proxy.Target = target

	proxy.Handler.Director = func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = proxy.TargetURL().Host
	}
}

func (proxy *Proxy) Start() {
	go http.ListenAndServe(proxy.Listen, proxy.Handler)
}

func (proxy *Proxy) TargetURL() url.URL {
	t := fmt.Sprintf("http://%v", proxy.Target)

	target, err := url.Parse(t)
	if err != nil {
		panic(err)
	}

	return *target
}
