package model

import (
	"time"

	"github.com/google/uuid"
)

type Like struct {
	UserId  uuid.UUID
	DealId  uuid.UUID
	Created time.Time
}
