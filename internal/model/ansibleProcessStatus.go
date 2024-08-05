// package model

// import (
// 	"context"
// 	"fmt"
// 	"net/http"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"go.mongodb.org/mongo-driver/mongo/options"
// 	"gopkg.in/mgo.v2/bson"

// 	config "iteasy.wrappedAnsible/configs"
// 	"iteasy.wrappedAnsible/internal/ansible"
// )

// type AnsibleProcessStatusDocument struct {
// 	ID        string                 `bson:"_id,omitempty"`
// 	Type      string                 `bson:"type"`
// 	IPs       []string               `bson:"ips"`
// 	Name      string                 `bson:"name"`
// 	Email     string                 `bson:"email"`
// 	Account   string                 `bson:"account"`
// 	Status    bool                   `bson:"status"`
// 	Payload   string                 `bson:"payload"`
// 	Duration  int64                  `bson:"duration"`  // Changed to int64 for BSON compatibility
// 	Timestamp int64                  `bson:"timestamp"` // Unix time field
// 	Options   map[string]interface{} `bson:"options,omitempty"`

// 	isOr       bool
// 	isDesc     bool
// 	comparison string
// 	ctx        context.Context
// }

// func _newAnsibleProcessStatus(a *ansible.AnsibleProcessStatus) *AnsibleProcessStatusDocument {
// 	return &AnsibleProcessStatusDocument{
// 		Type:      a.Type,
// 		IPs:       a.IPs,
// 		Name:      a.Name,
// 		Email:     a.Email,
// 		Account:   a.Account,
// 		Status:    a.Status,
// 		Payload:   a.Payload,
// 		Duration:  int64(a.Duration),
// 		Timestamp: time.Now().Unix(),
// 		Options:   a.Options,
// 	}
// }

// func _newHttpRequest(r *http.Request) *AnsibleProcessStatusDocument {
// 	query := r.URL.Query()

// 	var a AnsibleProcessStatusDocument

// 	// IPs 파라미터를 가져옴 (콤마로 구분된 문자열로 가정)
// 	ipsParam := query.Get("ips")
// 	var ips []string
// 	if ipsParam != "" {
// 		ips = strings.Split(ipsParam, ",")
// 		a.IPs = ips
// 	}

// 	// isOr 파라미터를 가져옴 (기본값은 false로 설정)
// 	isOrParam := query.Get("isOr")
// 	isOr := false
// 	if isOrParam == "true" {
// 		isOr = true
// 	}

// 	isDescParam := query.Get("isDesc")
// 	isDesc := true // 기본값: 내림차순
// 	if isDescParam == "false" {
// 		isDesc = false
// 	}

// 	typeParam := query.Get("type")
// 	if typeParam != "" {
// 		a.Type = typeParam
// 	}
// 	nameParam := query.Get("name")
// 	if nameParam != "" {
// 		a.Name = nameParam
// 	}
// 	accountParam := query.Get("account")
// 	if accountParam != "" {
// 		a.Account = accountParam
// 	}
// 	statusParam := query.Get("status")
// 	a.Status = true
// 	if statusParam != "" && statusParam == "false" {
// 		a.Status = false
// 	}
// 	durationParam := query.Get("duration")
// 	if durationParam != "" {
// 		duration, err := strconv.ParseInt(durationParam, 10, 64)
// 		if err == nil {
// 			comparisonParam := query.Get("comparison")
// 			if comparisonParam != "" {
// 				if comparisonParam == "gt" || comparisonParam == "lt" || comparisonParam == "gte" || comparisonParam == "lte" {
// 					a.comparison = fmt.Sprintf("$%s", comparisonParam)
// 				}
// 			} else {
// 				a.comparison = "$gte"
// 			}
// 			a.Duration = duration
// 		}
// 	}

// 	a.isDesc = isDesc
// 	a.isOr = isOr
// 	a.ctx = r.Context()

// 	return &a
// }

// func NewAnsibleProcessStatusDocument(v interface{}) *AnsibleProcessStatusDocument {
// 	switch t := v.(type) {
// 	case *ansible.AnsibleProcessStatus:
// 		return _newAnsibleProcessStatus(t)
// 	case *http.Request:
// 		return _newHttpRequest(t)
// 	default:
// 		return &AnsibleProcessStatusDocument{}
// 	}
// }

