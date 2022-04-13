package controller

import (
	"github.com/alob-mtc/wallet-engine/internal/common/constant"
	"github.com/alob-mtc/wallet-engine/internal/common/log"
	"github.com/alob-mtc/wallet-engine/internal/common/util"
	"github.com/alob-mtc/wallet-engine/internal/controller/types"
	"github.com/alob-mtc/wallet-engine/internal/model"
	"github.com/alob-mtc/wallet-engine/internal/response"
	"github.com/alob-mtc/wallet-engine/internal/service"
	serviceType "github.com/alob-mtc/wallet-engine/internal/service/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

type IWalletController interface {
	CreateWallet(ctx *gin.Context)
	DebitWallet(ctx *gin.Context)
	CreditWallet(ctx *gin.Context)
	SetWalletState(ctx *gin.Context)
}

type walletController struct {
	walletService service.IWalletService
}

func NewWalletController(walletService service.IWalletService) IWalletController {
	return &walletController{walletService: walletService}
}

//CreateWallet TODO
func (ctl *walletController) CreateWallet(ctx *gin.Context) {
	var body types.CreateWalletRequest

	logger := log.WithFields(log.FromContext(ctx).Fields)
	requestIdentifier := util.UniqueStringIdentifier()
	logger = logger.WithField(util.RequestIdentifier, requestIdentifier)
	logger.Info("CreateWallet request")

	if err := ctx.ShouldBindJSON(&body); err != nil {
		logger.Error("error while parsing request body: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &response.ErrorResponse{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	//TODO: Transfer service call
	transferReq := serviceType.CreateWalletRequest{
		Currency:   constant.TransactionCurrency(body.Currency),
		CustomerID: body.CustomerId,
	}
	newWallet, err := ctl.walletService.CreateWallet(ctx, logger, transferReq)
	if err != nil {
		logger.Error("error while initiating fund transfer: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &response.ErrorResponse{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	logger.Info("Wallet created successfully")
	ctx.JSON(http.StatusCreated, &response.SuccessResponse{Status: true, Message: "Wallet created successfully", Data: types.CreateWalletResponse{WalletId: newWallet.ID, CustomerId: newWallet.CustomerId, AvailableBalance: newWallet.AvailableBalance, Currency: newWallet.Currency}})
	return

}

//SetWalletState TODO
func (ctl *walletController) SetWalletState(ctx *gin.Context) {

	logger := log.WithFields(log.FromContext(ctx).Fields)
	requestIdentifier := util.UniqueStringIdentifier()
	logger = logger.WithField(util.RequestIdentifier, requestIdentifier)
	logger.Info("SetWalletState request")

	walletId := ctx.Param("walletId")

	if walletId == "" {
		logger.Error("wallet id was not passed")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &response.ErrorResponse{
			Status:  false,
			Message: "Must pass wallet Id",
		})
		return
	}

	query := ctx.Query("status")

	if query == "" {
		logger.Info("wallet status to be set not passed")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &response.ErrorResponse{
			Status:  false,
			Message: "You must provide the wallet status to be set",
		})
		return
	}

	_, err := ctl.walletService.SetWalletState(ctx, logger, walletId, query)
	if err != nil {
		logger.Error("error while initiating fund transfer: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &response.ErrorResponse{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	logger.Info("Request initiated successfully")
	ctx.JSON(http.StatusCreated, &response.SuccessResponse{Status: true, Message: "Request initiated successfully"})
	return

}

//DebitWallet TODO
func (ctl *walletController) DebitWallet(ctx *gin.Context) {
	var body types.InitiateTransactionRequest

	logger := log.WithFields(log.FromContext(ctx).Fields)
	requestIdentifier := util.UniqueStringIdentifier()
	logger = logger.WithField(util.RequestIdentifier, requestIdentifier)
	logger.Info("DebitWallet request")

	if err := ctx.ShouldBindJSON(&body); err != nil {
		logger.Error("error while parsing request body: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &response.ErrorResponse{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	//TODO: Transfer service call
	transferReq := serviceType.PerformTransactionRequest{
		WalletID:   body.WalletId,
		Amount:     body.Amount,
		Narration:  body.Narration,
		Currency:   constant.TransactionCurrency(body.Currency),
		CustomerID: body.CustomerId,
		Entry:      model.DebitEntry,
		Meta:       body.Meta,
	}
	transactionRes, err := ctl.walletService.PerformTransaction(ctx, logger, transferReq)
	if err != nil {
		logger.Error("error while initiating fund transfer: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &response.ErrorResponse{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	if transactionRes.Successful {
		logger.Info("Wallet debit initiated successfully")
		ctx.JSON(http.StatusCreated, &response.SuccessResponse{Status: true, Message: "Wallet debit initiated successfully", Data: types.CreateTransferResponse{WalletId: body.WalletId, Amount: body.Amount, Reference: transactionRes.Ref}})
		return
	}

	logger.Info("Wallet debit not initiated successfully")
	ctx.JSON(http.StatusCreated, &response.SuccessResponse{Status: false, Message: "Failed to initiate Wallet debit"})

}

//CreditWallet TODO
func (ctl *walletController) CreditWallet(ctx *gin.Context) {
	var body types.InitiateTransactionRequest

	logger := log.WithFields(log.FromContext(ctx).Fields)
	requestIdentifier := util.UniqueStringIdentifier()
	logger = logger.WithField(util.RequestIdentifier, requestIdentifier)
	logger.Info("DebitWallet request")

	if err := ctx.ShouldBindJSON(&body); err != nil {
		logger.Error("error while parsing request body: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &response.ErrorResponse{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	//TODO: Transfer service call
	transferReq := serviceType.PerformTransactionRequest{
		WalletID:   body.WalletId,
		Amount:     body.Amount,
		Narration:  body.Narration,
		Currency:   constant.TransactionCurrency(body.Currency),
		CustomerID: body.CustomerId,
		Entry:      model.CreditEntry,
		Meta:       body.Meta,
	}
	transactionRes, err := ctl.walletService.PerformTransaction(ctx, logger, transferReq)
	if err != nil {
		logger.Error("error while initiating fund transfer: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &response.ErrorResponse{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	if transactionRes.Successful {
		logger.Info("Wallet credit initiated successfully")
		ctx.JSON(http.StatusCreated, &response.SuccessResponse{Status: true, Message: "Wallet credit initiated successfully", Data: types.CreateTransferResponse{WalletId: body.WalletId, Amount: body.Amount, Reference: transactionRes.Ref}})
		return
	}

	logger.Info("Wallet credit not initiated successfully")
	ctx.JSON(http.StatusCreated, &response.SuccessResponse{Status: false, Message: "Failed to initiate Wallet credit"})

}
