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
	ALREADY_INIT_KEY = "already_init"
)

func (g *Gateway) loaddb() error {
	isInit := false
	g.dprintf("check database file")
	if stat, err := os.Stat(database_name); err != nil {
		g.dprintf("database file not exists, using default system config")
		// not exists
	} else if stat.IsDir() {
		g.dprintf("database file exists, load config from database")
		isInit = true
	} else {
		return errors.New("'data' is not a nutsdb format")
	}
	opt := nutsdb.DefaultOptions
	opt.Dir = database_name
	db, err := nutsdb.Open(opt)
	if err != nil {
		return err
	}
	g.db = db
	g.systemConfig = g.defaultSystemConfig()
	if isInit {
		return g.loadSystemConfig()
	} else {
		return g.saveConfigToDisk()
	}
}

func (g *Gateway) loadSystemConfig() error {
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
