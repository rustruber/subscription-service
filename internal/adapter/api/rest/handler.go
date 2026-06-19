package rest

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rustruber/subscription-service/internal/application/subscription"
)

type Handler struct {
	useCase *subscription.UseCase
}

func NewHandler(useCase *subscription.UseCase) *Handler {
	return &Handler{useCase: useCase}
}

// Create godoc
// @Summary      Создать подписку
// @Description  Создаёт новую запись о подписке для пользователя
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        request body CreateSubscriptionRequest true "Данные подписки"
// @Success      201  {object}  domain.Subscription
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /subscriptions [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	startDate, _ := time.Parse("01-2006", req.StartDate)
	var endDate *time.Time
	if req.EndDate != nil && *req.EndDate != "" {
		t, _ := time.Parse("01-2006", *req.EndDate)
		endDate = &t
	}

	sub, err := h.useCase.Create(r.Context(), subscription.CreateRequest{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sub)
}

// List godoc
// @Summary      Получить список подписок
// @Description  Возвращает список подписок с пагинацией
// @Tags         subscriptions
// @Produce      json
// @Param        page  query int false "Номер страницы" default(1)
// @Param        limit query int false "Количество на странице" default(10)
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /subscriptions [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	subs, total, err := h.useCase.List(r.Context(), page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"subscriptions": subs,
		"total":         total,
		"page":          page,
		"limit":         limit,
	})
}

// GetByID godoc
// @Summary      Получить подписку по ID
// @Description  Возвращает данные подписки по её идентификатору
// @Tags         subscriptions
// @Produce      json
// @Param        id path string true "ID подписки (UUID)"
// @Success      200  {object}  domain.Subscription
// @Failure      404  {object}  map[string]string
// @Router       /subscriptions/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	sub, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
}

// Update godoc
// @Summary      Обновить подписку
// @Description  Обновляет данные существующей подписки
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id      path string true "ID подписки (UUID)"
// @Param        request body UpdateSubscriptionRequest true "Данные для обновления"
// @Success      200  {object}  domain.Subscription
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /subscriptions/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	startDate, _ := time.Parse("01-2006", req.StartDate)
	var endDate *time.Time
	if req.EndDate != nil && *req.EndDate != "" {
		t, _ := time.Parse("01-2006", *req.EndDate)
		endDate = &t
	}

	sub, err := h.useCase.Update(r.Context(), subscription.UpdateRequest{
		ID:          id,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		StartDate:   startDate,
		EndDate:     endDate,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
}

// Delete godoc
// @Summary      Удалить подписку
// @Description  Удаляет подписку по её идентификатору
// @Tags         subscriptions
// @Produce      json
// @Param        id path string true "ID подписки (UUID)"
// @Success      204  "No Content"
// @Failure      404  {object}  map[string]string
// @Router       /subscriptions/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.useCase.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetTotalCost godoc
// @Summary      Подсчёт суммарной стоимости
// @Description  Считает общую стоимость всех подписок за указанный период с фильтрацией
// @Tags         subscriptions
// @Produce      json
// @Param        user_id      query string false "ID пользователя для фильтрации"
// @Param        service_name query string false "Название сервиса для фильтрации"
// @Param        start_date   query string true "Начало периода (MM-YYYY)"
// @Param        end_date     query string true "Конец периода (MM-YYYY)"
// @Success      200  {object}  map[string]int
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /subscriptions/total-cost [get]
func (h *Handler) GetTotalCost(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	serviceName := r.URL.Query().Get("service_name")
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	if startDateStr == "" || endDateStr == "" {
		http.Error(w, "start_date and end_date are required", http.StatusBadRequest)
		return
	}

	startDate, _ := time.Parse("01-2006", startDateStr)
	endDate, _ := time.Parse("01-2006", endDateStr)

	total, err := h.useCase.GetTotalCost(r.Context(), subscription.TotalCostRequest{
		UserID:      userID,
		ServiceName: serviceName,
		StartDate:   startDate,
		EndDate:     endDate,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"total_cost": total})
}
