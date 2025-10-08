package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"net/http"
	"subs/subservice/internal/dto"
	"subs/subservice/internal/service"
	"time"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{service.InitMethod(db)}
}

// CreateSub godoc
// @Summary Create new subscription
// @Description Create a new subscription for user
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param input body dto.CreateSubscriptionRequest true "Subscription data"
// @Success 201 {object} dto.Subscription
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /create_sub [post]
func (h *Handler) HandleCreateSub(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSubscriptionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uid, err := uuid.Parse(req.UserID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	t, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m := int(t.Month())
	y := t.Year()

	sub := &dto.Subscription{
		ID:          uuid.New(),
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      uid,
		StartMonth:  m,
		StartYear:   y,
	}

	if err := h.svc.Create(*sub); err != nil {
		slog.Info("Create")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sub)
}

// HandleSumSubs godoc
// @Summary Get total subscription sum
// @Description Returns the total amount spent on subscriptions for a given user, service, and period.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param input body dto.SumSubsRequest true "Sum subscriptions request"
// @Success 200 {object} map[string]int64 "Total sum"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /sum_subs [post]
func (h *Handler) HandleSumSubs(w http.ResponseWriter, r *http.Request) {
	var req dto.SumSubsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.From == "" || req.To == "" {
		http.Error(w, "from and to required", http.StatusBadRequest)
		return
	}

	fromTime, err := time.Parse("2006-01", req.From)
	if err != nil {
		http.Error(w, "invalid from format", http.StatusBadRequest)
		return
	}

	toTime, err := time.Parse("2006-01", req.To)
	if err != nil {
		http.Error(w, "invalid to format", http.StatusBadRequest)
		return
	}

	total, err := h.svc.SumSubs(
		req.UserID,
		req.ServiceName,
		fromTime.Year(),
		int(fromTime.Month()),
		toTime.Year(),
		int(toTime.Month()),
	)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int64{"total": total})
}

// GetSubs godoc
// @Summary Get all subscriptions
// @Description Returns list of all subscriptions
// @Tags subscriptions
// @Produce json
// @Success 200 {array} dto.Subscription
// @Failure 500 {string} string "Internal Server Error"
// @Router /get_all_subs [get]
func (h *Handler) GetSubs(w http.ResponseWriter, r *http.Request) {
	subs, err := h.svc.GetSubs()
	if err != nil {
		http.Error(w, "failed to get subscriptions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subs)
}

// UpdateSub godoc
// @Summary Update subscription
// @Description Update subscription fields by ID (JSON body)
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param input body dto.UpdateSubscriptionRequest true "Updated data"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /update_sub [put]
func (h *Handler) UpdateSub(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		http.Error(w, "invalid or missing id", http.StatusBadRequest)
		return
	}

	endTime, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		http.Error(w, "invalid end_date format", http.StatusBadRequest)
		return
	}

	endMonth := int(endTime.Month())
	endYear := endTime.Year()

	sub := dto.Subscription{
		ID:          id,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		EndMonth:    &endMonth,
		EndYear:     &endYear,
	}

	if err := h.svc.Update(sub); err != nil {
		http.Error(w, "failed to update", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// DeleteSub godoc
// @Summary Delete subscription
// @Description Delete subscription by ID (from JSON body)
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param input body dto.DeleteSubscriptionRequest true "Subscription ID"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /delete_sub [delete]
func (h *Handler) DeleteSub(w http.ResponseWriter, r *http.Request) {
	var req dto.DeleteSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	if req.ID == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}

	if err := h.svc.Delete(id); err != nil {
		http.Error(w, "failed to delete", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}
