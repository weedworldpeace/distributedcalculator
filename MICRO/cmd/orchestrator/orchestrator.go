package orchestrator

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"slices"
	"strconv"
	"sync"

	ev "github.com/weedworldpeace/distributedcalculator/cmd/env_var"
	pb "github.com/weedworldpeace/distributedcalculator/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/weedworldpeace/distributedcalculator/cmd/sql"
)

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

var mu = sync.Mutex{}
var miniglobalid int
var miniexpressions = make(map[int]*fortask)
var	miniresults = make(map[int]*result)

type Server struct {
	pb.CalculatorServiceServer
} 

func NewServer() *Server {
	return &Server{}
}

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
	res := &result{}
	res.Id = int(in.Id)
	res.Result = in.Resultat
	mu.Lock()
	miniresults[res.Id] = res
	mu.Unlock()
	return nil, nil
}

func GrpcServer() {
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

func Orchestrator(str []string, id int64) {
	wg := sync.WaitGroup{}
	supArr := []string{"+", "-", "*", "/"}

	for {
		for i := 0; i < len(str) - 2; i++ {
			if !slices.Contains(supArr, str[i]) && !slices.Contains(supArr, str[i + 1]) && slices.Contains(supArr, str[i + 2]) {
				wg.Add(1)
				go func(ind int) {
					fort := &fortask{}
					if str[ind + 2] == "+" {
						fort.Operation_time = ev.OperVars("+")
					} else if str[ind + 2] == "-" {
						fort.Operation_time = ev.OperVars("-")
					} else if str[ind + 2] == "*" {
						fort.Operation_time = ev.OperVars("*")
					} else {
						fort.Operation_time = ev.OperVars("/")
					}
					eblan := []string{str[ind], str[ind + 1], str[ind + 2]}

					fort.Arg1 = str[ind]
					fort.Arg2 = str[ind + 1]
					fort.Operation = str[ind + 2]
					mu.Lock()
					fort.Id = miniglobalid
					miniglobalid += 1
					miniexpressions[fort.Id] = fort
					mu.Unlock()

					for {
						mu.Lock()
						v, b := miniresults[fort.Id]
						delete(miniresults, fort.Id)
						mu.Unlock() 
						if b {
							eblan = ReplaceFirstSequence(str, eblan, v.Result)
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
		sql.MyDB.UpdateExpression(finres, int64(id))
		mu.Unlock()
	}
}

func findSequenceIndex(slice []string, sequence []string) int {
	for i := 0; i <= len(slice)-len(sequence); i++ {
	  match := true
	  for j := 0; j < len(sequence); j++ {
		if slice[i+j] != sequence[j] {
		  match = false
		  break
		}
	  }
	  if match {
		return i
	  }
	}
	return -1
}
  
func ReplaceFirstSequence(slice []string, sequence []string, replacement string) []string {
	index := findSequenceIndex(slice, sequence)
	if index == -1 {
	  return slice
	}
	newSlice := append(slice[:index], append([]string{replacement}, slice[index+len(sequence):]...)...)
	return newSlice
}