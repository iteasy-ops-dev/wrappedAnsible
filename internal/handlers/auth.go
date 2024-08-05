package handlers

import (
	"fmt"
	"net/http"
	"time"

	config "iteasy.wrappedAnsible/configs"
	"iteasy.wrappedAnsible/internal/model"
	"iteasy.wrappedAnsible/pkg/utils"
)

type UserReq struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

type ResetPasswordReq struct {
	Email string `json:"email"`
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
	b, err := utils.ParseRequestBody[UserReq](r)
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
	if err := sendVerificationEmail(b.Email, verificationToken); err != nil {
		http.Error(w, "Failed to send verification email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// mail 메일 인증
func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	if err := AllowMethod(w, r, http.MethodGet); err != nil {
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

func Login(w http.ResponseWriter, r *http.Request) {
	if err := AllowMethod(w, r, http.MethodPost); err != nil {
		return
	}
	b, err := utils.ParseRequestBody[UserReq](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a := model.NewAuth(r.Context())
	a.SetEmail(b.Email)
	a.SetPassword(b.Password)

	l, err := a.Login()
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
		default:
			http.Error(w, "Failed to login user", http.StatusInternalServerError)
			return
		}
	}

	IssueJWT(w, r, l)
}

// TODO: 로그아웃시 프론트에서 쿠키 제거하는 함수 사용중.
// 해당 함수를 어떻게 처리 할 것인가
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
	b, err := utils.ParseRequestBody[ResetPasswordReq](r)
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
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	a := model.NewAuth(r.Context())
	a.SetEmail(b.Email)
	a.SetPassword(string(hashedPassword))

	// Send temporary password to user's email
	err = sendResetPasswordEmail(b.Email, tempPassword)
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

	w.WriteHeader(http.StatusOK)
}
