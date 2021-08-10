package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/miRemid/yuki/response"
)

type SystemConfig struct {
	CQHTTPAddress string   `json:"cqhttp_address" form:"cqhttp_address" binding:"required"`
	Secret        string   `json:"secret" form:"secret" binding:"required"`
	AdminQQ       string   `json:"admin_qq" form:"admin_qq" binding:"required"`
	Prefix        []string `json:"prefix" form:"prefix" binding:"required"`
}

func (g *Gateway) defaultSystemConfig() *SystemConfig {
	var config = new(SystemConfig)
	config.CQHTTPAddress = "127.0.0.1:5600"
	config.AdminQQ = "1234567890"
	config.Prefix = []string{"!"}
	config.Secret = "yuki"
	return config
}

// ModifyConfig
// @Summary Modify config
// @Description Modify gateway's system config, include: "CQHTTP_API_ADDRESS", "CQHTTP_SECRET", "ADMIN_QQ", "CMD PREFIX"
// @Tags config
// @Accept json
// @Produce json
// @Param config body main.SystemConfig true "config struct"
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
	g.systemConfig.Secret = config.Secret
	// save to the database
	if err := g.saveConfigToDisk(); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    response.StatusSaveDiskError,
			Message: "modify success",
		})
		return
	}
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
