package model

import (
	"time"

	"github.com/google/uuid"
)

type Reservation struct {
	Id          uuid.UUID
	ApplianceId uuid.UUID
	UserId      string
	StartTime   time.Time
	EndTime     time.Time
}
