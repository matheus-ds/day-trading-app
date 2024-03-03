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

type StockMatch struct {
	Order       StockTransaction `json:"order" bson:"order"`               // original order; though matching engine will change OrderStatus
	QuantityTx  int              `json:"quantityTx" bson:"quantity_tx"`    // quantity actually transacted
	PriceTx     int              `json:"priceTx" bson:"price_tx"`          // price actually transacted
	CostTotalTx int              `json:"costTotalTx" bson:"cost_total_tx"` // total cost transacted; needed for parent tx
	IsParent    bool             `json:"isParent" bson:"is_parent"`        // true if transaction has created a child
	Killed      bool             `json:"killed" bson:"killed"`             // expired or cancelled
}
