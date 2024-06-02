package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type forreqget struct {
	Id int `json:"id"`
	Arg1 string `json:"arg1"`
	Arg2 string `json:"arg2"`
	Operation string `json:"operation"`
	Operation_time int `json:"operation_time"`
}

type forreqpost struct {
	Id int `json:"id"`
	Resultat string `json:"result"`
}

func NewDataGet() *forreqget{
	return &forreqget{}
}

func NewDataPost() *forreqpost{
	return &forreqpost{}
}

func Agent() {
	comp1 := os.Getenv("COMPUTING_POWER")
	if comp1 == "" {
		os.Setenv("COMPUTING_POWER", "10")
	}
	gocount, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil {
		fmt.Println(err)
	}
	cl := &http.Client{}
	reqget, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/internal/task", nil)

	for i := 0; i < gocount; i++ {
		go func(cli *http.Client) {
			for {
				time.Sleep(time.Second)
				repget, err := cli.Do(reqget)
				if err != nil {
					continue
				}
				if repget.StatusCode == 404 || repget.StatusCode == 500 {
					continue // empty list
				}
				dataget := NewDataGet()
				bdget, err := io.ReadAll(repget.Body)
				if err != nil {
					fmt.Println(err)
					continue
				}
				repget.Body.Close()
				err = json.Unmarshal(bdget, dataget)
				if err != nil {
					fmt.Println(err)
					continue
				}
				arg1, err := strconv.ParseFloat(dataget.Arg1, 64) 
				if err != nil {
					fmt.Println(err)
					continue
				}
				arg2, err := strconv.ParseFloat(dataget.Arg2, 64)
				if err != nil {
					fmt.Println(err)
					continue
				}
				result := resolve(arg1, arg2, dataget.Operation, dataget.Operation_time)
				datapost := NewDataPost()
				datapost.Id = dataget.Id
				datapost.Resultat = strconv.FormatFloat(result, 'g', -1, 64)
				jsonData, err := json.Marshal(datapost)
				if err != nil {
					fmt.Println(err) 
					continue
				}
				bdpost := bytes.NewReader(jsonData)
				reqpost, err := http.NewRequest(http.MethodPost, "http://localhost:8080/internal/task", bdpost)
				if err != nil {
					fmt.Println(err)
					continue
				}
				reppost, err := cli.Do(reqpost)
				if err != nil {
					fmt.Println(err)
					continue
				}
				if reppost.StatusCode != 200 {
					fmt.Println(err) 
					continue
				}
			}
		}(cl)
	}
}

func resolve(a, b float64, op string, optime int) float64 {
	time.Sleep(time.Duration(optime) * time.Millisecond)
	if op == "+" {
		return a + b
	} else if op == "-" {
		return a - b
	} else if op == "*" {
		return a * b
	} else {
		return a / b
	}
}