package models

type User struct {
	UserName     string              `bson:"user_name"`
	PasswordHash string              `bson:"password_hash"`
	Name         string              `bson:"name"`
	Balance      int                 `bson:"balance"`
	Stocks       []PortfolioItem     `bson:"stocks"`
	WalletTxns   []WalletTransaction `bson:"wallet_txns"`
}
