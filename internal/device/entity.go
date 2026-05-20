package device

import (
	"time"

	"mertani/internal/shared/id"
)

type Device struct {
	ID        id.ID
	Name      string
	Location  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
