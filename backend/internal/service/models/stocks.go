package models

type StockCreated struct {
	ID string `json:"stock_id" bson:"stock_id"`
}

type PortfolioItem struct {
	ID            string `json:"stock_id" bson:"stock_id"`
	Name          string `json:"stock_name" bson:"stock_name"`
	QuantityOwned int    `json:"quantity_owned" bson:"quantity_owned"`
}

type StockTransaction struct {
	StockTxID       string  `json:"stock_tx_id" bson:"stock_tx_id"`
	ParentStockTxID *string `json:"parent_stock_tx_id" bson:"parent_stock_tx_id"`
	StockID         string  `json:"stock_id" bson:"stock_id"`
	WalletTxID      string  `json:"wallet_tx_id" bson:"wallet_tx_id"`
	OrderStatus     string  `json:"order_status" bson:"order_status"`
	IsBuy           bool    `json:"is_buy" bson:"is_buy"`
	OrderType       string  `json:"order_type" bson:"order_type"`
	StockPrice      float64 `json:"stock_price" bson:"stock_price"`
	Quantity        int     `json:"quantity" bson:"quantity"`
	TimeStamp       int64   `json:"time_stamp" bson:"time_stamp"`
}

type StockPrice struct {
	StockID      string  `json:"stock_id" bson:"stock_id"`
	StockName    string  `json:"stock_name" bson:"stock_name"`
	CurrentPrice float64 `json:"current_price" bson:"current_price"`
}
