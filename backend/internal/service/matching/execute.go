package matching

import "day-trading-app/backend/internal/service/models"

func ExecuteOrders(stockTxCommitQueue []models.StockMatch) {
	for _, tx := range stockTxCommitQueue {
		if isParent(tx) {
			// Update stock transaction as completed

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

		// Delete stock transaction

		// Delete wallet transaction

	} else if tx.Order.OrderStatus == "PARTIALLY_FULFILLED" {
		// Refund remaining wallet amount

		// Delete wallet transaction

	} else if tx.Order.OrderStatus == "COMPLETED" {
		// Add stock quantity to user portfolio

		// Update stock transaction status to completed

	} else if tx.Order.OrderStatus == "" {
		// Error: empty string
	} else {
		// Error: spelling probably
	}
}

func executeSell(tx models.StockMatch) {
	if tx.Order.OrderStatus == "IN_PROGRESS" { // unfulfilled
		// Add stock quantity back to user portfolio

		// Delete stock transaction

	} else if tx.Order.OrderStatus == "PARTIALLY_FULFILLED" {
		// Add remaining stock quantity back to user portfolio (remaining = Order.Quantity - QuantityTx)

	} else if tx.Order.OrderStatus == "COMPLETED" {
		// Update stock transaction to completed

		// Add money to wallet

		// Insert wallet transaction

	} else if tx.Order.OrderStatus == "" {
		// Error: empty string
	} else {
		// Error: spelling probably
	}
}

func isParent(tx models.StockMatch) bool {
	return (tx.Order.OrderStatus != "IN_PROGRESS") && (tx.PriceTx == 0)
}
