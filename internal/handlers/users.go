package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"iteasy.wrappedAnsible/internal/model"
)

func Users(w http.ResponseWriter, r *http.Request) {
	if err := AllowMethod(w, r, http.MethodGet); err != nil {
		return
	}
	if err := ValidateToken(w, r); err != nil {
		return
	}
	n := model.NewAuth(r.Context())
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

func UpdateUserActive(w http.ResponseWriter, r *http.Request) {
	if err := AllowMethod(w, r, http.MethodPost); err != nil {
		return
	}
	if err := ValidateToken(w, r); err != nil {
		return
	}

	type Body struct {
		Email    string `json:"Email"`
		IsActive bool   `json:"IsActive"`
	}

	var b Body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to Error: %v", err), http.StatusInternalServerError)
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &b); err != nil {
		http.Error(w, fmt.Sprintf("failed to Error: %v", err), http.StatusInternalServerError)
	}

	n := model.NewAuth(r.Context())
	n.SetEmail(b.Email)
	err = n.UpdateUserActive(b.IsActive)

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to Query with DB: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
