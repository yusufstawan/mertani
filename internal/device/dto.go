package device

import "time"

type CreateDeviceRequest struct {
	Name     string `json:"name" validate:"required"`
	Location string `json:"location" validate:"required"`
}

type PatchDeviceRequest struct {
	Name     *string `json:"name"`
	Location *string `json:"location"`
}

type DeviceResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToResponse(device Device) DeviceResponse {
	return DeviceResponse{
		ID:        device.ID.String(),
		Name:      device.Name,
		Location:  device.Location,
		CreatedAt: device.CreatedAt,
		UpdatedAt: device.UpdatedAt,
	}
}

func ToResponses(devices []Device) []DeviceResponse {
	responses := make([]DeviceResponse, 0, len(devices))
	for _, device := range devices {
		responses = append(responses, ToResponse(device))
	}

	return responses
}
