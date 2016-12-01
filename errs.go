package kong

import (
	"errors"
)

var (
	ErrMissingParameter = errors.New("Missing parameter")
	ErrCreateConsumerFailed = errors.New("Create consumer failed")
	ErrDeleteConsumerFailed = errors.New("Delete consumer failed")
	ErrCreateAPIKeyFailed = errors.New("Create API key failed")
	ErrUnknowMedthod = errors.New("Unknown method")
)