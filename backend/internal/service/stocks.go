package service

import (
	"day-trading-app/backend/internal/service/matching"
	"errors"
	"github.com/google/uuid"
	"strings"
	"time"

	"day-trading-app/backend/internal/service/models"
)

func (s serviceImpl) CreateStock(stockName string) (models.StockCreated, error) {
	return s.db.CreateStock(stockName)
}

func (s serviceImpl) AddStockToUser(userName string, stockID string, quantity int) error {
	if quantity <= 0 {
		return errors.New("invalid quantity")
	}
	_, stkExists, err := s.doesStockExist(stockID)
	if err != nil {
		return err
	}
	if !stkExists {
		return errors.New("stock does not exist")
	}

	return s.db.AddStockToUser(userName, stockID, quantity)
}

func (s serviceImpl) GetStockPortfolio(userName string) ([]models.PortfolioItem, error) {
	return s.db.GetStockPortfolio(userName)
}

func (s serviceImpl) GetStockTransactions(userName string) ([]models.StockTransaction, error) {
	return s.db.GetStockTransactions(userName)
}

func (s serviceImpl) GetStockPrices() ([]models.StockPrice, error) {
	return s.db.GetStockPrices()
}

func (s serviceImpl) PlaceStockOrder(userName string, stockID string, isBuy bool, orderType string, quantity int, price int) error {
	if orderType != "MARKET" && orderType != "LIMIT" {
		return errors.New("invalid order type")
	}

	if quantity <= 0 {
		return errors.New("invalid quantity")
	}

	if price < 0 {
		return errors.New("negative price")
	} else if orderType == "MARKET" && price != 0 {
		return errors.New("market order with non-zero price")
	}

	stockPrice, err := s.db.GetStockPrice(stockID)
	if err != nil {
		return err
	}

	if orderType == "LIMIT" && stockPrice == 0 {
		err = s.db.UpdateStockPrice(stockID, price)
		if err != nil {
			return err
		}
	}

	if orderType == "MARKET" && isBuy {
		price = stockPrice
	}

	balance, err := s.db.GetWalletBalance(userName)
	if err != nil {
		return err
	}
	newBalance := balance
	if isBuy {
		// check if user has enough money to buy
		if balance < quantity*price {
			return errors.New("insufficient balance")
		}
		newBalance = balance - quantity*price
	} else {
		// check if user has enough stock to sell
		userStocksOwned, err := s.db.GetStockQuantityFromUser(userName, stockID)
		if err != nil {
			return err
		}
		if userStocksOwned < quantity {
			return errors.New("insufficient stock quantity")
		} else {
			// deduct stocks from user
			newUserStockQuantity := userStocksOwned - quantity
			if newUserStockQuantity == 0 {
				err = s.db.DeleteStockToUser(userName, stockID)
				if err != nil {
					return err
				}
			} else {
				err = s.db.UpdateStockToUser(userName, stockID, newUserStockQuantity)
				if err != nil {
					return err
				}
			}
		}
	}

	// add string "Tx" inbetween stockID's name, for example, "googleStockId" becomes "googleStockTxId"
	index := strings.Index(stockID, "Stock")
	stockTxID := stockID[:index+len("Stock")] + "Tx" + stockID[index+len("Stock"):] + uuid.New().String()
	walletTxID := ""
	if isBuy {
		walletTxID = strings.Replace(stockTxID, "StockTxId", "WalletTxId", 1)
	}

	if balance != newBalance {
		err = s.db.SetWalletBalance(userName, newBalance)
		if err != nil {
			return err
		}

		err = s.db.AddWalletTransaction(userName, walletTxID, stockTxID, isBuy, quantity*price, time.Now().UnixNano())
		if err != nil {
			return err
		}
	}

	transaction, err := s.db.PlaceStockOrder(userName, stockID, isBuy, orderType, quantity, price, stockTxID, walletTxID)
	err = matching.Match(transaction)
	return err
}

func (s serviceImpl) CancelStockTransaction(userName string, stockTxID string) error {
	txs, err := s.db.GetStockTransactions(userName)
	if err != nil {
		return err
	}

	for _, item := range txs {
		if item.UserName == userName && item.StockTxID == stockTxID {
			if item.OrderType == "MARKET" {
				return errors.New("cannot cancel a market transaction")
			} else if item.OrderStatus == "COMPLETED" {
				return errors.New("cannot cancel a completed transaction")
			} else if time.Now().UnixNano() >= item.TimeStamp+(15*time.Minute).Nanoseconds() {
				return errors.New("cannot cancel an expired transaction")
			} else {
				wasCancelled, err := matching.CancelOrder(item)
				if err != nil {
					return err
				} else if !wasCancelled {
					return errors.New("cannot cancel transaction; was not found in matching engine")
				} else {
					return nil
				}
			}
		}
	}
	return errors.New("stock transaction not found for given user")
}

func (s serviceImpl) doesStockExist(stockID string) (models.StockPrice, bool, error) {
	stk := models.StockPrice{}
	stocks, err := s.db.GetStockPrices()
	if err != nil {
		return stk, false, err
	}
	for _, stock := range stocks {
		if stock.StockID == stockID {
			return stock, true, nil
		}
	}
	return stk, false, nil
}

func (s serviceImpl) getStockFromUser(userName, stockID string) (models.PortfolioItem, error) {
	portfolio, err := s.db.GetStockPortfolio(userName)
	if err != nil {
		return models.PortfolioItem{}, err
	}
	for _, item := range portfolio {
		if item.StockID == stockID {
			return item, nil
		}
	}
	return models.PortfolioItem{}, errors.New("stock not found")
}
