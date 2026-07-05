package telegram

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestSetBioUsesBusinessEndpoint(t *testing.T) {
	var path string
	httpClient := &http.Client{Transport: roundTrip(func(r *http.Request) *http.Response {
		path = r.URL.Path
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("business_connection_id") != "business" {
			t.Fatalf("business_connection_id = %q", r.Form.Get("business_connection_id"))
		}
		return jsonResponse(http.StatusOK, `{"ok":true,"result":true}`)
	})}

	client := NewClientWithBaseURL("token", "business", "https://telegram.test", httpClient)
	if err := client.SetBio(context.Background(), "bio"); err != nil {
		t.Fatal(err)
	}
	if path != "/bottoken/setBusinessAccountBio" {
		t.Fatalf("path = %q", path)
	}
}

func TestRetryAfterError(t *testing.T) {
	httpClient := &http.Client{Transport: roundTrip(func(r *http.Request) *http.Response {
		return jsonResponse(http.StatusTooManyRequests, `{"ok":false,"description":"Too Many Requests","parameters":{"retry_after":7}}`)
	})}

	client := NewClientWithBaseURL("token", "business", "https://telegram.test", httpClient)
	err := client.SetBio(context.Background(), "bio")

	var retry interface{ RetryAfter() time.Duration }
	if !errors.As(err, &retry) {
		t.Fatalf("err = %v, want retry error", err)
	}
	if retry.RetryAfter() != 7*time.Second {
		t.Fatalf("retry = %s", retry.RetryAfter())
	}
}

func TestAvatarMethodsUsePrivateMainProfilePhoto(t *testing.T) {
	paths := make([]string, 0, 2)
	httpClient := &http.Client{Transport: roundTrip(func(r *http.Request) *http.Response {
		paths = append(paths, r.URL.Path)
		if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			if err := r.ParseMultipartForm(1024 * 1024); err != nil {
				t.Fatal(err)
			}
			if r.MultipartForm.Value["is_public"][0] != "false" {
				t.Fatalf("set is_public = %q", r.MultipartForm.Value["is_public"][0])
			}
		} else {
			if err := r.ParseForm(); err != nil {
				t.Fatal(err)
			}
			if r.Form.Get("is_public") != "false" {
				t.Fatalf("remove is_public = %q", r.Form.Get("is_public"))
			}
		}
		return jsonResponse(http.StatusOK, `{"ok":true,"result":true}`)
	})}

	client := NewClientWithBaseURL("token", "business", "https://telegram.test", httpClient)
	if err := client.RemoveAvatar(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := client.SetAvatar(context.Background(), []byte("jpg")); err != nil {
		t.Fatal(err)
	}
	if paths[0] != "/bottoken/removeBusinessAccountProfilePhoto" || paths[1] != "/bottoken/setBusinessAccountProfilePhoto" {
		t.Fatalf("paths = %v", paths)
	}
}

func TestRemovePublicAvatarUsesPublicFlag(t *testing.T) {
	httpClient := &http.Client{Transport: roundTrip(func(r *http.Request) *http.Response {
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("is_public") != "true" {
			t.Fatalf("is_public = %q", r.Form.Get("is_public"))
		}
		return jsonResponse(http.StatusOK, `{"ok":true,"result":true}`)
	})}

	client := NewClientWithBaseURL("token", "business", "https://telegram.test", httpClient)
	if err := client.RemovePublicAvatar(context.Background()); err != nil {
		t.Fatal(err)
	}
}

type roundTrip func(*http.Request) *http.Response

func (r roundTrip) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req), nil
}

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Status:     http.StatusText(status),
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}
