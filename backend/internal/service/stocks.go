package service

import (
	"errors"

	"day-trading-app/backend/internal/service/models"
)

func (s serviceImpl) CreateStock(stockName string) (models.StockCreated, error) {
	return s.db.CreateStock(stockName)
}

func (s serviceImpl) AddStockToUser(userName string, stockID string, quantity int) error {
	if quantity <= 0 {
		return errors.New("invalid quantity")
	}
	stkExists, err := s.doesStockExist(stockID)
	if err != nil {
		return err
	}
	if !stkExists {
		return errors.New("stock does not exist")
	}

	// check if user already has that stock, if so, update the quantity
	stock, err := s.getStockFromUser(userName, stockID)
	if err != nil {
		// user does not have that stock, add it
		return s.db.AddStockToUser(userName, stockID, quantity)
	}

	return s.db.AddStockToUser(userName, stockID, quantity+stock.Quantity)
}

func (s serviceImpl) GetStockPortfolio(userName string) ([]models.PortfolioItem, error) {
	return s.db.GetStockPortfolio(userName)
}

func (s serviceImpl) GetStockTransactions() ([]models.StockTransaction, error) {
	return s.db.GetStockTransactions()
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

	if price <= 0 {
		return errors.New("invalid price")
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
		stock, err := s.getStockFromUser(userName, stockID)
		if err != nil {
			return err
		}
		if stock.Quantity < quantity {
			return errors.New("insufficient stock quantity")
		}
	}

	if balance != newBalance {
		err = s.db.SetWalletBalance(userName, newBalance)
		if err != nil {
			return err
		}
	}

	// TODO: call matching engine to place the order

	return s.db.PlaceStockOrder(userName, stockID, isBuy, orderType, quantity, price)
}

func (s serviceImpl) CancelStockTransaction(userName string, stockTxID string) error {
	txs, err := s.db.GetStockTransactions()
	if err != nil {
		return err
	}

	for _, item := range txs {
		if item.UserName == userName && item.StockTxID == stockTxID {
			if item.OrderStatus != "COMPLETED" {
				return errors.New("cannot cancel a completed transaction")
			}
			return s.db.CancelStockTransaction(userName, stockTxID)
		}
	}
	return errors.New("stock transaction not found for given user")
}

func (s serviceImpl) doesStockExist(stockID string) (bool, error) {
	stocks, err := s.db.GetStockPrices()
	if err != nil {
		return false, err
	}
	for _, stock := range stocks {
		if stock.StockID == stockID {
			return true, nil
		}
	}
	return false, nil
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
