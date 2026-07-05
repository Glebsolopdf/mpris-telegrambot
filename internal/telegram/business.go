package telegram

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

var (
	ErrBusinessConnectionNotFound = errors.New("business connection update not found")
	ErrBusinessConnectionDisabled = errors.New("business connection is disabled")
)

type BusinessConnection struct {
	ID        string
	UserID    int64
	FirstName string
	LastName  string
	IsEnabled bool
}

func (c *Client) FindBusinessConnection(ctx context.Context, userID int64) (BusinessConnection, error) {
	form := url.Values{}
	form.Set("timeout", "0")
	form.Set("allowed_updates", `["business_connection"]`)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.methodURL("getUpdates"), strings.NewReader(form.Encode()))
	if err != nil {
		return BusinessConnection{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var updates []update
	if err := c.do(req, &updates); err != nil {
		return BusinessConnection{}, err
	}
	return findConnection(updates, userID)
}

func findConnection(updates []update, userID int64) (BusinessConnection, error) {
	var latest update
	found := false
	for _, item := range updates {
		conn := item.BusinessConnection
		if conn == nil || conn.User.ID != userID {
			continue
		}
		if !found || item.UpdateID > latest.UpdateID {
			latest = item
			found = true
		}
	}
	if !found {
		return BusinessConnection{}, ErrBusinessConnectionNotFound
	}
	conn := latest.BusinessConnection
	if !conn.IsEnabled {
		return BusinessConnection{}, ErrBusinessConnectionDisabled
	}
	return BusinessConnection{
		ID:        conn.ID,
		UserID:    conn.User.ID,
		FirstName: conn.User.FirstName,
		LastName:  conn.User.LastName,
		IsEnabled: true,
	}, nil
}

type update struct {
	UpdateID           int                 `json:"update_id"`
	BusinessConnection *businessConnection `json:"business_connection"`
}

type businessConnection struct {
	ID        string `json:"id"`
	User      user   `json:"user"`
	IsEnabled bool   `json:"is_enabled"`
}

type user struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
