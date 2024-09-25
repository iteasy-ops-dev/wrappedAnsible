package model

import (
	"context"
	"errors"
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
	/*
		*
		 MongoDB 클라이언트 옵션 설정
			MongoDB Go 드라이버는 내부적으로 커넥션 풀을 자동으로 관리.
			그러나, 이를 적절히 설정하고 최적화하는 것이 중요.
			MongoDB 드라이버에서 커넥션 풀을 설정하고 관리하는 방법 예시.
	*/
	// clientOptions := options.Client().
	// 	ApplyURI(config.GlobalConfig.MongoDB.URL).
	// 	SetMaxPoolSize(100).                // 커넥션 풀 최대 크기 설정
	// 	SetMinPoolSize(10).                 // 커넥션 풀 최소 크기 설정
	// 	SetMaxConnIdleTime(5 * time.Minute) // 커넥션 풀 내 유휴 커넥션 유지 시간

	once.Do(func() {
		log.Println("⚙️ Creating mongo Client now.")
		// client, err := mongo.Connect(context.Background(), clientOptions)
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
		log.Println("✅ Mongo Client created.")
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
		log.Println("✅ Successfully connected to MongoDB!")
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
