package errs

import "net/http"

type ErrorMessage interface {
	Error() string
	StatusCode() int
	Message() string
}

type ErrorData struct {
	ErrError      string `json:"error"`
	ErrStatusCode int    `json:"statusCode"`
	ErrMessage    string `json:"message"`
}

func (e *ErrorData) Error() string {
	return e.ErrError
}

func (e *ErrorData) StatusCode() int {
	return e.ErrStatusCode
}

func (e *ErrorData) Message() string {
	return e.ErrMessage
}

func NewBadRequestError(message string) ErrorMessage {
	return &ErrorData{
		ErrError:      "BAD_REQUEST",
		ErrStatusCode: http.StatusBadRequest,
		ErrMessage:    message,
	}
}

func NewUnauthenticatedError(message string) ErrorMessage {
	return &ErrorData{
		ErrError:      "NOT_AUTHENTICATED",
		ErrStatusCode: http.StatusUnauthorized,
		ErrMessage:    message,
	}
}

func NewUnauthorizedError(message string) ErrorMessage {
	return &ErrorData{
		ErrError:      "NOT_AUTHORIZED",
		ErrStatusCode: http.StatusForbidden,
		ErrMessage:    message,
	}
}

func NewNotFoundError(message string) ErrorMessage {
	return &ErrorData{
		ErrError:      "NOT_FOUND",
		ErrStatusCode: http.StatusNotFound,
		ErrMessage:    message,
	}
}

func NewMethodNotAllowedError(message string) ErrorMessage {
	return &ErrorData{
		ErrError:      "METHOD_NOT_ALLOWED",
		ErrStatusCode: http.StatusMethodNotAllowed,
		ErrMessage:    message,
	}
}

func NewUnsupportedMediaTypeError(message string) ErrorMessage {
	return &ErrorData{
		ErrError:      "UNSUPPORTED_MEDIA_TYPE",
		ErrStatusCode: http.StatusUnsupportedMediaType,
		ErrMessage:    message,
	}
}

func NewUnprocessableEntityError(message string) ErrorMessage {
	return &ErrorData{
		ErrError:      "INVALID_REQUEST_BODY",
		ErrStatusCode: http.StatusUnprocessableEntity,
		ErrMessage:    message,
	}
}

func NewInternalServerError(message string) ErrorMessage {
	return &ErrorData{
		ErrError:      "INTERNAL_SERVER_ERROR",
		ErrStatusCode: http.StatusInternalServerError,
		ErrMessage:    message,
	}
}
