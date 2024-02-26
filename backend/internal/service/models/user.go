package models

type User struct {
	UserName     string `bson:"user_name"`
	PasswordHash string `bson:"password_hash"`
	Stocks       []PortfolioItem
	Balance      float32 `bson:"balance"`
	WalletTxns   []WalletTransaction
}
