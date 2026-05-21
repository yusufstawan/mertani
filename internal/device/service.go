package device

import (
	"context"
	"errors"
	"strings"
	"time"

	"mertani/internal/shared/apperror"
	"mertani/internal/shared/id"
	"mertani/internal/shared/response"
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

func (s *Service) FindAll(ctx context.Context, params ListParams) ([]Device, response.Pagination, error) {
	devices, total, err := s.repository.FindAll(ctx, params)
	if err != nil {
		return nil, response.Pagination{}, apperror.Internal(err)
	}

	return devices, response.NewPagination(params.Page, params.Limit, total), nil
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

func (s *Service) Patch(ctx context.Context, deviceID id.ID, request PatchDeviceRequest) (Device, error) {
	validationErrors := make(map[string]string)
	hasUpdate := false

	if request.Name != nil {
		*request.Name = strings.TrimSpace(*request.Name)
		hasUpdate = true
		if *request.Name == "" {
			validationErrors["name"] = "name is required"
		}
	}
	if request.Location != nil {
		*request.Location = strings.TrimSpace(*request.Location)
		hasUpdate = true
		if *request.Location == "" {
			validationErrors["location"] = "location is required"
		}
	}
	if !hasUpdate {
		validationErrors["body"] = "at least one field is required"
	}
	if len(validationErrors) > 0 {
		return Device{}, apperror.BadRequest("Validation error", validationErrors)
	}

	device, err := s.repository.FindByID(ctx, deviceID)
	if errors.Is(err, ErrNotFound) {
		return Device{}, apperror.NotFound("Device not found")
	}
	if err != nil {
		return Device{}, apperror.Internal(err)
	}

	if request.Name != nil {
		device.Name = *request.Name
	}
	if request.Location != nil {
		device.Location = *request.Location
	}
	device.UpdatedAt = time.Now().UTC()

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
