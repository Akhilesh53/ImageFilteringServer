package errors

import "errors"

var (
	ErrFirestoreConnectionFailed = errors.New("failed to connect to the firestore database")
	ErrFirestoreProjectIDMissing = errors.New("firestore project id is missing")
	ErrURLNotPresent             = errors.New("URL not present in request")
	ErrEnvFilePathNotFound       = errors.New("environment file path not found")
)

var (
	InternalError = Error{
		ErrorCode:    "API-001-E",
		ErrorMessage: "Internal Error",
		StatusCode:   500,
	}

	RequestProcessSuccess = Error{
		ErrorCode:    "API-001-S",
		ErrorMessage: "Request processed successfully",
		StatusCode:   200,
	}

	URLNotPresent = Error{
		ErrorCode:    "API-002-E",
		ErrorMessage: "URL not present",
		StatusCode:   400,
	}

	InvalidRequest = Error{
		ErrorCode:    "API-003-E",
		ErrorMessage: "Invalid request",
		StatusCode:   400,
	}
)
