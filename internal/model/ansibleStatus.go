package model

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
	config "iteasy.wrappedAnsible/configs"
	"iteasy.wrappedAnsible/internal/ansible"
)

type AnsibleProcessStatus struct {
	ID        string                 `bson:"_id,omitempty"`
	Type      string                 `bson:"type"`
	IPs       []string               `bson:"ips"`
	Name      string                 `bson:"name"`
	Account   string                 `bson:"account"`
	Status    bool                   `bson:"status"`
	Payload   string                 `bson:"payload"`
	Duration  int64                  `bson:"duration"`  // Changed to int64 for BSON compatibility
	Timestamp int64                  `bson:"timestamp"` // Unix time field
	Options   map[string]interface{} `bson:"options,omitempty"`
}

func NewAnsibleProcessStatus(a *ansible.AnsibleProcessStatus) *AnsibleProcessStatus {
	return &AnsibleProcessStatus{
		Type:      a.Type,
		IPs:       a.IPs,
		Name:      a.Name,
		Account:   a.Account,
		Status:    a.Status,
		Payload:   a.Payload,
		Duration:  int64(a.Duration),
		Timestamp: time.Now().Unix(),
		Options:   a.Options,
	}
}

func (a *AnsibleProcessStatus) Put() {
	collection := db.Collection(config.COLLECTION_ANSIBLE_PROCESS_STATUS)
	_, err := collection.InsertOne(context.TODO(), a)
	if err != nil {
		log.Fatalf("Failed to insert user: %v", err)
	}
	fmt.Println("User created:", a)
}

func FindByIPs(ips []string, isOr bool, isDesc bool) ([]AnsibleProcessStatus, error) {
	collection := db.Collection(config.COLLECTION_ANSIBLE_PROCESS_STATUS)
	result := []AnsibleProcessStatus{}
	var orderKey string = "timestamp"

	var filter bson.M
	if len(ips) == 0 || !isOr {
		filter = bson.M{"ips": bson.M{"$all": ips}}
	} else {
		filter = bson.M{"ips": bson.M{"$in": ips}}
	}

	var sort bson.M
	if !isDesc {
		// 오름차순
		sort = bson.M{orderKey: 1}
	} else {
		sort = bson.M{orderKey: -1}
	}

	// 쿼리 옵션 설정
	options := options.Find().SetSort(sort)

	// 쿼리 수행
	cursor, err := collection.Find(context.Background(), filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// 결과 처리
	for cursor.Next(context.Background()) {
		var status AnsibleProcessStatus
		if err := cursor.Decode(&status); err != nil {
			return nil, err
		}
		result = append(result, status)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// func FindByIPs(ips []string, isOr bool, isDesc bool) ([]AnsibleProcessStatus, error) {
// 	collection := db.Collection(config.COLLECTION_ANSIBLE_PROCESS_STATUS)
// 	result := []AnsibleProcessStatus{}
// 	var orderKey string = "timestamp"

// 	var filter bson.D
// 	if len(ips) == 0 || !isOr {
// 		filter = bson.D{{Name: "ips", Value: bson.D{{Name: "$all", Value: ips}}}}
// 	} else {
// 		filter = bson.D{{Name: "ips", Value: bson.D{{Name: "$all", Value: ips}}}}
// 	}

// 	var sort bson.D
// 	if !isDesc {
// 		// 오름차순
// 		sort = bson.D{{Name: orderKey, Value: 1}}
// 	} else {
// 		sort = bson.D{{Name: orderKey, Value: -1}}
// 	}

// 	// 쿼리 옵션 설정
// 	options := options.Find().SetSort(sort)

// 	// 쿼리 수행
// 	cursor, err := collection.Find(context.Background(), filter, options)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(context.Background())

// 	// 결과 처리
// 	for cursor.Next(context.Background()) {
// 		var status AnsibleProcessStatus
// 		if err := cursor.Decode(&status); err != nil {
// 			return nil, err
// 		}
// 		result = append(result, status)
// 	}
// 	if err := cursor.Err(); err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }

// func CreateUser(a AnsibleProcessStatus) {
// 	collectionUser := db.Collection("ansible_status")
// 	_, err := collectionUser.InsertOne(context.TODO(), a)
// 	if err != nil {
// 		log.Fatalf("Failed to insert user: %v", err)
// 	}
// 	fmt.Println("User created:", a)
// }

// func ReadUser(id string) *AnsibleProcessStatus {
// 	var a AnsibleProcessStatus
// 	collectionUser := db.Collection("ansible_status")
// 	err := collectionUser.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&a)
// 	if err != nil {
// 		log.Fatalf("Failed to read user: %v", err)
// 	}
// 	return &a
// }

// func UpdateUser(id string, update bson.M) {
// 	collectionUser := db.Collection("ansible_status")
// 	_, err := collectionUser.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": update})
// 	if err != nil {
// 		log.Fatalf("Failed to update user: %v", err)
// 	}
// 	fmt.Println("User updated:", id)
// }

// func DeleteUser(id string) {
// 	collectionUser := db.Collection("ansible_status")
// 	_, err := collectionUser.DeleteOne(context.TODO(), bson.M{"_id": id})
// 	if err != nil {
// 		log.Fatalf("Failed to delete user: %v", err)
// 	}
// 	fmt.Println("User deleted:", id)
// }
