package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	config "iteasy.wrappedAnsible/configs"
	"iteasy.wrappedAnsible/internal/model"
)

type User struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	if err := AllowMethod(w, r, http.MethodPost); err != nil {
		return
	}
	ctx := r.Context()
	var s User
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(s.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	a := model.NewAuth(ctx)
	a.SetEmail(s.Email)
	a.SetName(s.Name)
	a.SetPassword(string(hashedPassword))

	if err := a.SignUp(); err != nil {
		switch err.(type) {
		case *model.AlreadyExistsError:
			http.Error(w, err.Error(), http.StatusConflict)
			return
		default:
			http.Error(w, "Failed to sign up user", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if err := AllowMethod(w, r, http.MethodPost); err != nil {
		return
	}
	ctx := r.Context()
	var s User
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	l := model.NewAuth(ctx)
	l.SetEmail(s.Email)
	l.SetPassword(s.Password)

	if err := l.Login(); err != nil {
		switch err.(type) {
		case *model.UserNotFoundError:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		case *model.IncorrectPasswordError:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		default:
			http.Error(w, "Failed to login user", http.StatusInternalServerError)
			return
		}
	}

	IssueJWT(w, r, s)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// TODO: jwt 만료시키기
	if err := AllowMethod(w, r, http.MethodPost); err != nil {
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: config.GlobalConfig.JWTTokenName,
		// Name:     config.JWT_TOKEN_NAME,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
}
