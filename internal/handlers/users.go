package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"iteasy.wrappedAnsible/internal/model"
	"iteasy.wrappedAnsible/pkg/utils"
)

type updateActiveReq struct {
	Email    string `json:"Email"`
	IsActive bool   `json:"IsActive"`
}

type updatePasswordReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func users(w http.ResponseWriter, r *http.Request) {
	if err := _allowMethod(w, r, http.MethodPost); err != nil {
		return
	}
	// if err := _validateToken(w, r); err != nil {
	// 	return
	// }

	type req struct {
		Filter   bson.M `json:"filter"`   // JSON의 "filter" 필드를 매핑
		Page     int    `json:"page"`     // 페이지 번호
		PageSize int    `json:"pageSize"` // 페이지 당 항목 수
	}

	var data req
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 페이지 번호와 페이지 크기 검증
	if data.Page < 1 {
		data.Page = 1
	}
	if data.PageSize < 1 {
		data.PageSize = 10 // 기본값 설정
	}

	a := model.NewAuth(r.Context())
	// result, err := a.Get()
	results, totalPages, err := a.Get(data.Filter, data.Page, data.PageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("failed to Query with DB: %v", err), http.StatusInternalServerError)
	// 	return
	// }

	// _httpResponse(w, http.StatusOK, result)

	response := struct {
		Data       []model.AuthReq `json:"data"`
		TotalPages int             `json:"totalPages"`
	}{
		Data:       results,
		TotalPages: totalPages,
	}

	_httpResponse(w, http.StatusOK, response)
}

// TODO: profile 에 이름을 비롯한 추가적인 정보에 대해서 업데이트 필요
// func Profile(w http.ResponseWriter, r *http.Request) {
// 	if err := _allowMethod(w, r, http.MethodGet); err != nil {
// 		return
// 	}
// 	if err := _validateToken(w, r); err != nil {
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

func updateUserActive(w http.ResponseWriter, r *http.Request) {
	if err := _allowMethod(w, r, http.MethodPost); err != nil {
		return
	}
	if err := _validateToken(w, r); err != nil {
		return
	}
	b, err := utils.ParseRequestBody[updateActiveReq](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	a := model.NewAuth(r.Context())
	a.SetEmail(b.Email)
	err = a.UpdateUserActive(b.IsActive)

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to Query with DB: %v", err), http.StatusInternalServerError)
		return
	}

	_httpResponse(w, http.StatusOK, nil)
}

func updateUserPassword(w http.ResponseWriter, r *http.Request) {
	if err := _allowMethod(w, r, http.MethodPost); err != nil {
		return
	}
	if err := _validateToken(w, r); err != nil {
		return
	}
	b, err := utils.ParseRequestBody[updatePasswordReq](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	_httpResponse(w, http.StatusOK, nil)
}
