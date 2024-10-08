package handlers

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	config "iteasy.wrappedAnsible/configs"
	"iteasy.wrappedAnsible/internal/model"
)

const (
	JWT_EXPIRE_TIME = 12 * 60 * time.Minute // 12시간
)

var JWT_KEY = []byte(config.GlobalConfig.JWT.Key)

// TODO: 구조체를 어떻게 정리할지 생각해야함
type Claims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	jwt.StandardClaims
}

func _issueJWT(w http.ResponseWriter, s *model.Auth) {
	expirationTime := time.Now().Add(JWT_EXPIRE_TIME)
	claims := &Claims{
		Name:  s.Name,
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
		Name:    config.GlobalConfig.JWT.TokenName,
		Value:   tokenString,
		Expires: expirationTime,
		Path:    "/",
		// HttpOnly: true,
		// SameSite: http.SameSiteNoneMode, // 중요: Cross-Site 요청에서 쿠키 사용 가능하게 설정
		// Secure:   true,                  // 중요: HTTPS에서만 쿠키 전송
	})
}

// ExtendJWT handles the extension of the existing JWT token
func extendJWT(w http.ResponseWriter, r *http.Request) {
	if err := _allowMethod(w, r, http.MethodPost); err != nil {
		return
	}

	// Get the existing token from the cookie
	cookie, err := r.Cookie(config.GlobalConfig.JWT.TokenName)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	tokenString := cookie.Value

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JWT_KEY, nil
	})
	if err != nil || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Create a new token with extended expiration
	claims, ok := token.Claims.(*Claims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(JWT_EXPIRE_TIME)
	claims.ExpiresAt = expirationTime.Unix()

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newTokenString, err := newToken.SignedString(JWT_KEY)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the new token in the cookie
	http.SetCookie(w, &http.Cookie{
		Name:    config.GlobalConfig.JWT.TokenName,
		Value:   newTokenString,
		Expires: expirationTime,
		Path:    "/",
		// HttpOnly: true,
		// SameSite: http.SameSiteNoneMode, // Cross-Site requests
		// Secure:   true,                  // HTTPS only
	})
}
