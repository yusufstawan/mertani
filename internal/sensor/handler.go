package sensor

import (
	"errors"
	"strconv"
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
	group.POST("/sensors", h.Create)
	group.GET("/sensors", h.FindAll)
	group.GET("/sensors/:id", h.FindByID)
	group.PATCH("/sensors/:id", h.Patch)
	group.DELETE("/sensors/:id", h.Delete)
}

func (h *Handler) Create(c *echo.Context) error {
	var request CreateSensorRequest
	if err := bindAndValidate(c, &request); err != nil {
		return err
	}

	sensor, err := h.service.Create(c.Request().Context(), request)
	if err != nil {
		return err
	}

	return response.Created(c, "Sensor created successfully", ToResponse(sensor))
}

func (h *Handler) FindAll(c *echo.Context) error {
	params, err := parseListParams(c)
	if err != nil {
		return err
	}

	sensors, pagination, err := h.service.FindAll(c.Request().Context(), params)
	if err != nil {
		return err
	}

	return response.OK(c, "Sensors retrieved successfully", response.NewPaginatedData(ToResponses(sensors), pagination))
}

func (h *Handler) FindByID(c *echo.Context) error {
	sensorID, err := parsePathID(c)
	if err != nil {
		return err
	}

	sensor, err := h.service.FindByID(c.Request().Context(), sensorID)
	if err != nil {
		return err
	}

	return response.OK(c, "Sensor retrieved successfully", ToResponse(sensor))
}

func (h *Handler) Patch(c *echo.Context) error {
	sensorID, err := parsePathID(c)
	if err != nil {
		return err
	}

	var request PatchSensorRequest
	if err := bindAndValidate(c, &request); err != nil {
		return err
	}

	sensor, err := h.service.Patch(c.Request().Context(), sensorID, request)
	if err != nil {
		return err
	}

	return response.OK(c, "Sensor updated successfully", ToResponse(sensor))
}

func (h *Handler) Delete(c *echo.Context) error {
	sensorID, err := parsePathID(c)
	if err != nil {
		return err
	}

	if err := h.service.Delete(c.Request().Context(), sensorID); err != nil {
		return err
	}

	return response.OK(c, "Sensor deleted successfully", nil)
}

func parseListParams(c *echo.Context) (ListParams, error) {
	page, err := parsePositiveQueryInt(c.QueryParam("page"), 1, "page")
	if err != nil {
		return ListParams{}, err
	}

	limit, err := parsePositiveQueryInt(c.QueryParam("limit"), 10, "limit")
	if err != nil {
		return ListParams{}, err
	}
	if limit > 100 {
		return ListParams{}, apperror.BadRequest("Invalid query parameter", map[string]string{
			"limit": "limit must be less than or equal to 100",
		})
	}

	return ListParams{
		Page:   page,
		Limit:  limit,
		Search: strings.TrimSpace(c.QueryParam("search")),
	}, nil
}

func parsePositiveQueryInt(value string, defaultValue int, field string) (int, error) {
	if value == "" {
		return defaultValue, nil
	}

	parsedValue, err := strconv.Atoi(value)
	if err != nil || parsedValue < 1 {
		return 0, apperror.BadRequest("Invalid query parameter", map[string]string{
			field: field + " must be a positive number",
		})
	}

	return parsedValue, nil
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
		field := requestFieldName(fieldError.Field())
		result[field] = validationMessage(field, fieldError)
	}

	return result
}

func requestFieldName(field string) string {
	switch field {
	case "DeviceID":
		return "device_id"
	case "SensorType":
		return "sensor_type"
	default:
		return strings.ToLower(field)
	}
}

func validationMessage(field string, fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return field + " is required"
	case "uuid":
		return field + " must be a valid UUID"
	default:
		return field + " is invalid"
	}
}
