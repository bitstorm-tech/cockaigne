package like

import (
	"github.com/google/uuid"
	"time"
)

type Like struct {
	UserId  uuid.UUID `gorm:"type:uuid;not null"`
	DealId  uuid.UUID `gorm:"type:uuid;not null"`
	Created time.Time `gorm:"type:timestamp with time zone;not null;default:now()"`
}
