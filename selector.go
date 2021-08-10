package main

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xujiajun/nutsdb"

	"github.com/miRemid/yuki/response"
	"github.com/miRemid/yuki/selector"
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
	var node selector.Node
	if err := ctx.ShouldBind(&node); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    response.StatusBindError,
			Message: "add node failed",
		})
		return
	}
	if err := g.selector.Check(node.RemoteAddr); err == nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code: response.StatusAlreadyExist,
		})
		return
	}
	g.selector.Add(&node)
	// save to the disk
	if err := g.db.Update(func(tx *nutsdb.Tx) error {
		data, _ := json.Marshal(&node)
		return tx.Put(NODE_BUCKET, []byte(node.RemoteAddr), data, 0)
	}); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    response.StatusSaveDiskError,
			Message: "add node failed",
		})
		return
	}
	ctx.JSON(http.StatusOK, response.Response{
		Code:    response.StatusOK,
		Message: "add node success",
	})
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
		ctx.JSON(http.StatusOK, response.Response{
			Code:    response.StatusGetError,
			Message: "get failed",
		})
		return
	}
	ctx.JSON(http.StatusOK, response.Response{
		Code: response.StatusOK,
		Data: gin.H{
			"selector_name": g.selector.Name(),
			"nodes":         nodes,
		},
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
		ctx.JSON(http.StatusOK, response.Response{
			Code:    response.StatusBindError,
			Message: "del node failed",
		})
		return
	}
	if err := g.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Delete(NODE_BUCKET, []byte(node.RemoteAddr))
	}); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    response.StatusDelError,
			Message: "del node failed",
		})
		return
	}
	if err := g.selector.Delete(node.RemoteAddr); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    response.StatusDelError,
			Message: "del node failed",
		})
		return
	}
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
		ctx.JSON(http.StatusOK, response.Response{
			Code:    response.StatusBindError,
			Message: "modify failed",
		})
		return
	}
	g.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Put(selector.SELECTOR_BUCKET, []byte(selector.SELECTOR_KEY), []byte(fun.FuncName), 0)
	})
	nodes, _ := g.selector.Getall()
	s, err := g.resetSelector(fun.FuncName, nodes...)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    response.StatusBindError,
			Message: "modify failed",
		})
		return
	}
	g.selector = s
	ctx.JSON(http.StatusOK, response.Response{
		Code: response.StatusOK,
	})
}
