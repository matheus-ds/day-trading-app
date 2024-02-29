package models

type WalletTransaction struct {
	UserName   string `json:"userName"     bson:"user_name"`
	WalletTxID string `json:"walletTxId"   bson:"wallet_tx_id"`
	StockID    string `json:"stockId"       bson:"stock_id"`
	Is_debit   bool   `json:"isDebit"       bson:"is_debit"`
	Amount     int    `json:"amount"       bson:"amount"`
	TimeStamp  int64  `json:"timeStamp"    bson:"time_stamp"`
}
