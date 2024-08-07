// package model

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"time"

// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// 	"golang.org/x/crypto/bcrypt"
// 	"gopkg.in/mgo.v2/bson"

// 	config "iteasy.wrappedAnsible/configs"
// )

// // TODO: 회원 가입 및 로그인에 대한 모델 작업 추가해야함.
// type Auth struct {
// 	ID       string `bson:"_id,omitempty"`
// 	Email    string `bson:"email"`
// 	Name     string `bson:"name"`
// 	Password string `bson:"password"`
// 	AtDate   int64  `bson:"atDate"`
// 	IsActive bool   `bson:"isActive"`
// 	IsAdmin  bool   `bson:"IsAdmin"`

// 	// mail 메일 인증
// 	VerificationToken string `bson:"verificationToken"`
// 	Verified          bool   `bson:"verified"`

// 	LoginFailureCount int `bson:"loginFailureCount"`

// 	isOr   bool
// 	isDesc bool

// 	ctx context.Context
// }

// func NewAuth(ctx context.Context) *Auth {
// 	return &Auth{
// 		IsActive: true,
// 		IsAdmin:  false,
// 		AtDate:   time.Now().Unix(),

// 		ctx: ctx,
// 	}
// }

// func (a *Auth) SetEmail(s string) {
// 	a.Email = s
// }
// func (a *Auth) SetName(s string) {
// 	a.Name = s
// }
// func (a *Auth) SetPassword(s string) {
// 	a.Password = s
// }

// // mail 메일 인증
// func (a *Auth) SetVerificationToken(token string) {
// 	a.VerificationToken = token
// }

// func (a *Auth) checkExistingUser() (bool, error) {
// 	col := db.Collection(config.GlobalConfig.MongoDB.Collections.Auth)
// 	// col := db.Collection(config.COLLECTION_AUTH)
// 	filter := bson.M{"email": a.Email}
// 	var existingUser Auth
// 	err := col.FindOne(a.ctx, filter).Decode(&existingUser)
// 	if err == nil {
// 		return true, nil
// 	} else if err == mongo.ErrNoDocuments {
// 		return false, nil
// 	} else {
// 		return false, err
// 	}
// }

// // func (a *Auth) FindUser() (bool, error) {

// // }

// func (a *Auth) signUp() error {
// 	col := db.Collection(config.GlobalConfig.MongoDB.Collections.Auth)
// 	_, err := col.InsertOne(a.ctx, a)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (a *Auth) SignUp() error {
// 	isExist, err := a.checkExistingUser()
// 	if err != nil {
// 		return err
// 	}
// 	if isExist {
// 		return NewAlreadyExistsError(a.Email)
// 	}
// 	return a.signUp()
// }

// func (a *Auth) _incrementLoginFailure() error {
// 	col := db.Collection(config.GlobalConfig.MongoDB.Collections.Auth)
// 	filter := bson.M{"email": a.Email}
// 	update := bson.M{
// 		"$inc": bson.M{"loginFailureCount": 1}, // 실패 횟수 증가
// 	}

// 	// 로그인 실패 횟수 업데이트
// 	_, err := col.UpdateOne(a.ctx, filter, update)
// 	if err != nil {
// 		return fmt.Errorf("failed to update login failure count: %v", err)
// 	}

// 	// 현재 로그인 실패 횟수 가져오기
// 	sr := col.FindOne(a.ctx, filter)
// 	user, err := EvaluateAndDecodeSingleResult[Auth](sr)
// 	if err != nil {
// 		return fmt.Errorf("failed to fetch user after login failure increment: %v", err)
// 	}

// 	// 5번 실패 시 계정 비활성화
// 	if user.LoginFailureCount >= 5 {
// 		update = bson.M{"$set": bson.M{"isActive": false}}
// 		_, err := col.UpdateOne(a.ctx, filter, update)
// 		if err != nil {
// 			return fmt.Errorf("failed to deactivate account after 5 login failures: %v", err)
// 		}
// 	}

// 	return nil
// }

// func (a *Auth) Login() (*Auth, error) {
// 	col := db.Collection(config.GlobalConfig.MongoDB.Collections.Auth)
// 	filter := bson.M{"email": a.Email}
// 	sr := col.FindOne(a.ctx, filter)
// 	r, err := EvaluateAndDecodeSingleResult[Auth](sr)

// 	if err != nil {
// 		return r, NewUserNotFoundError(a.Email)
// 	}

// 	// Check if the user is active
// 	if !r.IsActive {
// 		return r, NewUserNotActiveError(a.Email)
// 	}

// 	// Check if the user is verified
// 	if !r.Verified {
// 		return r, NewUserNotVerifiedError(a.Email)
// 	}

// 	err = bcrypt.CompareHashAndPassword([]byte(r.Password), []byte(a.Password))
// 	if err != nil {
// 		// 비밀번호가 일치하지 않으면 로그인 실패 처리
// 		if incrementErr := a._incrementLoginFailure(); incrementErr != nil {
// 			return nil, incrementErr
// 		}
// 		return r, NewIncorrectPasswordError()
// 	}

// 	update := bson.M{"$set": bson.M{"loginFailureCount": 0}}
// 	_, err = col.UpdateOne(a.ctx, filter, update)
// 	if err != nil {
// 		return r, fmt.Errorf("failed to reset login failure count: %v", err)
// 	}

// 	return r, nil
// }

// // mail 메일 인증
// func (a *Auth) VerifyEmail() error {
// 	col := db.Collection(config.GlobalConfig.MongoDB.Collections.Auth)

// 	filter := bson.M{"verificationToken": a.VerificationToken}
// 	update := bson.M{"$set": bson.M{"verified": true, "verificationToken": ""}}

