package main

import (
	"fmt"
	"strconv"

	"github.com/alob-mtc/wallet-engine/internal/common/log"
	"github.com/alob-mtc/wallet-engine/internal/config"
	"github.com/alob-mtc/wallet-engine/internal/controller"
	"github.com/alob-mtc/wallet-engine/internal/database"
	"github.com/alob-mtc/wallet-engine/internal/middleware"
	"github.com/alob-mtc/wallet-engine/internal/repository"
	"github.com/alob-mtc/wallet-engine/internal/service"
	wirepay_router "github.com/alob-mtc/wallet-engine/router"
	"github.com/gin-gonic/gin"
)

func main() {
	err := run()
	if err != nil {
		log.Error("error: %v", err)
		return
	}
}

// initializes modules and starts the server
func run() error {
	err := config.Load()
	if err != nil {
		return err
	}

	// Set app log level
	level := log.ParseLevel("info")
	log.SetLevel(level)
	logger := log.New(level, 3)

	log.Info("Starting the WirePay API Service!")

	logger.Error("Connecting to postgres")
	dbConnection := database.ConnectDB(config.Instance.DatabaseURL)
	err = database.MigrateAll(dbConnection)
	if err != nil {
		return err
	}

	defer func() {
		sqlDB, _ := dbConnection.DB()
		err := sqlDB.Close()
		if err != nil {
			log.Error("error: %v", err)
			return
		}
	}()

	// Initialize repositories
	var (
		walletRepository       = repository.NewWalletRepository(dbConnection)
		walletLedgerRepository = repository.NewWalletLedgerRepository(dbConnection)
		transactionRepository  = repository.NewTransactionRepo(dbConnection)
		unitOfWork             = repository.NewGormUnitOfWork(dbConnection)
	)

	// Setup services
	var (
		walletService = service.NewWalletService(walletRepository, transactionRepository, walletLedgerRepository, unitOfWork)
	)

	// Setup controllers
	var (
		generalController = controller.NewGeneralController()
		walletController  = controller.NewWalletController(walletService)
	)

	router := gin.Default()

	// Setup middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	// setting up routes
	wirepay_router.SetUpRoutes(router, generalController, walletController)

	port := 3000

	if config.Instance.Port != nil {
		p, err := strconv.Atoi(*config.Instance.Port)
		if err != nil {
			return fmt.Errorf("error parsing port, must be numeric: %v", err)
		}

		port = p
	}

	if err := router.Run(fmt.Sprintf("0.0.0.0:%d", port)); err != nil {
		log.Error("Could not run infrastructure -> %v", err)
		return err
	}

	return nil
}
