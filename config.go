package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miRemid/cqhttp-gateway/response"
)

type SystemConfig struct {
	CQHTTPAddress string   `json:"cqhttp_address" form:"cqhttp_address" binding:"required"`
	AdminQQ       string   `json:"admin_qq" form:"admin_qq" binding:"required"`
	Prefix        []string `json:"prefix" form:"prefix" binding:"required"`
}

func defaultSystemConfig() *SystemConfig {
	var config = new(SystemConfig)
	config.CQHTTPAddress = "127.0.0.1:5600"
	config.AdminQQ = "1234567890"
	config.Prefix = []string{"!"}
	return config
}

// ModifyConfig
// @Summary Modify config
// @Description Modify gateway's system config
// @Tags config
// @Accept json
// @Produce json
// @Param config body main.SystemConfig true "Modify Config"
// @Success 200 {object} response.Response
// @Router /api/config/modify [post]
func (g *Gateway) ModifyConfig(ctx *gin.Context) {
	var config SystemConfig
	if err := ctx.ShouldBindJSON(&config); err != nil {
		log.Println(err)
		ctx.AbortWithStatusJSON(http.StatusOK, response.Response{
			Code:    response.StatusBindError,
			Message: "modify failed",
		})
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	g.systemConfig.AdminQQ = config.AdminQQ
	g.systemConfig.CQHTTPAddress = config.CQHTTPAddress
	g.systemConfig.Prefix = config.Prefix
	ctx.JSON(http.StatusOK, response.Response{
		Code:    response.StatusOK,
		Message: "modify success",
	})
}

// GetConfig
// @Summary Get system config
// @Description Get gateway's system config
// @Tags config
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/config/get [get]
func (g *Gateway) GetConfig(ctx *gin.Context) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	ctx.JSON(http.StatusOK, response.Response{
		Code: response.StatusOK,
		Data: g.systemConfig,
	})
}
