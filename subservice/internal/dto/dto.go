package dto

import "time"

type SubsMethods interface {
	CreateSubscription()
	UpdateSubscription(id int)
	GetSubscription()
	DeleteSubscription()
}

type Subscription struct {
	Name       string    `json:"service_name"`
	Price      int       `json:"price"`
	UserID     int       `json:"user_id"`
	StartDate  time.Time `json:"start_date"`
	EndingDate time.Time `json:"ending_ate"` // todo опционально
}
