package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/weedworldpeace/distributedcalculator/cmd/orchestrator"
	"github.com/weedworldpeace/distributedcalculator/cmd/postfix"
	"github.com/weedworldpeace/distributedcalculator/cmd/sql"
	"github.com/weedworldpeace/distributedcalculator/cmd/token"
)

type respBody struct { 
	Token string `json:"token"`
	Expression string `json:"expression"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
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
		w.Write([]byte("try another login"))
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
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
		tok, err := token.NewToken(u.Login)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("token error"))
			log.Println(err)
			return
		}
		w.Write([]byte(tok))
	}
}

func ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	rb := respBody{}
	by, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("body err"))
		return
	}
	err = json.Unmarshal(by, &rb)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("json err"))
		return
	}

	login, err := token.Validation(rb.Token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	id, err := sql.MyDB.SelectID(login.(string))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	exps, err := sql.MyDB.SelectExpressions(int64(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(exps)
}

func ExpressionsIDHandler(w http.ResponseWriter, r *http.Request) {
	rb := respBody{}
	by, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("body err"))
		return
	}
	err = json.Unmarshal(by, &rb)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("json err"))
		return
	}

	login, err := token.Validation(rb.Token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	userid, err := sql.MyDB.SelectID(login.(string))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	expID, err := strconv.Atoi(strings.Split(r.URL.Path, "expressions/")[1])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "invalid id")
	} else {
		exp, err := sql.MyDB.SelectExpressionById(int64(expID), userid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		json.NewEncoder(w).Encode(exp)
	}
}

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	rb := respBody{}
	by, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("body err"))
		return
	}
	err = json.Unmarshal(by, &rb)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("json err"))
		return
	}

	login, err := token.Validation(rb.Token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	id, err := sql.MyDB.SelectID(login.(string))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	
	defer r.Body.Close()

	newStr, err := postfix.ToPostfix(rb.Expression)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	} else {
		expID, err := sql.MyDB.InsertExpression(int64(id))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}
		w.WriteHeader(201)
		resp := "id = " + strconv.Itoa(int(expID))
		fmt.Fprint(w, "accepted, " + resp)
		orchestrator.Orchestrator(newStr, expID)
	}
}

func ClearHandler(w http.ResponseWriter, r *http.Request) {
	err := sql.MyDB.Clear()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
	w.Write([]byte("database is clear"))
}