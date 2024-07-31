package model

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	// mail 메일 인증
	VerificationToken string `bson:"verificationToken"`
	Verified          bool   `bson:"verified"`

	isOr   bool
	isDesc bool

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

// mail 메일 인증
func (a *Auth) SetVerificationToken(token string) {
	a.VerificationToken = token
}

func (a *Auth) checkExistingUser() (bool, error) {
	col := db.Collection(config.GlobalConfig.MongoDB.Collections.Auth)
	// col := db.Collection(config.COLLECTION_AUTH)
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
	col := db.Collection(config.GlobalConfig.MongoDB.Collections.Auth)
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
	col := db.Collection(config.GlobalConfig.MongoDB.Collections.Auth)
	filter := bson.M{"email": a.Email}
	sr := col.FindOne(a.ctx, filter)
	r, err := EvaluateAndDecodeSingleResult[Auth](sr)

	if err != nil {
		return NewUserNotFoundError(a.Email)
	}

	// Check if the user is active
	if !r.IsActive {
		return NewUserNotActiveError(a.Email)
	}

	// Check if the user is verified
	if !r.Verified {
		return NewUserNotVerifiedError(a.Email)
	}

	err = bcrypt.CompareHashAndPassword([]byte(r.Password), []byte(a.Password))
	if err != nil {
		return NewIncorrectPasswordError()
	}

	return nil
}

// mail 메일 인증
func (a *Auth) VerifyEmail() error {
	col := db.Collection(config.GlobalConfig.MongoDB.Collections.Auth)

	filter := bson.M{"verificationToken": a.VerificationToken}
	update := bson.M{"$set": bson.M{"verified": true, "verificationToken": ""}}

	result, err := col.UpdateOne(a.ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("invalid or expired token")
	}
	return nil
}

func (a *Auth) Get() ([]Auth, error) {
	col := db.Collection(config.GlobalConfig.MongoDB.Collections.Auth)

	var orderKey string = "atDate"

	filter := bson.M{}

	if a.Email != "" {
		filter["email"] = a.Email
	}

	// 정렬 옵션 설정
	sort := bson.M{orderKey: 1}
	if a.isDesc {
		sort = bson.M{orderKey: -1}
	}

	// 쿼리 옵션 설정
	options := options.Find().SetSort(sort)

	// 쿼리 수행
	cursor, err := col.Find(a.ctx, filter, options)

	if err != nil {
		return nil, err
	}

	return DecodeCursor[Auth](cursor)
}

func (a *Auth) UpdateUserActive(s bool) error {
	col := db.Collection(config.GlobalConfig.MongoDB.Collections.Auth)

	filter := bson.M{"email": a.Email}
	update := bson.M{"$set": bson.M{"isActive": !s}}

	result, err := col.UpdateOne(a.ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("invalid or expired token")
	}
	return nil
}

func (a *Auth) updatePassword() error {
	col := db.Collection(config.GlobalConfig.MongoDB.Collections.Auth)
	filter := bson.M{"email": a.Email}
	update := bson.M{"$set": bson.M{"password": a.Password}}

	_, err := col.UpdateOne(a.ctx, filter, update)

	return err
}

func (a *Auth) UpdatePassword() error {
	isExist, err := a.checkExistingUser()
	if err != nil {
		return err
	}
	if !isExist {
		return NewUserNotFoundError(a.Email)
	}
	return a.updatePassword()
}
