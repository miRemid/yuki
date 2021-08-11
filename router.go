package main

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/miRemid/yuki/docs"
)

func (g *Gateway) Router() *gin.Engine {
	var route = gin.New()
	route.Use(gin.Logger())
	route.Use(gin.Recovery())

	// reverse proxy
	route.POST("/", g.checkSignature(), g.reverseProxy)

	// web static page route
	route.Use(static.Serve("/", static.LocalFile("./web/dist", true)))
	// route.GET("/", func(ctx *gin.Context) {
	// 	ctx.String(http.StatusOK, "hello, welcome to the yuki!")
	// })

	// swagger page
	url := ginSwagger.URL("http://192.168.1.106:8080/swagger/doc.json") // The url pointing to API definition
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

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
