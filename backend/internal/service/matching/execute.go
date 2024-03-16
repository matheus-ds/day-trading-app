package matching

import (
	"day-trading-app/backend/internal/service/store"
	"errors"
)

var mh *store.MongoHandler

func ExecuteOrders(txCommitQueue []StockMatch) (err error) {
	mh = store.GetMongoHandler()
	for _, tx := range txCommitQueue {
		if tx.IsParent && !tx.Killed {
			// Update stock transaction status
			err = mh.UpdateStockOrder(tx.Order)
		} else { // child or non-parent
			if tx.Order.IsBuy {
				err = executeBuy(tx)
			} else {
				err = executeSell(tx)
			}
		}
	}
	return err
}

func executeBuy(tx StockMatch) (err error) {
	if tx.Order.OrderStatus == "IN_PROGRESS" { // unfulfilled
		// Add deducted money back to wallet
		var deducted = tx.Order.Quantity * tx.Order.StockPrice
		walletBalance, _ := mh.GetWalletBalance(tx.Order.UserName)
		err = mh.SetWalletBalance(tx.Order.UserName, walletBalance+deducted)
		if err != nil {
			return err
		}

		// Delete stock transaction
		err = mh.DeleteStockTransaction(tx.Order.StockTxID)
		if err != nil {
			return err
		}

		// Delete wallet transaction
		err = mh.DeleteWalletTransaction(tx.Order.UserName, tx.Order.WalletTxID)
		if err != nil {
			return err
		}

	} else if tx.Order.OrderStatus == "PARTIAL_FULFILLED" {
		// Refund remaining wallet amount
		var deducted = tx.Order.Quantity * tx.Order.StockPrice
		walletBalance, _ := mh.GetWalletBalance(tx.Order.UserName)
		err = mh.SetWalletBalance(tx.Order.UserName, walletBalance+(deducted-tx.CostTotalTx))
		if err != nil {
			return err
		}

		// Delete wallet transaction
		err = mh.DeleteWalletTransaction(tx.Order.UserName, tx.Order.WalletTxID)
		if err != nil {
			return err
		}

		if tx.Order.OrderType == "MARKET" && !tx.IsParent {
			tx.Order.StockPrice = tx.PriceTx
			tx.Order.Quantity = tx.QuantityTx
			err = mh.UpdateStockOrder(tx.Order)
			if err != nil {
				return err
			}
		}

	} else if tx.Order.OrderStatus == "COMPLETED" {
		// Add stock quantity to user portfolio
		err = mh.AddStockToUser(tx.Order.UserName, tx.Order.StockID, tx.Order.Quantity)
		if err != nil {
			return err
		}

		if tx.Order.ParentStockTxID != nil { // child
			// Add child tx to db
			_, err = mh.PlaceStockOrder(tx.Order.UserName, tx.Order.StockID, tx.Order.IsBuy, tx.Order.OrderType, tx.Order.Quantity, tx.Order.StockPrice)
			if err != nil {
				return err
			}
			err = mh.UpdateStockOrder(tx.Order) // to update parentID as well. Optimize later.
			if err != nil {
				return err
			}
		} else if tx.IsParent {
			// Update stock transaction status to completed
			err = mh.UpdateStockOrder(tx.Order)
			if err != nil {
				return err
			}
		} else {
			tx.Order.StockPrice = tx.PriceTx
			err = mh.UpdateStockOrder(tx.Order)
			if err != nil {
				return err
			}
		}

	} else if tx.Order.OrderStatus == "" {
		return errors.New("order status is empty string")
	} else {
		return errors.New("order status string is invalid")
	}
	return nil
}

func executeSell(tx StockMatch) (err error) {
	if tx.Order.OrderStatus == "IN_PROGRESS" { // unfulfilled
		// Add stock quantity back to user portfolio
		err = mh.AddStockToUser(tx.Order.UserName, tx.Order.StockID, tx.Order.Quantity)
		if err != nil {
			return err
		}

		// Delete stock transaction
		err = mh.DeleteStockTransaction(tx.Order.StockTxID)
		if err != nil {
			return err
		}

	} else if tx.Order.OrderStatus == "PARTIAL_FULFILLED" {
		// Add remaining stock quantity back to user portfolio (remaining = Order.Quantity - QuantityTx)
		err = mh.AddStockToUser(tx.Order.UserName, tx.Order.StockID, tx.Order.Quantity-tx.QuantityTx)
		if err != nil {
			return err
		}

		if tx.Order.OrderType == "MARKET" && !tx.IsParent {
			tx.Order.StockPrice = tx.PriceTx
			tx.Order.Quantity = tx.QuantityTx
			err = mh.UpdateStockOrder(tx.Order)
			if err != nil {
				return err
			}
		}

	} else if tx.Order.OrderStatus == "COMPLETED" {
		// Add money to wallet
		walletBalance, _ := mh.GetWalletBalance(tx.Order.UserName)
		err = mh.SetWalletBalance(tx.Order.UserName, walletBalance+tx.CostTotalTx)
		if err != nil {
			return err
		}

		// Insert wallet transaction
		err = mh.DeleteWalletTransaction(tx.Order.UserName, tx.Order.WalletTxID)
		if err != nil {
			return err
		}

		if tx.Order.ParentStockTxID != nil { // child
			// Add child tx to db
			_, err = mh.PlaceStockOrder(tx.Order.UserName, tx.Order.StockID, tx.Order.IsBuy, tx.Order.OrderType, tx.Order.Quantity, tx.Order.StockPrice)
			if err != nil {
				return err
			}
			err = mh.UpdateStockOrder(tx.Order) // to update parentID as well. Optimize later.
			if err != nil {
				return err
			}
		} else if tx.IsParent {
			// Update stock transaction status to completed
			err = mh.UpdateStockOrder(tx.Order)
			if err != nil {
				return err
			}
		} else {
			tx.Order.StockPrice = tx.PriceTx
			err = mh.UpdateStockOrder(tx.Order)
			if err != nil {
				return err
			}
		}

	} else if tx.Order.OrderStatus == "" {
		return errors.New("order status is empty string")
	} else {
		return errors.New("order status string is invalid")
	}
	return nil
}
