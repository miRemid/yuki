package main

import "github.com/gin-gonic/gin"

type Target struct {
	RemoteAddr string // reverse proxy target address
}

// AddTarget will add a proxy node
func (g *Gateway) AddTarget(ctx *gin.Context) {

}
