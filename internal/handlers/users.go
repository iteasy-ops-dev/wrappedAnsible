package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"iteasy.wrappedAnsible/internal/model"
	"iteasy.wrappedAnsible/pkg/utils"
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

func UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	if err := AllowMethod(w, r, http.MethodPost); err != nil {
		return
	}
	if err := ValidateToken(w, r); err != nil {
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Hash the temporary password
	hashedPassword, err := utils.HashingPassword(req.Password)
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	auth := model.NewAuth(ctx)
	auth.SetEmail(req.Email)
	auth.SetPassword(string(hashedPassword))

	// Update user's password in the database
	if err := auth.UpdatePassword(); err != nil {
		switch err.(type) {
		case *model.UserNotFoundError:
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
