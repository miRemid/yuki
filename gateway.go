package main

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/xujiajun/nutsdb"
)

type Gateway struct {
	mu sync.RWMutex

	Addr string

	nodes        map[string]*Target
	rules        map[string]*Rule
	systemConfig *SystemConfig

	db *nutsdb.DB

	Debug bool
}

// NewGateway returns a pointer of cqhttp-gateway struct
func NewGateway(addr string) (*Gateway, error) {
	var g = new(Gateway)
	g.Addr = addr
	g.nodes = make(map[string]*Target)
	g.rules = make(map[string]*Rule)
	g.mu = sync.RWMutex{}
	g.Debug = true
	// load database
	if err := g.loaddb(); err != nil {
		return nil, err
	}
	return g, nil
}

func (g *Gateway) ListenAndServe() error {
	r := g.Router()
	server := &http.Server{
		Addr:           g.Addr,
		Handler:        r,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("Listening and serving HTTP on %s", g.Addr)
	l, err := net.Listen("tcp", g.Addr)
	if err != nil {
		return err
	}
	return server.Serve(l)
}
