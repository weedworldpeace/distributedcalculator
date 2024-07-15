package agent

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	pb "github.com/weedworldpeace/distributedcalculator/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Agent() {
	comp1 := os.Getenv("COMPUTING_POWER")
	if comp1 == "" {
		os.Setenv("COMPUTING_POWER", "10")
	}
	gocount, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil {
		fmt.Println(err)
	}

	host := "localhost"
	port := "5000"
	addr := fmt.Sprintf("%s:%s", host, port)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println("could not connect to grpc server: ", err)
		os.Exit(1)
	}
 
	grpcClient := pb.NewCalculatorServiceClient(conn)

	for i := 0; i < gocount; i++ {
		go func() {
			for {
				time.Sleep(time.Second)
				res, err := grpcClient.TaskGet(context.TODO(), nil)
				if err != nil {
					continue
				}
				arg1, err := strconv.ParseFloat(res.Arg1, 64) 
				if err != nil {
					fmt.Println(err)
					continue
				}
				arg2, err := strconv.ParseFloat(res.Arg2, 64)
				if err != nil {
					fmt.Println(err)
					continue
				}
				result := resolve(arg1, arg2, res.Operation, int(res.OperationTime))
				grpcClient.TaskPost(context.TODO(), &pb.TaskPostRequest{Id: res.Id, Resultat: strconv.FormatFloat(result, 'g', -1, 64)})
			}
		}()
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