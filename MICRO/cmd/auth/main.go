package auth

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/weedworldpeace/distributedcalculator/cmd/sql"
)

const hmacSampleSecret = "super_secret_signature"

func NewToken(login string) (string, error) {
	const hmacSampleSecret = "super_secret_signature"
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": login,
		"nbf": now.Unix(),
		"exp": now.Add(5 * time.Minute).Unix(),
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

func Register(w http.ResponseWriter, r *http.Request) {
	u := sql.User{}

	defer r.Body.Close()
	bd, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("jsonData error"))
		log.Println(err)
		return
	}

	err = json.Unmarshal(bd, &u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("jsonData error"))
		log.Println(err)
		return
	}

	res, err := sql.MyDB.LoginExists(u.Login)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("sql error"))
		log.Println(err)
		return
	}
	if res {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("try autre login"))
		return
	}
	_, err = sql.MyDB.InsertUser(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("sql error"))
		log.Println(err)
		return
	}
	w.Write([]byte("user registered"))
}

func Login(w http.ResponseWriter, r *http.Request) {
	u := sql.User{}

	defer r.Body.Close()
	bd, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("jsonData error"))
		log.Println(err)
		return
	}

	err = json.Unmarshal(bd, &u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("jsonData error"))
		log.Println(err)
		return
	}

	pass, err := sql.MyDB.SelectPassword(u.Login)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if pass != u.Password {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("wrong password"))
		return
	} else {
		tok, err := NewToken(u.Login)
		if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("token error"))
		log.Println(err)
		return
		}
		w.Write([]byte(tok))
	}

}