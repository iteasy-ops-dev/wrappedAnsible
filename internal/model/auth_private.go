package model

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	config "iteasy.wrappedAnsible/configs"
)

func (a *Auth) _collection() *mongo.Collection {
	return db.Collection(config.GlobalConfig.MongoDB.Collections.Auth)
}

// 이메일로 사용자를 조회하는 헬퍼 함수
func (a *Auth) _findUserByEmail() (*Auth, error) {
	col := a._collection()
	filter := bson.M{"email": a.Email}
	var user Auth
	err := col.FindOne(a.ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// MongoDB 업데이트 연산을 수행하는 헬퍼 함수
func (a *Auth) _updateUser(filter, update bson.M) error {
	col := a._collection()
	_, err := col.UpdateOne(a.ctx, filter, update)
	return err
}

// 사용자 존재 여부 체크
func (a *Auth) _checkExistingUser() (bool, error) {
	user, err := a._findUserByEmail()
	if err != nil {
		return false, err
	}
	return user != nil, nil
}

func (a *Auth) _signUp() error {
	col := a._collection()
	a.SetInitAccessLog()
	_, err := col.InsertOne(a.ctx, a)

	return err
}

func (a *Auth) _incrementLoginFailure() error {
	filter := bson.M{"email": a.Email}
	update := bson.M{"$inc": bson.M{"loginFailureCount": 1}} // 실패 횟수 증가

	if err := a._updateUser(filter, update); err != nil {
		return fmt.Errorf("failed to update login failure count: %v", err)
	}

	user, err := a._findUserByEmail()
	if err != nil {
		return fmt.Errorf("failed to fetch user after login failure increment: %v", err)
	}

	if user.LoginFailureCount >= 5 {
		update = bson.M{"$set": bson.M{"isActive": false}}
		if err := a._updateUser(filter, update); err != nil {
			return fmt.Errorf("failed to deactivate account after 5 login failures: %v", err)
		}
	}

	return nil
}

func (a *Auth) _addAccessLog(agent, ip string) error {
	accessLog := AccessLog{
		Agent:      agent,
		Ip:         ip,
		AccessTime: time.Now().Unix(),
	}
	update := bson.M{"$push": bson.M{"accessLog": accessLog}}

	return a._updateUser(bson.M{"email": a.Email}, update)
}

func (a *Auth) _updatePassword() error {
	filter := bson.M{"email": a.Email}
	update := bson.M{"$set": bson.M{"password": a.Password}}

	return a._updateUser(filter, update)
}
