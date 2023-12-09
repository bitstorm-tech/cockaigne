package dealer

import (
	"github.com/google/uuid"
	"time"
)

type DealerRating struct {
	DealerId uuid.UUID `db:"dealer_id"`
	UserId   uuid.UUID `db:"user_id"`
	Stars    int
	Text     string
	Created  time.Time
}
