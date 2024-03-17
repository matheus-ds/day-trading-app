package models

type WalletTransaction struct {
	UserName   string `json:"user_name"     bson:"user_name"`
	WalletTxID string `json:"wallet_tx_id"   bson:"wallet_tx_id"`
	StockTxID  string `json:"stock_tx_id"       bson:"stock_tx_id"`
	Is_debit   bool   `json:"is_debit"       bson:"is_debit"`
	Amount     int    `json:"amount"       bson:"amount"`
	TimeStamp  int64  `json:"time_stamp"    bson:"time_stamp"`
}
