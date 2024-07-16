package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	config "iteasy.wrappedAnsible/configs"
)

const (
	JWT_EXPIRE_TIME  = 5 * time.Minute
	REFRESH_JWT_TIMW = 1 * time.Minute
)

var JWT_KEY = []byte(config.JWT_KEY)

// TODO: 구조체를 어떻게 정리할지 생각해야함
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type Credentials struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

// TODO: DB로 넘길것이기 때문에 다 되면 이거 삭제해야함.
var users = map[string]string{} // map to store user data

// TODO: 회원 가입 및 로그인에 대한 모델 작업 추가해야함.
func SignUp(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	users[creds.Email] = string(hashedPassword)
	w.WriteHeader(http.StatusCreated)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	storedPassword, ok := users[creds.Email]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(JWT_EXPIRE_TIME)
	claims := &Claims{
		Email: creds.Email,
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

func ValidateToken(w http.ResponseWriter, r *http.Request) (*Claims, error) {
	// 클라이언트로부터 쿠키에서 토큰을 가져옵니다.
	c, err := r.Cookie("token")
	if err != nil {
		// 쿠키가 없는 경우 401 Unauthorized를 반환합니다.
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return nil, fmt.Errorf("unauthorized: no cookie found")
		}
		// 쿠키 오류가 발생한 경우 400 Bad Request를 반환합니다.
		w.WriteHeader(http.StatusBadRequest)
		return nil, fmt.Errorf("bad request: %v", err)
	}

	// 토큰 값을 가져옵니다.
	tknStr := c.Value

	// 토큰을 파싱하여 클레임(클라임 데이터)을 추출합니다.
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		// 토큰 서명을 검증하기 위해 사용할 key를 반환합니다.
		return JWT_KEY, nil
	})
	if err != nil {
		// 토큰이 유효하지 않은 경우 401 Unauthorized를 반환합니다.
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return nil, fmt.Errorf("unauthorized: signature invalid")
		}
		// 토큰 해석 중 오류가 발생한 경우 400 Bad Request를 반환합니다.
		w.WriteHeader(http.StatusBadRequest)
		return nil, fmt.Errorf("bad request: %v", err)
	}

	// 토큰이 유효하지 않은 경우 401 Unauthorized를 반환합니다.
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return nil, fmt.Errorf("unauthorized: token is not valid")
	}

	// 토큰의 만료 시간을 확인합니다.
	if time.Until(time.Unix(claims.ExpiresAt, 0)) < REFRESH_JWT_TIMW {
		// 토큰의 만료 시간이 1분 미만이면 새로운 토큰을 발급하여 클라이언트에게 전달합니다.
		newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		newExpirationTime := time.Now().Add(JWT_EXPIRE_TIME) // 새로운 토큰의 만료 시간 설정
		newTokenString, err := newToken.SignedString(JWT_KEY)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return nil, fmt.Errorf("internal server error: %v", err)
		}

		// 새로운 토큰을 클라이언트에게 전달하기 위해 쿠키를 설정합니다.
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   newTokenString,
			Expires: newExpirationTime,
			Path:    "/",
		})

		// 클레임과 함께 새로운 토큰을 반환합니다.
		return claims, nil
	}

	// 토큰의 만료 시간이 1분 이상 남은 경우 클레임을 반환합니다.
	return claims, nil
}
