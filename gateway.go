package main

import (
	"errors"
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
func NewGateway(addr string, debug bool) (*Gateway, error) {
	var g = new(Gateway)
	g.Addr = addr
	g.mu = sync.RWMutex{}
	g.Debug = debug
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
		if err := g.loadRules(ex); err != nil {
			return nil, err
		}
	}
	return g, nil
}

func (g *Gateway) getLocalIP() (ipv4 string, err error) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet
		isIpNet bool
	)
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}
	for _, addr = range addrs {
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String()
				return
			}
		}
	}
	err = errors.New("local ip not found")
	return
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
	g.log("Listening and serving HTTP on %s", g.Addr)
	l, err := net.Listen("tcp", g.Addr)
	if err != nil {
		return err
	}
	return server.Serve(l)
}
