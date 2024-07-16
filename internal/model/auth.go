package model

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"

	config "iteasy.wrappedAnsible/configs"
)

// TODO: 회원 가입 및 로그인에 대한 모델 작업 추가해야함.
type Auth struct {
	ID        string `bson:"_id,omitempty"`
	Email     string `bson:"email"`
	Name      string `bson:"name"`
	Password  string `bson:"password"`
	TimeStamp int64  `bson:"timestamp"`
	IsActive  bool   `bson:"isActive"`
}

func NewFunc() *Auth {
	return &Auth{}
}

func CheckExistingUser(email string) (bool, error) {
	collection := db.Collection(config.COLLECTION_AUTH)
	filter := bson.M{"email": email}
	var existingUser Auth
	err := collection.FindOne(context.Background(), filter).Decode(&existingUser)
	if err == nil {
		// 사용자가 이미 존재함
		return true, nil
	} else if err == mongo.ErrNoDocuments {
		// 사용자가 존재하지 않음
		return false, nil
	} else {
		// 기타 오류
		return false, err
	}
}
