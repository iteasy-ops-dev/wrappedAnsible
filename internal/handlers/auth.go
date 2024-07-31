package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	config "iteasy.wrappedAnsible/configs"
	"iteasy.wrappedAnsible/internal/model"
	"iteasy.wrappedAnsible/pkg/utils"
)

type User struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

// mail 메일 인증
func sendVerificationEmail(to, token string) error {
	subject := "Email Verification"
	verificationLink := fmt.Sprintf("%s/verify?token=%s", config.GlobalConfig.Default.Host, token)
	mailBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<body>
			<p>Please verify your email using the following link:</p>
			<p><a href="%s">Verify Email</a></p>
		</body>
		</html>`, verificationLink)

	if err := utils.SendEmail(to, subject, mailBody); err != nil {
		return err
	}
	return nil
}

func sendResetPasswordEmail(to, tempPassword string) error {
	subject := "Password Reset"
	mailBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<body>
			<p>Your temporary password is: <b>%s</b></p>
		</body>
		</html>`, tempPassword)

	if err := utils.SendEmail(to, subject, mailBody); err != nil {
		return err
	}
	return nil
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

	hashedPassword, err := utils.HashingPassword(s.Password)
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(s.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// mail 메일 인증
	verificationToken := utils.GenerateToken()

	a := model.NewAuth(ctx)
	a.SetEmail(s.Email)
	a.SetName(s.Name)
	a.SetPassword(string(hashedPassword))
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
	if err := sendVerificationEmail(s.Email, verificationToken); err != nil {
		http.Error(w, "Failed to send verification email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// mail 메일 인증
func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	a := model.NewAuth(ctx)
	a.SetVerificationToken(token)

	if err := a.VerifyEmail(); err != nil {
		http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email verified successfully"))
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
		// Name: config.GlobalConfig.JWTTokenName,
		Name:     config.GlobalConfig.JWT.TokenName,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	if err := AllowMethod(w, r, http.MethodPost); err != nil {
		return
	}
	var req struct {
		Email string `json:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
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
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	auth := model.NewAuth(ctx)
	auth.SetEmail(req.Email)
	auth.SetPassword(string(hashedPassword))

	// Send temporary password to user's email
	err = sendResetPasswordEmail(req.Email, tempPassword)
	if err != nil {
		http.Error(w, "Error sending email", http.StatusInternalServerError)
		return
	}

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
