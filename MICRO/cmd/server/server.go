package server

import (
	"net/http"

	"github.com/weedworldpeace/distributedcalculator/cmd/handlers"
)

func ServerHTTP() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/calculate", handlers.CalculateHandler)

	mux.HandleFunc("/api/v1/expressions", handlers.ExpressionsHandler)

	mux.HandleFunc("/api/v1/expressions/", handlers.ExpressionsIDHandler)

	mux.HandleFunc("/api/v1/register", handlers.RegisterHandler)

	mux.HandleFunc("/api/v1/login", handlers.LoginHandler)

	mux.HandleFunc("/api/v1/clear", handlers.ClearHandler)

	http.ListenAndServe(":8080", mux)
}