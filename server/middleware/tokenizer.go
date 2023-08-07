package middleware 

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var signInKey = "jdnfksdmfksd"

func GenerateToken(userId uint64, username string) (string, error) {
	var err error 
	os.Setenv("ACCESS_SECRET", signInKey)

	claims := jwt.MapClaims{}
	claims["authorized"] = true 
	claims["user_id"] = userId 
	claims["user_name"] = username 
	claims["exp"] = time.Now().Add(time.Minute*15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", errors.New("An error occured during the token creation")
	}

	fmt.Println("jwt map --> ", claims)
	return token, nil 
}
