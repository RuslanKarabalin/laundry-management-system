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

type CreateReservation struct {
	UserId    string    `json:"userId"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}
