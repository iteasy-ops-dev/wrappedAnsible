package model

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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

func _newAnsibleProcessStatus(a *ansible.AnsibleProcessStatus) *AnsibleProcessStatusDocument {
	return &AnsibleProcessStatusDocument{
		Type:      a.Type,
		IPs:       a.IPs,
		Name:      a.Name,
		Email:     a.Email,
		Account:   a.Account,
		Status:    a.Status,
		Payload:   a.Payload,
		Duration:  int64(a.Duration),
		Timestamp: time.Now().Unix(),
		Options:   a.Options,
	}
}

func _newHttpRequest(r *http.Request) *AnsibleProcessStatusDocument {
	query := r.URL.Query()

	var a AnsibleProcessStatusDocument

	// IPs 파라미터를 가져옴
	ipsParam := query.Get("ips")
	var ips []string
	if ipsParam != "" {
		ips = strings.Split(ipsParam, ",")
		a.IPs = ips
	}

	// isOr 파라미터를 가져옴
	isOrParam := query.Get("isOr")
	isOr := false
	if isOrParam == "true" {
		isOr = true
	}

	isDescParam := query.Get("isDesc")
	isDesc := true
	if isDescParam == "false" {
		isDesc = false
	}

	typeParam := query.Get("type")
	if typeParam != "" {
		a.Type = typeParam
	}
	nameParam := query.Get("name")
	if nameParam != "" {
		a.Name = nameParam
	}
	accountParam := query.Get("account")
	if accountParam != "" {
		a.Account = accountParam
	}
	statusParam := query.Get("status")
	a.Status = true
	if statusParam != "" && statusParam == "false" {
		a.Status = false
	}
	durationParam := query.Get("duration")
	if durationParam != "" {
		duration, err := strconv.ParseInt(durationParam, 10, 64)
		if err == nil {
			comparisonParam := query.Get("comparison")
			if comparisonParam != "" {
				if comparisonParam == "gt" || comparisonParam == "lt" || comparisonParam == "gte" || comparisonParam == "lte" {
					a.comparison = fmt.Sprintf("$%s", comparisonParam)
				}
			} else {
				a.comparison = "$gte"
			}
			a.Duration = duration
		}
	}

	a.isDesc = isDesc
	a.isOr = isOr
	a.ctx = r.Context()

	return &a
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

func (a *AnsibleProcessStatusDocument) Get() ([]AnsibleProcessStatusDocument, error) {
	col := db.Collection(config.GlobalConfig.MongoDB.Collections.AnsibleProcessStatus)

	var orderKey string = "timestamp"

	filter := bson.M{}
	if len(a.IPs) > 0 {
		if a.isOr {
			filter["ips"] = bson.M{"$in": a.IPs}
		} else {
			filter["ips"] = bson.M{"$all": a.IPs}
		}
	}
	if a.Type != "" {
		filter["type"] = a.Type
	}
	if a.Name != "" {
		filter["name"] = a.Name
	}
	if a.Email != "" {
		filter["email"] = a.Email
	}
	if a.Account != "" {
		filter["account"] = a.Account
	}

	if a.Status {
		filter["$or"] = []bson.M{
			{"status": true},
			{"status": bson.M{"$exists": false}},
		}
	} else {
		filter["$or"] = []bson.M{
			{"status": false},
			{"status": bson.M{"$exists": false}},
		}
	}

	if a.Duration > 0 {
		filter["duration"] = bson.M{a.comparison: a.Duration}
	}

	sort := bson.M{orderKey: 1}
	if a.isDesc {
		sort = bson.M{orderKey: -1}
	}

	opts := options.Find().SetSort(sort)

	cursor, err := col.Find(a.ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(a.ctx)

	var results []AnsibleProcessStatusDocument
	if err := cursor.All(a.ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (a *AnsibleProcessStatusDocument) Dashboard() (map[string]interface{}, error) {
	col := db.Collection(config.GlobalConfig.MongoDB.Collections.AnsibleProcessStatus)

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
