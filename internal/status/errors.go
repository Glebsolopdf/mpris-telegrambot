package status

import (
	"errors"
	"time"
)

type RetryAfter interface {
	RetryAfter() time.Duration
}

type BusinessConnectionInvalid interface {
	BusinessConnectionInvalid() bool
}

func retryDelay(err error) time.Duration {
	if retry, ok := err.(RetryAfter); ok {
		return retry.RetryAfter()
	}
	return 0
}

func IsBusinessConnectionInvalid(err error) bool {
	var invalid BusinessConnectionInvalid
	return errors.As(err, &invalid) && invalid.BusinessConnectionInvalid()
}
