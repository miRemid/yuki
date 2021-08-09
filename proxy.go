package main

import (
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
)

func (g *Gateway) ReverseProxy(ctx *gin.Context) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = g.systemConfig.CQHTTPAddress
		req.Host = g.systemConfig.CQHTTPAddress
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}
