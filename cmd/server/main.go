// cmd/server/main.go
package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/rustruber/subscription-service/docs"
	"github.com/rustruber/subscription-service/internal/infrastructure/config"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Subscription Service API
// @version         0.1.1
// @description     REST-сервис для агрегации данных об онлайн подписках пользователей
// @host            localhost:8080
// @BasePath        /api/v1

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	r := mux.NewRouter()

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	addr := ":" + cfg.ServerPort
	fmt.Println("Server started at ", addr)
	url := fmt.Sprintf("Swagger UI: http://localhost:%s/swagger/index.html", addr)
	fmt.Println(url)

	if err := http.ListenAndServe(addr, r); err != nil {
		panic("Server failed: " + err.Error())
	}
}
