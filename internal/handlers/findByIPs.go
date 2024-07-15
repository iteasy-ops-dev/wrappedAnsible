package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"iteasy.wrappedAnsible/internal/model"
)

func FindByIPs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 쿼리 파라미터에서 값을 가져옴
	query := r.URL.Query()

	// IPs 파라미터를 가져옴 (콤마로 구분된 문자열로 가정)
	ipsParam := query.Get("ips")
	var ips []string
	if ipsParam != "" {
		ips = strings.Split(ipsParam, ",")
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

	result, err := model.FindByIPs(ips, isOr, isDesc)

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to Query with DB: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
}
