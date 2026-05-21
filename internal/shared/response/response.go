package response

import (
	"errors"
	stdhttp "net/http"

	"github.com/labstack/echo/v5"

	"mertani/internal/shared/apperror"
)

const (
	StatusOK                  = stdhttp.StatusOK
	StatusCreated             = stdhttp.StatusCreated
	StatusBadRequest          = stdhttp.StatusBadRequest
	StatusNotFound            = stdhttp.StatusNotFound
	StatusInternalServerError = stdhttp.StatusInternalServerError
)

type Body struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type PaginatedData struct {
	Items      any        `json:"items"`
	Pagination Pagination `json:"pagination"`
}

func Success(message string, data any) Body {
	return Body{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func Error(message string, errors any) Body {
	return Body{
		Success: false,
		Message: message,
		Errors:  errors,
	}
}

func NewPagination(page int, limit int, total int) Pagination {
	totalPages := 0
	if limit > 0 {
		totalPages = (total + limit - 1) / limit
	}

	return Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}

func NewPaginatedData(items any, pagination Pagination) PaginatedData {
	return PaginatedData{
		Items:      items,
		Pagination: pagination,
	}
}

func OK(c *echo.Context, message string, data any) error {
	return c.JSON(StatusOK, Success(message, data))
}

func Created(c *echo.Context, message string, data any) error {
	return c.JSON(StatusCreated, Success(message, data))
}

func BadRequest(c *echo.Context, message string, errors any) error {
	return c.JSON(StatusBadRequest, Error(message, errors))
}

func NotFound(c *echo.Context, message string) error {
	return c.JSON(StatusNotFound, Error(message, nil))
}

func InternalServerError(c *echo.Context) error {
	return c.JSON(StatusInternalServerError, Error("Internal server error", nil))
}

func FromError(c *echo.Context, err error) error {
	status, body := ErrorResponse(err)
	return c.JSON(status, body)
}

func ErrorHandler(c *echo.Context, err error) {
	if response, _ := echo.UnwrapResponse(c.Response()); response != nil && response.Committed {
		return
	}

	status, body := ErrorResponse(err)
	var responseErr error
	if c.Request().Method == stdhttp.MethodHead {
		responseErr = c.NoContent(status)
	} else {
		responseErr = c.JSON(status, body)
	}
	if responseErr != nil {
		c.Logger().Error("failed to send error response", "error", responseErr)
	}
}

func ErrorResponse(err error) (int, Body) {
	if appErr, ok := apperror.As(err); ok {
		return appErrorResponse(appErr)
	}

	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		message := httpErr.Message
		if message == "" {
			message = stdhttp.StatusText(httpErr.Code)
		}

		return httpErr.Code, Error(message, nil)
	}

	if status := echo.StatusCode(err); status != 0 {
		return status, Error(stdhttp.StatusText(status), nil)
	}

	return StatusInternalServerError, Error("Internal server error", nil)
}

func appErrorResponse(err *apperror.Error) (int, Body) {
	switch err.Code {
	case apperror.CodeBadRequest:
		return StatusBadRequest, Error(err.Message, err.Errors)
	case apperror.CodeNotFound:
		return StatusNotFound, Error(err.Message, err.Errors)
	case apperror.CodeInternal:
		return StatusInternalServerError, Error("Internal server error", nil)
	default:
		return StatusInternalServerError, Error("Internal server error", nil)
	}
}
