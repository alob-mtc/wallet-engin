package controller

import (
	"net/http"
	"time"

	"github.com/alob-mtc/wallet-engine/internal/response"
	"github.com/gin-gonic/gin"
)

type IGeneralController interface {
	Health(c *gin.Context)
	Welcome(c *gin.Context)
}

type generalController struct {
}

func NewGeneralController() IGeneralController {
	return &generalController{}
}

type WelcomeResponse struct {
	Service string    `json:"service"`
	Version string    `json:"version"`
	Date    time.Time `json:"date"`
}

func (ctl *generalController) Health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Up 🚀",
	})
}

func (ctl *generalController) Welcome(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, &response.SuccessResponse{
		Status:  true,
		Message: "Welcome to Opay Wallet Engine",
		Data:    &WelcomeResponse{Service: "opay-wallet-engine ", Version: "1.0.0", Date: time.Now()},
	})
}
