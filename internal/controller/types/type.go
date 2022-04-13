package types

import "github.com/alob-mtc/wallet-engine/internal/common/constant"

type (
	CreateWalletRequest struct {
		CustomerId string `json:"customer_id" binding:"required"`
		Currency   string `json:"currency" binding:"required"`
	}

	CreateWalletResponse struct {
		WalletId         string                       `json:"wallet_id"`
		CustomerId       string                       `json:"customer_id"`
		AvailableBalance int64                        `json:"available_balance"`
		Currency         constant.TransactionCurrency `json:"currency"`
	}

	InitiateTransactionRequest struct {
		WalletId   string      `json:"wallet_id" binding:"required"`
		CustomerId string      `json:"customer_id" binding:"required"`
		Amount     int64       `json:"amount" binding:"required"`
		Narration  string      `json:"reason"`
		Currency   string      `json:"currency" binding:"required"`
		Meta       interface{} `json:"meta"`
	}

	CreateTransferResponse struct {
		WalletId  string `json:"wallet_id"`
		Amount    int64  `json:"amount"`
		Fee       int64  `json:"fee"`
		Reference string `json:"reference"`
	}
)
