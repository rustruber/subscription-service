// cmd/server/main.go
package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	_ "github.com/rustruber/subscription-service/docs"
	httpSwagger "github.com/swaggo/http-swagger"

	"net/http"

	"github.com/rustruber/subscription-service/internal/adapter/api/rest"
	"github.com/rustruber/subscription-service/internal/adapter/logger"
	"github.com/rustruber/subscription-service/internal/adapter/repository/postgres"
	"github.com/rustruber/subscription-service/internal/application/subscription"
	"github.com/rustruber/subscription-service/internal/infrastructure/config"
	_ "net/http/pprof"
)

// @title           Subscription Service API
// @version         0.3.1
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

	// 3. Создаём БД, если её нет
	if err := ensureDatabaseExists(cfg, log); err != nil {
		log.Fatal("Failed to ensure database exists", "error", err)
	}

	// 4. Подключаемся к БД
	db, err := sql.Open("postgres", cfg.GetDBConnectionString())
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}
	defer db.Close()

	// 5. Проверяем подключение
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database", "error", err)
	}
	log.Info("Database connected")

	// 6. Автоматически накатываем миграции
	if err := runMigrations(db, log); err != nil {
		log.Fatal("Failed to run migrations", "error", err)
	}
	log.Info("Migrations applied successfully")

	// 7. Создаём репозиторий
	repo := postgres.NewPostgresRepository(db, log)

	// 8. Создаём Use Case
	useCase := subscription.NewUseCase(repo, log)

	// 9. Создаём Handler
	handler := rest.NewHandler(useCase)

	// 10. Создаём роутер
	r := mux.NewRouter()

	// 11. Регистрируем роуты
	r.HandleFunc("/subscriptions/total-cost", handler.GetTotalCost).Methods("GET")
	r.HandleFunc("/subscriptions/{id}", handler.GetByID).Methods("GET")
	r.HandleFunc("/subscriptions/{id}", handler.Update).Methods("PUT")
	r.HandleFunc("/subscriptions/{id}", handler.Delete).Methods("DELETE")
	r.HandleFunc("/subscriptions", handler.List).Methods("GET")
	r.HandleFunc("/subscriptions", handler.Create).Methods("POST")

	// 12. Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	// 13. Swagger UI
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// 14. Запускаем сервер
	addr := ":" + cfg.ServerPort
	log.Info("Server starting", "port", cfg.ServerPort)
	log.Info("Swagger UI", "url", "http://localhost"+addr+"/swagger/index.html")
	go foo()
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal("Server failed", "error", err)
	}
}

// ensureDatabaseExists проверяет, существует ли база данных, и создаёт её при необходимости
func ensureDatabaseExists(cfg *config.Config, log *logger.Logger) error {
	// Строка подключения к стандартной БД postgres (без указания имени БД)
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBSSLMode)

	// Подключаемся к БД postgres
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres database: %w", err)
	}
	defer db.Close()

	// Проверяем, существует ли БД
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = '%s')", cfg.DBName)
	err = db.QueryRow(query).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check database existence: %w", err)
	}

	if exists {
		log.Info("Database already exists", "db", cfg.DBName)
		return nil
	}

	// Если БД нет — создаём
	log.Info("Creating database", "db", cfg.DBName)
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DBName))
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	log.Info("Database created successfully", "db", cfg.DBName)
	return nil
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

const (
	addr    = ":8080"  // адрес сервера
	maxSize = 10000000 // будем растить слайс до 10 миллионов элементов
)

func foo() {
	// полезная нагрузка
	for {
		var s []int
		for i := 0; i < maxSize; i++ {
			s = append(s, i)
		}
	}
}
