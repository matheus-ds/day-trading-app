package store

import (
	"context"
	"fmt"
	"sync"
	"time"

	"day-trading-app/backend/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
			handler = &mongoHandler{}
		} else {
			fmt.Println("Mongo single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return handler
}

func NewTxInterface() mongoHandler {
	_, client, _, _ := ConnectMongoDB(&config.Config{})
	return mongoHandler{
		client: client,
	}
}

func (mh *mongoHandler) BeginMongoTransaction(ctx context.Context, callback func(mongo.SessionContext) (interface{}, error)) (interface{}, error) {
	session, err := mh.client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)
	result, err := session.WithTransaction(ctx, callback)
	if err != nil {
		return nil, err
	}
	return result, err
}

func ConnectMongoDB(cfg *config.Config) (*mongo.Database, *mongo.Client, func() error, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	// got rid of cgf.Mongo.User and cfg.Mongo.Password for now. Will add back in later
	connString := fmt.Sprintf("mongodb://%s:%s/?authSource=admin&readPreference=primary&retryWrites=true&w=majority", cfg.Mongo.Host, cfg.Mongo.Port)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connString).SetMinPoolSize(20).SetHeartbeatInterval(1*time.Second))
	if err != nil {
		return nil, nil, nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, nil, nil, err
	}
	db := client.Database(cfg.Mongo.Dbname)
	disconnect := func() error {
		err = client.Disconnect(ctx)
		if err != nil {
			return err
		}
		return nil
	}
	return db, client, disconnect, nil
}
