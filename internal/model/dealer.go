package model

import (
	"time"

	"github.com/google/uuid"
)

type DealerRating struct {
	DealerId uuid.UUID `db:"dealer_id"`
	UserId   uuid.UUID `db:"user_id"`
	Stars    int
	Text     string
	Username string
	CanEdit  bool
	Created  time.Time
}
