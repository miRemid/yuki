package main

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
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

// AddRule will add a command rule into the gateway
// and save to the disk
// @Summary Add Command Rule
// @Description Add a command rule
// @Tags Rule
// @Accept json
// @Produce json
// @Param node body main.Rule true "Proxy command rule's cmd and remote address, eg: remote_addr: 127.0.0.1:8081, cmd: help"
// @Success 200 {object} response.Response
// @Router /api/rule/add [post]
func (g *Gateway) AddRule(ctx echo.Context) error {
	var rule = new(Rule)
	if err := ctx.Bind(rule); err != nil {
		return response.BindError(ctx, "add rule failed: bind failed")
	}
	u, t := tools.CheckValidURL(rule.RemoteAddr)
	if !t {
		g.dprintf("%s is an invalid url address", rule.RemoteAddr)
		return response.InvalidURLFormatError(ctx, "add rule failed: invalid remote address")
	}
	rule.RemoteAddr = u.String()
	if rule.Regex != "" {
		reg, err := regexp.Compile(rule.Regex)
		if err != nil {
			g.dprintf("regexp compile error: %v", err)
			return response.RegexpCompileError(ctx, "add rule failed: regexp invalid")
		}
		rule.reg = reg
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	// 1. check rule exist
	g.dprintf("check rule from cache")
	if _, ok := g.rules[rule.Cmd]; ok {
		g.dprintf("%s rule already exist", rule.Cmd)
		return response.AlreadyExisterror(ctx, "add rule failed: already exist")
	}
	// 2. check remote exist
	g.dprintf("check %s proxy node", rule.RemoteAddr)
	if err := g.selector.Check(rule.RemoteAddr); err != nil {
		g.dprintf("proxy node %s not found", rule.RemoteAddr)
		return response.NotExistError(ctx, "add rule failed: proxy node not found")
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
		return response.DatabaseAddError(ctx, "add rule failed: save to the database failed")
	}
	g.dprintf("add %s cmd into the cache", rule.Cmd)
	g.rules[rule.Cmd] = rule
	return response.OK(ctx, "add success", nil)
}

type DelRuleReq struct {
	Cmd string `json:"cmd" form:"cmd" binding:"required"`
}

// DeleteRule will delete a command rule
// @Summary Delete Rule
// @Description Delete a command rule
// @Tags Rule
// @Accept json
// @Produce json
// @Param node body main.DelRuleReq true "Rule's cmd"
// @Success 200 {object} response.Response
// @Router /api/rule/ [delete]
func (g *Gateway) DeleteRule(ctx echo.Context) error {
	var req = new(DelRuleReq)
	if err := ctx.Bind(req); err != nil {
		return response.BindError(ctx, "del rule failed")
	}
	g.mu.Lock()
	defer g.mu.Unlock()

	node := g.rules[req.Cmd]

	// 1. check exists
	if _, ok := g.rules[node.Cmd]; !ok {
		g.dprintf("remove " + node.Cmd + " rule failed, not exist")
		return response.NotExistError(ctx, "delete rule failed: not exist")
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
		return response.DatabaseDelError(ctx, "del rule failed: database delete failed")
	}
	// 3. del map
	delete(g.rules, req.Cmd)
	return response.OK(ctx, "del rule success", nil)
}

// GetRules will return all rules
// @Summary Get all rules
// @Description Get all command rules
// @Tags Rule
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/rule/ [get]
func (g *Gateway) GetRules(ctx echo.Context) error {
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
		return response.DatabaseGetError(ctx, "get rules failed")
	}

	rules := make([]interface{}, 0)
	for key := range g.rules {
		rules = append(rules, g.rules[key])
	}

	return response.OK(ctx, "", gin.H{
		"rules": rules,
		"ex":    data,
	})
}

// ModifyRule will modify a command rule
// @Summary Modify Rule
// @Description Modify a command rule
// @Tags Rule
// @Accept json
// @Produce json
// @Param node body main.Rule true "Rule's struct"
// @Success 200 {object} response.Response
// @Router /api/rule/ [patch]
func (g *Gateway) ModifyRule(ctx echo.Context) error {
	var rule = new(Rule)
	if err := ctx.Bind(rule); err != nil {
		return response.BindError(ctx, "modify rule failed")
	}

	if rule.Regex != "" {
		if r, err := regexp.Compile(rule.Regex); err != nil {
			return response.RegexpCompileError(ctx, "rule regexp compile failed")
		} else {
			rule.reg = r
		}
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	// 1. check exists
	if _, ok := g.rules[rule.Cmd]; !ok {
		g.dprintf("modify " + rule.Cmd + " rule failed, not exist")
		return response.NotExistError(ctx, "modify rule failed")
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
	return response.OK(ctx, "modify rule success", nil)
}
