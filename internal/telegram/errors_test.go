package telegram

import "testing"

func TestAPIErrorDetectsBusinessConnectionInvalid(t *testing.T) {
	err := APIError{Description: "Bad Request: BUSINESS_CONNECTION_INVALID"}

	if !err.BusinessConnectionInvalid() {
		t.Fatal("business connection invalid was not detected")
	}
}
