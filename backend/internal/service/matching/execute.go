package matching

import "day-trading-app/backend/internal/service/models"

var s serviceImpl

func ExecuteOrders(stockTxCommitQueue []models.StockMatch) {
	for _, tx := range stockTxCommitQueue {
		if isParent(tx) {
			// Update stock transaction as completed
			s.db.UpdateStockOrderStatus(tx.Order.UserName, tx.Order.StockTxID, tx.Order.OrderStatus)
		} else { // child or non-parent
			if tx.Order.IsBuy {
				executeBuy(tx)
			} else {
				executeSell(tx)
			}
		}
	}
}

func executeBuy(tx models.StockMatch) {
	if tx.Order.OrderStatus == "IN_PROGRESS" { // unfulfilled
		// Add deducted money back to wallet
		var deducted = tx.Order.Quantity * tx.Order.StockPrice
		walletBalance, _ := s.db.GetWalletBalance(tx.Order.UserName)
		s.db.SetWalletBalance(tx.Order.UserName, walletBalance+deducted)

		// Delete stock transaction todo

		// Delete wallet transaction todo

	} else if tx.Order.OrderStatus == "PARTIALLY_FULFILLED" {
		// Refund remaining wallet amount
		var deducted = tx.Order.Quantity * tx.Order.StockPrice
		walletBalance, _ := s.db.GetWalletBalance(tx.Order.UserName)
		s.db.SetWalletBalance(tx.Order.UserName, walletBalance+(deducted-tx.CostTotalTx))

		// Delete wallet transaction todo

	} else if tx.Order.OrderStatus == "COMPLETED" {
		// Add stock quantity to user portfolio
		s.db.AddStockToUser(tx.Order.UserName, tx.Order.StockID, tx.Order.Quantity)

		// Update stock transaction status to completed
		s.db.UpdateStockOrderStatus(tx.Order.UserName, tx.Order.StockTxID, tx.Order.OrderStatus)

	} else if tx.Order.OrderStatus == "" {
		// Error: empty string
	} else {
		// Error: spelling probably
	}
}

func executeSell(tx models.StockMatch) {
	if tx.Order.OrderStatus == "IN_PROGRESS" { // unfulfilled
		// Add stock quantity back to user portfolio
		s.db.AddStockToUser(tx.Order.UserName, tx.Order.StockID, tx.Order.Quantity)

		// Delete stock transaction // todo

	} else if tx.Order.OrderStatus == "PARTIALLY_FULFILLED" {
		// Add remaining stock quantity back to user portfolio (remaining = Order.Quantity - QuantityTx)
		s.db.AddStockToUser(tx.Order.UserName, tx.Order.StockID, tx.Order.Quantity-tx.QuantityTx)

	} else if tx.Order.OrderStatus == "COMPLETED" {
		// Update stock transaction to completed
		s.db.UpdateStockOrderStatus(tx.Order.UserName, tx.Order.StockTxID, tx.Order.OrderStatus)

		// Add money to wallet
		walletBalance, _ := s.db.GetWalletBalance(tx.Order.UserName)
		s.db.SetWalletBalance(tx.Order.UserName, walletBalance+tx.CostTotalTx)

		// Insert wallet transaction todo

	} else if tx.Order.OrderStatus == "" {
		// Error: empty string
	} else {
		// Error: spelling probably
	}
}

func isParent(tx models.StockMatch) bool {
	return (tx.Order.OrderStatus != "IN_PROGRESS") && (tx.PriceTx == 0)
}

type Database interface {
	// users
	RegisterUser(userName, password, name string) error
	GetUserByUserName(userName string) (models.User, error)

	// stocks
	CreateStock(stockName string) (models.StockCreated, error)
	AddStockToUser(userName string, stockID string, quantity int) error
	GetStockPortfolio(userName string) ([]models.PortfolioItem, error)
	GetStockTransactions() ([]models.StockTransaction, error)
	GetStockPrices() ([]models.StockPrice, error)
	PlaceStockOrder(userName string, stockID string, isBuy bool, orderType string, quantity int, price int) error
	UpdateStockOrderStatus(userName string, stockTxID string, orderStatus string) error //for the matching engine to update the status of the order
	CancelStockTransaction(userName string, stockTxID string) error

	// wallet
	SetWalletBalance(userName string, newBalance int) error
	GetWalletBalance(userName string) (int, error)
	GetWalletTransactions(userName string) ([]models.WalletTransaction, error)
}

type serviceImpl struct {
	db Database
}
