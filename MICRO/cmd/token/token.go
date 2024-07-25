package token

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const hmacSampleSecret = "super_secret_signature"

func NewToken(login string) (string, error) {
	const hmacSampleSecret = "super_secret_signature"
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": login,
		"nbf": now.Unix(),
		"exp": now.Add(10 * time.Minute).Unix(),
		"iat": now.Unix(),
	})
	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		log.Println(err)
		return "", nil
	}
	return tokenString, nil
}

func Validation(tokenString string) (interface{}, error) {
	tokenFromString, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error){
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("wrong sign method", token.Header["alg"])
		}

		return []byte(hmacSampleSecret), nil
	})

	if err != nil {
		log.Println(err)
		return "", err
	}

	if claims, ok := tokenFromString.Claims.(jwt.MapClaims); ok {
		return claims["login"], nil
	} else {
		log.Println("claims error")
		return nil, errors.New("claims error")
	}
}