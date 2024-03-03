package store

import (
	"context"
	"day-trading-app/backend/internal/service/models"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// look up context and what it does
// global context with a timeout of 20 seconds
// var ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)

// Tested
func (mh *mongoHandler) CreateStock(stockName string) (models.StockCreated, error) {
	// generate stock id by appending string "StockId" to stockName while making stockname all lowercase
	// stock_name:"Google", stock_id: <googleStockId>
	stockID := strings.ToLower(stockName) + "StockId" + uuid.New().String()
	stock := models.StockCreated{
		ID:           stockID,
		StockName:    stockName,
		CurrentPrice: 0.0, // initial stock price is 0
	}
	// todo: create stock in db
	collection := mh.client.Database("day-trading-app").Collection("stocks")
	_, err := collection.InsertOne(context.Background(), stock)
	if err != nil {
		return models.StockCreated{}, err
	}
	return stock, nil
}

// Tested
func (mh *mongoHandler) AddStockToUser(userName string, stockID string, quantity int) error {
	collection := mh.client.Database("day-trading-app").Collection("users")

	//test use only:
	//_, err := collection.UpdateOne(context.Background(), bson.M{"user_name": "VanguardETF"}, bson.M{"$push": bson.M{"stocks": bson.M{"stock_id": "googleStockId", "quantity": 550}}})

	//Uncomment this line and comment the above line for production
	_, err := collection.UpdateOne(context.Background(), bson.M{"user_name": userName}, bson.M{"$push": bson.M{"stocks": bson.M{"stock_id": stockID, "quantity": quantity}}})
	if err != nil {
		return err
	}
	return nil
}

// Tested
func (mh *mongoHandler) GetStockPortfolio(userName string) ([]models.PortfolioItem, error) {
	// Access the collection where user portfolio data is stored
	collection := mh.client.Database("day-trading-app").Collection("users")

	// Find the user by their username
	var user models.User
	//test use only:
	//err := collection.FindOne(context.Background(), bson.M{"user_name": "VanguardETF"}).Decode(&user)
	//Uncomment line below and comment the above line for production
	err := collection.FindOne(context.Background(), bson.M{"user_name": userName}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user.Stocks, nil
}

// Tested
func (mh *mongoHandler) GetStockTransactions() ([]models.StockTransaction, error) {
	// Access the collection where user portfolio data is stored
	collection := mh.client.Database("day-trading-app").Collection("stock_transactions")

	// Create a cursor for the find operation
	cur, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	// Decode the documents into a slice of StockTransaction
	var transactions []models.StockTransaction
	for cur.Next(context.Background()) {
		var transaction models.StockTransaction
		err := cur.Decode(&transaction)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	// Check if the cursor encountered any errors while iterating
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

// Tested
func (mh *mongoHandler) GetStockPrices() ([]models.StockPrice, error) {
	// Access the collection where stock price data is stored
	collection := mh.client.Database("day-trading-app").Collection("stocks")

	// Create a cursor for the find operation
	cur, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	// Decode the documents into a slice of StockPrice
	var prices []models.StockPrice
	for cur.Next(context.Background()) {
		var price models.StockPrice
		err := cur.Decode(&price)
		if err != nil {
			return nil, err
		}
		prices = append(prices, price)
	}
	// Check if the cursor encountered any errors while iterating
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return prices, nil
}

// Tested
func (mh *mongoHandler) PlaceStockOrder(userName string, stockID string, isBuy bool, orderType string, quantity int, price int) error {
	collection := mh.client.Database("day-trading-app").Collection("stock_transactions")
	// add string "Tx" inbetween stockID's name, for example, "googleStockId" becomes "googleStockTxId"
	index := strings.Index(stockID, "Stock")
	stockTxID := stockID[:index+len("Stock")] + "Tx" + stockID[index+len("Stock"):] + uuid.New().String()
	// replace "StockId" with "WalletTxId" in stockID
	walletTxID := strings.Replace(stockID, "StockId", "WalletTxId", 1)

	// Create a new stock transaction
	// //test use only:
	// transaction := models.StockTransaction{
	// 	UserName:        "VanguardETF",
	// 	StockTxID:       "googleStockTxID",
	// 	ParentStockTxID: nil, // ParentStockTxID is nil for the first transaction but how do we handle it for subsequent transactions?
	// 	StockID:         "googleStockId",
	// 	WalletTxID:      "googleWalletTxId", // WalletTxID
	// 	OrderStatus:     "IN_PROGRESS",      // initial status of the order is "IN_PROGRESS" needs to be updated to "COMPLETED" or "CANCELLED" later
	// 	IsBuy:           isBuy,
	// 	OrderType:       orderType,
	// 	StockPrice:      price,
	// 	Quantity:        quantity,
	// 	TimeStamp:       time.Now().Unix(), // Use the current time as the timestamp
	// }

	//Uncomment this line and comment the above line for production
	transaction := models.StockTransaction{
		UserName:        userName,
		StockTxID:       stockTxID,
		ParentStockTxID: nil, // ParentStockTxID is nil for the first transaction but how do we handle it for subsequent transactions?
		StockID:         stockID,
		WalletTxID:      walletTxID,    // WalletTxID
		OrderStatus:     "IN_PROGRESS", // initial status of the order is "IN_PROGRESS" needs to be updated to "COMPLETED" or "CANCELLED" later
		IsBuy:           isBuy,
		OrderType:       orderType,
		StockPrice:      price,
		Quantity:        quantity,
		TimeStamp:       time.Now().Unix(), // Use the current time as the timestamp
	}

	// Insert the new stock transaction into the collection
	_, err := collection.InsertOne(context.Background(), transaction)
	if err != nil {
		return err
	}
	return nil
}

// NOT TESTED
func (mh *mongoHandler) UpdateStockOrderStatus(userName string, stockTxID string, orderStatus string) error {
	// UpdateStockOrder updates the status of a stock transaction with the given stockTxID to have the status "COMPLETED" or "PARTIALLY_FULFILLED"
	collection := mh.client.Database("day-trading-app").Collection("stock_transactions")
	// Update the stock transaction with the given stockTxID to have the status "COMPLETED" or "PARTIALLY_FULFILLED"
	_, err := collection.UpdateOne(context.Background(), bson.M{"stock_tx_id": stockTxID}, bson.M{"$set": bson.M{"order_status": orderStatus}})
	if err != nil {
		return err
	}
	return nil

}

// TESTED
func (mh *mongoHandler) CancelStockTransaction(userName string, stockTxID string) error {
	collection := mh.client.Database("day-trading-app").Collection("stock_transactions")
	// Update the stock transaction with the given stockTxID to have the status "CANCELLED"
	// need to check first if transaction is IN_PROGRESS or PARTIALLY_FULFILLED and abort if not. Because some transactions might be too late to cancel.
	var transaction models.StockTransaction
	err := collection.FindOne(context.Background(), bson.M{"stock_tx_id": stockTxID}).Decode(&transaction)
	if err != nil {
		return err
	}
	if transaction.OrderStatus != "IN_PROGRESS" && transaction.OrderStatus != "PARTIALLY_FULFILLED" {
		_, err := collection.UpdateOne(context.Background(), bson.M{"stock_tx_id": stockTxID}, bson.M{"$set": bson.M{"order_status": "CANCELLED"}})
		if err != nil {
			return err
		}
	} else {
		return errors.New("transaction cannot be cancelled because it is not in progress or partially fulfilled")
	}

	return nil
}

// NOT TESTED.
func (mh *mongoHandler) DeleteStockTransaction(stockTxID string) error {
	collection := mh.client.Database("day-trading-app").Collection("stock_transactions")

	_, err := collection.DeleteOne(context.Background(), bson.M{"stock_tx_id": stockTxID})
	if err != nil {
		return err
	}

	return nil
}
