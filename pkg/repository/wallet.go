package repository

import (
	"context"
	"time"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
)

// find wallet by userID
func (c *OrderDatabase) FindWalletByUserID(ctx context.Context, userID uint) (wallet domain.Wallet, err error) {

	query := `SELECT * FROM wallets WHERE user_id = $1`
	err = c.DB.Raw(query, userID).Scan(&wallet).Error

	return
}

// create a new wallet for user
func (c *OrderDatabase) SaveWallet(ctx context.Context, userID uint) (walletID uint, err error) {

	query := `INSERT INTO wallets (user_id,total_amount) VALUES ($1, $2) RETURNING id`
	err = c.DB.Raw(query, userID, 0).Scan(&walletID).Error

	return
}
func (c *OrderDatabase) UpdateWallet(ctx context.Context, walletID, upateTotalAmount uint) error {

	query := `UPDATE wallets SET total_amount = $1 WHERE id = $2`
	err := c.DB.Exec(query, upateTotalAmount, walletID).Error

	return err
}

func (c *OrderDatabase) SaveWalletTransaction(ctx context.Context, walletTrx domain.Transaction) error {

	trxDate := time.Now()
	query := `INSERT INTO transactions (wallet_id, transaction_date, amount, transaction_type) 
	VALUES ($1, $2, $3, $4)`
	err := c.DB.Exec(query, walletTrx.WalletID, trxDate, walletTrx.Amount, walletTrx.TransactionType).Error

	return err
}

// find wallet transaction history

func (c *OrderDatabase) FindWalletTransactions(ctx context.Context, walletID uint,
	pagination request.Pagination) (transaction []domain.Transaction, err error) {

	limit := pagination.Count
	offset := (pagination.PageNumber - 1) * limit

	query := `SELECT * FROM transactions WHERE wallet_id = $1
	ORDER BY transaction_date DESC LIMIT $2 OFFSET $3`

	err = c.DB.Raw(query, walletID, limit, offset).Scan(&transaction).Error

	return
}