// 	result, err := col.UpdateOne(a.ctx, filter, update)
// 	if err != nil {
// 		return err
// 	}
// 	if result.MatchedCount == 0 {
// 		return errors.New("invalid or expired token")
// 	}
// 	return nil
// }

// // TODO: 네이밍이 이게 뭐냐? 제대로 바꾸자
// func (a *Auth) Get() ([]Auth, error) {
// 	col := db.Collection(config.GlobalConfig.MongoDB.Collections.Auth)

// 	var orderKey string = "atDate"

// 	filter := bson.M{}

// 	if a.Email != "" {
// 		filter["email"] = a.Email
// 	}

// 	// 정렬 옵션 설정
// 	sort := bson.M{orderKey: 1}
// 	if a.isDesc {
// 		sort = bson.M{orderKey: -1}
// 	}

// 	// 쿼리 옵션 설정
// 	options := options.Find().SetSort(sort)

// 	// 쿼리 수행
// 	cursor, err := col.Find(a.ctx, filter, options)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return DecodeCursor[Auth](cursor)
// }

// func (a *Auth) UpdateUserActive(s bool) error {
// 	col := db.Collection(config.GlobalConfig.MongoDB.Collections.Auth)

// 	filter := bson.M{"email": a.Email}
// 	update := bson.M{"$set": bson.M{"isActive": !s}}

// 	result, err := col.UpdateOne(a.ctx, filter, update)
// 	if err != nil {
// 		return err
// 	}
// 	if result.MatchedCount == 0 {
// 		return errors.New("invalid or expired token")
// 	}
// 	return nil
// }

// func (a *Auth) updatePassword() error {
// 	col := db.Collection(config.GlobalConfig.MongoDB.Collections.Auth)
// 	filter := bson.M{"email": a.Email}
// 	update := bson.M{"$set": bson.M{"password": a.Password}}

// 	_, err := col.UpdateOne(a.ctx, filter, update)

// 	return err
// }

// func (a *Auth) UpdatePassword() error {
// 	isExist, err := a.checkExistingUser()
// 	if err != nil {
// 		return err
// 	}
// 	if !isExist {
// 		return NewUserNotFoundError(a.Email)
// 	}
// 	return a.updatePassword()
// }

package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	config "iteasy.wrappedAnsible/configs"
)

// Auth는 인증 관련 정보를 담는 구조체입니다.
type Auth struct {
	ID                string `bson:"_id,omitempty"`
	Email             string `bson:"email"`
	Name              string `bson:"name"`
	Password          string `bson:"password"`
	AtDate            int64  `bson:"atDate"`
	IsActive          bool   `bson:"isActive"`
	IsAdmin           bool   `bson:"isAdmin"`
	VerificationToken string `bson:"verificationToken"`
	Verified          bool   `bson:"verified"`
	LoginFailureCount int    `bson:"loginFailureCount"`

	isOr   bool
	isDesc bool

	ctx context.Context
}

// NewAuth는 새로운 Auth 구조체를 초기화합니다.
func NewAuth(ctx context.Context) *Auth {
	return &Auth{
		IsActive: true,
		IsAdmin:  false,
		AtDate:   time.Now().Unix(),
		ctx:      ctx,
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
func (a *Auth) SetVerified(b bool) {
	a.Verified = b
}

func (a *Auth) SetVerificationToken(token string) {
	a.VerificationToken = token
}

func (a *Auth) checkExistingUser() (bool, error) {
	col := db.Collection(config.GlobalConfig.MongoDB.Collections.Auth)
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

func (a *Auth) _incrementLoginFailure() error {
	col := db.Collection(config.GlobalConfig.MongoDB.Collections.Auth)
	filter := bson.M{"email": a.Email}
	update := bson.M{
		"$inc": bson.M{"loginFailureCount": 1}, // 실패 횟수 증가
	}

	_, err := col.UpdateOne(a.ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update login failure count: %v", err)
	}

	var user Auth
	sr := col.FindOne(a.ctx, filter)
	if err := sr.Decode(&user); err != nil {
		return fmt.Errorf("failed to fetch user after login failure increment: %v", err)
	}

	if user.LoginFailureCount >= 5 {
		update = bson.M{"$set": bson.M{"isActive": false}}
		_, err := col.UpdateOne(a.ctx, filter, update)
		if err != nil {
			return fmt.Errorf("failed to deactivate account after 5 login failures: %v", err)
		}
	}

	return nil
}

func (a *Auth) Login() (*Auth, error) {
	col := db.Collection(config.GlobalConfig.MongoDB.Collections.Auth)
	filter := bson.M{"email": a.Email}
	sr := col.FindOne(a.ctx, filter)
	var r Auth
	if err := sr.Decode(&r); err != nil {
		return nil, NewUserNotFoundError(a.Email)
	}

	if !r.IsActive {
		return nil, NewUserNotActiveError(a.Email)
	}

	if !r.Verified {
		return nil, NewUserNotVerifiedError(a.Email)
	}

	err := bcrypt.CompareHashAndPassword([]byte(r.Password), []byte(a.Password))
	if err != nil {
		if incrementErr := a._incrementLoginFailure(); incrementErr != nil {
			return nil, incrementErr
		}
		return nil, NewIncorrectPasswordError()
	}

	update := bson.M{"$set": bson.M{"loginFailureCount": 0}}
	_, err = col.UpdateOne(a.ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to reset login failure count: %v", err)
	}

	return &r, nil
}

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

	sort := bson.M{orderKey: 1}
	if a.isDesc {
		sort = bson.M{orderKey: -1}
	}

	findOptions := options.Find().SetSort(sort)
	cursor, err := col.Find(a.ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	var results []Auth
	if err := cursor.All(a.ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
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
		return errors.New("user not found")
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
