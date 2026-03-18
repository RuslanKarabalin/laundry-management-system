package model

import "github.com/google/uuid"

type Appliance struct {
	Id   uuid.UUID
	Name string
	Type string
}
