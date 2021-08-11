package main

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/miRemid/yuki/response"
	"github.com/miRemid/yuki/tools"
	"github.com/xujiajun/nutsdb"
)

const (
	RULE_BUCKET              = "rule_bucket"
	RULE_REMOTE_ARRAY_BUCKET = "rule_remote_arr_bucket"
)

func remote_key(remote string) string {
	return fmt.Sprintf("%s_bucket", remote)
}

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
	u, t := tools.CheckValidURL(rule.RemoteAddr)
	if !t {
		g.dprintf("%s is an invalid url address", rule.RemoteAddr)
		response.InvalidURLFormatError(ctx, "add rule failed: invalid remote address")
		return
	}
	rule.RemoteAddr = u.String()
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
	// 1. check rule exist
	g.dprintf("check rule from cache")
	if _, ok := g.rules[rule.Cmd]; ok {
		g.dprintf("%s rule already exist", rule.Cmd)
		response.AlreadyExisterror(ctx, "add rule failed: already exist")
		return
	}
	// 2. check remote exist
	g.dprintf("check %s proxy node", rule.RemoteAddr)
	if err := g.selector.Check(rule.RemoteAddr); err != nil {
		g.dprintf("proxy node %s not found", rule.RemoteAddr)
		response.NotExistError(ctx, "add rule failed: proxy node not found")
		return
	}
	if err := g.db.Update(func(tx *nutsdb.Tx) error {
		data, err := json.Marshal(rule)
		if err != nil {
			return err
		}
		// add rule's cmd into the database
		g.dprintf("push cmd %v into the %s in %s", rule.Cmd, rule.RemoteAddr, RULE_REMOTE_ARRAY_BUCKET)
		if err := tx.Put(remote_key(rule.RemoteAddr), []byte(rule.Cmd), []byte(""), 0); err != nil {
			return err
		}
		g.dprintf("put cmd %v into the %s", rule.Cmd, RULE_BUCKET)
		return tx.Put(RULE_BUCKET, []byte(rule.Cmd), data, 0)
	}); err != nil {
		g.dprintf("save rule to database failed: %v", err)
		response.DatabaseAddError(ctx, "add rule failed: save to the database failed")
		return
	}
	g.dprintf("add %s cmd into the cache", rule.Cmd)
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

	node := g.rules[req.Cmd]

	// 1. check exists
	if _, ok := g.rules[node.Cmd]; !ok {
		g.dprintf("remove " + node.Cmd + " rule failed, not exist")
		response.NotExistError(ctx, "delete rule failed: not exist")
		return
	}

	// 2. del database
	if err := g.db.Update(func(tx *nutsdb.Tx) error {
		// 2.1 delete rule's remote arr
		g.dprintf("delete %s cmd from the %s key", node.Cmd, node.RemoteAddr)
		if err := tx.Delete(remote_key(node.RemoteAddr), []byte(node.Cmd)); err != nil {
			return nil
		}
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

	data := make(map[string][]string)

	if err := g.db.View(func(tx *nutsdb.Tx) error {
		// get all remote
		g.dprintf("get nodes from database")
		es, err := tx.GetAll(NODE_BUCKET)
		if err != nil {
			return err
		}
		g.dprintf("get remote_address's cmd rule from database")
		for _, e := range es {
			var d = make([]string, 0)
			var key = remote_key(string(e.Key))
			// e.key is the remote address
			// get the remote's rule
			g.dprintf("get %s node's cmd rules", key)
			ss, err := tx.GetAll(key)
			if err != nil {
				return err
			}
			g.dprintf("find %d rules in %s", len(ss), key)
			for _, s := range ss {
				d = append(d, string(s.Key))
			}
			data[string(e.Key)] = d
		}
		return nil
	}); err != nil {
		g.dprintf("get ex error: %v", err)
	}

	response.OK(ctx, "", gin.H{
		"rules": g.rules,
		"ex":    data,
	})
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
