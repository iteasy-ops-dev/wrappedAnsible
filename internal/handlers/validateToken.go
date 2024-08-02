package handlers

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	config "iteasy.wrappedAnsible/configs"
)

// func ValidateToken(w http.ResponseWriter, r *http.Request) error {
// 	// 클라이언트로부터 쿠키에서 토큰을 가져옵니다.
// 	c, err := r.Cookie(JWT_TOKEN_NAME)
// 	if err != nil {
// 		// 쿠키가 없는 경우 401 Unauthorized를 반환합니다.
// 		if err == http.ErrNoCookie {
// 			w.WriteHeader(http.StatusUnauthorized)
// 			return fmt.Errorf("unauthorized: no cookie found")
// 		}
// 		// 쿠키 오류가 발생한 경우 400 Bad Request를 반환합니다.
// 		w.WriteHeader(http.StatusBadRequest)
// 		return fmt.Errorf("bad request: %v", err)
// 	}

// 	// 토큰 값을 가져옵니다.
// 	tknStr := c.Value

// 	// 토큰을 파싱하여 클레임(클라임 데이터)을 추출합니다.
// 	claims := &Claims{}
// 	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
// 		// 토큰 서명을 검증하기 위해 사용할 key를 반환합니다.
// 		return JWT_KEY, nil
// 	})
// 	if err != nil {
// 		// 토큰이 유효하지 않은 경우 401 Unauthorized를 반환합니다.
// 		if err == jwt.ErrSignatureInvalid {
// 			w.WriteHeader(http.StatusUnauthorized)
// 			return fmt.Errorf("unauthorized: signature invalid")
// 		}
// 		// 토큰 해석 중 오류가 발생한 경우 400 Bad Request를 반환합니다.
// 		w.WriteHeader(http.StatusBadRequest)
// 		return fmt.Errorf("bad request: %v", err)
// 	}

// 	// 토큰이 유효하지 않은 경우 401 Unauthorized를 반환합니다.
// 	if !tkn.Valid {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		return fmt.Errorf("unauthorized: token is not valid")
// 	}

// 	// 토큰의 만료 시간을 확인합니다.
// 	if time.Until(time.Unix(claims.ExpiresAt, 0)) < REFRESH_JWT_TIMW {
// 		// 토큰의 만료 시간이 1분 미만이면 새로운 토큰을 발급하여 클라이언트에게 전달합니다.
// 		newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 		newExpirationTime := time.Now().Add(JWT_EXPIRE_TIME) // 새로운 토큰의 만료 시간 설정
// 		newTokenString, err := newToken.SignedString(JWT_KEY)
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return fmt.Errorf("internal server error: %v", err)
// 		}

// 		// 새로운 토큰을 클라이언트에게 전달하기 위해 쿠키를 설정합니다.
// 		http.SetCookie(w, &http.Cookie{
// 			Name:    JWT_TOKEN_NAME,
// 			Value:   newTokenString,
// 			Expires: newExpirationTime,
// 			Path:    "/",
// 		})

// 		// 클레임과 함께 새로운 토큰을 반환합니다.
// 		return nil
// 	}

// 	// 토큰의 만료 시간이 1분 이상 남은 경우 클레임을 반환합니다.
// 	return nil
// }

func ValidateToken(w http.ResponseWriter, r *http.Request) error {
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
