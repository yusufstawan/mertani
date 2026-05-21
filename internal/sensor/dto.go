package sensor

import "time"

type CreateSensorRequest struct {
	DeviceID   string   `json:"device_id" validate:"required,uuid"`
	SensorType string   `json:"sensor_type" validate:"required"`
	Value      *float64 `json:"value" validate:"required"`
}

type PatchSensorRequest struct {
	DeviceID   *string  `json:"device_id"`
	SensorType *string  `json:"sensor_type"`
	Value      *float64 `json:"value"`
}

type SensorResponse struct {
	ID         string    `json:"id"`
	DeviceID   string    `json:"device_id"`
	SensorType string    `json:"sensor_type"`
	Value      float64   `json:"value"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func ToResponse(sensor Sensor) SensorResponse {
	return SensorResponse{
		ID:         sensor.ID.String(),
		DeviceID:   sensor.DeviceID.String(),
		SensorType: sensor.SensorType,
		Value:      sensor.Value,
		CreatedAt:  sensor.CreatedAt,
		UpdatedAt:  sensor.UpdatedAt,
	}
}

func ToResponses(sensors []Sensor) []SensorResponse {
	responses := make([]SensorResponse, 0, len(sensors))
	for _, sensor := range sensors {
		responses = append(responses, ToResponse(sensor))
	}

	return responses
}
