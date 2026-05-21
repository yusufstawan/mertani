package response

import (
	"errors"
	"testing"

	"github.com/labstack/echo/v5"

	"mertani/internal/shared/apperror"
)

func TestErrorResponseMapsAppErrors(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		status  int
		message string
	}{
		{
			name:    "bad request",
			err:     apperror.BadRequest("Validation error", map[string]string{"name": "required"}),
			status:  StatusBadRequest,
			message: "Validation error",
		},
		{
			name:    "not found",
			err:     apperror.NotFound("Device not found"),
			status:  StatusNotFound,
			message: "Device not found",
		},
		{
			name:    "internal",
			err:     apperror.Internal(errors.New("database timeout")),
			status:  StatusInternalServerError,
			message: "Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, body := ErrorResponse(tt.err)

			if status != tt.status {
				t.Fatalf("expected status %d, got %d", tt.status, status)
			}
			if body.Success {
				t.Fatal("expected error response")
			}
			if body.Message != tt.message {
				t.Fatalf("expected message %q, got %q", tt.message, body.Message)
			}
		})
	}
}

func TestErrorResponseMapsEchoErrors(t *testing.T) {
	status, body := ErrorResponse(echo.NewHTTPError(StatusBadRequest, "Invalid JSON body"))

	if status != StatusBadRequest {
		t.Fatalf("expected status %d, got %d", StatusBadRequest, status)
	}
	if body.Message != "Invalid JSON body" {
		t.Fatalf("expected echo error message, got %q", body.Message)
	}
}

func TestErrorResponseDefaultsUnknownErrorToInternalServerError(t *testing.T) {
	status, body := ErrorResponse(errors.New("unexpected failure"))

	if status != StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", StatusInternalServerError, status)
	}
	if body.Message != "Internal server error" {
		t.Fatalf("expected generic internal message, got %q", body.Message)
	}
}

func TestNewPaginationCalculatesTotalPages(t *testing.T) {
	pagination := NewPagination(2, 10, 25)

	if pagination.Page != 2 {
		t.Fatalf("expected page 2, got %d", pagination.Page)
	}
	if pagination.Limit != 10 {
		t.Fatalf("expected limit 10, got %d", pagination.Limit)
	}
	if pagination.Total != 25 {
		t.Fatalf("expected total 25, got %d", pagination.Total)
	}
	if pagination.TotalPages != 3 {
		t.Fatalf("expected total pages 3, got %d", pagination.TotalPages)
	}
}
