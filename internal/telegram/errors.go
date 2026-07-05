package telegram

import (
	"fmt"
	"strings"
	"time"
)

type APIError struct {
	Description string
}

func (e APIError) Error() string {
	return fmt.Sprintf("telegram api: %s", e.Description)
}

func (e APIError) BusinessConnectionInvalid() bool {
	return strings.Contains(strings.ToUpper(e.Description), "BUSINESS_CONNECTION_INVALID")
}

type RetryError struct {
	Delay time.Duration
	Err   error
}

func (e RetryError) Error() string {
	return e.Err.Error()
}

func (e RetryError) Unwrap() error {
	return e.Err
}

func (e RetryError) RetryAfter() time.Duration {
	return e.Delay
}

func newAPIError(description string, retryAfter int) error {
	err := APIError{Description: description}
	if retryAfter <= 0 {
		return err
	}
	return RetryError{Delay: time.Duration(retryAfter) * time.Second, Err: err}
}
