package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xujiajun/nutsdb"

	"github.com/miRemid/yuki/response"
	"github.com/miRemid/yuki/selector"
	"github.com/miRemid/yuki/tools"
)

func (g *Gateway) loadSelector(e bool) (selector.Selector, error) {
	if !e {
		// set funcName
		g.log("using default nodes selector")
		err := g.db.Update(func(tx *nutsdb.Tx) error {
			return tx.Put(selector.SELECTOR_BUCKET, []byte(selector.SELECTOR_KEY), []byte(selector.RANDOM_SELECTOR), 0)
		})
		return selector.NewRandomSelector(), err
	} else {
		// 1. get selector funciton, and nodes
		g.log("load nodes from database")
		var funcName string
		var nodes = make([]*selector.Node, 0)
		if err := g.db.View(func(tx *nutsdb.Tx) error {
			if entry, err := tx.Get(selector.SELECTOR_BUCKET, []byte(selector.SELECTOR_KEY)); err != nil {
				if err == nutsdb.ErrBucketAndKey(selector.SELECTOR_BUCKET, []byte(selector.SELECTOR_KEY)) {
					funcName = selector.RANDOM_SELECTOR
				} else {
					g.derrorf(err)
					return err
				}
			} else {
				funcName = string(entry.Value)
			}

			entries, err := tx.GetAll(NODE_BUCKET)
			if err != nil {
				if err == nutsdb.ErrBucketEmpty {
					return nil
				}
				return err
			}

			for _, entry := range entries {
				var node = new(selector.Node)
				json.Unmarshal(entry.Value, node)
				nodes = append(nodes, node)
			}

			return nil
		}); err != nil {
			if err == nutsdb.ErrBucketNotFound {
				funcName = selector.RANDOM_SELECTOR
			}
			return nil, err
		}
		return g.resetSelector(funcName, nodes...)
	}
}

func (g *Gateway) resetSelector(funcName string, nodes ...*selector.Node) (selector.Selector, error) {
	g.dprintf("using %s selector", funcName)
	switch funcName {
	case selector.RANDOM_SELECTOR:
		return selector.NewRandomSelector(nodes...), nil
	case selector.ROUND_ROBIN_SELECTOR:
		return selector.NewRoundRobinSelector(nodes...), nil
	case selector.HASH_SELECTOR:
		return selector.NewHashSelector(nodes...), nil
	case selector.WEIGHT_SELECTOR:
		return selector.NewWeightRoundRobinSelector(nodes...), nil
	default:
		g.dprintf("%s is not a valid selector method, using random instead", funcName)
		return selector.NewRandomSelector(nodes...), nil
	}
}

// AddNode will add a proxy node into the gateway
// and save to the disk
// @Summary Add Proxy Node
// @Description Add a proxy remote node into the gateway's selector
// @Tags Selector
// @Accept json
// @Produce json
// @Param node body selector.Node true "Proxy node's remote address, eg: 127.0.0.1:8081"
// @Success 200 {object} response.Response
// @Router /api/node/add [post]
func (g *Gateway) AddNode(ctx *gin.Context) {
	var node = new(selector.Node)
	if err := ctx.ShouldBind(node); err != nil {
		g.dprintf("add proxy node failed: %v", err)
		response.BindError(ctx, "add node failed: binding failed")
		return
	}
	// check remote valid
	u, ok := tools.CheckValidURL(node.RemoteAddr)
	if !ok {
		g.dprintf("%s is an invaild url address", node.RemoteAddr)
		response.InvalidURLFormatError(ctx, "add node failed: invalid remote address")
		return
	}
	node.RemoteAddr = u.String()
	// check remote add exist
	if err := g.selector.Check(node.RemoteAddr); err == nil {
		g.dprintf("%s already exist", node.RemoteAddr)
		response.AlreadyExisterror(ctx, "add node failed: node already exist")
		return
	}
	// save to the disk
	if err := g.db.Update(func(tx *nutsdb.Tx) error {
		data, _ := json.Marshal(node)
		g.dprintf("save %s into the %s", node.RemoteAddr, NODE_BUCKET)
		return tx.Put(NODE_BUCKET, []byte(node.RemoteAddr), data, 0)
	}); err != nil {
		g.dprintf("add proxy node to database failed: %v", err)
		response.DatabaseAddError(ctx, "add node failed: save to the database failed")
		return
	}
	node.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	g.dprintf("add %s node into the selector", node.RemoteAddr)
	g.selector.Add(node)
	response.OK(ctx, "add node success", nil)
}

