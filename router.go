package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/miRemid/yuki/docs"
)

func (g *Gateway) Router() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// reverse proxy
	e.POST("/", g.reverseProxy, g.checkSignature())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Static("/", "web/dist")

	api := e.Group("/api")
	{
		config := api.Group("/config")
		{
			config.POST("", g.ModifyConfig)
			config.GET("", g.GetConfig)
		}
		node := api.Group("/node")
		{
			node.POST("", g.AddNode)
			node.DELETE("", g.DeleteNode)
			node.GET("", g.GetAllNodes)
			// update
			node.PATCH("", g.ModifySelector)
		}
		rule := api.Group("/rule")
		{
			rule.POST("", g.AddRule)
			rule.GET("", g.GetRules)
			rule.PATCH("", g.ModifyRule)
			rule.DELETE("", g.DeleteRule)
		}
	}
	data, _ := json.MarshalIndent(e.Routes(), "", " ")
	ioutil.WriteFile("routes.json", data, 0644)
	return e
}
