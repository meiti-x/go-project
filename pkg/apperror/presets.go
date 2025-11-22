package apperror

import "net/http"

var (
	ErrNotFound     = New("NOT_FOUND", "Not found", http.StatusNotFound)
	ErrValidation   = New("VALIDATION_ERROR", "Validation error", http.StatusBadRequest)
	ErrUnauthorized = New("UNAUTHORIZED", "Unauthorized Request", http.StatusUnauthorized)
	ErrBadRequest   = New("BAD_REQUEST", "Invalid request param/body", http.StatusBadRequest)
	ErrForbidden    = New("FORBIDDEN", "Forbidden request", http.StatusForbidden)
	ErrServer       = New("SERVER", "Internal Server error", http.StatusInternalServerError)
)

func ResolveError(statusCode int) error {

	if statusCode == http.StatusNotFound {
		return ErrNotFound
	}
	if statusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if statusCode == http.StatusForbidden {
		return ErrForbidden
	}
	if statusCode >= http.StatusInternalServerError {
		return ErrServer
	}

	return ErrBadRequest
}
