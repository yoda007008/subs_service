package service

import (
	"github.com/jmoiron/sqlx"
	"subs/subservice/internal/dto"
)

type Service struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Create(sub dto.Subscription) error {
	q := `INSERT INTO subscriptions (id, service_name, price, user_id, start_month, start_year, end_month, end_year, created_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,now())`
	_, err := s.db.Exec(q, sub.ID, sub.ServiceName, sub.Price, sub.UserID, sub.StartMonth, sub.StartYear, sub.EndMonth, sub.EndYear)
	return err
}
