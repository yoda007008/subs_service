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
