package model

import "github.com/alob-mtc/wallet-engine/internal/common/constant"

type TransactionStatus string
type TransactionEntry string

const (
	PendingStatus TransactionStatus = "PENDING"
	SuccessStatus TransactionStatus = "SUCCESS"
	FailedStatus  TransactionStatus = "FAILED"

	DebitEntry  TransactionEntry = "DEBIT"
	CreditEntry TransactionEntry = "CREDIT"
)

type Transaction struct {
	Base
	WalletID            string                       `json:"wallet_id" gorm:"index;not null"`
	SourceCurrency      constant.TransactionCurrency `json:"source_currency" gorm:"not null"`
	DestinationCurrency constant.TransactionCurrency `json:"destination_currency" gorm:"not null"`
	Status              TransactionStatus            `json:"status" gorm:"index;not null;"`
	Entry               TransactionEntry             `json:"type" gorm:"index;not null;"`
	Reason              string                       `json:"reason" gorm:""`
	SourceAmount        int64                        `json:"source_amount" gorm:"not null"`
	DestinationAmount   int64                        `json:"destination_amount" gorm:"not null"`
	Fee                 int64                        `json:"fee" gorm:"not null;"`
	Meta                JSONMap                      `json:"meta" gorm:"type:json"`
}
