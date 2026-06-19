package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/rustruber/subscription-service/internal/application/port"
	"github.com/rustruber/subscription-service/internal/domain"
)

type PostgresRepository struct {
	db     *sql.DB
	logger port.Logger
}

func NewPostgresRepository(db *sql.DB, logger port.Logger) port.SubscriptionRepository {
	return &PostgresRepository{
		db:     db,
		logger: logger,
	}
}

// Create — сохраняет подписку
func (r *PostgresRepository) Create(ctx context.Context, sub *domain.Subscription) error {
	r.logger.Debug("Creating subscription", "id", sub.ID)

	query := `
        INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `
	_, err := r.db.ExecContext(ctx, query,
		sub.ID,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
		sub.CreatedAt,
		sub.UpdatedAt,
	)
	return err
}

// GetByID — получает подписку по ID
func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*domain.Subscription, error) {
	r.logger.Debug("Getting subscription", "id", id)

	query := `
        SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
        FROM subscriptions WHERE id = $1
    `
	var sub domain.Subscription
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

// Update — обновляет подписку
func (r *PostgresRepository) Update(ctx context.Context, sub *domain.Subscription) error {
	r.logger.Debug("Updating subscription", "id", sub.ID)

	query := `
        UPDATE subscriptions 
        SET service_name = $1, price = $2, start_date = $3, end_date = $4, updated_at = $5
        WHERE id = $6
    `
	result, err := r.db.ExecContext(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.StartDate,
		sub.EndDate,
		time.Now(),
		sub.ID,
	)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// Delete — удаляет подписку
func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	r.logger.Debug("Deleting subscription", "id", id)

	result, err := r.db.ExecContext(ctx, "DELETE FROM subscriptions WHERE id = $1", id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// List — возвращает список с пагинацией
func (r *PostgresRepository) List(ctx context.Context, limit, offset int) ([]*domain.Subscription, int64, error) {
	r.logger.Debug("Listing subscriptions", "limit", limit, "offset", offset)

	var total int64
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM subscriptions").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
        SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
        FROM subscriptions
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var subs []*domain.Subscription
	for rows.Next() {
		var sub domain.Subscription
		err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
			&sub.CreatedAt,
			&sub.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		subs = append(subs, &sub)
	}

	return subs, total, nil
}

// GetTotalCost — считает сумму за период
func (r *PostgresRepository) GetTotalCost(ctx context.Context, userID, serviceName string, startDate, endDate time.Time) (int, error) {
	r.logger.Debug("Calculating total cost",
		"user_id", userID,
		"service", serviceName,
		"start", startDate,
		"end", endDate,
	)

	query := `SELECT COALESCE(SUM(price), 0) FROM subscriptions WHERE start_date >= $1 AND start_date <= $2`
	args := []interface{}{startDate, endDate}
	argIndex := 3

	if userID != "" {
		query += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, userID)
		argIndex++
	}
	if serviceName != "" {
		query += fmt.Sprintf(" AND LOWER(service_name) = LOWER($%d)", argIndex)
		args = append(args, serviceName)
	}

	var total int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&total)
	return total, err
}