// GetAllNodes will return all proxy nodes
// @Summary Get all nodes
// @Description Get all proxy nodes and current selector's function
// @Tags Selector
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/node/getAll [get]
func (g *Gateway) GetAllNodes(ctx *gin.Context) {
	nodes, err := g.selector.Getall()
	if err != nil {
		g.dprintf("get all nodes failed: %v", err)
		response.GetError(ctx, "get all nodes failed")
		return
	}
	response.OK(ctx, "", gin.H{
		"method": g.selector.Name(),
		"nodes":  nodes,
	})
}

// DeleteNode will delete a proxy node
// @Summary Delete Node
// @Description Delete a proxy node
// @Tags Selector
// @Accept json
// @Produce json
// @Param node body selector.Node true "Node's remote address"
// @Success 200 {object} response.Response
// @Router /api/node/remove [post]
func (g *Gateway) DeleteNode(ctx *gin.Context) {
	var node selector.Node
	if err := ctx.ShouldBind(&node); err != nil {
		g.dprintf("delete node failed: %v", err)
		response.BindError(ctx, "delete node failed: binding failed")
		return
	}
	_, ok := tools.CheckValidURL(node.RemoteAddr)
	if !ok {
		response.InvalidURLFormatError(ctx, "delete proxy node failed: invalid remote address")
		return
	}
	rules := make([]string, 0)
	if err := g.db.Update(func(tx *nutsdb.Tx) error {
		// 1. get all cmd rule from database
		key := remote_key(node.RemoteAddr)
		g.dprintf("get %s from %s bucket", node.RemoteAddr, RULE_REMOTE_ARRAY_BUCKET)
		cmds, err := tx.GetAll(key)
		if err != nil && err != nutsdb.ErrBucketEmpty {
			return err
		}
		// 2. delete cmd rule
		g.dprintf("delete rule from database")
		for _, cmd := range cmds {
			g.dprintf("delete %s from %s", string(cmd.Key), RULE_BUCKET)
			// delete rule record
			tx.Delete(RULE_BUCKET, cmd.Key)
			g.dprintf("delete %s from %s", string(cmd.Key), key)
			// delete remote_arr record
			tx.Delete(key, cmd.Key)
			rules = append(rules, string(cmd.Key))
		}
		// 3. delete node
		g.dprintf("delete %s from %s", node.RemoteAddr, NODE_BUCKET)
		return tx.Delete(NODE_BUCKET, []byte(node.RemoteAddr))
	}); err != nil {
		g.dprintf("delete proxy node failed: %v", err)
		response.DelError(ctx, "delete node failed")
		return
	}
	// delete map
	g.mu.Lock()
	g.dprintf("delete selector's proxy node")
	g.selector.Delete(node.RemoteAddr)
	g.dprintf("delete proxy node's cache")
	for _, r := range rules {
		g.dprintf("delete %s rule in cache", r)
		delete(g.rules, r)
	}
	g.mu.Unlock()
	ctx.JSON(http.StatusOK, response.Response{
		Code: response.StatusOK,
	})
}

type SelectorFuncName struct {
	FuncName string `json:"func_name" form:"func_name" binding:"required"`
}

// ModifySelector will change gateway's selector method
// @Summary Modify Selector
// @Description Change gateway's selector method, supported: "random", "round_robin", "hash", "weight"
// @Tags Selector
// @Accept json
// @Produce json
// @Param selector body main.SelectorFuncName true "Load balance algorithm name"
// @Success 200 {object} response.Response
// @Router /api/node/modifySelector [post]
func (g *Gateway) ModifySelector(ctx *gin.Context) {
	var fun SelectorFuncName
	if err := ctx.ShouldBind(&fun); err != nil {
		g.dprintf("modify selector failed: %v", err)
		response.BindError(ctx, "modify selector faield: binding error")
		return
	}
	if err := g.db.Update(func(tx *nutsdb.Tx) error {
		nodes, _ := g.selector.Getall()
		s, err := g.resetSelector(fun.FuncName, nodes...)
		if err != nil {
			return err
		}
		g.selector = s
		return tx.Put(selector.SELECTOR_BUCKET, []byte(selector.SELECTOR_KEY), []byte(fun.FuncName), 0)
	}); err != nil {
		g.dprintf("modify selector failed: %v", err)
		response.ModError(ctx, "modify selector failed")
	}
	response.OK(ctx, "modify selector success", nil)
}
