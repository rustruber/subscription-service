package port

import (
	"context"
	"github.com/rustruber/subscription-service/internal/domain"
	"time"
)

// SubscriptionRepository порт в домен
type SubscriptionRepository interface {
	Create(ctx context.Context, sub *domain.Subscription) error
	GetByID(ctx context.Context, id string) (*domain.Subscription, error)
	Update(ctx context.Context, sub *domain.Subscription) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*domain.Subscription, int64, error)
	GetTotalCost(ctx context.Context, userID, serviceName string, startDate, endDate time.Time) (int, error)
}
