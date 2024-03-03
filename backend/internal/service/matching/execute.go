package matching

import (
	"day-trading-app/backend/internal/service/models"
	"day-trading-app/backend/internal/service/store"
)

var mh = store.GetMongoHandler()

func ExecuteOrders(stockTxCommitQueue []models.StockMatch) {
	for _, tx := range stockTxCommitQueue {
		if isParent(tx) && !tx.Killed {
			// Update stock transaction as completed
			mh.UpdateStockOrderStatus(tx.Order.UserName, tx.Order.StockTxID, tx.Order.OrderStatus)
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
		walletBalance, _ := mh.GetWalletBalance(tx.Order.UserName)
		mh.SetWalletBalance(tx.Order.UserName, walletBalance+deducted)

		// Delete stock transaction
		mh.DeleteStockTransaction(tx.Order.StockTxID)

		// Delete wallet transaction
		mh.DeleteWalletTransaction(tx.Order.UserName, tx.Order.WalletTxID)

	} else if tx.Order.OrderStatus == "PARTIALLY_FULFILLED" {
		// Refund remaining wallet amount
		var deducted = tx.Order.Quantity * tx.Order.StockPrice
		walletBalance, _ := mh.GetWalletBalance(tx.Order.UserName)
		mh.SetWalletBalance(tx.Order.UserName, walletBalance+(deducted-tx.CostTotalTx))

		// Delete wallet transaction
		mh.DeleteWalletTransaction(tx.Order.UserName, tx.Order.WalletTxID)

	} else if tx.Order.OrderStatus == "COMPLETED" {
		// Add stock quantity to user portfolio
		mh.AddStockToUser(tx.Order.UserName, tx.Order.StockID, tx.Order.Quantity)

		// Update stock transaction status to completed
		mh.UpdateStockOrderStatus(tx.Order.UserName, tx.Order.StockTxID, tx.Order.OrderStatus)

	} else if tx.Order.OrderStatus == "" {
		// Error: empty string
	} else {
		// Error: spelling probably
	}
}

func executeSell(tx models.StockMatch) {
	if tx.Order.OrderStatus == "IN_PROGRESS" { // unfulfilled
		// Add stock quantity back to user portfolio
		mh.AddStockToUser(tx.Order.UserName, tx.Order.StockID, tx.Order.Quantity)

		// Delete stock transaction
		mh.DeleteStockTransaction(tx.Order.StockTxID)

	} else if tx.Order.OrderStatus == "PARTIALLY_FULFILLED" {
		// Add remaining stock quantity back to user portfolio (remaining = Order.Quantity - QuantityTx)
		mh.AddStockToUser(tx.Order.UserName, tx.Order.StockID, tx.Order.Quantity-tx.QuantityTx)

	} else if tx.Order.OrderStatus == "COMPLETED" {
		// Update stock transaction to completed
		mh.UpdateStockOrderStatus(tx.Order.UserName, tx.Order.StockTxID, tx.Order.OrderStatus)

		// Add money to wallet
		walletBalance, _ := mh.GetWalletBalance(tx.Order.UserName)
		mh.SetWalletBalance(tx.Order.UserName, walletBalance+tx.CostTotalTx)

		// Insert wallet transaction
		mh.DeleteWalletTransaction(tx.Order.UserName, tx.Order.WalletTxID)

	} else if tx.Order.OrderStatus == "" {
		// Error: empty string
	} else {
		// Error: spelling probably
	}
}

func isParent(tx models.StockMatch) bool {
	return (tx.Order.OrderStatus != "IN_PROGRESS") && (tx.PriceTx == 0)
}
