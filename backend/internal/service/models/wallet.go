package models

type WalletTransaction struct {
	StockID      string  `json:"stockId"       bson:"stock_id"`
	StockName    string  `json:"stockName"     bson:"stock_name"`
	CurrentPrice float32 `json:"currentPrice"  bson:"current_price"`
}
