package service

import (
	"day-trading-app/backend/internal/service/models"
)

type Service interface {
	AuthenticateUser(userName, password string) (string, error)
	RegisterUser(userName, password, name string) error

	// stocks
	CreateStock(stockName string) (models.StockCreated, error)
	AddStockToUser(userName, stockID string, quantity int) error
	GetStockPortfolio(userName string) ([]models.PortfolioItem, error)
	GetStockTransactions(userName string) ([]models.StockTransaction, error)
	GetStockPrices() ([]models.StockPrice, error)
	PlaceStockOrder(userName, stockID string, isBuy bool, orderType string, quantity int, price int) error
	CancelStockTransaction(userName, stockTxID string) error

	// wallet
	AddMoneyToWallet(userName string, amount int) error
	GetWalletBalance(userName string) (int, error)
	GetWalletTransactions(userName string) ([]models.WalletTransaction, error)
}

type Database interface {
	// users
	RegisterUser(userName, password, name string) error
	GetUserByUserName(userName string) (models.User, error)

	// stocks
	CreateStock(stockName string) (models.StockCreated, error)
	AddStockToUser(userName string, stockID string, quantity int) error
	UpdateStockToUser(userName string, stockID string, quantity int) error
	DeleteStockToUser(userName string, stockID string) error
	GetStockQuantityFromUser(userName string, stockID string) (int, error)
	ManageUserStock(userName string, stockID string, quantity int) error
	GetStockPortfolio(userName string) ([]models.PortfolioItem, error)
	GetStockTransactions(userName string) ([]models.StockTransaction, error)
	GetStockPrices() ([]models.StockPrice, error) // get prices from all stocks in the stocks collection
	GetStockPrice(stockID string) (int, error)    // get price of a specific stock
	UpdateStockPrice(stockID string, newPrice int) error
	PlaceStockOrder(userName string, stockID string, isBuy bool, orderType string, quantity int, price int, stockTxID string, walletTxID string) (models.StockTransaction, error)
	UpdateStockOrder(models.StockTransaction) error //for the matching engine to update the status of the order
	CancelStockTransaction(userName string, stockTxID string) error

	// wallet
	SetWalletBalance(userName string, newBalance int) error
	GetWalletBalance(userName string) (int, error)
	ManageUserWalletBalance(userName string, amountToAdd int) error
	GetWalletTransactions(userName string) ([]models.WalletTransaction, error)
	AddWalletTransaction(userName string, walletTxID string, stockTxID string, is_debit bool, amount int, timeStamp int64) error
	DeleteWalletTransaction(userName string, walletTxID string) error
}

type serviceImpl struct {
	db Database
}

func New(db Database) Service {
	return &serviceImpl{
		db: db,
	}
}
