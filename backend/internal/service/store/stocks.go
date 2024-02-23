package store

import (
	"day-trading-app/backend/internal/service/models"

	"errors"
)

func (mh *mongoHandler) CreateStock(stockName string) (models.StockCreated, error) {
	stock := models.StockCreated{}
	// todo: create stock in db
	return stock, nil
}

func (mh *mongoHandler) AddStockToUser(userName string, stockID string, quantity int) error {
	return errors.New("not implemented")
}

func (mh *mongoHandler) GetStockPortfolio(userName string) ([]models.PortfolioItem, error) {
	return nil, errors.New("not implemented")
}

func (mh *mongoHandler) GetStockTransactions() ([]models.StockTransaction, error) {
	return nil, errors.New("not implemented")
}

func (mh *mongoHandler) GetStockPrices() ([]models.StockPrice, error) {
	return nil, errors.New("not implemented")
}

func (mh *mongoHandler) PlaceStockOrder(userName string, stockID string, isBuy bool, orderType string, quantity int, price float32) error {
	return errors.New("not implemented")
}

func (mh *mongoHandler) CancelStockTransaction(userName string, stockTxID string) error {
	return errors.New("not implemented")
}
