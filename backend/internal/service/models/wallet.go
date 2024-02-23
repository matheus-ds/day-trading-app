package models

type WalletTransaction struct {
	StockID      string  `json:"stock_id"       bson:"stock_id"`
	StockName    string  `json:"stock_name"     bson:"stock_name"`
	CurrentPrice float32 `json:"current_price"  bson:"current_price"`
}
