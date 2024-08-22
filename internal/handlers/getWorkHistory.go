package handlers

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"iteasy.wrappedAnsible/internal/model"
)

func getWorkHistory(w http.ResponseWriter, r *http.Request) {
	if err := _allowMethod(w, r, http.MethodPost); err != nil {
		return
	}
	if err := _validateToken(w, r); err != nil {
		return
	}

	type req struct {
		Filter   bson.M `json:"filter"`   // JSON의 "filter" 필드를 매핑
		Page     int    `json:"page"`     // 페이지 번호
		PageSize int    `json:"pageSize"` // 페이지 당 항목 수
	}

	// 요청 바디에서 JSON 데이터를 파싱합니다.
	var data req
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 페이지 번호와 페이지 크기 검증
	if data.Page < 1 {
		data.Page = 1
	}
	if data.PageSize < 1 {
		data.PageSize = 10 // 기본값 설정
	}

	m := model.NewWorkHistory()
	results, totalPages, err := m.Get(data.Filter, data.Page, data.PageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 응답 데이터 구조
	response := struct {
		Data       []model.WorkHistory `json:"data"`
		TotalPages int                 `json:"totalPages"`
	}{
		Data:       results,
		TotalPages: totalPages,
	}

	_httpResponse(w, http.StatusOK, response)
}
