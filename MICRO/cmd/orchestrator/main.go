package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"

	pb "github.com/weedworldpeace/distributedcalculator/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/weedworldpeace/distributedcalculator/cmd/postfix"
)

type expression struct {
	Id int
	Status string
	Result float64
}

type calculate struct {
	Expression string `json:"expression"`
}

type result struct {
	Id int `json:"id"`
	Result string `json:"result"`
}

type fortask struct {
	Id int `json:"id"`
	Arg1 string `json:"arg1"`
	Arg2 string `json:"arg2"`
	Operation string `json:"operation"`
	Operation_time int `json:"operation_time"`
}

func newCalculate() *calculate {
	return &calculate{}
}

func newResult() *result {
	return &result{}
}

func newForTask() *fortask {
	return &fortask{}
}

var globalid int
var miniglobalid int

var miniexpressions = make(map[int]*fortask)
var	miniresults = make(map[int]*result)

type Server struct {
	pb.CalculatorServiceServer
} 

func NewServer() *Server {
	return &Server{}
}

var mu = sync.Mutex{}

func (s *Server) TaskGet(ctx context.Context, _ *emptypb.Empty) (*pb.TaskGetResponse, error) {
	mu.Lock()
	defer mu.Unlock()
	for i, t := range miniexpressions{ 
		defer delete(miniexpressions, i)
		return &pb.TaskGetResponse{Id: int32(t.Id), Arg1: t.Arg1, Arg2: t.Arg2, Operation: t.Operation, OperationTime: int32(t.Operation_time)}, nil
	}
	return nil, status.Error(codes.OutOfRange, "no data")
}

func (s *Server) TaskPost(ctx context.Context, in *pb.TaskPostRequest) (*emptypb.Empty, error) {
	res := newResult()
	res.Id = int(in.Id)
	res.Result = in.Resultat
	mu.Lock()
	miniresults[res.Id] = res
	mu.Unlock()
	return nil, nil
}

func grpcSer() {
	host := "localhost"
	port := "5000" // to change

	addr := fmt.Sprintf("%s:%s", host, port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println("error starting tcp listener: ", err)
		os.Exit(1)
	}

	log.Println("tcp listener started at port: ", port)

	grpcServer := grpc.NewServer()
	calcServiceServer := NewServer()
	pb.RegisterCalculatorServiceServer(grpcServer, calcServiceServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Println("error serving grpc: ", err)
		os.Exit(1)
	}
}

func forExpression(expressions *map[int]*expression, miniexpressions *map[int]*fortask, miniresults *map[int]*result, str []string, id int, mu *sync.Mutex) {
	wg := sync.WaitGroup{}
	exp := *expressions
	minexp := *miniexpressions
	minres := *miniresults
	supArr := []string{"+", "-", "*", "/"}

	for {
		for i := 0; i < len(str) - 2; i++ {
			if !slices.Contains(supArr, str[i]) && !slices.Contains(supArr, str[i + 1]) && slices.Contains(supArr, str[i + 2]) {
				wg.Add(1)
				go func(ind int) {
					fort := newForTask()
					if str[ind + 2] == "+" {
						addenv, err := strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
						if err != nil {
							fmt.Println("addenv error")
						}
						fort.Operation_time = addenv
					} else if str[ind + 2] == "-" {
						subenv, err := strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS"))
						if err != nil {
							fmt.Println("subenv error")
						}
						fort.Operation_time = subenv
					} else if str[ind + 2] == "*" {
						mulenv, err := strconv.Atoi(os.Getenv("TIME_MULTIPLICATIONS_MS"))
						if err != nil {
							fmt.Println("mulenv error")
						}
						fort.Operation_time = mulenv
					} else if str[ind + 2] == "/" {
						divenv, err := strconv.Atoi(os.Getenv("TIME_DIVISIONS_MS"))
						if err != nil {
							fmt.Println("divenv error")
						}
						fort.Operation_time = divenv
					}
					eblan := []string{str[ind], str[ind + 1], str[ind + 2]}
					fort.Arg1 = str[ind]
					fort.Arg2 = str[ind + 1]
					fort.Operation = str[ind + 2]
					mu.Lock()
					fort.Id = miniglobalid
					miniglobalid += 1
					minexp[fort.Id] = fort
					mu.Unlock()
					fmt.Println(fort.Operation_time)
					for {
						mu.Lock()
						v, b := minres[fort.Id]
						delete(minres, fort.Id)
						mu.Unlock() 
						if b {
							eblan = postfix.ReplaceFirstSequence(str, eblan, v.Result)
							mu.Lock()
							str = eblan
							mu.Unlock()
							break
						}
					}
					wg.Done()
				}(i)
			}
		}
		wg.Wait()
		if len(str) == 1 {
			break
		}
	}
	finres, err := strconv.ParseFloat(string(str[0]), 64)
	if err != nil {
		fmt.Println(err)
	} else {
		mu.Lock()
		exp[id].Result = finres
		exp[id].Status = "resolved"
		mu.Unlock()
	}
}

