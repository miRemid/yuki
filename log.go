package main

import "log"

func init() {
	log.SetPrefix("[YUKI] ")
}

func (g *Gateway) dprintf(format string, args ...interface{}) {
	if !g.Debug {
		return
	}
	log.Printf(format, args...)
}
