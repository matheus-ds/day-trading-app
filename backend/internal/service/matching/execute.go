package matching

import (
	"day-trading-app/backend/internal/service/models"
	"day-trading-app/backend/internal/service/store"
)

var mh = store.GetMongoHandler()

func ExecuteOrders(txCommitQueue []models.StockMatch) {
	for _, tx := range txCommitQueue {
		if tx.IsParent && !tx.Killed {
			// Update stock transaction status
			mh.UpdateStockOrder(tx.Order)
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

		if tx.Order.OrderType == "MARKET" && !tx.IsParent {
			tx.Order.StockPrice = tx.PriceTx
			tx.Order.Quantity = tx.QuantityTx
			mh.UpdateStockOrder(tx.Order)
		}

	} else if tx.Order.OrderStatus == "COMPLETED" {
		// Add stock quantity to user portfolio
		mh.AddStockToUser(tx.Order.UserName, tx.Order.StockID, "", tx.Order.Quantity)

		if tx.Order.ParentStockTxID != nil { // child
			// Add child tx to db
			mh.PlaceStockOrder(tx.Order.UserName, tx.Order.StockID, tx.Order.IsBuy, tx.Order.OrderType, tx.Order.Quantity, tx.Order.StockPrice)
			mh.UpdateStockOrder(tx.Order) // to update parentID as well. Optimize later.
		} else if tx.IsParent {
			// Update stock transaction status to completed
			mh.UpdateStockOrder(tx.Order)
		} else {
			tx.Order.StockPrice = tx.PriceTx
			mh.UpdateStockOrder(tx.Order)
		}

	} else if tx.Order.OrderStatus == "" {
		// Error: empty string
	} else {
		// Error: spelling probably
	}
}

func executeSell(tx models.StockMatch) {
	if tx.Order.OrderStatus == "IN_PROGRESS" { // unfulfilled
		// Add stock quantity back to user portfolio
		mh.AddStockToUser(tx.Order.UserName, tx.Order.StockID, "", tx.Order.Quantity)

		// Delete stock transaction
		mh.DeleteStockTransaction(tx.Order.StockTxID)

	} else if tx.Order.OrderStatus == "PARTIALLY_FULFILLED" {
		// Add remaining stock quantity back to user portfolio (remaining = Order.Quantity - QuantityTx)
		mh.AddStockToUser(tx.Order.UserName, tx.Order.StockID, "", tx.Order.Quantity-tx.QuantityTx)

		if tx.Order.OrderType == "MARKET" && !tx.IsParent {
			tx.Order.StockPrice = tx.PriceTx
			tx.Order.Quantity = tx.QuantityTx
			mh.UpdateStockOrder(tx.Order)
		}

	} else if tx.Order.OrderStatus == "COMPLETED" {
		// Add money to wallet
		walletBalance, _ := mh.GetWalletBalance(tx.Order.UserName)
		mh.SetWalletBalance(tx.Order.UserName, walletBalance+tx.CostTotalTx)

		// Insert wallet transaction
		mh.DeleteWalletTransaction(tx.Order.UserName, tx.Order.WalletTxID)

		if tx.Order.ParentStockTxID != nil { // child
			// Add child tx to db
			mh.PlaceStockOrder(tx.Order.UserName, tx.Order.StockID, tx.Order.IsBuy, tx.Order.OrderType, tx.Order.Quantity, tx.Order.StockPrice)
			mh.UpdateStockOrder(tx.Order) // to update parentID as well. Optimize later.
		} else if tx.IsParent {
			// Update stock transaction status to completed
			mh.UpdateStockOrder(tx.Order)
		} else {
			tx.Order.StockPrice = tx.PriceTx
			mh.UpdateStockOrder(tx.Order)
		}

	} else if tx.Order.OrderStatus == "" {
		// Error: empty string
	} else {
		// Error: spelling probably
	}
}
