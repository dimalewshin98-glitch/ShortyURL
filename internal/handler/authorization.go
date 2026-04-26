package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

var userID = 0

const TOKEN_EXP = time.Hour * 3
const SECRET_KEY = "supersecretkey"

func BuildJWTString() (string, error) {
	userID += 1
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		UserID: userID,
	})
	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetUserID(tokenString string) int {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SECRET_KEY), nil
		})
	if err != nil {
		return -1
	}
	if !token.Valid {
		return -1
	}
	return claims.UserID
}

func AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string
		cookie, err := r.Cookie("token")
		if err != nil {
			if err.Error() == "http: named cookie not present" {
				tokenString, err = BuildJWTString()
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
				}
			} else {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			tokenString = cookie.Value
		}
		userID := GetUserID(tokenString)
		switch userID {
		case 0:
			http.Error(w, "", http.StatusUnauthorized)
			return
		case -1:
			tokenString, err = BuildJWTString()
		}
		http.SetCookie(w, &http.Cookie{
			Name:  "token",
			Value: tokenString,
		})
		ctx := context.WithValue(r.Context(), "userID", strconv.Itoa(userID))
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
