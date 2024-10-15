package errors

import "errors"

var (
	ErrFirestoreConnectionFailed = errors.New("failed to connect to the firestore database")
	ErrFirestoreProjectIDMissing = errors.New("firestore project id is missing")
)
