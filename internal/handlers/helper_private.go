package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	config "iteasy.wrappedAnsible/configs"
)

func _allowMethod(w http.ResponseWriter, r *http.Request, method string) error {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("invalid request method")
	}

	return nil
}

func _validateToken(w http.ResponseWriter, r *http.Request) error {
	// 클라이언트로부터 쿠키에서 토큰을 가져옵니다.
	c, err := r.Cookie(config.GlobalConfig.JWT.TokenName)
	if err != nil {
		// 쿠키가 없는 경우 401 Unauthorized를 반환합니다.
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return fmt.Errorf("unauthorized: no cookie found")
		}
		// 쿠키 오류가 발생한 경우 400 Bad Request를 반환합니다.
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("bad request: %v", err)
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
			return fmt.Errorf("unauthorized: signature invalid")
		}
		// 토큰 해석 중 오류가 발생한 경우 400 Bad Request를 반환합니다.
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("bad request: %v", err)
	}

	// 토큰이 유효하지 않은 경우 401 Unauthorized를 반환합니다.
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return fmt.Errorf("unauthorized: token is not valid")
	}

	// 클레임을 반환합니다.
	return nil
}

func _httpResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// If data is nil, just return with the status code
	if data == nil {
		return
	}

	// Otherwise, encode the data as JSON and write it to the response
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
