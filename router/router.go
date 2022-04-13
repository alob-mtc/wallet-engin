package router

import (
	"github.com/alob-mtc/wallet-engine/internal/controller"
	"github.com/gin-gonic/gin"
)

func SetUpRoutes(
	router *gin.Engine,
	generalController controller.IGeneralController,
	walletController controller.IWalletController,

) {

	// Setup routes
	router.GET("/", generalController.Welcome)
	router.GET("/health", generalController.Health)

	v1 := router.Group("/api/v1")

	wallet := v1.Group("/wallet")
	wallet.POST("/create", walletController.CreateWallet)
	wallet.POST("/credit", walletController.CreditWallet)
	wallet.POST("/debit", walletController.DebitWallet)
	wallet.GET("/set-status/:walletId", walletController.SetWalletState)

}
