// cmd/server/main.go
package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/rustruber/subscription-service/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Subscription Service API
// @version         0.1.0
// @description     REST-сервис для агрегации данных об онлайн подписках пользователей
// @host            localhost:8080
// @BasePath        /api/v1

func main() {
	r := mux.NewRouter()

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	fmt.Println("Server started at :8080")
	fmt.Println("Swagger UI: http://localhost:8080/swagger/index.html")

	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		panic(err)
	}
}
