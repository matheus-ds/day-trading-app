package store

import (
	"context"
	"fmt"
	"os"
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
			uri := fmt.Sprintf("mongodb://%s:27017", os.Getenv("MONGO_HOST"))
			client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
			// Do we need to disconnect once we are done?
			if err != nil {
				fmt.Println("Error connecting to mongo: ", err)
			}
			handler = &MongoHandler{}
			handler.client = client
			//Test USE ONLY
			//handler.ManageUserWalletBalance("VanguardETF", 1000)
			//handler.ManageUserWalletBalance("lis8fhithvx6s0pjpp", -1000)
			// handler.ManageUserStock("VanguardETF", "appleStockId36714945-db6e-4ed9-92f3-58dcc98b214a", -100)
			// handler.ManageUserStock("VanguardETF", "googleStockId7a592f16-64e0-47d3-8110-03b2ff568337", 100)
			// handler.ManageUserStock("s2v8g00zbmkqbu", "googleStockId7a592f16-64e0-47d3-8110-03b2ff568337", -1)
			// handler.ManageUserStock("s2v8g00zbmkqbu", "appleStockId36714945-db6e-4ed9-92f3-58dcc98b214a", 1)
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
