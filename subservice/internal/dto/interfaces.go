package dto

import (
	"github.com/google/uuid"
	"net/http"
)

type SubsHandlersMethods interface {
	HandleCreateSub(w http.ResponseWriter, r *http.Request)
	HandleSumSubs(w http.ResponseWriter, r *http.Request)
	GetSubs(w http.ResponseWriter, r *http.Request)
	UpdateSub(w http.ResponseWriter, r *http.Request)
	DeleteSub(w http.ResponseWriter, r *http.Request)
}

type SubsLogicMethods interface {
	Create(sub Subscription) error
	SumSubs(userID, serviceName string, fromYear, fromMonth, toYear, toMonth int) (int64, error)
	GetSubs() ([]Subscription, error)
	Update(sub Subscription) error
	Delete(id uuid.UUID) error
}
