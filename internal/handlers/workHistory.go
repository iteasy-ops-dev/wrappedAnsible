package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"iteasy.wrappedAnsible/internal/erpparser"
)

func workHistory(w http.ResponseWriter, r *http.Request) {
	if err := _allowMethod(w, r, http.MethodGet); err != nil {
		return
	}
	if err := _validateToken(w, r); err != nil {
		return
	}

	workHistories, err := erpparser.WorkHistoryParsing()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(workHistories); err != nil {
		http.Error(w, fmt.Sprintf("JSON encoding error: %v", err), http.StatusInternalServerError)
	}
}
