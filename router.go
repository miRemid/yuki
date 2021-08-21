package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/miRemid/yuki/docs"
)

//go:embed web/dist
var local embed.FS

func (g *Gateway) frontEnd(ctx *gin.Context) {
	fsys := fs.FS(local)
	contentStatic, _ := fs.Sub(fsys, "web/dist")
	handler := http.FileServer(http.FS(contentStatic))
	handler.ServeHTTP(ctx.Writer, ctx.Request)
}

func (g *Gateway) Router() *gin.Engine {
	var route = gin.New()
	route.Use(gin.Logger())
	route.Use(gin.Recovery())

	// reverse proxy
	route.POST("/", g.checkSignature, g.reverseProxy)

	// swagger page
	ipv4, _ := g.getLocalIP()
	url := ginSwagger.URL(fmt.Sprintf("http://%s%s/swagger/doc.json", ipv4, g.Addr)) // The url pointing to API definition
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	route.NoRoute(g.frontEnd)

	// web api routes
	api := route.Group("/api")
	{

		config := api.Group("/config")
		{
			config.POST("/modify", g.ModifyConfig)
			config.GET("/get", g.GetConfig)
		}

		node := api.Group("/node")
		{
			node.POST("/add", g.AddNode)
			node.POST("/remove", g.DeleteNode)
			node.GET("/getAll", g.GetAllNodes)
			node.POST("/modifySelector", g.ModifySelector)
		}

		rule := api.Group("/rule")
		{
			rule.POST("/add", g.AddRule)
			rule.POST("/remove", g.DelRule)
			rule.POST("/modify", g.ModifyRule)
			rule.GET("/getAll", g.GetRules)
		}
	}
	return route
}
