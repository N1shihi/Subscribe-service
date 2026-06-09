package models

import (
	"github.com/google/uuid"
)

type Subscription struct {
	ServiceName string    `json:"service_name"`
	Price       int64     `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date,omitempty"`
}

type SubscriptionDelete struct {
	ServiceName string    `json:"service_name"`
	UserID      uuid.UUID `json:"user_id"`
}

type AggregateRequest struct {
	UserID      string `json:"user_id"`
	ServiceName string `json:"service_name"`
	StartDate   string `json:"from"` // "07-2025"
	EndDate     string `json:"to"`   // "08-2025"
}
