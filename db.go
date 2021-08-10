package main

import (
	"errors"
	"os"

	"github.com/xujiajun/nutsdb"
)

const (
	database_name = "data"
)

const (
	SYSTEM_BUCKET = "system"
	CQHTTP_KEY    = "cqhttp"
	SECRET_KEY    = "secret"
	ADMINQQ_KEY   = "adminqq"
	PREFIX_KEY    = "prefix"
)

const (
	NODE_BUCKET = "node"
)

const (
	ALREADY_INIT_KEY = "already_init"
)

func (g *Gateway) loadDatabase() (bool, error) {
	g.log("check database file")
	var e bool
	if stat, err := os.Stat(database_name); err != nil {
		g.log("database file not exists")
		// not exists
		e = false
	} else if stat.IsDir() {
		g.log("database file exists")
		e = true
	} else {
		return false, errors.New("'data' is not a nutsdb format")
	}
	opt := nutsdb.DefaultOptions
	opt.Dir = database_name
	db, err := nutsdb.Open(opt)
	if err != nil {
		return e, err
	}
	g.db = db
	return e, nil
}

func (g *Gateway) loadSystemConfig(e bool) error {
	g.systemConfig = g.defaultSystemConfig()
	if e {
		g.log("load config from database")
		return g.loadSystemConfigFromDisk()
	} else {
		g.log("using default system config")
		return g.saveConfigToDisk()
	}
}

func (g *Gateway) loadSystemConfigFromDisk() error {
	return g.db.View(func(tx *nutsdb.Tx) error {
		bucket := SYSTEM_BUCKET
		entry, _ := tx.Get(bucket, []byte(CQHTTP_KEY))
		g.systemConfig.CQHTTPAddress = string(entry.Value)
		entry, _ = tx.Get(bucket, []byte(SECRET_KEY))
		g.systemConfig.Secret = string(entry.Value)
		entry, _ = tx.Get(bucket, []byte(ADMINQQ_KEY))
		g.systemConfig.AdminQQ = string(entry.Value)
		items, _ := tx.LRange(bucket, []byte(PREFIX_KEY), 0, -1)
		g.systemConfig.Prefix = make([]string, 0)
		for _, item := range items {
			g.systemConfig.Prefix = append(g.systemConfig.Prefix, string(item))
		}
		return nil
	})
}

func (g *Gateway) saveConfigToDisk() error {
	return g.db.Update(func(tx *nutsdb.Tx) (err error) {
		bucket := SYSTEM_BUCKET
		if err = tx.Put(bucket, []byte(CQHTTP_KEY), []byte(g.systemConfig.CQHTTPAddress), 0); err != nil {
			return err
		}
		if err = tx.Put(bucket, []byte(SECRET_KEY), []byte(g.systemConfig.Secret), 0); err != nil {
			return err
		}
		if err = tx.Put(bucket, []byte(ADMINQQ_KEY), []byte(g.systemConfig.Secret), 0); err != nil {
			return err
		}
		tx.Delete(bucket, []byte(PREFIX_KEY))
		for _, p := range g.systemConfig.Prefix {
			if err = tx.LPush(bucket, []byte(PREFIX_KEY), []byte(p)); err != nil {
				return err
			}
		}
		return nil
	})
}
