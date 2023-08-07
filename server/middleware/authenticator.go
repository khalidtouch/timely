package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenBearer := r.Header.Get("Authorization")
		if !strings.Contains(tokenBearer, "Bearer") {
			ctx := context.WithValue(r.Context(), "props", jwt.MapClaims{"user_name": ""})
			next.ServeHTTP(w, r.WithContext(ctx))
			return 
		}

		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		fmt.Printf("authHeader -> %s and len -> %d\n", authHeader, len(authHeader))

		if len(authHeader) != 2 || authHeader[0] == "null" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
			log.Fatal("Malformed Token")
			return 
		}

		jwtToken := authHeader[1]
		token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(signInKey), nil 
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), "props", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
			return 
		}

		fmt.Println("token err -> ", err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("you are unauthorized or your token is expired"))
	})

}