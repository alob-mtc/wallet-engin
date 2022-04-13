package database

import (
	"fmt"

	"github.com/alob-mtc/wallet-engine/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println(
			err.Error(),
		)
		panic("failed to connect database")
	}

	fmt.Println("Established database connection")

	return db
}

func MigrateAll(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Wallet{},
		&model.Transaction{},
		&model.WalletLedger{},
	)
}
