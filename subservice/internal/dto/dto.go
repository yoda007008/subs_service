package dto

import (
	"database/sql"
	"github.com/google/uuid"
)

// для удобства тестирования функционала приложения
type SubsMethods interface {
	CreateSubscription()
	UpdateSubscription(id int)
	GetSubscription()
	DeleteSubscription(id int)
}

type Subscription struct {
	ID          uuid.UUID     `db:"id" json:"id"`
	ServiceName string        `db:"service_name" json:"service_name"`
	Price       int           `db:"price" json:"price"`
	UserID      uuid.UUID     `db:"user_id" json:"user_id"`
	StartYear   int           `db:"start_year" json:"start_year"`
	StartMonth  int           `db:"start_month" json:"start_month"`
	EndYear     sql.NullInt64 `db:"end_year" json:"end_year,omitempty"`
	EndMonth    sql.NullInt64 `db:"end_month" json:"end_month,omitempty"`
}

type CreateSubscriptionRequest struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date,omitempty"`
}
