package model

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"

	config "iteasy.wrappedAnsible/configs"
)

// TODO: 회원 가입 및 로그인에 대한 모델 작업 추가해야함.
type Auth struct {
	ID       string `bson:"_id,omitempty"`
	Email    string `bson:"email"`
	Name     string `bson:"name"`
	Password string `bson:"password"`
	AtDate   int64  `bson:"atDate"`
	IsActive bool   `bson:"isActive"`
	IsAdmin  bool   `bson:"IsAdmin"`

	ctx context.Context
}

func NewAuth(ctx context.Context) *Auth {
	return &Auth{
		IsActive: true,
		IsAdmin:  false,
		AtDate:   time.Now().Unix(),

		ctx: ctx,
	}
}

func (a *Auth) SetEmail(s string) {
	a.Email = s
}
func (a *Auth) SetName(s string) {
	a.Name = s
}
func (a *Auth) SetPassword(s string) {
	a.Password = s
}

func (a *Auth) checkExistingUser() (bool, error) {
	col := db.Collection(config.COLLECTION_AUTH)
	filter := bson.M{"email": a.Email}
	var existingUser Auth
	err := col.FindOne(a.ctx, filter).Decode(&existingUser)
	if err == nil {
		return true, nil
	} else if err == mongo.ErrNoDocuments {
		return false, nil
	} else {
		return false, err
	}
}

// func (a *Auth) FindUser() (bool, error) {

// }

func (a *Auth) signUp() error {
	col := db.Collection(config.COLLECTION_AUTH)
	_, err := col.InsertOne(a.ctx, a)
	if err != nil {
		return err
	}
	return nil
}

func (a *Auth) SignUp() error {
	isExist, err := a.checkExistingUser()
	if err != nil {
		return err
	}
	if isExist {
		return NewAlreadyExistsError(a.Email)
	}
	return a.signUp()
}

func (a *Auth) Login() error {
	col := db.Collection(config.COLLECTION_AUTH)
	filter := bson.M{"email": a.Email}
	sr := col.FindOne(a.ctx, filter)
	r, err := EvaluateAndDecodeSingleResult[Auth](sr)

	if err != nil {
		return NewUserNotFoundError(a.Email)
	}

	err = bcrypt.CompareHashAndPassword([]byte(r.Password), []byte(a.Password))
	if err != nil {
		return NewIncorrectPasswordError()
	}

	return nil
}
