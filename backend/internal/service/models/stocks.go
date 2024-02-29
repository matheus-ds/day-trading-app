package models

type StockCreated struct {
	ID           string `json:"stockId" bson:"stock_id"`
	StockName    string `json:"stockName" bson:"stock_name"`
	CurrentPrice int    `json:"currentPrice" bson:"current_price"`
}

type PortfolioItem struct {
	StockID  string `json:"stockId"      bson:"stock_id"`
	Quantity int    `json:"quantityOwned" bson:"quantity"`
}

type StockTransaction struct {
	UserName        string  `json:"userName"        bson:"user_name"`
	StockTxID       string  `json:"stockTxId"       bson:"stock_tx_id"`
	ParentStockTxID *string `json:"parentStockTxId" bson:"parent_stock_tx_id"`
	StockID         string  `json:"stockId"         bson:"stock_id"`
	WalletTxID      string  `json:"walletTxId"      bson:"wallet_tx_id"`
	OrderStatus     string  `json:"orderStatus"     bson:"order_status"`
	IsBuy           bool    `json:"isBuy"           bson:"is_buy"`
	OrderType       string  `json:"orderType"       bson:"order_type"`
	StockPrice      int     `json:"stockPrice"      bson:"stock_price"`
	Quantity        int     `json:"quantity"        bson:"quantity"`
	TimeStamp       int64   `json:"timeStamp"       bson:"time_stamp"`
}

type StockPrice struct {
	StockID      string `json:"stockId"      bson:"stock_id"`
	StockName    string `json:"stockName"    bson:"stock_name"`
	CurrentPrice int    `json:"currentPrice" bson:"current_price"`
}
