package model

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	config "iteasy.wrappedAnsible/configs"
)

var (
	once          sync.Once
	mongoInstance *dbInstance
	db            *mongo.Database
)

type dbInstance struct {
	client *mongo.Client
}

func getInstance() *dbInstance {
	once.Do(func() {
		fmt.Println("⚙️ Creating mongo Client now.")
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.GlobalConfig.MongoDB.URL))
		if err != nil {
			log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
		}

		// Ping the primary
		if err := client.Ping(context.Background(), nil); err != nil {
			log.Fatalf("❌ Failed to ping MongoDB: %v", err)
		}

		mongoInstance = &dbInstance{
			client: client,
		}
		fmt.Println("✅ Mongo Client created.")
	})
	return mongoInstance
}

func GetMongoInstance() *mongo.Client {
	instance := getInstance()
	if instance != nil {
		return instance.client
	}
	return nil
}

func PingMongoDB(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("❌ Failed to ping MongoDB: %v", err)
	} else {
		fmt.Println("✅ Successfully connected to MongoDB!")
		db = mongoInstance.client.Database(config.GlobalConfig.MongoDB.Database)
	}
}

// cursor
func DecodeCursor[T any](cursor *mongo.Cursor) ([]T, error) {
	defer cursor.Close(context.Background())
	var slice []T
	for cursor.Next(context.Background()) {
		var doc T
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		slice = append(slice, doc)
	}
	return slice, nil
}

// single result
func EvaluateAndDecodeSingleResult[T any](result *mongo.SingleResult) (*T, error) {
	if result == nil {
		return nil, errors.New("result is nil")
	}
	var v T
	if err := result.Decode(&v); err != nil {
		return nil, err
	}
	return &v, nil
}
