package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	config "iteasy.wrappedAnsible/configs"
	"iteasy.wrappedAnsible/pkg/utils"
)

func ListFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	_, err := ValidateToken(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Authentication error: %v", err), http.StatusUnauthorized)
		return
	}

	res := make([]string, 0)
	l := utils.GetFileList(config.PATH_STATIC_PLAYBOOK)
	for _, name := range l {
		if utils.CheckExtension(name, `.yml`) {
			// 초기화에 필요한 yml이므로 제외
			if name == "requirements.yml" || name == "init.yml" {
				continue
			}
			res = append(res, utils.TruncationExtension(name))
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
