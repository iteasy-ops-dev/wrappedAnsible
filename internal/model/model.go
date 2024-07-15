package model

// 문제점 및 보완 방법
// 문제점: 초기화 지연 (Lazy Initialization)

// 싱글턴 인스턴스를 처음 사용할 때까지 생성하지 않습니다. 이는 멀티스레드 환경에서 여러 스레드가 동시에 접근할 경우 인스턴스가 여러 번 생성될 수 있는 위험이 있습니다.
// 보완 방법: 이중 검사 잠금 (Double-Checked Locking)

// getInstance 메서드에서 인스턴스가 존재하는지 두 번 검사하여 성능을 최적화하고 스레드 안전성을 유지합니다.
// 첫 번째 검사: 락을 걸지 않고 mongoInstance가 nil인지 확인합니다.
// 두 번째 검사: 락을 걸고 mongoInstance가 nil인지 다시 확인합니다. 이는 첫 번째 검사를 통과한 후 다른 스레드가 인스턴스를 생성하지 않았는지 확인합니다.
import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	config "iteasy.wrappedAnsible/configs"
)

var lock = &sync.Mutex{}

type dbInstance struct {
	client *mongo.Client
}

var mongoInstance *dbInstance
var db *mongo.Database

func getInstance() *dbInstance {
	if mongoInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if mongoInstance == nil {
			fmt.Println("⚙️ Creating mongo Client now.")
			client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.MONGODB_URL))
			if err != nil {
				fmt.Println("❌ Failed to connect to MongoDB:", err)
				return nil
			}

			// Ping the primary
			if err := client.Ping(context.TODO(), nil); err != nil {
				fmt.Println("❌ Failed to ping MongoDB:", err)
				return nil
			}

			mongoInstance = &dbInstance{
				client: client,
			}
		} else {
			fmt.Println("❗️ mongo Client already created.")
		}
	} else {
		fmt.Println("❗️ mongo Client already created.")
	}

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
		panic(fmt.Sprintf("❌ Failed to ping MongoDB: %v", err))
	} else {
		fmt.Println("✅ Successfully connected to MongoDB!")
		db = mongoInstance.client.Database(config.MONGODB_DATABASE)
	}
}
