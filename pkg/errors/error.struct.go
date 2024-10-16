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

// SetErrorCode is a function that sets the error code
func (e *Error) SetErrorCode(errorCode string) *Error {
	e.ErrorCode = errorCode
	return e
}

// SetErrorMessage is a function that sets the error message
func (e *Error) SetErrorMessage(errorMessage string) *Error {
	e.ErrorMessage = errorMessage
	return e
}

// SetStatusCode is a function that sets the status code
func (e *Error) SetStatusCode(statusCode int) *Error {
	e.StatusCode = statusCode
	return e
}

// SetUUID is a function that sets the UUID
func (e *Error) SetUUID(uuid string) *Error {
	e.UUID = uuid
	return e
}


