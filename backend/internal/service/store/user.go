package store

import (
	"day-trading-app/backend/internal/service/models"

	"errors"
)

func (mh *mongoHandler) RegisterUser(userName, password string) error {
	return errors.New("not implemented")
}

func (mh *mongoHandler) GetUserByUserName(userName string) (models.User, error) {
	return models.User{}, errors.New("not implemented")
}

func (mh *mongoHandler) GetWalletTransactions(userName string) ([]models.WalletTransaction, error) {
	return nil, errors.New("not implemented")
}

func (mh *mongoHandler) GetWalletBalance(userName string) (float32, error) {
	return -1, errors.New("not implemented")
}

func (mh *mongoHandler) SetWalletBalance(userName string, newBalance float32) error {
	return errors.New("not implemented")
}
