package subscription

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rustruber/subscription-service/internal/application/port"
	"github.com/rustruber/subscription-service/internal/domain"
)

// UseCase — бизнес-логика для подписок
type UseCase struct {
	repo   port.SubscriptionRepository
	logger port.Logger
}

// NewUseCase создаёт Use Case
func NewUseCase(repo port.SubscriptionRepository, logger port.Logger) *UseCase {
	return &UseCase{
		repo:   repo,
		logger: logger,
	}
}

// CreateRequest — данные для создания подписки
type CreateRequest struct {
	ServiceName string
	Price       int
	UserID      string
	StartDate   time.Time
	EndDate     *time.Time
}

// UpdateRequest — данные для обновления подписки
type UpdateRequest struct {
	ID          string
	ServiceName string
	Price       int
	StartDate   time.Time
	EndDate     *time.Time
}

// TotalCostRequest — данные для подсчёта стоимости
type TotalCostRequest struct {
	UserID      string
	ServiceName string
	StartDate   time.Time
	EndDate     time.Time
}

// --------------------- МЕТОДЫ USE CASE -----------------------

// Create — создание подписки
func (uc *UseCase) Create(ctx context.Context, req CreateRequest) (*domain.Subscription, error) {
	uc.logger.Info("Creating subscription", "user_id", req.UserID, "service", req.ServiceName)

	// Создаём доменную сущность
	sub, err := domain.NewSubscription(
		req.ServiceName,
		req.Price,
		req.UserID,
		req.StartDate,
		req.EndDate,
	)
	if err != nil {
		uc.logger.Error("Validation failed", "error", err)
		return nil, err
	}

	// Генерируем ID
	sub.ID = uuid.New().String()

	// Сохраняем через репозиторий
	if err := uc.repo.Create(ctx, sub); err != nil {
		uc.logger.Error("Failed to save subscription", "error", err)
		return nil, err
	}

	uc.logger.Info("Subscription created", "id", sub.ID)
	return sub, nil
}

// GetByID — получение подписки по ID
func (uc *UseCase) GetByID(ctx context.Context, id string) (*domain.Subscription, error) {
	uc.logger.Debug("Getting subscription", "id", id)

	sub, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to get subscription", "id", id, "error", err)
		return nil, err
	}

	return sub, nil
}

// Update — обновление подписки
func (uc *UseCase) Update(ctx context.Context, req UpdateRequest) (*domain.Subscription, error) {
	uc.logger.Info("Updating subscription", "id", req.ID)

	// Проверяем, существует ли подписка
	sub, err := uc.repo.GetByID(ctx, req.ID)
	if err != nil {
		uc.logger.Error("Subscription not found", "id", req.ID, "error", err)
		return nil, err
	}

	// Обновляем поля
	sub.ServiceName = req.ServiceName
	sub.Price = req.Price
	sub.StartDate = req.StartDate
	sub.EndDate = req.EndDate
	sub.UpdatedAt = time.Now()

	// Сохраняем
	if err := uc.repo.Update(ctx, sub); err != nil {
		uc.logger.Error("Failed to update subscription", "id", req.ID, "error", err)
		return nil, err
	}

	uc.logger.Info("Subscription updated", "id", sub.ID)
	return sub, nil
}

// Delete — удаление подписки
func (uc *UseCase) Delete(ctx context.Context, id string) error {
	uc.logger.Info("Deleting subscription", "id", id)

	if err := uc.repo.Delete(ctx, id); err != nil {
		uc.logger.Error("Failed to delete subscription", "id", id, "error", err)
		return err
	}

	uc.logger.Info("Subscription deleted", "id", id)
	return nil
}

// List — список подписок с пагинацией
func (uc *UseCase) List(ctx context.Context, page, limit int) ([]*domain.Subscription, int64, error) {
	uc.logger.Debug("Listing subscriptions", "page", page, "limit", limit)

	offset := (page - 1) * limit
	subs, total, err := uc.repo.List(ctx, limit, offset)
	if err != nil {
		uc.logger.Error("Failed to list subscriptions", "error", err)
		return nil, 0, err
	}

	return subs, total, nil
}

// GetTotalCost — подсчёт суммарной стоимости за период
func (uc *UseCase) GetTotalCost(ctx context.Context, req TotalCostRequest) (int, error) {
	uc.logger.Info("Calculating total cost",
		"user_id", req.UserID,
		"service", req.ServiceName,
		"start", req.StartDate,
		"end", req.EndDate,
	)

	if req.StartDate.After(req.EndDate) {
		return 0, domain.ErrInvalidPeriod
	}

	total, err := uc.repo.GetTotalCost(ctx, req.UserID, req.ServiceName, req.StartDate, req.EndDate)
	if err != nil {
		uc.logger.Error("Failed to calculate total cost", "error", err)
		return 0, err
	}

	uc.logger.Info("Total cost calculated", "total", total)
	return total, nil
}
