package dealer

import (
	"github.com/google/uuid"
	"time"
)

type Rating struct {
	DealerId uuid.UUID `db:"dealer_id"`
	UserId   uuid.UUID `db:"user_id"`
	Stars    int
	Text     string
	Username string
	CanEdit  bool
	Created  time.Time
}
