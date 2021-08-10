package main

import (
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	var debug bool
	var port int64

	flag.BoolVar(&debug, "d", false, "debug mode")
	flag.Int64Var(&port, "p", 8080, "server's port")

	flag.Parse()
	g, err := NewGateway(fmt.Sprintf(":%d", port), debug)
	if err != nil {
		panic(err)
	}
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	if err := g.ListenAndServe(); err != nil {
		panic(err)
	}
}
