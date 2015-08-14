package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

type API struct {
  Proxies *ProxyCollection
}

func NewAPI() *API {
	return &API{Proxies: NewProxyCollection()}
}

func (api *API) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		api.RenderSummary(rw)
		return
	}

    params := &ProxyParams{
        Listen: r.FormValue("listen"),
        Target: r.FormValue("target"),
    }

	if !params.ValidateForMethod(r.Method) {
		log.Printf("ERROR: Invalid params: %v", params.Errors)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

    api.UpdateProxy(params, rw)
}

func (api *API) RenderSummary(rw http.ResponseWriter) {
	var buffer bytes.Buffer
	buffer.WriteString("SUMMARY\n\n")

    if len(api.Proxies.Items) == 0 {
        buffer.WriteString("Nothing going on.\n")
    }

	for listen, proxy := range api.Proxies.Items {
		buffer.WriteString(fmt.Sprintf("%v => %v\n", listen, proxy.Target))
	}

	rw.Write(buffer.Bytes())
}

func (api *API) UpdateProxy(params *ProxyParams, rw http.ResponseWriter) {
    var proxy Proxy
    var logAction string

    api.Proxies.Lock()

    if existing, ok := api.Proxies.Items[params.Listen]; ok {
        proxy = existing
        logAction = "UPDATED"
    } else {
        proxy = *NewProxy(params)
        proxy.Start()
        logAction = "CREATED"
    }

    proxy.SetTarget(params.Target)
    api.Proxies.Items[params.Listen] = proxy

    api.Proxies.Unlock()

    log.Printf("%v: %v => %v", logAction, proxy.Listen, proxy.Target)
}
