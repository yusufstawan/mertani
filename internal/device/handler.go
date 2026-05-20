package device

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"

	"mertani/internal/shared/apperror"
	"mertani/internal/shared/id"
	"mertani/internal/shared/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(group *echo.Group) {
	group.POST("/devices", h.Create)
	group.GET("/devices", h.FindAll)
	group.GET("/devices/:id", h.FindByID)
	group.PUT("/devices/:id", h.Update)
	group.DELETE("/devices/:id", h.Delete)
}

func (h *Handler) Create(c *echo.Context) error {
	var request CreateDeviceRequest
	if err := bindAndValidate(c, &request); err != nil {
		return err
	}

	device, err := h.service.Create(c.Request().Context(), request)
	if err != nil {
		return err
	}

	return response.Created(c, "Device created successfully", ToResponse(device))
}

func (h *Handler) FindAll(c *echo.Context) error {
	devices, err := h.service.FindAll(c.Request().Context())
	if err != nil {
		return err
	}

	return response.OK(c, "Devices retrieved successfully", ToResponses(devices))
}

func (h *Handler) FindByID(c *echo.Context) error {
	deviceID, err := parsePathID(c)
	if err != nil {
		return err
	}

	device, err := h.service.FindByID(c.Request().Context(), deviceID)
	if err != nil {
		return err
	}

	return response.OK(c, "Device retrieved successfully", ToResponse(device))
}

func (h *Handler) Update(c *echo.Context) error {
	deviceID, err := parsePathID(c)
	if err != nil {
		return err
	}

	var request UpdateDeviceRequest
	if err := bindAndValidate(c, &request); err != nil {
		return err
	}

	device, err := h.service.Update(c.Request().Context(), deviceID, request)
	if err != nil {
		return err
	}

	return response.OK(c, "Device updated successfully", ToResponse(device))
}

func (h *Handler) Delete(c *echo.Context) error {
	deviceID, err := parsePathID(c)
	if err != nil {
		return err
	}

	if err := h.service.Delete(c.Request().Context(), deviceID); err != nil {
		return err
	}

	return response.OK(c, "Device deleted successfully", nil)
}

func bindAndValidate(c *echo.Context, request any) error {
	if err := c.Bind(request); err != nil {
		return apperror.BadRequest("Invalid request body", nil)
	}
	if err := c.Validate(request); err != nil {
		return apperror.BadRequest("Validation error", validationErrors(err))
	}

	return nil
}

func parsePathID(c *echo.Context) (id.ID, error) {
	parsedID, err := id.Parse(c.Param("id"))
	if err != nil {
		return id.ID{}, apperror.BadRequest("Invalid id parameter", map[string]string{
			"id": "id must be a valid UUID",
		})
	}

	return parsedID, nil
}

func validationErrors(err error) map[string]string {
	var validatorErrors validator.ValidationErrors
	if !errors.As(err, &validatorErrors) {
		return map[string]string{
			"validation": err.Error(),
		}
	}

	result := make(map[string]string, len(validatorErrors))
	for _, fieldError := range validatorErrors {
		field := strings.ToLower(fieldError.Field())
		result[field] = validationMessage(field, fieldError)
	}

	return result
}

func validationMessage(field string, fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return field + " is required"
	default:
		return field + " is invalid"
	}
}
