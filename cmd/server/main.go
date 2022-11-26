package main

import (
	"fmt"
	"net/http"
	"os"
	"sofa-logs-servers/infra/zincsearch"
	"sofa-logs-servers/routes/requests"
	"sofa-logs-servers/routes/transactions"
	"sofa-logs-servers/utils"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/cors"
)

func main() {
	router := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodPost,
			http.MethodGet,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	zincClient, err := zincsearch.Init()
	utils.PanicErr(err)

	router.HandleFunc("/api/logs/create", utils.Middleware(requests.Create, zincClient))
	router.HandleFunc("/api/logs/update", utils.Middleware(requests.Update, zincClient))
	router.HandleFunc("/api/logs/delete", utils.Middleware(requests.Delete, zincClient))
	router.HandleFunc("/api/logs/id", utils.Middleware(requests.FindById, zincClient))
	router.HandleFunc("/api/logs/all", utils.Middleware(requests.FindAll, zincClient))
	router.HandleFunc("/api/transactions/create", utils.Middleware(transactions.Create, zincClient))
	router.HandleFunc("/api/transactions/update", utils.Middleware(transactions.Update, zincClient))
	router.HandleFunc("/api/transactions/delete", utils.Middleware(transactions.Delete, zincClient))
	router.HandleFunc("/api/transactions/all", utils.Middleware(transactions.FindAll, zincClient))
	router.HandleFunc("/api/transactions/id", utils.Middleware(transactions.FindById, zincClient))

	fmt.Println("server started at " + port)
	err = http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, c.Handler(router)))
	if err != nil {
		panic(err)
	}
}
