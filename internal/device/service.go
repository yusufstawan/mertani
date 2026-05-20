package device

import (
	"context"
	"errors"
	"strings"
	"time"

	"mertani/internal/shared/apperror"
	"mertani/internal/shared/id"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(ctx context.Context, request CreateDeviceRequest) (Device, error) {
	request.Name = strings.TrimSpace(request.Name)
	request.Location = strings.TrimSpace(request.Location)
	if validationErrors := validateDeviceInput(request.Name, request.Location); len(validationErrors) > 0 {
		return Device{}, apperror.BadRequest("Validation error", validationErrors)
	}

	now := time.Now().UTC()
	device := Device{
		ID:        id.New(),
		Name:      request.Name,
		Location:  request.Location,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repository.Create(ctx, &device); err != nil {
		return Device{}, apperror.Internal(err)
	}

	return device, nil
}

func (s *Service) FindAll(ctx context.Context) ([]Device, error) {
	devices, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, apperror.Internal(err)
	}

	return devices, nil
}

func (s *Service) FindByID(ctx context.Context, deviceID id.ID) (Device, error) {
	device, err := s.repository.FindByID(ctx, deviceID)
	if errors.Is(err, ErrNotFound) {
		return Device{}, apperror.NotFound("Device not found")
	}
	if err != nil {
		return Device{}, apperror.Internal(err)
	}

	return device, nil
}

func (s *Service) Update(ctx context.Context, deviceID id.ID, request UpdateDeviceRequest) (Device, error) {
	request.Name = strings.TrimSpace(request.Name)
	request.Location = strings.TrimSpace(request.Location)
	if validationErrors := validateDeviceInput(request.Name, request.Location); len(validationErrors) > 0 {
		return Device{}, apperror.BadRequest("Validation error", validationErrors)
	}

	device := Device{
		ID:        deviceID,
		Name:      request.Name,
		Location:  request.Location,
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.repository.Update(ctx, &device); errors.Is(err, ErrNotFound) {
		return Device{}, apperror.NotFound("Device not found")
	} else if err != nil {
		return Device{}, apperror.Internal(err)
	}

	return device, nil
}

func (s *Service) Delete(ctx context.Context, deviceID id.ID) error {
	if err := s.repository.Delete(ctx, deviceID); errors.Is(err, ErrNotFound) {
		return apperror.NotFound("Device not found")
	} else if err != nil {
		return apperror.Internal(err)
	}

	return nil
}

func validateDeviceInput(name string, location string) map[string]string {
	errors := make(map[string]string)
	if name == "" {
		errors["name"] = "name is required"
	}
	if location == "" {
		errors["location"] = "location is required"
	}

	return errors
}
