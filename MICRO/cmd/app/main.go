package main

import (
	"github.com/weedworldpeace/distributedcalculator/cmd/agent"
	"github.com/weedworldpeace/distributedcalculator/cmd/orchestrator"
	"github.com/weedworldpeace/distributedcalculator/cmd/server"
	"github.com/weedworldpeace/distributedcalculator/cmd/sql"
)

func main() {
	sql.SqlUp()
	go orchestrator.GrpcServer()
	go agent.Agent()
	server.ServerHTTP()
}