package models

type StockCreated struct {
	ID           string `json:"stock_id" bson:"stock_id"`
	StockName    string `json:"-" bson:"stock_name"`
	CurrentPrice int    `json:"-" bson:"current_price"`
}

type PortfolioItem struct {
	StockID   string `json:"stock_id"      bson:"stock_id"`
	StockName string `json:"stock_name"      bson:"stock_name"`
	Quantity  int    `json:"quantity" bson:"quantity"`
}

type StockTransaction struct {
	UserName        string  `json:"user_name"        bson:"user_name"`
	StockTxID       string  `json:"stock_tx_id"       bson:"stock_tx_id"`
	ParentStockTxID *string `json:"parent_stock_tx_id" bson:"parent_stock_tx_id"`
	StockID         string  `json:"stock_id"         bson:"stock_id"`
	WalletTxID      string  `json:"wallet_tx_id"      bson:"wallet_tx_id"`
	OrderStatus     string  `json:"order_status"     bson:"order_status"`
	IsBuy           bool    `json:"is_buy"           bson:"is_buy"`
	OrderType       string  `json:"order_type"       bson:"order_type"`
	StockPrice      int     `json:"stock_price"      bson:"stock_price"`
	Quantity        int     `json:"quantity"        bson:"quantity"`
	TimeStamp       int64   `json:"time_stamp"       bson:"time_stamp"`
}

type StockPrice struct {
	StockID      string `json:"stock_id"      bson:"stock_id"`
	StockName    string `json:"stock_name"    bson:"stock_name"`
	CurrentPrice int    `json:"current_price" bson:"current_price"`
}

type StockMatch struct {
	Order       StockTransaction `json:"order" bson:"order"`                 // original order; though matching engine will change OrderStatus
	QuantityTx  int              `json:"quantity_tx" bson:"quantity_tx"`     // quantity actually transacted
	PriceTx     int              `json:"price_tx" bson:"price_tx"`           // price actually transacted
	CostTotalTx int              `json:"cost_total_tx" bson:"cost_total_tx"` // total cost transacted; needed for parent tx
	IsParent    bool             `json:"is_parent" bson:"is_parent"`         // true if transaction has created a child
	Killed      bool             `json:"killed" bson:"killed"`               // expired or cancelled
}
