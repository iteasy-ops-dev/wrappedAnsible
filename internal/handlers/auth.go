package handlers

import (
	"net/http"

	"iteasy.wrappedAnsible/internal/model"
	"iteasy.wrappedAnsible/pkg/utils"
)

type userReq struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

type resetPasswordReq struct {
	Email string `json:"email"`
}
type logoutReq struct {
	Email string `json:"email"`
}

func signUp(w http.ResponseWriter, r *http.Request) {
	if err := _allowMethod(w, r, http.MethodPost); err != nil {
		return
	}
	b, err := utils.ParseRequestBody[userReq](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hashedPassword, err := utils.HashingPassword(b.Password)
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(s.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// mail 메일 인증
	verificationToken := utils.GenerateToken()

	a := model.NewAuth(r.Context())
	a.SetEmail(b.Email)
	a.SetName(b.Name)
	a.SetPassword(string(hashedPassword))
	a.SetActive(true)
	// mail 메일 인증
	a.SetVerificationToken(verificationToken)

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

	// mail 메일 인증
	if err := _sendVerificationEmail(b.Email, verificationToken); err != nil {
		http.Error(w, "Failed to send verification email", http.StatusInternalServerError)
		return
	}

	_httpResponse(w, http.StatusCreated, nil)
}

// mail 메일 인증
func verifyEmail(w http.ResponseWriter, r *http.Request) {
	if err := _allowMethod(w, r, http.MethodGet); err != nil {
		return
	}
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	a := model.NewAuth(r.Context())
	a.SetVerificationToken(token)

	if err := a.VerifyEmail(); err != nil {
		http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email verified successfully"))
}

func login(w http.ResponseWriter, r *http.Request) {
	if err := _allowMethod(w, r, http.MethodPost); err != nil {
		return
	}
	b, err := utils.ParseRequestBody[userReq](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a := model.NewAuth(r.Context())
	a.SetEmail(b.Email)
	a.SetPassword(b.Password)

	userAgent := r.Header.Get("User-Agent")
	ipAddress := r.RemoteAddr
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		ipAddress = forwardedFor
	}

	l, err := a.Login(userAgent, ipAddress)
	if err != nil {
		switch err.(type) {
		case *model.UserNotFoundError:
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		case *model.IncorrectPasswordError:
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		case *model.UserNotActiveError:
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		case *model.UserNotVerifiedError:
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		case *model.AlreadyLoginError:
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	_issueJWT(w, l)
}

func logout(w http.ResponseWriter, r *http.Request) {
	if err := _allowMethod(w, r, http.MethodPost); err != nil {
		return
	}
	b, err := utils.ParseRequestBody[logoutReq](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a := model.NewAuth(r.Context())
	a.SetEmail(b.Email)

	if err := a.Logout(); err != nil {
		switch err.(type) {
		case *model.UserLogoutError:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	_httpResponse(w, http.StatusOK, nil)
}

func resetPassword(w http.ResponseWriter, r *http.Request) {
	if err := _allowMethod(w, r, http.MethodPost); err != nil {
		return
	}
	b, err := utils.ParseRequestBody[resetPasswordReq](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate temporary password
	tempPassword, err := utils.GenerateTempPassword()
	if err != nil {
		http.Error(w, "Error generating temporary password", http.StatusInternalServerError)
		return
	}

	// Hash the temporary password
	hashedPassword, err := utils.HashingPassword(tempPassword)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	a := model.NewAuth(r.Context())
	a.SetEmail(b.Email)
	a.SetPassword(string(hashedPassword))

	// Send temporary password to user's email
	err = _sendResetPasswordEmail(b.Email, tempPassword)
	if err != nil {
		http.Error(w, "Error sending email", http.StatusInternalServerError)
		return
	}

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
