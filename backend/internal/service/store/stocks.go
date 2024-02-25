package store

import (
	"context"
	"day-trading-app/backend/internal/service/models"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (mh *mongoHandler) CreateStock(stockName string) (models.StockCreated, error) {
	// generate stock id by appending string "StockId" to stockName while making stockname all lowercase
	// stock_name:"Google", stock_id: <googleStockId>
	stockID := strings.ToLower(stockName) + "StockId"
	stock := models.StockCreated{
		ID: stockID,
	}
	// todo: create stock in db
	collection := mh.client.Database("day-trading-app").Collection("stocks")
	_, err := collection.InsertOne(context.Background(), stock)
	if err != nil {
		return models.StockCreated{}, err
	}
	return stock, nil
}

func (mh *mongoHandler) AddStockToUser(userName string, stockID string, stockName string, quantity int) error {
	collection := mh.client.Database("day-trading-app").Collection("users")
	_, err := collection.UpdateOne(context.Background(), bson.M{"user_name": userName}, bson.M{"$push": bson.M{"stocks": bson.M{"stock_id": stockID, "stock_name": stockName, "quantity": quantity}}})
	if err != nil {
		return err
	}
	return nil
}

func (mh *mongoHandler) GetStockPortfolio(userName string) ([]models.PortfolioItem, error) {
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

func (mh *mongoHandler) GetStockTransactions() ([]models.StockTransaction, error) {
	// Access the collection where user portfolio data is stored
	collection := mh.client.Database("day-trading-app").Collection("stock_transactions")
	// Create a context for the operation
	ctx := context.Background()

	// Create a cursor for the find operation
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	// Decode the documents into a slice of StockTransaction
	var transactions []models.StockTransaction
	for cur.Next(ctx) {
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

func (mh *mongoHandler) GetStockPrices() ([]models.StockPrice, error) {
	// Access the collection where stock price data is stored
	collection := mh.client.Database("day-trading-app").Collection("stock_prices")

	// Create a context for the operation
	ctx := context.Background()

	// Create a cursor for the find operation
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	// Decode the documents into a slice of StockPrice
	var prices []models.StockPrice
	for cur.Next(ctx) {
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

func (mh *mongoHandler) PlaceStockOrder(userName string, stockID string, isBuy bool, orderType string, quantity int, price float32) error {
	collection := mh.client.Database("day-trading-app").Collection("stock_transactions")

	// add string "Tx" inbetween stockID's name, for example, "googleStockId" becomes "googleStockTxId"
	index := strings.Index(stockID, "Stock")
	stockTxID := stockID[:index+len("Stock")] + "Tx" + stockID[index+len("Stock"):]
	// Create a new stock transaction
	transaction := models.StockTransaction{
		StockTxID:       stockTxID,
		ParentStockTxID: nil, // ParentStockTxID is nil for the first transaction but how do we handle it for subsequent transactions?
		StockID:         stockID,
		OrderStatus:     "IN_PROGRESS", // initial status of the order is "IN_PROGRESS" needs to be updated to "COMPLETED" or "CANCELLED" later
		IsBuy:           isBuy,
		OrderType:       orderType,
		StockPrice:      float64(price),
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

func (mh *mongoHandler) CancelStockTransaction(userName string, stockTxID string) error {
	collection := mh.client.Database("day-trading-app").Collection("stock_transactions")
	// Update the stock transaction with the given stockTxID to have the status "CANCELLED"
	// need to check first if transaction is IN_PROGRESS or PARTIALLY_FULFILLED and abort if not. Because some transactions might be too late to cancel.
	_, err := collection.UpdateOne(context.Background(), bson.M{"stock_tx_id": stockTxID}, bson.M{"$set": bson.M{"order_status": "CANCELLED"}})
	if err != nil {
		return err
	}
	return nil
}
