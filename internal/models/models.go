package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	Service_name string
	Price        int64
	User_id      uuid.UUID
	Start_date   time.Time // В тз месяц-год 07-2025
	End_date     time.Time //
}

type SubscriptionDelete struct {
	Service_name string
	User_id      uuid.UUID
}
