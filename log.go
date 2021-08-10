package main

import (
	"log"
)

func init() {
	log.SetPrefix("[YUKI] ")
}

func (g *Gateway) log(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (g *Gateway) dprintf(format string, args ...interface{}) {
	if !g.Debug {
		return
	}
	log.Printf(format, args...)
}

func (g *Gateway) derrorf(err error) {
	if !g.Debug {
		return
	}
	log.Printf("[ERR] %v", err)
}
