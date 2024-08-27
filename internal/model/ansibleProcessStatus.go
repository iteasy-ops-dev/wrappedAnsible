package model

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	config "iteasy.wrappedAnsible/configs"
	"iteasy.wrappedAnsible/internal/ansible"
)

type AnsibleProcessStatusDocument struct {
	ID        string                 `bson:"_id,omitempty"`
	Type      string                 `bson:"type"`
	IPs       []string               `bson:"ips"`
	Name      string                 `bson:"name"`
	Email     string                 `bson:"email"`
	Account   string                 `bson:"account"`
	Status    bool                   `bson:"status"`
	Payload   string                 `bson:"payload"`
	Duration  int64                  `bson:"duration"`
	Timestamp int64                  `bson:"timestamp"`
	Options   map[string]interface{} `bson:"options,omitempty"`

	isOr       bool
	isDesc     bool
	comparison string
	ctx        context.Context
}

func NewAnsibleProcessStatusDocument(v interface{}) *AnsibleProcessStatusDocument {
	switch t := v.(type) {
	case *ansible.AnsibleProcessStatus:
		return _newAnsibleProcessStatus(t)
	case *http.Request:
		return _newHttpRequest(t)
	default:
		return &AnsibleProcessStatusDocument{}
	}
}

func (a *AnsibleProcessStatusDocument) Put() error {
	col := db.Collection(config.GlobalConfig.MongoDB.Collections.AnsibleProcessStatus)
	_, err := col.InsertOne(a.ctx, a)
	if err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func (a *AnsibleProcessStatusDocument) Get(filter bson.M, page int, pageSize int) ([]AnsibleProcessStatusDocument, int, error) {
	col := a._collection()

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
	results, err := DecodeCursor[AnsibleProcessStatusDocument](cursor)
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

func (a *AnsibleProcessStatusDocument) Dashboard() (map[string]interface{}, error) {
	col := a._collection()

	// 전체 문서 통계 집계 파이프라인
	overallPipeline := mongo.Pipeline{
		{
			{Key: "$group", Value: bson.M{
				"_id":             nil,                                                                // 전체 문서 하나로 그룹화
				"totalCount":      bson.M{"$sum": 1},                                                  // 총 문서 수
				"totalDuration":   bson.M{"$sum": "$duration"},                                        // 총 duration 합계
				"successCount":    bson.M{"$sum": bson.M{"$cond": bson.A{"$status", 1, 0}}},           // 성공 상태(true)의 문서 수
				"successDuration": bson.M{"$sum": bson.M{"$cond": bson.A{"$status", "$duration", 0}}}, // 성공 상태(true)의 duration 합계
				"failureCount":    bson.M{"$sum": bson.M{"$cond": bson.A{"$status", 0, 1}}},           // 실패 상태(false)의 문서 수
				"failureDuration": bson.M{"$sum": bson.M{"$cond": bson.A{"$status", 0, "$duration"}}}, // 실패 상태(false)의 duration 합계
			}},
		},
	}

	// 타입별 문서 통계 집계 파이프라인
	typePipeline := mongo.Pipeline{
		{
			{Key: "$group", Value: bson.M{
				"_id":             "$type",                                                            // 타입별로 그룹화
				"totalCount":      bson.M{"$sum": 1},                                                  // 총 문서 수
				"totalDuration":   bson.M{"$sum": "$duration"},                                        // 총 duration 합계
				"successCount":    bson.M{"$sum": bson.M{"$cond": bson.A{"$status", 1, 0}}},           // 성공 상태(true)의 문서 수
				"successDuration": bson.M{"$sum": bson.M{"$cond": bson.A{"$status", "$duration", 0}}}, // 성공 상태(true)의 duration 합계
				"failureCount":    bson.M{"$sum": bson.M{"$cond": bson.A{"$status", 0, 1}}},           // 실패 상태(false)의 문서 수
				"failureDuration": bson.M{"$sum": bson.M{"$cond": bson.A{"$status", 0, "$duration"}}}, // 실패 상태(false)의 duration 합계
			}},
		},
		{
			{Key: "$project", Value: bson.M{
				"type":            "$_id", // _id를 type으로 이름 변경
				"totalCount":      1,
				"totalDuration":   1,
				"successCount":    1,
				"successDuration": 1,
				"failureCount":    1,
				"failureDuration": 1,
			}},
		},
		{
			{Key: "$sort", Value: bson.M{"type": 1}}, // 선택 사항: 타입별로 정렬 (필요 시)
		},
	}

	// 전체 문서 통계 가져오기
	overallCursor, err := col.Aggregate(a.ctx, overallPipeline)
	if err != nil {
		return nil, fmt.Errorf("전체 문서 통계 집계 오류: %w", err)
	}
	defer overallCursor.Close(a.ctx)

	var overallResult []map[string]interface{}
	if err := overallCursor.All(a.ctx, &overallResult); err != nil {
		return nil, fmt.Errorf("전체 문서 통계 커서 디코딩 오류: %w", err)
	}

	// 타입별 문서 통계 가져오기
	typeCursor, err := col.Aggregate(a.ctx, typePipeline)
	if err != nil {
		return nil, fmt.Errorf("타입별 문서 통계 집계 오류: %w", err)
	}
	defer typeCursor.Close(a.ctx)

	var typeResult []map[string]interface{}
	if err := typeCursor.All(a.ctx, &typeResult); err != nil {
		return nil, fmt.Errorf("타입별 문서 통계 커서 디코딩 오류: %w", err)
	}

	// 결과 결합
	result := map[string]interface{}{
		"overall": overallResult,
		"types":   typeResult,
	}

	return result, nil
}
