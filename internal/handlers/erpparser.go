package handlers

import (
	"log"
	"net/http"

	"iteasy.wrappedAnsible/internal/erpparser"
	"iteasy.wrappedAnsible/pkg/utils"
)

type ErpParserReq struct {
	URL string `json:"url"` // JSON의 "url" 필드를 매핑
}

func ErpParser(w http.ResponseWriter, r *http.Request) {
	if err := AllowMethod(w, r, http.MethodPost); err != nil {
		return
	}
	if err := ValidateToken(w, r); err != nil {
		return
	}
	data, err := utils.ParseRequestBody[ErpParserReq](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// JSON 데이터 파싱
	// var data ErpParserReq
	// if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
	// 	http.Error(w, "Error decoding JSON", http.StatusBadRequest)
	// 	return
	// }

	// URL 값 추출
	url := data.URL

	if url == "" {
		http.Error(w, "URL not provided", http.StatusBadRequest)
		return
	}

	// 추출된 URL 사용 (예: 로그에 기록)
	log.Printf("Received URL: %s", url)

	e := erpparser.New(url)
	e.Parsing()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(e.ToBytes()); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
