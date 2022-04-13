package types

import (
	"github.com/alob-mtc/wallet-engine/internal/common/constant"
	"github.com/alob-mtc/wallet-engine/internal/model"
)

type (
	CreateWalletRequest struct {
		CustomerID string
		Currency   constant.TransactionCurrency
	}

	PerformTransactionRequest struct {
		WalletID   string
		CustomerID string
		Narration  string
		Currency   constant.TransactionCurrency
		Amount     int64
		Entry      model.TransactionEntry
		Meta       interface{}
	}

	PerformTransactionResponse struct {
		Successful bool
		Ref        string
		Fee        int64
	}
)
