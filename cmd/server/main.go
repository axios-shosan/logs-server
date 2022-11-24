package main

import (
	"fmt"
	"net/http"
	"os"
	"sofa-logs-servers/infra/elastic"
	"sofa-logs-servers/routes/requests"
	"sofa-logs-servers/utils"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/cors"
)

func main() {
	router := mux.NewRouter()

	cors := cors.New(cors.Options{
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

	client, err := elastic.Init()
	utils.PanicErr(err)

	router.HandleFunc("/api/logs/create", utils.Middleware(requests.Create, client))
	router.HandleFunc("/api/logs/update", utils.Middleware(requests.Update, client))
	router.HandleFunc("/api/logs/delete", utils.Middleware(requests.Delete, client))
	router.HandleFunc("/api/logs/id", utils.Middleware(requests.FindById, client))
	router.HandleFunc("/api/logs/all", utils.Middleware(requests.FindAll, client))

	fmt.Println("server started at " + port)
	err = http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, cors.Handler(router)))
	if err != nil {
		panic(err)
	}
}
