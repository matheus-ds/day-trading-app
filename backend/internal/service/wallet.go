package service

import (
	"errors"

	"day-trading-app/backend/internal/service/models"
)

func (s serviceImpl) AddMoneyToWallet(userName string, amount int) error {

	if amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	currentBalance, err := s.db.GetWalletBalance(userName)
	if err != nil {
		return err
	}
	newBalance := currentBalance + amount

	return s.db.SetWalletBalance(userName, newBalance)
}

func (s serviceImpl) GetWalletBalance(userName string) (int, error) {
	return s.db.GetWalletBalance(userName)
}

func (s serviceImpl) GetWalletTransactions(userName string) ([]models.WalletTransaction, error) {
	return s.db.GetWalletTransactions(userName)
}
