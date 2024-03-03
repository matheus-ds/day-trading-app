package store

import (
	"context"
	"day-trading-app/backend/internal/service/models"

	"go.mongodb.org/mongo-driver/bson"
)

// Not Tested due to Register not implemented
func (mh *mongoHandler) RegisterUser(userName, password, name string) error {
	//create user in db
	collection := mh.client.Database("day-trading-app").Collection("users")
	_, err := collection.InsertOne(context.Background(), models.User{UserName: userName, PasswordHash: password, Name: name})
	if err != nil {
		return err
	}
	return nil
}

// Not Tested, No Postman Collection
func (mh *mongoHandler) GetUserByUserName(userName string) (models.User, error) {
	// Access the collection where user data is stored
	collection := mh.client.Database("day-trading-app").Collection("users")

	// Find the user by their username
	var user models.User
	err := collection.FindOne(context.Background(), bson.M{"user_name": userName}).Decode(&user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

// Tested
func (mh *mongoHandler) GetWalletTransactions(userName string) ([]models.WalletTransaction, error) {
	//For testing purposes only:
	//userName = "VanguardETF"
	// Access the collection where user data is stored
	collection := mh.client.Database("day-trading-app").Collection("wallet_transactions")

	// return every transaction in the wallet_transactions collection
	cursor, err := collection.Find(context.Background(), bson.M{"user_name": userName})
	if err != nil {
		return nil, err
	}
	var walletTransactions []models.WalletTransaction
	if err = cursor.All(context.Background(), &walletTransactions); err != nil {
		return nil, err
	}
	return walletTransactions, nil
}

// Tested
func (mh *mongoHandler) GetWalletBalance(userName string) (int, error) {
	//For testing purposes only:
	//userName = "VanguardETF"
	//access the collection where user data is stored
	collection := mh.client.Database("day-trading-app").Collection("users")

	//find the user by their username
	var user models.User
	err := collection.FindOne(context.Background(), bson.M{"user_name": userName}).Decode(&user)
	if err != nil {
		return 0, err
	}
	return user.Balance, nil
}

// Tested
func (mh *mongoHandler) SetWalletBalance(userName string, newBalance int) error {
	//For testing purposes only:
	//userName = "VanguardETF"
	//newBalance = 100000
	// Access the collection where user data is stored
	collection := mh.client.Database("day-trading-app").Collection("users")

	// Find the user by their username
	var user models.User
	err := collection.FindOne(context.Background(), bson.M{"user_name": userName}).Decode(&user)
	if err != nil {
		return err
	}

	// Update the user's balance
	_, err = collection.UpdateOne(context.Background(), bson.M{"user_name": userName}, bson.M{"$set": bson.M{"balance": newBalance}})
	if err != nil {
		return err
	}
	return nil
}

// Not Tested.
func (mh *mongoHandler) AddWalletTransaction(userName string, walletTxID string, stockID string, is_debit bool, amount int, timeStamp int64) error {
	var walletTx models.WalletTransaction = models.WalletTransaction{
		UserName:   userName,
		WalletTxID: walletTxID,
		StockID:    stockID,
		Is_debit:   is_debit,
		Amount:     amount,
		TimeStamp:  timeStamp,
	}

	// Add to 'wallet_transactions' collection
	collection := mh.client.Database("day-trading-app").Collection("wallet_transactions")
	_, err := collection.InsertOne(context.Background(), bson.M{"stock_tx_id": walletTx})
	if err != nil {
		return err
	}

	// * Add to user's entry in 'users' collection *

	// Access the collection where user data is stored
	collection = mh.client.Database("day-trading-app").Collection("users")

	// Find the user by their username
	var user models.User
	err = collection.FindOne(context.Background(), bson.M{"user_name": userName}).Decode(&user)
	if err != nil {
		return err
	}

	// Add to the user's wallet transactions
	//_, err = collection.InsertOne() // todo
	if err != nil {
		return err
	}

	return nil
}

// Not Tested.
func (mh *mongoHandler) DeleteWalletTransaction(userName string, WalletTxID string) error {
	// Remove from 'wallet_transactions' collection
	collection := mh.client.Database("day-trading-app").Collection("wallet_transactions")
	_, err := collection.DeleteOne(context.Background(), bson.M{"stock_tx_id": WalletTxID})
	if err != nil {
		return err
	}

	// Remove from user's entry in 'users' collection

	// Access the collection where user data is stored
	collection = mh.client.Database("day-trading-app").Collection("users")

	// Find the user by their username
	var user models.User
	err = collection.FindOne(context.Background(), bson.M{"user_name": userName}).Decode(&user)
	if err != nil {
		return err
	}

	// Add to the user's wallet transactions
	//_, err = collection.DeleteOne() // todo
	if err != nil {
		return err
	}

	return nil
}
