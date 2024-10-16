package errors

type Error struct {
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	StatusCode   int    `json:"status_code"`
	UUID         string `json:"uuid"`
}

// NewError is a function that returns a new instance of Error
func NewError(errorCode string, errorMessage string, statusCode int, uuid string) *Error {
	return &Error{
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
		StatusCode:   statusCode,
		UUID:         uuid,
	}
}

// GetErrorCode is a function that returns the error code
func (e *Error) GetErrorCode() string {
	return e.ErrorCode
}

// GetErrorMessage is a function that returns the error message
func (e *Error) GetErrorMessage() string {
	return e.ErrorMessage
}

// GetStatusCode is a function that returns the status code
func (e *Error) GetStatusCode() int {
	return e.StatusCode
}

// GetUUID is a function that returns the UUID
func (e *Error) GetUUID() string {
	return e.UUID
}
