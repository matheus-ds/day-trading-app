package store

import (
	"context"
	"fmt"
	"time"

	"github.com/matheus-ds/day-trading-app/backend/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoTxInterface struct {
	client *mongo.Client
}

func NewTxInterface(c *mongo.Client) MongoTxInterface {
	return MongoTxInterface{
		client: c,
	}
}

func (tx *MongoTxInterface) BeginMongoTransaction(ctx context.Context, callback func(mongo.SessionContext) (interface{}, error)) (interface{}, error) {
	session, err := tx.client.StartSession()
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
	connString := fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin&readPreference=primary&retryWrites=true&w=majority", cfg.Mongo.User, cfg.Mongo.Password, cfg.Mongo.Host, cfg.Mongo.Port)
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
