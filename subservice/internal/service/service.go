package service

import (
	"github.com/jmoiron/sqlx"
	"strconv"
	"subs/subservice/internal/dto"
)

type Service struct {
	db *sqlx.DB
}

func InitMethod(db *sqlx.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Create(sub dto.Subscription) error {
	q := `INSERT INTO subscriptions (id, service_name, price, user_id, start_month, start_year, end_month, end_year, created_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,now())`
	_, err := s.db.Exec(q, sub.ID, sub.ServiceName, sub.Price, sub.UserID, sub.StartMonth, sub.StartYear, sub.EndMonth, sub.EndYear)
	return err
}

func (s *Service) SumSubs(userID, serviceName string, fromYear, fromMonth, toYear, toMonth int) (int64, error) {
	q := `SELECT price, start_month, start_year, end_month, end_year FROM subscriptions WHERE 1=1`
	args := []interface{}{}

	if userID != "" {
		q += " AND user_id = $" + strconv.Itoa(len(args)+1)
		args = append(args, userID)
	}
	if serviceName != "" {
		q += " AND service_name ILIKE $" + strconv.Itoa(len(args)+1)
		args = append(args, "%"+serviceName+"%")
	}

	rows := []dto.RowUpdate{}
	if err := s.db.Select(&rows, q, args...); err != nil {
		return 0, err
	}

	fromIdx := fromYear*12 + fromMonth - 1
	toIdx := toYear*12 + toMonth - 1

	var total int64 = 0
	for _, rec := range rows {
		startIdx := rec.StartYear*12 + rec.StartMonth - 1
		endIdx := startIdx
		if rec.EndYear.Valid && rec.EndMonth.Valid {
			endIdx = int(rec.EndYear.Int64)*12 + int(rec.EndMonth.Int64) - 1
		}

		overlapStart := max(startIdx, fromIdx)
		overlapEnd := min(endIdx, toIdx)
		if overlapEnd >= overlapStart {
			months := overlapEnd - overlapStart + 1
			total += int64(months) * int64(rec.Price)
		}
	}

	return total, nil
}

func (s *Service) GetSubs() ([]dto.Subscription, error) {
	q := `SELECT id, service_name, price, user_id, start_month, start_year, end_month, end_year, created_at 
	      FROM subscriptions`

	var subs []dto.Subscription
	if err := s.db.Select(&subs, q); err != nil {
		return nil, err
	}
	return subs, nil
}
