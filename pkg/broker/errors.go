package broker

import "errors"

var (
	// ErrUnavailable Represents an error for requests made after the server has started shutting down.
	ErrUnavailable = errors.New("service is unavailable")
	// ErrInvalidID Indicates that the message with the provided ID is not valid or was never published.
	ErrInvalidID = errors.New("message with id provided is not valid or never published")
	// ErrExpiredID Indicates that the message with the provided ID has expired and is no longer available.
	ErrExpiredID = errors.New("message with id provided is expired")
)
