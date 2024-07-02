package handlers

import (
	"encoding/json"
	"net/http"

	"iteasy.wrappedAnsible/pkg/utils"
)

func ListFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	res := make([]string, 0)
	l := utils.GetFileList("static/playbooks")
	for _, name := range l {
		if utils.CheckExtension(name, `.yml`) {
			res = append(res, utils.TruncationExtension(name))
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
