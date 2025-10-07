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

// todo refactor this method
//func parseMonthYear(s string) (int, int, error) {
//
//
//	parts := strings.Split(s, "-")
//	if len(parts) != 2 {
//		return 0, 0, fmt.Errorf("bad")
//	}
//	m, _ := strconv.Atoi(parts[0])
//	y, _ := strconv.Atoi(parts[1])
//	return m, y, nil
//}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{service.NewHandler(db)}
}

func (h *Handler) CreateSub(w http.ResponseWriter, r *http.Request) {
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

	//m, y, _ := parseMonthYear(req.StartDate)

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
