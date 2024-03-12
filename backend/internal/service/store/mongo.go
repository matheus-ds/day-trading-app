package store

import (
	"context"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var lock = &sync.Mutex{}

type MongoHandler struct {
	client *mongo.Client
}

var handler *MongoHandler

func GetMongoHandler() *MongoHandler {
	if handler == nil {
		lock.Lock()
		defer lock.Unlock()
		if handler == nil {
			fmt.Println("Creating mongo single instance now.")
			// Please connect.
			// host.docker.internal only works in docker container
			// need to come up with a better way to handle this for local development vs production
			client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://host.docker.internal:27017"))
			// Do we need to disconnect once we are done?
			if err != nil {
				fmt.Println("Error connecting to mongo: ", err)
			}
			handler = &MongoHandler{}
			handler.client = client
			//Test USE ONLY
			//handler.AddWalletTransaction("TESTonPOSTMAN_after", "testWalletTxId", "teststockID", true, 888, 8888)
			//handler.DeleteWalletTransaction("TESTonPOSTMAN_after", "testWalletTxId")
		} else {
			fmt.Println("Mongo single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return handler
}
