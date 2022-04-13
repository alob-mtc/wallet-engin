package service

import (
	"context"
	"fmt"
	"github.com/alob-mtc/wallet-engine/internal/common/constant"
	"github.com/alob-mtc/wallet-engine/internal/common/log"
	"github.com/alob-mtc/wallet-engine/internal/model"
	"github.com/alob-mtc/wallet-engine/internal/repository"
	"github.com/alob-mtc/wallet-engine/internal/service/types"
)

type IWalletService interface {
	CreateWallet(ctx context.Context, logger log.Entry, data types.CreateWalletRequest) (*model.Wallet, error)
	SetWalletState(ctx context.Context, logger log.Entry, WalletID string, status string) (bool, error)
	PerformTransaction(ctx context.Context, logger log.Entry, data types.PerformTransactionRequest) (*types.PerformTransactionResponse, error)
}

type walletService struct {
	walletRepository repository.IWalletRepository
	transactionRepo  repository.ITransactionRepository
	walletLedgerRepo repository.IWalletLedgerRepository
	unitOfWork       repository.UnitOfWork
}

// NewWalletService returns a new instance of the wallet service
func NewWalletService(wr repository.IWalletRepository, tr repository.ITransactionRepository, wl repository.IWalletLedgerRepository, ut repository.UnitOfWork) IWalletService {
	return &walletService{
		walletRepository: wr,
		transactionRepo:  tr,
		walletLedgerRepo: wl,
		unitOfWork:       ut,
	}
}

func (ws *walletService) CreateWallet(ctx context.Context, logger log.Entry, data types.CreateWalletRequest) (*model.Wallet, error) {

	if data.Currency != constant.NGN {
		return nil, fmt.Errorf("wallet currency must be %s", constant.NGN)
	}

	newWallet := &model.Wallet{Currency: data.Currency, CustomerId: data.CustomerID}

	if err := ws.walletRepository.CreateWallet(ctx, newWallet); err != nil {
		return nil, err
	}

	return newWallet, nil
}

// SetWalletState TODO
func (ws *walletService) SetWalletState(ctx context.Context, logger log.Entry, WalletID string, status string) (bool, error) {

	// Get the Wallet
	wallet, err := ws.walletRepository.GetWalletById(ctx, WalletID)

	if err != nil {
		return false, err
	}

	if status == "activate" {
		if wallet.Active {
			logger.Info("customer wallet is already active")
			return false, fmt.Errorf("customer wallet is already active")
		}
		wallet.Active = true
	} else if status == "de-activate" {
		if !wallet.Active {
			logger.Info("customer wallet is already de-activated")
			return false, fmt.Errorf("customer wallet is already de-activated")
		}
		wallet.Active = false
	} else {
		logger.Info("got an unsupported status")
		return false, fmt.Errorf("please pass a valid status: activate | de-activate")
	}

	if err := ws.walletRepository.UpdateWalletState(ctx, wallet); err != nil {
		return false, err
	}

	return false, nil
}

