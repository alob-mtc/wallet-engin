package model

import (
	"github.com/alob-mtc/wallet-engine/internal/common/constant"
)

type BalanceType string

const (
	AvailableBalanceType BalanceType = "AVAILABLE"
	LedgerBalanceType    BalanceType = "LEDGER"
)

type Wallet struct {
	Base
	CustomerId       string                       `json:"customer_id" gorm:"not null"`
	Currency         constant.TransactionCurrency `json:"currency" gorm:"not null"`
	LedgerBalance    int64                        `json:"ledger_balance" gorm:"not null;default:0"`
	AvailableBalance int64                        `json:"available_balance" gorm:"not null;default:0"`
	Active           bool                         `json:"active" gorm:"not null;default:true"`
	WalletLedger     []WalletLedger               `json:"wallet_ledger" gorm:"foreignkey:WalletID"`
	Transaction      []Transaction                `json:"transaction" gorm:"foreignkey:WalletID"`
}

type WalletActionRequest struct {
	Amount      int64
	WalletID    string
	CustomerID  string
	Entry       TransactionEntry
	BalanceType BalanceType
	Transaction *Transaction
}

type CreditWalletResponse struct {
}
