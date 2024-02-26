package models

type StockCreated struct {
	ID string `json:"stockId" bson:"stock_id"`
}

type PortfolioItem struct {
	ID            string `json:"id" bson:"id"`
	Name          string `json:"name" bson:"name"`
	QuantityOwned int    `json:"quantityOwned" bson:"quantity_owned"`
}

type StockTransaction struct {
	StockTxID       string  `json:"stockTxId"       bson:"stock_tx_id"`
	ParentStockTxID *string `json:"parentStockTxId" bson:"parent_stock_tx_id"`
	StockID         string  `json:"stockId"         bson:"stock_id"`
	WalletTxID      string  `json:"walletTxId"      bson:"wallet_tx_id"`
	OrderStatus     string  `json:"orderStatus"     bson:"order_status"`
	IsBuy           bool    `json:"isBuy"           bson:"is_buy"`
	OrderType       string  `json:"orderType"       bson:"order_type"`
	StockPrice      float64 `json:"stockPrice"      bson:"stock_price"`
	Quantity        int     `json:"quantity"        bson:"quantity"`
	TimeStamp       int64   `json:"timeStamp"       bson:"time_stamp"`
}

type StockPrice struct {
	StockID      string  `json:"stockId"      bson:"stock_id"`
	StockName    string  `json:"stockName"    bson:"stock_name"`
	CurrentPrice float64 `json:"currentPrice" bson:"current_price"`
}
