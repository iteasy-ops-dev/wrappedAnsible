package handlers

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	config "iteasy.wrappedAnsible/configs"
)

const (
	JWT_EXPIRE_TIME  = 60 * time.Minute
	REFRESH_JWT_TIMW = 1 * time.Minute //TODO: 안쓸거면 지우자
)

var JWT_KEY = []byte(config.JWT_KEY)

// TODO: 구조체를 어떻게 정리할지 생각해야함
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func IssueJWT(w http.ResponseWriter, r *http.Request, s User) {
	expirationTime := time.Now().Add(JWT_EXPIRE_TIME)
	claims := &Claims{
		Email: s.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWT_KEY)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
		Path:    "/",
		// HttpOnly: true,
		// SameSite: http.SameSiteNoneMode, // 중요: Cross-Site 요청에서 쿠키 사용 가능하게 설정
		// Secure:   true,                  // 중요: HTTPS에서만 쿠키 전송
	})
}
