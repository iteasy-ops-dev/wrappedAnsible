package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type AccessLog struct {
	Agent      string `bson:"agent"`
	Ip         string `bson:"ip"`
	AccessTime int64  `bson:"accessTime"`
}

type AuthReq struct {
	Email             string `bson:"email"`
	Name              string `bson:"name"`
	AtCreateDate      int64  `bson:"atCreateDate"`
	IsActive          bool   `bson:"isActive"`
	Verified          bool   `bson:"verified"`
	LoginFailureCount int    `bson:"loginFailureCount"`

	AccessLog []AccessLog `bson:"accessLog"`
}

// Auth는 인증 관련 정보를 담는 구조체입니다.
type Auth struct {
	ID                string `bson:"_id,omitempty"`
	Email             string `bson:"email"`
	Name              string `bson:"name"`
	Password          string `bson:"password"`
	AtCreateDate      int64  `bson:"atCreateDate"`
	IsActive          bool   `bson:"isActive"`
	IsAdmin           bool   `bson:"isAdmin"`
	VerificationToken string `bson:"verificationToken"`
	Verified          bool   `bson:"verified"`
	LoginFailureCount int    `bson:"loginFailureCount"`

	AccessLog []AccessLog `bson:"accessLog"`

	// isOr   bool
	isDesc bool

	ctx context.Context
}

// NewAuth는 새로운 Auth 구조체를 초기화합니다.
func NewAuth(ctx context.Context) *Auth {
	return &Auth{
		AtCreateDate: time.Now().Unix(),
		ctx:          ctx,
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
func (a *Auth) SetActive(b bool) {
	a.IsActive = b
}
func (a *Auth) SetIsAdmin(b bool) {
	a.IsAdmin = b
}
func (a *Auth) SetInitAccessLog() {
	a.AccessLog = []AccessLog{}
}
func (a *Auth) SetVerificationToken(token string) {
	a.VerificationToken = token
}

func (a *Auth) SignUp() error {
	isExist, err := a._checkExistingUser()
	if err != nil {
		return err
	}
	if isExist {
		return NewAlreadyExistsError(a.Email)
	}
	return a._signUp()
}

// func InitializeAccessLogs() error {
// 	col := db.Collection(config.GlobalConfig.MongoDB.Collections.Auth)
// 	filter := bson.M{"accessLog": bson.M{"$exists": false}}
// 	update := bson.M{"$set": bson.M{"accessLog": []AccessLog{}}}

// 	_, err := col.UpdateMany(context.Background(), filter, update)
// 	if err != nil {
// 		return fmt.Errorf("failed to initialize accessLog fields: %v", err)
// 	}
// 	return nil
// }

func (a *Auth) Login(agent, ip string) (*Auth, error) {
	user, err := a._findUserByEmail()
	if err != nil {
		return nil, NewUserNotFoundError(a.Email)
	}
	if user == nil {
		return nil, NewUserNotFoundError(a.Email)
	}
	if !user.IsActive {
		return nil, NewUserNotActiveError(a.Email)
	}
	if !user.Verified {
		return nil, NewUserNotVerifiedError(a.Email)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(a.Password))
	if err != nil {
		if incrementErr := a._incrementLoginFailure(); incrementErr != nil {
			return nil, incrementErr
		}
		return nil, NewIncorrectPasswordError()
	}

	if err := a._updateUser(bson.M{"email": a.Email}, bson.M{"$set": bson.M{"loginFailureCount": 0}}); err != nil {
		return nil, fmt.Errorf("failed to reset login failure count: %v", err)
	}

	if err := a._addAccessLog(agent, ip); err != nil {
		return nil, fmt.Errorf("failed to add access log: %v", err)
	}

	return user, nil
}

func (a *Auth) VerifyEmail() error {
	filter := bson.M{"verificationToken": a.VerificationToken}
	update := bson.M{"$set": bson.M{"verified": true, "verificationToken": ""}}

	result, err := a._collection().UpdateOne(a.ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("invalid or expired token")
	}
	return nil
}

func (a *Auth) Get(filter bson.M, page int, pageSize int) ([]AuthReq, int, error) {
	col := a._collection()
	// filter := bson.M{}
	// if a.Email != "" {
	// 	filter["email"] = a.Email
	// }

	// sort := bson.M{"AtCreateDate": 1}
	// if a.isDesc {
	// 	sort = bson.M{"AtCreateDate": -1}
	// }

	// cursor, err := col.Find(a.ctx, filter, options.Find().SetSort(sort))
	// if err != nil {
	// 	return nil, err
	// }

	// var results []AuthReq
	// if err := cursor.All(a.ctx, &results); err != nil {
	// 	return nil, err
	// }

	// return results, nil
	// 페이지 번호와 페이지 크기 검증
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10 // 기본값 설정
	}

	// Skip 및 Limit 설정
	skip := (page - 1) * pageSize

	// 정렬 옵션 추가
	sortOptions := bson.D{
		{Key: "timestamp", Value: -1}, // timestamp 내림차순
	}

	// 데이터 조회
	cursor, err := col.Find(
		context.Background(),
		filter,
		options.Find().
			SetSkip(int64(skip)).
			SetLimit(int64(pageSize)).
			SetSort(sortOptions), // 정렬 옵션 설정
	)
	if err != nil {
		return nil, 0, err
	}

	// 커서를 디코딩하여 결과 반환
	results, err := DecodeCursor[AuthReq](cursor)
	if err != nil {
		return nil, 0, err
	}

	// 전체 데이터 수 조회
	totalCount, err := col.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, err
	}

	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize)) // 총 페이지 수 계산

	return results, totalPages, nil
}

func (a *Auth) UpdateUserActive(s bool) error {
	filter := bson.M{"email": a.Email}
	update := bson.M{"$set": bson.M{"isActive": !s}}

	result, err := a._collection().UpdateOne(a.ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return NewUserNotFoundError(a.Email)
	}
	return nil
}

func (a *Auth) UpdatePassword() error {
	isExist, err := a._checkExistingUser()
	if err != nil {
		return err
	}
	if !isExist {
		return NewUserNotFoundError(a.Email)
	}
	return a._updatePassword()
}

// TODO: 보류
func (a *Auth) Logout() error {
	filter := bson.M{"email": a.Email}
	update := bson.M{"$set": bson.M{"isLogin": false}}

	result, err := a._collection().UpdateOne(a.ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return NewUserLogoutError(a.Email)
	}
	return nil
}
