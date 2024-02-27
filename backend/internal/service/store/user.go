package store

import (
	"context"
	"day-trading-app/backend/internal/service/models"

	"go.mongodb.org/mongo-driver/bson"
)

func (mh *mongoHandler) RegisterUser(userName, password string) error {
	//create user in db
	collection := mh.client.Database("day-trading-app").Collection("users")
	_, err := collection.InsertOne(context.Background(), models.User{UserName: userName, PasswordHash: password})
	if err != nil {
		return err
	}
	return nil
}

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

func (mh *mongoHandler) GetWalletTransactions(userName string) ([]models.WalletTransaction, error) {
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

func (mh *mongoHandler) GetWalletBalance(userName string) (float32, error) {
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

func (mh *mongoHandler) SetWalletBalance(userName string, newBalance float32) error {
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
