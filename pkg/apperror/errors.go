package apperror

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

type ErrorWithStatus struct {
	Err        *AppError
	StatusCode int
}

func (e *ErrorWithStatus) Error() string {
	return e.Err.Error()
}

func (e *ErrorWithStatus) Unwrap() error {
	return e.Err
}

func New(code, message string, status int) *ErrorWithStatus {
	return &ErrorWithStatus{
		Err: &AppError{
			Code:    code,
			Message: message,
		},
		StatusCode: status,
	}
}
