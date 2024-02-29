package store

import (
	"context"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var lock = &sync.Mutex{}

type mongoHandler struct {
	client *mongo.Client
}

var handler *mongoHandler

func GetMongoHandler() *mongoHandler {
	if handler == nil {
		lock.Lock()
		defer lock.Unlock()
		if handler == nil {
			fmt.Println("Creating mongo single instance now.")
			// Please connect.
			client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
			// Do we need to disconnect once we are done?
			if err != nil {
				fmt.Println("Error connecting to mongo: ", err)
			}
			handler = &mongoHandler{}
			handler.client = client
		} else {
			fmt.Println("Mongo single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return handler
}
