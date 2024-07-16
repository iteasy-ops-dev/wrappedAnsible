package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"iteasy.wrappedAnsible/internal/model"
)

func Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	_, err := ValidateToken(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Authentication error: %v", err), http.StatusUnauthorized)
		return
	}

	n := model.NewAnsibleProcessStatusDocument(r)
	result, err := n.Get()

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
