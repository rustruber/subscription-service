package rest

// CreateSubscriptionRequest — запрос на создание подписки
type CreateSubscriptionRequest struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}

// UpdateSubscriptionRequest — запрос на обновление подписки
type UpdateSubscriptionRequest struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}
