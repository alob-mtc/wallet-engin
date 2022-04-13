package model

type WalletLedger struct {
	Base
	TransactionID   string           `json:"transaction_id" gorm:"index;not null"`
	WalletID        string           `json:"wallet_id" gorm:"index;not null"`
	Amount          int64            `gorm:"not null;default:0"`
	Entry           TransactionEntry `json:"type" gorm:"index;not null;"`
	PreviousBalance int64            `gorm:"not null"`
	CurrentBalance  int64            `gorm:"not null"`
	Reversal        bool             `gorm:"not null;default:false"`
}
