package telegram

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) SetAvatar(ctx context.Context, jpeg []byte) error {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	_ = writer.WriteField("business_connection_id", c.businessConnectionID)
	_ = writer.WriteField("is_public", "false")
	_ = writer.WriteField("photo", `{"type":"static","photo":"attach://avatar"}`)

	part, err := writer.CreateFormFile("avatar", "avatar.jpg")
	if err != nil {
		return err
	}
	if _, err := part.Write(jpeg); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.methodURL("setBusinessAccountProfilePhoto"), &body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return c.do(req, nil)
}

func (c *Client) RemoveAvatar(ctx context.Context) error {
	form := url.Values{}
	form.Set("business_connection_id", c.businessConnectionID)
	form.Set("is_public", "false")

	return c.removeAvatar(ctx, form)
}

func (c *Client) RemovePublicAvatar(ctx context.Context) error {
	form := url.Values{}
	form.Set("business_connection_id", c.businessConnectionID)
	form.Set("is_public", "true")

	return c.removeAvatar(ctx, form)
}

func (c *Client) removeAvatar(ctx context.Context, form url.Values) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.methodURL("removeBusinessAccountProfilePhoto"), strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return c.do(req, nil)
}
