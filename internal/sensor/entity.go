package sensor

import (
	"time"

	"mertani/internal/shared/id"
)

type Sensor struct {
	ID         id.ID
	DeviceID   id.ID
	SensorType string
	Value      float64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
