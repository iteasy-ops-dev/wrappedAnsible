package model

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"

	config "iteasy.wrappedAnsible/configs"
	"iteasy.wrappedAnsible/internal/ansible"
)

type AnsibleProcessStatusDocument struct {
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

	// IPs 파라미터를 가져옴 (콤마로 구분된 문자열로 가정)
	ipsParam := query.Get("ips")
	var ips []string
	if ipsParam != "" {
		ips = strings.Split(ipsParam, ",")
		a.IPs = ips
	}

	// isOr 파라미터를 가져옴 (기본값은 false로 설정)
	isOrParam := query.Get("isOr")
	isOr := false
	if isOrParam == "true" {
		isOr = true
	}

	isDescParam := query.Get("isDesc")
	isDesc := true // 기본값
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
	if statusParam != "" {
		if statusParam == "true" {
			a.Status = true
		}
		if statusParam == "false" {
			a.Status = false
		}
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

func (a *AnsibleProcessStatusDocument) Put() {
	col := db.Collection(config.COLLECTION_ANSIBLE_PROCESS_STATUS)
	_, err := col.InsertOne(context.Background(), a)
	if err != nil {
		log.Fatalf("Failed to insert user: %v", err)
	}
	fmt.Println("User created:", a)
}

func (a *AnsibleProcessStatusDocument) Get() ([]AnsibleProcessStatusDocument, error) {
	col := db.Collection(config.COLLECTION_ANSIBLE_PROCESS_STATUS)
	// result := []AnsibleProcessStatus{}
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
	if a.Account != "" {
		filter["account"] = a.Account
	}
	// if a.Status {
	// 	filter["status"] = a.Status
	// }

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
		// filter["duration"] = bson.M{"$gte": a.Duration} // duration 이상
		// filter["duration"] = bson.M{"$lte": a.Duration} // duration 이하 또는 같음
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
	// cursor, err := collection.Find(context.Background(), filter, options)
	if err != nil {
		return nil, err
	}

	return DecodeCursor[AnsibleProcessStatusDocument](cursor)
}
