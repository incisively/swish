package main

import (
  "sync"
)

type ProxyCollection struct {
  sync.RWMutex
  Items map[string]Proxy
}

func NewProxyCollection() *ProxyCollection {
    return &ProxyCollection{
        Items: make(map[string]Proxy),
    }
}