// func (a *AnsibleProcessStatusDocument) Put() error {
// 	col := db.Collection(config.GlobalConfig.MongoDB.Collections.AnsibleProcessStatus)
// 	// col := db.Collection(config.COLLECTION_ANSIBLE_PROCESS_STATUS)
// 	_, err := col.InsertOne(a.ctx, a)
// 	if err != nil {
// 		return fmt.Errorf("failed to insert: %w", err)
// 	}
// 	return nil
// }

// // TODO: 페이징을 할 것인가
// // TODO: 네이밍이 이게 뭐냐? 제대로 바꾸자
// func (a *AnsibleProcessStatusDocument) Get() ([]AnsibleProcessStatusDocument, error) {
// 	col := db.Collection(config.GlobalConfig.MongoDB.Collections.AnsibleProcessStatus)

// 	var orderKey string = "timestamp"

// 	filter := bson.M{}
// 	if len(a.IPs) > 0 {
// 		if a.isOr {
// 			filter["ips"] = bson.M{"$in": a.IPs}
// 		} else {
// 			filter["ips"] = bson.M{"$all": a.IPs}
// 		}
// 	}
// 	if a.Type != "" {
// 		filter["type"] = a.Type
// 	}
// 	if a.Name != "" {
// 		filter["name"] = a.Name
// 	}
// 	if a.Email != "" {
// 		filter["email"] = a.Email
// 	}
// 	if a.Account != "" {
// 		filter["account"] = a.Account
// 	}

// 	// TODO: 해당 옵션이 없을 경우 모두 나올 수 있게 하는 방법 찾기
// 	if a.Status {
// 		filter["$or"] = []bson.M{
// 			{"status": true},
// 			{"status": bson.M{"$exists": false}},
// 		}
// 	} else {
// 		filter["$or"] = []bson.M{
// 			{"status": false},
// 			{"status": bson.M{"$exists": false}},
// 		}
// 	}

// 	if a.Duration > 0 {
// 		filter["duration"] = bson.M{a.comparison: a.Duration}
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
// 	// cursor, err := collection.Find(context.Background(), filter, options)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return DecodeCursor[AnsibleProcessStatusDocument](cursor)
// }

// func (a *AnsibleProcessStatusDocument) GetStats() (map[string]interface{}, error) {
// 	col := db.Collection(config.GlobalConfig.MongoDB.Collections.AnsibleProcessStatus)

// 	// Aggregation pipeline
// 	pipeline := []bson.M{
// 		// Match filter: apply the same filter as in Get method
// 		{
// 			"$match": bson.M{
// 				"ips":      bson.M{"$all": a.IPs},
// 				"type":     a.Type,
// 				"name":     a.Name,
// 				"email":    a.Email,
// 				"account":  a.Account,
// 				"duration": bson.M{a.comparison: a.Duration},
// 				"$or": []bson.M{
// 					{"status": a.Status},
// 					{"status": bson.M{"$exists": false}},
// 				},
// 			},
// 		},
// 		// Group stage: calculate the stats
// 		{
// 			"$group": bson.M{
// 				"_id":             nil,
// 				"totalCount":      bson.M{"$sum": 1},
// 				"totalDuration":   bson.M{"$sum": "$duration"},
// 				"successCount":    bson.M{"$sum": bson.M{"$cond": []interface{}{bson.M{"$eq": []interface{}{"$status", true}}, 1, 0}}},
// 				"successDuration": bson.M{"$sum": bson.M{"$cond": []interface{}{bson.M{"$eq": []interface{}{"$status", true}}, "$duration", 0}}},
// 				"failureCount":    bson.M{"$sum": bson.M{"$cond": []interface{}{bson.M{"$eq": []interface{}{"$status", false}}, 1, 0}}},
// 				"failureDuration": bson.M{"$sum": bson.M{"$cond": []interface{}{bson.M{"$eq": []interface{}{"$status", false}}, "$duration", 0}}},
// 			},
// 		},
// 	}

// 	var result []map[string]interface{}
// 	err := col.Pipe(pipeline).All(&result)
// 	if err != nil {
// 		return nil, fmt.Errorf("aggregation error: %w", err)
// 	}

