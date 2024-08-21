package handlers

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"iteasy.wrappedAnsible/internal/model"
)

func Insight(w http.ResponseWriter, r *http.Request) {
	if err := _allowMethod(w, r, http.MethodPost); err != nil {
		return
	}
	if err := _validateToken(w, r); err != nil {
		return
	}

	type req struct {
		Filter bson.M `json:"filter"` // JSON의 "filter" 필드를 매핑
	}

	// 요청 바디에서 JSON 데이터를 파싱합니다.
	var data req
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m := model.NewWorkHistory()
	results, err := m.Dashboard(data.Filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
