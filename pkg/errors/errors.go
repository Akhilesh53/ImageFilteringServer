package errors

import "errors"

var (
	ErrFirestoreConnectionFailed = errors.New("failed to connect to the firestore database")
	ErrFirestoreProjectIDMissing = errors.New("firestore project id is missing")
)

var (
	InternalError = Error{
		ErrorCode:    "API-001",
		ErrorMessage: "Internal Error",
		StatusCode:   500,
	}
)
