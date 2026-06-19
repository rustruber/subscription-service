package domain

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

// Subscription представляет подписку пользователя на сервис.
// Содержит информацию о названии сервиса, цене и периоде действия
type Subscription struct {
	ID          string `json:"id,omitempty"`
	ServiceName string `json:"service_name"`
	// БИЗНЕС-ПРАВИЛО: Стоимость подписки — целое число рублей (копейки не учитываются)
	// Price - стоимость подписки в рублях (целое число)
	Price     int        `json:"price"`
	UserID    string     `json:"user_id"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// NewSubscription конструктор подписки
func NewSubscription(
	serviceName string,
	price int,
	userID string,
	startDate time.Time,
	endDate *time.Time,
) (*Subscription, error) {
	// Валидация (бизнес-правила)
	if serviceName == "" {
		return nil, errors.New("service name is required")
	}
	if price <= 0 {
		return nil, errors.New("price must be greater than 0")
	}
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if startDate.IsZero() {
		return nil, errors.New("start date is required")
	}
	if endDate != nil && endDate.Before(startDate) {
		return nil, errors.New("end date must be after start date")
	}

	now := time.Now()
	return &Subscription{
		ID:          uuid.New().String(), // Генерируем UUID
		ServiceName: serviceName,
		Price:       price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// IsActive вычисляет активна ли подписка прямо сейчас
func (s *Subscription) IsActive() bool {
	now := time.Now()

	// Если end_date не указан — подписка бессрочная (всегда активна)
	if s.EndDate == nil {
		return true
	}

	// Активна, если сейчас между start_date и end_date
	return now.After(s.StartDate) && now.Before(*s.EndDate)
}

// IsExpired просрочена ли подписка
func (s *Subscription) IsExpired() bool {
	if s.EndDate == nil {
		return false
	}
	return time.Now().After(*s.EndDate)
}

// CanBeCancelled можно ли отменить
func (s *Subscription) CanBeCancelled() bool {
	return s.IsActive()
}

// CanBeExtended можно ли продлить
func (s *Subscription) CanBeExtended() bool {
	if s.EndDate == nil {
		return true // Бессрочную можно продлить
	}
	return !s.IsExpired() || time.Now().Before(s.EndDate.AddDate(0, 0, 30))
}

// Extend продление подписки
func (s *Subscription) Extend(months int) error {
	if months <= 0 {
		return errors.New("extension period must be positive")
	}
	if s.EndDate == nil {
		return errors.New("cannot extend perpetual subscription")
	}

	newEndDate := s.EndDate.AddDate(0, months, 0)
	s.EndDate = &newEndDate
	return nil
}

// Cancel отмена подписки
func (s *Subscription) Cancel() error {
	if !s.CanBeCancelled() {
		return errors.New("subscription cannot be cancelled")
	}
	// Логика отмены (например, установка EndDate = сегодня)
	now := time.Now()
	s.EndDate = &now
	return nil
}

// Validate стоимость подписки должна быть положительной
func (s *Subscription) Validate() error {
	if s.ServiceName == "" {
		return errors.New("service name is required")
	}
	if s.Price <= 0 {
		return errors.New("price must be greater than 0")
	}
	if s.UserID == "" {
		return errors.New("user ID is required")
	}
	if s.StartDate.IsZero() {
		return errors.New("start date is required")
	}
	if s.EndDate != nil && s.EndDate.Before(s.StartDate) {
		return errors.New("end date must be after start date")
	}
	return nil
}

// DaysUntilExpiration возвращает количество дней до окончания
func (s *Subscription) DaysUntilExpiration() int {
	if s.EndDate == nil {
		return -1 // Бессрочная
	}
	days := int(time.Until(*s.EndDate).Hours() / 24)
	if days < 0 {
		return 0 // Уже просрочена
	}
	return days
}
