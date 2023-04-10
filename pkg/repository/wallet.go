package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
)

// find wallet by userID
func (c *OrderDatabase) FindWalletByUserID(ctx context.Context, userID uint) (wallet domain.Wallet, err error) {
	query := `SELECT * FROM wallets WHERE user_id = $1`

	err = c.DB.Raw(query, userID).Scan(&wallet).Error

	if err != nil {
		return wallet, fmt.Errorf("faild to find wallet for wallet_id %v", userID)
	}

	return wallet, nil
}

// create a new wallet for user
func (c *OrderDatabase) SaveWallet(ctx context.Context, userID uint) (walletID uint, err error) {

	query := `INSERT INTO wallets (user_id,total_amount) VALUES ($1, $2) RETURNING wallet_id`

	var wallet domain.Wallet
	err = c.DB.Raw(query, userID, 0).Scan(&wallet).Error

	if err != nil {
		return walletID, fmt.Errorf("faild to save wallet for user_id %v", userID)
	}

	walletID = wallet.WalletID

	return walletID, nil
}

// trancation type that already defined on domain Debit or Credit
func (c *OrderDatabase) UpdateWallet(ctx context.Context, walletID, amount uint, transactionType domain.TransactionType) error {

	trx := c.DB.Begin()

	// get the wallet total amount
	var totalAmount uint
	query := `SELECT total_amount FROM wallets WHERE wallet_id = $1`
	err := trx.Raw(query, walletID).Scan(&totalAmount).Error

	if err != nil {
		return fmt.Errorf("faild to find total price of wallet of wallert_id %v", walletID)
	}

	// check the transaction type is debit and the amount is less than total amount or not
	if transactionType == domain.Debit && totalAmount < amount {
		trx.Rollback()
		return fmt.Errorf("can't update the total amount total amount in wallet %v lesser than the given amount %v on debit", totalAmount, amount)
	}
	fmt.Println(totalAmount, amount)
	// calculate the total_amount according to trancation_type
	switch transactionType {
	case domain.Credit:
		totalAmount += amount
	case domain.Debit:
		totalAmount -= amount
	default:
		trx.Rollback()
		return fmt.Errorf("invalid transaction type")
	}

	// there is no conflict then update the amount with total_price
	query = `UPDATE wallets SET total_amount = $1 WHERE wallet_id = $2`
	err = trx.Exec(query, totalAmount, walletID).Error

	if err != nil {
		trx.Rollback()
		return fmt.Errorf("faild to update user wallet for wallet_id %v", walletID)
	}

	// update the transaction for wallet
	query = `INSERT INTO transactions (wallet_id,transaction_date,amount,transaction_type) 
	VALUES ($1, $2, $3, $4)`

	transactionDate := time.Now()
	err = trx.Exec(query, walletID, transactionDate, amount, transactionType).Error

	if err != nil {
		trx.Rollback()
		return fmt.Errorf("faild to save wallet transaction for wallet_id %v", walletID)
	}

	// complete the transaction
	err = trx.Commit().Error
	if err != nil {
		return fmt.Errorf("faild to complete the updation of wallet for wallet_id %v", walletID)
	}

	return nil
}

// find wallet transaction history

func (c *OrderDatabase) FindWalletTransactions(ctx context.Context, walletID uint, pagination req.ReqPagination) (transaction []domain.Transaction, err error) {

	limit := pagination.Count
	offset := (pagination.PageNumber - 1) * limit

	query := `SELECT * FROM transactions WHERE wallet_id = $1
	ORDER BY transaction_date DESC LIMIT $2 OFFSET $3`

	err = c.DB.Raw(query, walletID, limit, offset).Scan(&transaction).Error

	if err != nil {
		return transaction, fmt.Errorf("faild get transactions of wallet_id %v", walletID)
	}

	return transaction, nil
}
