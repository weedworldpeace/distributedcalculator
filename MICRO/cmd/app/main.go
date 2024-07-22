package main

import (
	"fmt"
	"os"

	"github.com/weedworldpeace/distributedcalculator/cmd/agent"
	"github.com/weedworldpeace/distributedcalculator/cmd/orchestrator"
	"github.com/weedworldpeace/distributedcalculator/cmd/sql"
)

func main() {
	addenv := os.Getenv("TIME_ADDITION_MS")
	if addenv == "" {
		err := os.Setenv("TIME_ADDITION_MS", "100")
		if err != nil {
			fmt.Println("createenverr")
		}
	}
	subenv := os.Getenv("TIME_SUBTRACTION_MS")
	if subenv == "" {
		err := os.Setenv("TIME_SUBTRACTION_MS", "100")
		if err != nil {
			fmt.Println("createenverr")
		}
	}
	mulenv := os.Getenv("TIME_MULTIPLICATIONS_MS")
	if mulenv == "" {
		err := os.Setenv("TIME_MULTIPLICATIONS_MS", "100")
		if err != nil {
			fmt.Println("createenverr")
		}
	}
	divenv := os.Getenv("TIME_DIVISIONS_MS")
	if divenv == "" {
		err := os.Setenv("TIME_DIVISIONS_MS", "100")
		if err != nil {
			fmt.Println("createenverr")
		}
	}  

	sql.SqlUp()
	go agent.Agent()
	orchestrator.Orchestrator()
}