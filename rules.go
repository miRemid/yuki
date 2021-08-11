package main

import (
	"encoding/json"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/miRemid/yuki/response"
	"github.com/xujiajun/nutsdb"
)

const (
	RULE_BUCKET = "rule_bucket"
)

type Rule struct {
	RemoteAddr string `json:"remote_addr" form:"remote_addr" binding:"required"`
	Cmd        string `json:"cmd" form:"cmd" binding:"required"`
	Regex      string `json:"regex" form:"regex"`
	reg        *regexp.Regexp
}

func NewRule(cmd, regex, remote string) (*Rule, error) {
	var r = new(Rule)
	reg, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}
	r.reg = reg
	r.Cmd = cmd
	r.Regex = regex
	r.RemoteAddr = remote
	return r, nil
}

func (r *Rule) Match(cmd string) bool {
	return r.reg.MatchString(cmd)
}

func (g *Gateway) loadRules(e bool) error {
	var rules = make(map[string]*Rule)
	if e {
		g.log("load rules from database")
		if err := g.db.View(func(tx *nutsdb.Tx) error {
			entries, err := tx.GetAll(RULE_BUCKET)
			if err != nil {
				return err
			}
			for _, entry := range entries {
				var rule = new(Rule)
				json.Unmarshal(entry.Value, rule)
				if rule.Regex != "" {
					rule.reg = regexp.MustCompile(rule.Regex)
				}
				rules[rule.Cmd] = rule
			}
			return nil
		}); err != nil {
			if err != nutsdb.ErrBucketNotFound && err != nutsdb.ErrBucketEmpty {
				return err
			}
		}
	}
	g.rules = rules
	return nil
}

func (g *Gateway) AddRule(ctx *gin.Context) {
	var rule = new(Rule)
	if err := ctx.ShouldBind(rule); err != nil {
		response.BindError(ctx, "add rule failed: bind failed")
		return
	}
	if rule.Regex != "" {
		reg, err := regexp.Compile(rule.Regex)
		if err != nil {
			g.dprintf("regexp compile error: %v", err)
			response.RegexpCompileError(ctx, "add rule failed: regexp invalid")
			return
		}
		rule.reg = reg
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	// 1. check exist
	if _, ok := g.rules[rule.Cmd]; ok {
		response.AlreadyExisterror(ctx, "add rule failed: already exist")
		return
	}
	if err := g.db.Update(func(tx *nutsdb.Tx) error {
		data, err := json.Marshal(rule)
		if err != nil {
			return err
		}
		return tx.Put(RULE_BUCKET, []byte(rule.Cmd), data, 0)
	}); err != nil {
		g.dprintf("save rule to database failed: %v", err)
		response.DatabaseAddError(ctx, "add rule failed: save to the database failed")
		return
	}
	g.rules[rule.Cmd] = rule
	response.OK(ctx, "add success", nil)
}

type DelRuleReq struct {
	Cmd string `json:"cmd" form:"cmd"`
}

func (g *Gateway) DelRule(ctx *gin.Context) {
	var req = new(DelRuleReq)
	if err := ctx.ShouldBind(req); err != nil {
		response.BindError(ctx, "del rule failed")
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()

	// 1. check exists
	if _, ok := g.rules[req.Cmd]; !ok {
		g.dprintf("remove " + req.Cmd + " rule failed, not exist")
		response.NotExistError(ctx, "delete rule failed: not exist")
		return
	}

	// 2. del database
	if err := g.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Delete(RULE_BUCKET, []byte(req.Cmd))
	}); err != nil {
		g.dprintf("delete rule from database failed: %v", err)
		response.DatabaseDelError(ctx, "del rule failed: database delete failed")
		return
	}
	// 3. del map
	delete(g.rules, req.Cmd)
	response.OK(ctx, "del rule success", nil)
}

func (g *Gateway) GetRules(ctx *gin.Context) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	response.OK(ctx, "", g.rules)
}

func (g *Gateway) ModifyRule(ctx *gin.Context) {
	var rule = new(Rule)
	if err := ctx.ShouldBind(rule); err != nil {
		response.BindError(ctx, "modify rule failed")
		return
	}

	if rule.Regex != "" {
		if r, err := regexp.Compile(rule.Regex); err != nil {
			response.RegexpCompileError(ctx, "rule regexp compile failed")
			return
		} else {
			rule.reg = r
		}
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	// 1. check exists
	if _, ok := g.rules[rule.Cmd]; !ok {
		g.dprintf("modify " + rule.Cmd + " rule failed, not exist")
		response.NotExistError(ctx, "modify rule failed")
		return
	}

	// 2. modify database
	if err := g.db.Update(func(tx *nutsdb.Tx) error {
		data, _ := json.Marshal(rule)
		return tx.Put(RULE_BUCKET, []byte(rule.Cmd), data, 0)
	}); err != nil {
		g.dprintf("Modify rule from database failed: %v", err)
		response.DatabaseModError(ctx, "modify rule failed")
	}
	// 3. modify map
	g.rules[rule.Cmd] = rule
	response.OK(ctx, "modify rule success", nil)
}
