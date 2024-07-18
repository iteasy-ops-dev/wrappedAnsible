package handlers

import (
	"encoding/json"
	"net/http"

	config "iteasy.wrappedAnsible/configs"
	"iteasy.wrappedAnsible/pkg/utils"
)

func Functions(w http.ResponseWriter, r *http.Request) {
	if err := AllowMethod(w, r, http.MethodGet); err != nil {
		return
	}
	if err := ValidateToken(w, r); err != nil {
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