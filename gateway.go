package main

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/miRemid/yuki/selector"
	"github.com/xujiajun/nutsdb"
)

type Gateway struct {
	mu sync.RWMutex

	Addr string

	selector     selector.Selector
	rules        map[string]*Rule
	systemConfig *SystemConfig

	db *nutsdb.DB

	Debug bool
}

// NewGateway returns a pointer of cqhttp-gateway struct
func NewGateway(addr string) (*Gateway, error) {
	var g = new(Gateway)
	g.Addr = addr
	g.rules = make(map[string]*Rule)
	g.mu = sync.RWMutex{}
	g.Debug = true
	// load database
	if ex, err := g.loadDatabase(); err != nil {
		return nil, err
	} else {
		if err := g.loadSystemConfig(ex); err != nil {
			return nil, err
		}
		if s, err := g.loadSelector(ex); err != nil {
			return nil, err
		} else {
			g.selector = s
		}
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
