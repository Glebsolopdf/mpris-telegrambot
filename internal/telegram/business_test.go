package telegram

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestFindBusinessConnection(t *testing.T) {
	httpClient := &http.Client{Transport: roundTrip(func(r *http.Request) *http.Response {
		if r.URL.Path != "/bottoken/getUpdates" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		body := `{"ok":true,"result":[{"update_id":1,"business_connection":{"id":"bc_1","user":{"id":42},"is_enabled":true}}]}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}
	})}

	client := NewClientWithBaseURL("token", "", "https://telegram.test", httpClient)
	connection, err := client.FindBusinessConnection(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if connection.ID != "bc_1" || connection.UserID != 42 {
		t.Fatalf("connection = %+v", connection)
	}
}

func TestFindBusinessConnectionDisabled(t *testing.T) {
	httpClient := &http.Client{Transport: roundTrip(func(*http.Request) *http.Response {
		body := `{"ok":true,"result":[{"update_id":1,"business_connection":{"id":"bc_1","user":{"id":42},"is_enabled":false}}]}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}
	})}

	client := NewClientWithBaseURL("token", "", "https://telegram.test", httpClient)
	_, err := client.FindBusinessConnection(context.Background(), 42)
	if !errors.Is(err, ErrBusinessConnectionDisabled) {
		t.Fatalf("err = %v, want disabled", err)
	}
}

func TestFindBusinessConnectionUsesLatestUpdate(t *testing.T) {
	httpClient := &http.Client{Transport: roundTrip(func(*http.Request) *http.Response {
		body := `{"ok":true,"result":[` +
			`{"update_id":2,"business_connection":{"id":"new","user":{"id":42},"is_enabled":true}},` +
			`{"update_id":1,"business_connection":{"id":"old","user":{"id":42},"is_enabled":false}}` +
			`]}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}
	})}

	client := NewClientWithBaseURL("token", "", "https://telegram.test", httpClient)
	connection, err := client.FindBusinessConnection(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if connection.ID != "new" {
		t.Fatalf("connection id = %q, want new", connection.ID)
	}
}
