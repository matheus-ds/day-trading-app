package store

import (
	"context"
	"day-trading-app/backend/internal/service/models"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// look up context and what it does
// global context with a timeout of 20 seconds
// var ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)

// Tested
func (mh *MongoHandler) CreateStock(stockName string) (models.StockCreated, error) {
	// generate stock id by appending string "StockId" to stockName while making stockname all lowercase
	// stock_name:"Google", stock_id: <googleStockId>
	stockID := strings.ToLower(stockName) + "StockId" + uuid.New().String()
	stock := models.StockCreated{
		ID:           stockID,
		StockName:    stockName,
		CurrentPrice: 0, // initial stock price is 0
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
func (mh *MongoHandler) AddStockToUser(userName string, stockID string, quantity int) error {

	collection := mh.client.Database("day-trading-app").Collection("users")

	stockName := strings.Split(cases.Title(language.Make(stockID)).String(stockID), "stockid")[0]

	_, err := collection.UpdateOne(context.Background(), bson.M{"user_name": userName}, bson.M{"$push": bson.M{"stocks": bson.M{"stock_id": stockID, "stock_name": stockName, "quantity": quantity}}})
	if err != nil {
		return err
	}
	return nil
}

// TESTED
func (mh *MongoHandler) UpdateStockToUser(userName string, stockID string, quantity int) error {

	collection := mh.client.Database("day-trading-app").Collection("users")

	_, err := collection.UpdateOne(context.Background(), bson.M{"user_name": userName, "stocks": bson.M{"$elemMatch": bson.M{"stock_id": stockID}}}, bson.M{"$set": bson.M{"stocks.$.quantity": quantity}})
	if err != nil {
		return err
	}
	return nil
}

// TESTED
func (mh *MongoHandler) DeleteStockToUser(userName string, stockID string) error {

	collection := mh.client.Database("day-trading-app").Collection("users")

	_, err := collection.UpdateOne(context.Background(), bson.M{"user_name": userName}, bson.M{"$pull": bson.M{"stocks": bson.M{"stock_id": stockID}}})
	if err != nil {
		return err
	}
	return nil
}

// TESTED
func (mh *MongoHandler) GetStockQuantityFromUser(userName string, stockID string) (int, error) {

	collection := mh.client.Database("day-trading-app").Collection("users")

	var user models.User
	collection.FindOne(context.Background(), bson.M{"user_name": userName}).Decode(&user)
	for _, stock := range user.Stocks {
		if stock.StockID == stockID {
			return stock.Quantity, nil
		}
	}
	return 0, errors.New("stock not found in user's portfolio")
}

// Tested
func (mh *MongoHandler) GetStockPortfolio(userName string) ([]models.PortfolioItem, error) {
	// Access the collection where user portfolio data is stored
	collection := mh.client.Database("day-trading-app").Collection("users")

	// Find the user by their username
	var user models.User

	err := collection.FindOne(context.Background(), bson.M{"user_name": userName}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user.Stocks, nil
}

// Tested
func (mh *MongoHandler) GetStockTransactions(userName string) ([]models.StockTransaction, error) {
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
func (mh *MongoHandler) GetStockPrices() ([]models.StockPrice, error) {
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

func (mh *MongoHandler) GetStockPrice(stockID string) (int, error) {
	collection := mh.client.Database("day-trading-app").Collection("stocks")

	var stock models.StockPrice
	err := collection.FindOne(context.Background(), bson.M{"stock_id": stockID}).Decode(&stock)
	if err != nil {
		return 0, err
	}
	return stock.CurrentPrice, nil
}

func (mh *MongoHandler) UpdateStockPrice(stockID string, newPrice int) error {
	collection := mh.client.Database("day-trading-app").Collection("stocks")

	_, err := collection.UpdateOne(context.Background(), bson.M{"stock_id": stockID}, bson.M{"$set": bson.M{"current_price": newPrice}})
	if err != nil {
		return err
	}
	return nil
}

// Tested
func (mh *MongoHandler) PlaceStockOrder(userName string, stockID string, isBuy bool, orderType string, quantity int, price int, stockTxID string, walletTxID string) (models.StockTransaction, error) {
	collection := mh.client.Database("day-trading-app").Collection("stock_transactions")

	transaction := models.StockTransaction{
		UserName:        userName,
		StockTxID:       stockTxID,
		ParentStockTxID: nil, // ParentStockTxID is nil for the first transaction
		StockID:         stockID,
		WalletTxID:      nil,
		OrderStatus:     "IN_PROGRESS", // initial status of the order is "IN_PROGRESS" needs to be updated to "COMPLETED" or "CANCELLED" later
		IsBuy:           isBuy,
		OrderType:       orderType,
		StockPrice:      price,
		Quantity:        quantity,
		TimeStamp:       time.Now().UnixNano(), // Use the current time as the timestamp
	}

	if walletTxID != "" {
		transaction.WalletTxID = &walletTxID
	}

	// Insert the new stock transaction into the collection
	_, err := collection.InsertOne(context.Background(), transaction)
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}

// Tested
func (mh *MongoHandler) UpdateStockOrder(stockTransaction models.StockTransaction) error {
	collection := mh.client.Database("day-trading-app").Collection("stock_transactions")
	// update the stock transaction by stockTxID and replace it with models.StockTransaction

	_, err := collection.ReplaceOne(context.Background(), bson.M{"stock_tx_id": stockTransaction.StockTxID}, stockTransaction)
	if err != nil {
		return err
	}
	return nil
}

// TESTED
func (mh *MongoHandler) CancelStockTransaction(userName string, stockTxID string) error {
	collection := mh.client.Database("day-trading-app").Collection("stock_transactions")
	// Update the stock transaction with the given stockTxID to have the status "CANCELLED"
	// need to check first if transaction is IN_PROGRESS or PARTIAL_FULFILLED and abort if not. Because some transactions might be too late to cancel.
	var transaction models.StockTransaction
	err := collection.FindOne(context.Background(), bson.M{"stock_tx_id": stockTxID}).Decode(&transaction)
	if err != nil {
		return err
	}
	if transaction.OrderStatus != "IN_PROGRESS" && transaction.OrderStatus != "PARTIAL_FULFILLED" {
		_, err := collection.UpdateOne(context.Background(), bson.M{"stock_tx_id": stockTxID}, bson.M{"$set": bson.M{"order_status": "CANCELLED"}})
		if err != nil {
			return err
		}
	} else {
		return errors.New("transaction cannot be cancelled because it is not in progress or partially fulfilled")
	}

	return nil
}

// TESTED
func (mh *MongoHandler) DeleteStockTransaction(stockTxID string) error {
	collection := mh.client.Database("day-trading-app").Collection("stock_transactions")

	_, err := collection.DeleteOne(context.Background(), bson.M{"stock_tx_id": stockTxID})
	if err != nil {
		return err
	}

	return nil
}
