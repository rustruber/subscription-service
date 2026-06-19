// cmd/server/main.go
package main

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // 👈 ВАЖНО: пустой импорт драйвера!
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/rustruber/subscription-service/internal/adapter/api/rest"
	"github.com/rustruber/subscription-service/internal/adapter/logger"
	"github.com/rustruber/subscription-service/internal/adapter/repository/postgres"
	"github.com/rustruber/subscription-service/internal/application/subscription"
	"github.com/rustruber/subscription-service/internal/infrastructure/config"

	_ "github.com/rustruber/subscription-service/docs"
)

// @title           Subscription Service API
// @version         0.2.0
// @description     REST-сервис для агрегации данных об онлайн подписках пользователей
// @host            localhost:8080
// @BasePath        /api/v1
func main() {
	// 1. Загружаем конфиг
	cfg, err := config.Load()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// 2. Создаём логгер
	log := logger.NewLogger(cfg.LogLevel)

	// 3. Подключаем БД
	db, err := sql.Open("postgres", cfg.GetDBConnectionString())
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}
	defer db.Close()

	// Автоматически накатываем миграции
	if err := runMigrations(db, log); err != nil {
		log.Fatal("Failed to run migrations", "error", err)
	}
	log.Info("Migrations applied successfully")

	// 4. Проверяем подключение
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database", "error", err)
	}
	log.Info("Database connected")

	// 5. Создаём репозиторий
	repo := postgres.NewPostgresRepository(db, log)

	// 6. Создаём Use Case
	useCase := subscription.NewUseCase(repo, log)

	// 7. Создаём Handler
	handler := rest.NewHandler(useCase)

	// 8. Создаём роутер
	r := mux.NewRouter()

	// 9. Регистрируем роуты
	r.HandleFunc("/subscriptions/total-cost", handler.GetTotalCost).Methods("GET")
	r.HandleFunc("/subscriptions/{id}", handler.GetByID).Methods("GET")
	r.HandleFunc("/subscriptions/{id}", handler.Update).Methods("PUT")
	r.HandleFunc("/subscriptions/{id}", handler.Delete).Methods("DELETE")
	r.HandleFunc("/subscriptions", handler.List).Methods("GET")
	r.HandleFunc("/subscriptions", handler.Create).Methods("POST")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	// Swagger UI
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// 10. Запускаем сервер
	addr := ":" + cfg.ServerPort
	log.Info("Server starting", "port", cfg.ServerPort)
	log.Info("Swagger UI", "url", "http://localhost"+addr+"/swagger/index.html")

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal("Server failed", "error", err)
	}
}

func runMigrations(db *sql.DB, log *logger.Logger) error {
	// Проверяем, существует ли таблица
	var exists bool
	err := db.QueryRow(`
        SELECT EXISTS (
            SELECT 1 FROM information_schema.tables 
            WHERE table_name = 'subscriptions'
        )
    `).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		log.Info("Table subscriptions already exists")
		return nil
	}

	log.Info("Creating table subscriptions...")

	// Создаём таблицу
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS subscriptions (
            id UUID PRIMARY KEY,
            service_name VARCHAR(255) NOT NULL,
            price INTEGER NOT NULL CHECK (price > 0),
            user_id UUID NOT NULL,
            start_date TIMESTAMP NOT NULL,
            end_date TIMESTAMP,
            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		return err
	}

	// Создаём индексы
	_, err = db.Exec(`
        CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions(user_id);
        CREATE INDEX IF NOT EXISTS idx_subscriptions_service_name ON subscriptions(service_name);
        CREATE INDEX IF NOT EXISTS idx_subscriptions_start_date ON subscriptions(start_date);
    `)
	if err != nil {
		return err
	}

	log.Info("Table subscriptions created with indexes")
	return nil
}
