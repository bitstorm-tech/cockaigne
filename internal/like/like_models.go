package like

import (
	"github.com/google/uuid"
	"time"
)

type Like struct {
	UserId  uuid.UUID
	DealId  uuid.UUID
	Created time.Time
}
