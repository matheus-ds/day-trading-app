package models

type User struct {
	UserName     string `bson:"user_name"`
	PasswordHash string `bson:"password_hash"`
	Stocks       []struct {
		StockID   string `bson:"stock_id"`
		StockName string `bson:"stock_name"`
		Quantity  int    `bson:"quantity"`
	} `bson:"stocks"`
}
