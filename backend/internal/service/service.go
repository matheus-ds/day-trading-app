package service

import (
	"day-trading-app/backend/internal/service/models"
)

type Service interface {
	AuthenticateUser(userName, password string) (string, error)

	// stocks
	CreateStock(stockName string) (models.StockCreated, error)
	AddStockToUser(userName, stockID string, quantity int) error
	GetStockPortfolio(userName string) ([]models.PortfolioItem, error)
	GetStockTransactions() ([]models.StockTransaction, error)
	GetStockPrices() ([]models.StockPrice, error)
	PlaceStockOrder(userName, stockID string, isBuy bool, orderType string, quantity int, price float32) error
	CancelStockTransaction(userName, stockTxID string) error

	// wallet
	AddMoneyToWallet(userName string, amount float32) error
	GetWalletBalance(userName string) (float32, error)
	GetWalletTransactions(userName string) ([]models.WalletTransaction, error)
}

type Database interface {
	// users
	RegisterUser(userName, password string) error
	GetUserByUserName(userName string) (models.User, error)

	// stocks
	CreateStock(stockName string) (models.StockCreated, error)
	AddStockToUser(userName string, stockID string, quantity int) error
	GetStockPortfolio(userName string) ([]models.PortfolioItem, error)
	GetStockTransactions() ([]models.StockTransaction, error)
	GetStockPrices() ([]models.StockPrice, error)
	PlaceStockOrder(userName string, stockID string, isBuy bool, orderType string, quantity int, price float32) error
	UpdateStockOrder(userName string, stockTxID string, orderStatus string) error //for the matching engine to update the status of the order
	CancelStockTransaction(userName string, stockTxID string) error

	// wallet
	SetWalletBalance(userName string, newBalance float32) error
	GetWalletBalance(userName string) (float32, error)
	GetWalletTransactions(userName string) ([]models.WalletTransaction, error)
}

type serviceImpl struct {
	db Database
}

func New(db Database) Service {
	return &serviceImpl{
		db: db,
	}
}
