package domain

import "time"

// Subscription представляет подписку пользователя на сервис.
// Содержит информацию о названии сервиса, цене и периоде действия.
type Subscription struct {
	ID          string
	ServiceName string
	Price       int
	UserID      string
	StartDate   time.Time
	EndDate     *time.Time
}

// IsActive возвращает true, если подписка активна на текущий момент.
func (s *Subscription) IsActive() bool {
	// TODO: реализовать метод
	return false
}
