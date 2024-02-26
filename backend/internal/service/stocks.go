package service

import (
	"errors"

	"day-trading-app/backend/internal/service/models"
)

func (s serviceImpl) CreateStock(stockName string) (models.StockCreated, error) {
	return s.db.CreateStock(stockName)

	//return models.StockCreated{}, errors.New("not implemented")
}

func (s serviceImpl) AddStockToUser(userName string, stockID string, quantity int) error {
	if quantity <= 0 {
		return errors.New("invalid quantity")
	}
	// check if stockID exists

	// check if user already has that stock, if so, update the quantity

	// if not, add the stock to the user

	return errors.New("not implemented")
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

func (s serviceImpl) PlaceStockOrder(userName string, stockID string, isBuy bool, orderType string, quantity int, price float32) error {

	if orderType != "MARKET" && orderType != "LIMIT" {
		return errors.New("invalid order type")
	}

	if quantity <= 0 {
		return errors.New("invalid quantity")
	}

	if price <= 0 {
		return errors.New("invalid price")
	}

	return s.db.PlaceStockOrder(userName, stockID, isBuy, orderType, quantity, price)
}

func (s serviceImpl) CancelStockTransaction(userName string, stockTxID string) error {
	// TODO: confirm that the user is the owner of the transaction

	return s.db.CancelStockTransaction(userName, stockTxID)
}
