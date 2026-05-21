package sensor

import (
	"context"
	"errors"
	"strings"
	"time"

	"mertani/internal/device"
	"mertani/internal/shared/apperror"
	"mertani/internal/shared/id"
	"mertani/internal/shared/response"
)

type Service struct {
	repository       Repository
	deviceRepository device.Repository
}

func NewService(repository Repository, deviceRepository device.Repository) *Service {
	return &Service{
		repository:       repository,
		deviceRepository: deviceRepository,
	}
}

func (s *Service) Create(ctx context.Context, request CreateSensorRequest) (Sensor, error) {
	sensorInput, err := s.validateInput(ctx, request.DeviceID, request.SensorType, request.Value)
	if err != nil {
		return Sensor{}, err
	}

	now := time.Now().UTC()
	sensor := Sensor{
		ID:         id.New(),
		DeviceID:   sensorInput.deviceID,
		SensorType: sensorInput.sensorType,
		Value:      sensorInput.value,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.repository.Create(ctx, &sensor); err != nil {
		return Sensor{}, apperror.Internal(err)
	}

	return sensor, nil
}

func (s *Service) FindAll(ctx context.Context, params ListParams) ([]Sensor, response.Pagination, error) {
	sensors, total, err := s.repository.FindAll(ctx, params)
	if err != nil {
		return nil, response.Pagination{}, apperror.Internal(err)
	}

	return sensors, response.NewPagination(params.Page, params.Limit, total), nil
}

func (s *Service) FindByID(ctx context.Context, sensorID id.ID) (Sensor, error) {
	sensor, err := s.repository.FindByID(ctx, sensorID)
	if errors.Is(err, ErrNotFound) {
		return Sensor{}, apperror.NotFound("Sensor not found")
	}
	if err != nil {
		return Sensor{}, apperror.Internal(err)
	}

	return sensor, nil
}

func (s *Service) Patch(ctx context.Context, sensorID id.ID, request PatchSensorRequest) (Sensor, error) {
	validationErrors := make(map[string]string)
	hasUpdate := false

	var deviceID id.ID
	hasDeviceUpdate := false
	if request.DeviceID != nil {
		*request.DeviceID = strings.TrimSpace(*request.DeviceID)
		hasUpdate = true
		hasDeviceUpdate = true
		if *request.DeviceID == "" {
			validationErrors["device_id"] = "device_id is required"
		} else {
			parsedDeviceID, err := id.Parse(*request.DeviceID)
			if err != nil {
				validationErrors["device_id"] = "device_id must be a valid UUID"
			} else {
				deviceID = parsedDeviceID
			}
		}
	}

	if request.SensorType != nil {
		*request.SensorType = strings.TrimSpace(*request.SensorType)
		hasUpdate = true
		if *request.SensorType == "" {
			validationErrors["sensor_type"] = "sensor_type is required"
		}
	}
	if request.Value != nil {
		hasUpdate = true
	}
	if !hasUpdate {
		validationErrors["body"] = "at least one field is required"
	}
	if len(validationErrors) > 0 {
		return Sensor{}, apperror.BadRequest("Validation error", validationErrors)
	}

	sensor, err := s.repository.FindByID(ctx, sensorID)
	if errors.Is(err, ErrNotFound) {
		return Sensor{}, apperror.NotFound("Sensor not found")
	}
	if err != nil {
		return Sensor{}, apperror.Internal(err)
	}

	if hasDeviceUpdate {
		if _, err := s.deviceRepository.FindByID(ctx, deviceID); errors.Is(err, device.ErrNotFound) {
			return Sensor{}, apperror.NotFound("Device not found")
		} else if err != nil {
			return Sensor{}, apperror.Internal(err)
		}
		sensor.DeviceID = deviceID
	}
	if request.SensorType != nil {
		sensor.SensorType = *request.SensorType
	}
	if request.Value != nil {
		sensor.Value = *request.Value
	}
	sensor.UpdatedAt = time.Now().UTC()

	if err := s.repository.Update(ctx, &sensor); errors.Is(err, ErrNotFound) {
		return Sensor{}, apperror.NotFound("Sensor not found")
	} else if err != nil {
		return Sensor{}, apperror.Internal(err)
	}

	return sensor, nil
}

func (s *Service) Delete(ctx context.Context, sensorID id.ID) error {
	if err := s.repository.Delete(ctx, sensorID); errors.Is(err, ErrNotFound) {
		return apperror.NotFound("Sensor not found")
	} else if err != nil {
		return apperror.Internal(err)
	}

	return nil
}

type sensorInput struct {
	deviceID   id.ID
	sensorType string
	value      float64
}

func (s *Service) validateInput(ctx context.Context, deviceIDValue string, sensorType string, value *float64) (sensorInput, error) {
	deviceIDValue = strings.TrimSpace(deviceIDValue)
	sensorType = strings.TrimSpace(sensorType)

	validationErrors := make(map[string]string)
	if deviceIDValue == "" {
		validationErrors["device_id"] = "device_id is required"
	}
	if sensorType == "" {
		validationErrors["sensor_type"] = "sensor_type is required"
	}
	if value == nil {
		validationErrors["value"] = "value is required"
	}

	deviceID, err := id.Parse(deviceIDValue)
	if deviceIDValue != "" && err != nil {
		validationErrors["device_id"] = "device_id must be a valid UUID"
	}

	if len(validationErrors) > 0 {
		return sensorInput{}, apperror.BadRequest("Validation error", validationErrors)
	}

	if _, err := s.deviceRepository.FindByID(ctx, deviceID); errors.Is(err, device.ErrNotFound) {
		return sensorInput{}, apperror.NotFound("Device not found")
	} else if err != nil {
		return sensorInput{}, apperror.Internal(err)
	}

	return sensorInput{
		deviceID:   deviceID,
		sensorType: sensorType,
		value:      *value,
	}, nil
}