func Orchestrator() {
	expressions := make(map[int]*expression)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/calculate", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		bd, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "smth goes wrong")
		} else {
			data := newCalculate()
			err := json.Unmarshal(bd, data)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "smth goes wrong")
			} else {
				newStr, err := postfix.ToPostfix(data.Expression)
				if err != nil {
					w.WriteHeader(422)
					fmt.Fprint(w, "invalid data ", err)
				} else {
					mu.Lock()
					id := globalid
					globalid += 1
					expressions[id] = &expression{id,  "accepted", 0}
					mu.Unlock()
					w.WriteHeader(201)
					resp := "id = " + strconv.Itoa(id)
					fmt.Fprint(w, "accepted, " + resp)
					go forExpression(&expressions, &miniexpressions, &miniresults, newStr, id, &mu)
				}
			}
		}
	})

	mux.HandleFunc("/api/v1/expressions", func(w http.ResponseWriter, r *http.Request) {
		resp := make(map[string][]expression)
		resparr := []expression{}
		for i := range expressions {
			newexp := *expressions[i]
			resparr = append(resparr, newexp)
		}
		resp["expressions"] = resparr
		json.NewEncoder(w).Encode(resp)
	})

	mux.HandleFunc("/api/v1/expressions/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		id, err := strconv.Atoi(strings.Split(path, "expressions/")[1])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "invalid id")
		} else if len(expressions) <= id {
			w.WriteHeader(http.StatusNotFound) 
			fmt.Fprint(w, "bad id")
		} else {
			resp := make(map[string]expression)
			resp["expression"] = *expressions[id]
			json.NewEncoder(w).Encode(resp)
		}
	})

	// mux.HandleFunc("/internal/task", func(w http.ResponseWriter, r *http.Request) {
	// 	if r.Method == http.MethodGet {
	// 		mu.Lock()
	// 		if len(miniexpressions) == 0 {
	// 			w.WriteHeader(404)
	// 		} else {
	// 			for i := range miniexpressions{
	// 				err := json.NewEncoder(w).Encode(miniexpressions[i])
	// 				if err != nil {
	// 					fmt.Println(err)
	// 				}
	// 				delete(miniexpressions, i)
	// 				break
	// 			}
	// 		}
	// 		mu.Unlock()
	// 	} else if r.Method == http.MethodPost {
	// 		res := newResult()
	// 		defer r.Body.Close()
	// 		iodata, err := io.ReadAll(r.Body)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 			w.WriteHeader(500)
	// 			fmt.Fprint(w, "smth goes wrong")
	// 		} else {
	// 			err := json.Unmarshal(iodata, res)
	// 			if err != nil {
	// 				w.WriteHeader(422)
	// 				fmt.Fprint(w, "invalid data")
	// 			} else {
	// 				mu.Lock()
	// 				miniresults[res.Id] = res
	// 				mu.Unlock()
	// 				fmt.Fprint(w, "accepted")
	// 			}
	// 		}
	// 	} else {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		fmt.Fprint(w, "wrong request")
	// 	}
	// })

	go grpcSer()

	http.ListenAndServe(":8080", mux)
}