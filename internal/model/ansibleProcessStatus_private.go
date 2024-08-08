package model

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	config "iteasy.wrappedAnsible/configs"
	"iteasy.wrappedAnsible/internal/ansible"
)

func (a *AnsibleProcessStatusDocument) _collection() *mongo.Collection {
	return db.Collection(config.GlobalConfig.MongoDB.Collections.AnsibleProcessStatus)
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
