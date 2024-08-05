package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"iteasy.wrappedAnsible/internal/model"
	"iteasy.wrappedAnsible/pkg/utils"
)

type UpdateActiveReq struct {
	Email    string `json:"Email"`
	IsActive bool   `json:"IsActive"`
}

type UpdatePasswordReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Users(w http.ResponseWriter, r *http.Request) {
	if err := AllowMethod(w, r, http.MethodGet); err != nil {
		return
	}
	if err := ValidateToken(w, r); err != nil {
		return
	}

	a := model.NewAuth(r.Context())
	result, err := a.Get()

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

// TODO: profile 에 이름을 비롯한 추가적인 정보에 대해서 업데이트 필요
// func Profile(w http.ResponseWriter, r *http.Request) {
// 	if err := AllowMethod(w, r, http.MethodGet); err != nil {
// 		return
// 	}
// 	if err := ValidateToken(w, r); err != nil {
// 		return
// 	}

// 	type Body struct {
// 		Email string `json:"Email"`
// 	}

// 	var b Body
// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("failed to Error: %v", err), http.StatusInternalServerError)
// 	}
// 	defer r.Body.Close()

// 	if err := json.Unmarshal(body, &b); err != nil {
// 		http.Error(w, fmt.Sprintf("failed to Error: %v", err), http.StatusInternalServerError)
// 	}

// 	n := model.NewAuth(r.Context())
// 	n.SetEmail(b.Email)

// }

func UpdateUserActive(w http.ResponseWriter, r *http.Request) {
	if err := AllowMethod(w, r, http.MethodPost); err != nil {
		return
	}
	if err := ValidateToken(w, r); err != nil {
		return
	}
	b, err := utils.ParseRequestBody[UpdateActiveReq](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a := model.NewAuth(r.Context())
	a.SetEmail(b.Email)
	err = a.UpdateUserActive(b.IsActive)

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
	b, err := utils.ParseRequestBody[UpdatePasswordReq](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Hash the temporary password
	hashedPassword, err := utils.HashingPassword(b.Password)
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	a := model.NewAuth(ctx)
	a.SetEmail(b.Email)
	a.SetPassword(string(hashedPassword))

	// Update user's password in the database
	if err := a.UpdatePassword(); err != nil {
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