// 	if len(result) > 0 {
// 		return result[0], nil
// 	}

// 	return nil, nil
// }

// func (a *AnsibleProcessStatusDocument) GetStats() (map[string]interface{}, error) {
// 	col := db.Collection(config.GlobalConfig.MongoDB.Collections.AnsibleProcessStatus) // Use mongo-driver collection method

// 	// Aggregation pipeline
// 	pipeline := mongo.Pipeline{
// 		// Match filter: apply the same filter as in Get method
// 		{
// 			{Key: "$match", Value: bson.D{
// 				{Name: "ips", Value: bson.D{{Key: "$all", Value: a.IPs}}},
// 				{Name: "type", Value: a.Type},
// 				{Name: "name", Value: a.Name},
// 				{Name: "email", Value: a.Email},
// 				{Name: "account", Value: a.Account},
// 				{Name: "duration", Value: bson.D{{Key: a.comparison, Value: a.Duration}}},
// 				{Name: "$or", Value: bson.A{
// 					bson.D{{Key: "status", Value: a.Status}},
// 					bson.D{{Key: "status", Value: bson.D{{Key: "$exists", Value: false}}}},
// 				}},
// 			}},
// 		},
// 		// Group stage: calculate the stats
// 		{
// 			{Key: "$group", Value: bson.D{
// 				{Name: "_id", Value: nil},
// 				{Name: "totalCount", Value: bson.D{{Key: "$sum", Value: 1}}},
// 				{Name: "totalDuration", Value: bson.D{{Key: "$sum", Value: "$duration"}}},
// 				{Name: "successCount", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{"$status", 1, 0}}}}}},
// 				{Name: "successDuration", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{"$status", "$duration", 0}}}}}},
// 				{Name: "failureCount", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{"$status", 0, 1}}}}}},
// 				{Name: "failureDuration", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{"$status", 0, "$duration"}}}}}},
// 			}},
// 		},
// 	}

// 	// Run the aggregation pipeline
// 	cursor, err := col.Aggregate(context.Background(), pipeline)
// 	if err != nil {
// 		return nil, fmt.Errorf("aggregation error: %w", err)
// 	}
// 	defer cursor.Close(context.Background())

// 	var result []map[string]interface{}
// 	if err := cursor.All(context.Background(), &result); err != nil {
// 		return nil, fmt.Errorf("error decoding cursor: %w", err)
// 	}

// 	if len(result) > 0 {
// 		return result[0], nil
// 	}

// 	return nil, nil
// }

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

// func (a *AnsibleProcessStatusDocument) Dashboard() (map[string]interface{}, error) {
// 	col := db.Collection(config.GlobalConfig.MongoDB.Collections.AnsibleProcessStatus)

// 	// Aggregation pipeline without $match
// 	pipeline := mongo.Pipeline{
// 		{
// 			{Key: "$group", Value: bson.M{
// 				"_id":             nil,                                                                // Group all documents together
// 				"totalCount":      bson.M{"$sum": 1},                                                  // Count total number of documents
// 				"totalDuration":   bson.M{"$sum": "$duration"},                                        // Sum of all durations
// 				"successCount":    bson.M{"$sum": bson.M{"$cond": bson.A{"$status", 1, 0}}},           // Count documents with status true
// 				"successDuration": bson.M{"$sum": bson.M{"$cond": bson.A{"$status", "$duration", 0}}}, // Sum of durations where status is true
// 				"failureCount":    bson.M{"$sum": bson.M{"$cond": bson.A{"$status", 0, 1}}},           // Count documents with status false
// 				"failureDuration": bson.M{"$sum": bson.M{"$cond": bson.A{"$status", 0, "$duration"}}}, // Sum of durations where status is false
// 			}},
// 		},
// 	}

// 	cursor, err := col.Aggregate(a.ctx, pipeline)
// 	if err != nil {
// 		return nil, fmt.Errorf("aggregation error: %w", err)
// 	}
// 	defer cursor.Close(a.ctx)

// 	var result []map[string]interface{}
// 	if err := cursor.All(a.ctx, &result); err != nil {
// 		return nil, fmt.Errorf("error decoding cursor: %w", err)
// 	}

// 	if len(result) > 0 {
// 		return result[0], nil
// 	}

// 	return nil, nil
// }

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