//PerformTransaction TODO
func (ws *walletService) PerformTransaction(ctx context.Context, logger log.Entry, data types.PerformTransactionRequest) (*types.PerformTransactionResponse, error) {

	//Begin DB Transaction
	tx, err := ws.unitOfWork.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get the Wallet
	wallet, err := ws.walletRepository.GetWallet(ctx, data.WalletID, data.CustomerID)

	if err != nil {
		return nil, err
	}

	if !wallet.Active {
		logger.Info("customer wallet is not active")
		return nil, fmt.Errorf("customer wallet is not active")
	}

	if wallet.Currency != data.Currency {
		logger.Info("got an unsupported currency")
		return nil, fmt.Errorf("currency %s is not a supported currency on this wallet", data.Currency)
	}

	var PreviousBalance, CurrentBalance int64

	if data.Entry == model.DebitEntry {
		PreviousBalance, CurrentBalance, err = ws.debitBalance(wallet, model.AvailableBalanceType, data.Amount)

		if err != nil {
			return nil, err
		}
	} else if data.Entry == model.CreditEntry {
		PreviousBalance, CurrentBalance, err = ws.creditBalance(wallet, model.AvailableBalanceType, data.Amount)

		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("entry not supported")
	}

	// create transaction
	debitAmount := int64(data.Amount)
	newTransaction := &model.Transaction{
		WalletID:            data.WalletID,
		SourceCurrency:      constant.NGN,
		DestinationCurrency: constant.NGN,
		Status:              model.SuccessStatus,
		Entry:               model.DebitEntry,
		Reason:              data.Narration,
		SourceAmount:        debitAmount,
		DestinationAmount:   debitAmount,
		Fee:                 int64(0),
	}

	if err := ws.transactionRepo.WithTx(tx).Create(ctx, newTransaction); err != nil {
		logger.Error("error while creating wallet transaction: %s", err)
		if err = ws.unitOfWork.Rollback(tx); err != nil {
			logger.Error("error while rolling back: %s", err)
			return nil, err
		}
		return nil, err
	}

	if err := ws.walletRepository.WithTx(tx).UpdateWalletState(ctx, wallet); err != nil {
		logger.Error("error while updating wallet state: %s", err)
		if err = ws.unitOfWork.Rollback(tx); err != nil {
			logger.Error("error while rolling back: %s", err)
			return nil, err
		}
		return nil, err
	}

	wl := &model.WalletLedger{
		TransactionID:   newTransaction.ID,
		WalletID:        wallet.ID,
		Amount:          data.Amount,
		PreviousBalance: PreviousBalance,
		CurrentBalance:  CurrentBalance,
	}

	if err = ws.walletLedgerRepo.WithTx(tx).Create(ctx, wl); err != nil {
		logger.Error("error while creating wallet ledger: %s", err)
		if err = ws.unitOfWork.Rollback(tx); err != nil {
			logger.Error("error while rolling back: %s", err)
			return nil, err
		}
		return nil, err
	}

	if err = ws.unitOfWork.Commit(tx); err != nil {
		return nil, err
	}

	result := &types.PerformTransactionResponse{}
	result.Successful = true
	result.Ref = newTransaction.ID
	result.Fee = newTransaction.Fee
	return result, nil
}

//
func (ws *walletService) debitBalance(wallet *model.Wallet, balanceType model.BalanceType, amount int64) (pb, cb int64, err error) {
	switch balanceType {
	case model.AvailableBalanceType:
		if wallet.AvailableBalance < amount {
			return 0, 0, fmt.Errorf("insufficient balance")
		}
		pb = wallet.AvailableBalance
		wallet.AvailableBalance = wallet.AvailableBalance - amount
		cb = wallet.AvailableBalance
		return pb, cb, nil
	case model.LedgerBalanceType:
		if wallet.LedgerBalance < amount {
			return 0, 0, fmt.Errorf("insufficient balance")
		}
		pb = wallet.LedgerBalance
		wallet.LedgerBalance = wallet.LedgerBalance - amount
		cb = wallet.AvailableBalance
		return pb, cb, nil
	default:
		return 0, 0, fmt.Errorf("unknown balance type selected")
	}
}

//
func (ws *walletService) creditBalance(wallet *model.Wallet, balanceType model.BalanceType, amount int64) (pb, cb int64, err error) {
	switch balanceType {
	case model.AvailableBalanceType:
		pb = wallet.AvailableBalance
		wallet.AvailableBalance = wallet.AvailableBalance + amount
		cb = wallet.AvailableBalance
		return pb, cb, nil
	case model.LedgerBalanceType:
		pb = wallet.LedgerBalance
		wallet.LedgerBalance = wallet.LedgerBalance + amount
		cb = wallet.AvailableBalance
		return pb, cb, nil
	default:
		return 0, 0, fmt.Errorf("unknown balance type selected")
	}
}
