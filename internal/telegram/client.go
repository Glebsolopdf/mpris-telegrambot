package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const defaultBaseURL = "https://api.telegram.org"

type Client struct {
	token                string
	businessConnectionID string
	baseURL              string
	http                 *http.Client
}

func NewClient(token string, businessConnectionID string, httpClient *http.Client) *Client {
	return NewClientWithBaseURL(token, businessConnectionID, defaultBaseURL, httpClient)
}

func NewClientWithBaseURL(token, businessConnectionID, baseURL string, httpClient *http.Client) *Client {
	return &Client{
		token:                token,
		businessConnectionID: businessConnectionID,
		baseURL:              strings.TrimRight(baseURL, "/"),
		http:                 httpClient,
	}
}

func (c *Client) SetBio(ctx context.Context, bio string) error {
	form := url.Values{}
	form.Set("business_connection_id", c.businessConnectionID)
	form.Set("bio", bio)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.methodURL("setBusinessAccountBio"), strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return c.do(req, nil)
}

func (c *Client) SetName(ctx context.Context, firstName string, lastName string) error {
	form := url.Values{}
	form.Set("business_connection_id", c.businessConnectionID)
	form.Set("first_name", firstName)
	if lastName != "" {
		form.Set("last_name", lastName)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.methodURL("setBusinessAccountName"), strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return c.do(req, nil)
}

func (c *Client) methodURL(method string) string {
	return fmt.Sprintf("%s/bot%s/%s", c.baseURL, c.token, method)
}

func (c *Client) do(req *http.Request, result any) error {
	res, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	payload, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return readErr
	}

	var api response
	if len(payload) > 0 {
		_ = json.Unmarshal(payload, &api)
	}
	if res.StatusCode == http.StatusTooManyRequests {
		return retryFromResponse(res, api)
	}
	if res.StatusCode >= 300 || !api.OK {
		return newAPIError(description(api, res.Status), api.Parameters.RetryAfter)
	}
	if result != nil && len(api.Result) > 0 {
		return json.Unmarshal(api.Result, result)
	}
	return nil
}

func retryFromResponse(res *http.Response, api response) error {
	delay := api.Parameters.RetryAfter
	if delay == 0 {
		delay, _ = strconv.Atoi(res.Header.Get("Retry-After"))
	}
	return newAPIError(description(api, res.Status), delay)
}

func description(api response, fallback string) string {
	if api.Description != "" {
		return api.Description
	}
	return fallback
}
