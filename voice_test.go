package minimax

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/giztoy/minimax-go/internal/protocol"
	"github.com/giztoy/minimax-go/internal/transport"
)

func TestListVoices(t *testing.T) {
	t.Parallel()

	t.Run("success with default request", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Fatalf("method = %s, want POST", r.Method)
			}

			if r.URL.Path != defaultVoiceListPath {
				t.Fatalf("path = %s, want %s", r.URL.Path, defaultVoiceListPath)
			}

			if got := r.URL.Query().Get("voice_type"); got != defaultVoiceType {
				t.Fatalf("query.voice_type = %q, want %q", got, defaultVoiceType)
			}

			var payload listVoicesWireRequest
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("Decode(request body) error = %v", err)
			}

			if payload.VoiceType != defaultVoiceType {
				t.Fatalf("payload.voice_type = %q, want %q", payload.VoiceType, defaultVoiceType)
			}

			if payload.PageSize != nil {
				t.Fatalf("payload.page_size = %d, want nil", *payload.PageSize)
			}

			if payload.PageToken != "" {
				t.Fatalf("payload.page_token = %q, want empty", payload.PageToken)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"base_resp":{"status_code":0,"status_msg":"ok"},"voices":[{"voice_id":"voice-system-1","voice_name":"calm narrator","description":["calm"],"created_time":"2026-03-01","voice_type":"system","gender":"female"}],"next_page_token":"cursor-2","has_more":true,"request_id":"req-1"}`))
		}))
		defer srv.Close()

		client := newVoiceTestClient(t, srv, transport.RetryConfig{MaxAttempts: 1})
		resp, err := client.Voice.ListVoices(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListVoices() error = %v, want nil", err)
		}

		if got := len(resp.Voices); got != 1 {
			t.Fatalf("len(resp.Voices) = %d, want 1", got)
		}

		if resp.NextPageToken != "cursor-2" {
			t.Fatalf("resp.NextPageToken = %q, want %q", resp.NextPageToken, "cursor-2")
		}

		if !resp.HasMore {
			t.Fatal("resp.HasMore = false, want true")
		}

		voice := resp.Voices[0]
		if voice.VoiceID != "voice-system-1" || voice.VoiceType != "system" {
			t.Fatalf("voice = %+v, want voice_id=voice-system-1 voice_type=system", voice)
		}

		if _, ok := voice.Raw["gender"]; !ok {
			t.Fatalf("voice.Raw = %v, want gender field", voice.Raw)
		}

		if _, ok := resp.Raw["request_id"]; !ok {
			t.Fatalf("resp.Raw = %v, want request_id field", resp.Raw)
		}
	})

	t.Run("empty list returns empty slice", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"base_resp":{"status_code":0,"status_msg":"ok"},"voices":[]}`))
		}))
		defer srv.Close()

		client := newVoiceTestClient(t, srv, transport.RetryConfig{MaxAttempts: 1})
		resp, err := client.Voice.ListVoices(context.Background(), &ListVoicesRequest{VoiceType: "system"})
		if err != nil {
			t.Fatalf("ListVoices() error = %v, want nil", err)
		}

		if resp.Voices == nil {
			t.Fatal("resp.Voices = nil, want empty slice")
		}

		if got := len(resp.Voices); got != 0 {
			t.Fatalf("len(resp.Voices) = %d, want 0", got)
		}
	})

	t.Run("pagination and filter are forwarded", func(t *testing.T) {
		t.Parallel()

		pageSize := 25
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query()
			if query.Get("voice_type") != "voice_cloning" {
				t.Fatalf("query.voice_type = %q, want %q", query.Get("voice_type"), "voice_cloning")
			}

			if query.Get("page_size") != "25" {
				t.Fatalf("query.page_size = %q, want %q", query.Get("page_size"), "25")
			}

			if query.Get("page_token") != "cursor-2" {
				t.Fatalf("query.page_token = %q, want %q", query.Get("page_token"), "cursor-2")
			}

			var payload listVoicesWireRequest
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("Decode(request body) error = %v", err)
			}

			if payload.VoiceType != "voice_cloning" {
				t.Fatalf("payload.voice_type = %q, want %q", payload.VoiceType, "voice_cloning")
			}

			if payload.PageSize == nil || *payload.PageSize != pageSize {
				t.Fatalf("payload.page_size = %v, want %d", payload.PageSize, pageSize)
			}

			if payload.PageToken != "cursor-2" {
				t.Fatalf("payload.page_token = %q, want %q", payload.PageToken, "cursor-2")
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"base_resp":{"status_code":0,"status_msg":"ok"},"voices":[{"voice_id":"clone-1"}],"next_page_token":"cursor-3","has_more":true}`))
		}))
		defer srv.Close()

		client := newVoiceTestClient(t, srv, transport.RetryConfig{MaxAttempts: 1})
		resp, err := client.Voice.ListVoices(context.Background(), &ListVoicesRequest{
			VoiceType: "voice_cloning",
			PageSize:  &pageSize,
			PageToken: "cursor-2",
		})
		if err != nil {
			t.Fatalf("ListVoices() error = %v, want nil", err)
		}

		if resp.NextPageToken != "cursor-3" || !resp.HasMore {
			t.Fatalf("resp = %+v, want next_page_token=cursor-3 has_more=true", resp)
		}
	})

	t.Run("legacy response shape is normalized", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"base_resp":{"status_code":0,"status_msg":"ok"},"system_voice":[{"voice_id":"sys-1","voice_name":"sys"}],"voice_cloning":[{"voice_id":"clone-1"}],"voice_generation":[{"voice_id":"gen-1"}],"request_id":"legacy-1"}`))
		}))
		defer srv.Close()

		client := newVoiceTestClient(t, srv, transport.RetryConfig{MaxAttempts: 1})
		resp, err := client.Voice.ListVoices(context.Background(), &ListVoicesRequest{VoiceType: "all"})
		if err != nil {
			t.Fatalf("ListVoices() error = %v, want nil", err)
		}

		if got := len(resp.Voices); got != 3 {
			t.Fatalf("len(resp.Voices) = %d, want 3", got)
		}

		if resp.Voices[0].VoiceType != "system" || resp.Voices[1].VoiceType != "voice_cloning" || resp.Voices[2].VoiceType != "voice_generation" {
			t.Fatalf("voice types = [%s %s %s], want [system voice_cloning voice_generation]", resp.Voices[0].VoiceType, resp.Voices[1].VoiceType, resp.Voices[2].VoiceType)
		}

		if _, ok := resp.Raw["request_id"]; !ok {
			t.Fatalf("resp.Raw = %v, want request_id field", resp.Raw)
		}
	})

	t.Run("http error returns unified APIError", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":"unauthorized"}`))
		}))
		defer srv.Close()

		client := newVoiceTestClient(t, srv, transport.RetryConfig{MaxAttempts: 1})
		_, err := client.Voice.ListVoices(context.Background(), &ListVoicesRequest{VoiceType: "all"})
		if err == nil {
			t.Fatal("ListVoices() error = nil, want non-nil")
		}

		var apiErr *protocol.APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("ListVoices() error type = %T, want *protocol.APIError", err)
		}

		if apiErr.HTTPStatus != http.StatusUnauthorized {
			t.Fatalf("apiErr.HTTPStatus = %d, want %d", apiErr.HTTPStatus, http.StatusUnauthorized)
		}
	})

	t.Run("http 500 returns unified APIError", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"internal"}`))
		}))
		defer srv.Close()

		client := newVoiceTestClient(t, srv, transport.RetryConfig{MaxAttempts: 1})
		_, err := client.Voice.ListVoices(context.Background(), &ListVoicesRequest{VoiceType: "all"})
		if err == nil {
			t.Fatal("ListVoices() error = nil, want non-nil")
		}

		var apiErr *protocol.APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("ListVoices() error type = %T, want *protocol.APIError", err)
		}

		if apiErr.HTTPStatus != http.StatusInternalServerError {
			t.Fatalf("apiErr.HTTPStatus = %d, want %d", apiErr.HTTPStatus, http.StatusInternalServerError)
		}
	})

	t.Run("base_resp error returns unified APIError", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"base_resp":{"status_code":2013,"status_msg":"invalid voice_type"}}`))
		}))
		defer srv.Close()

		client := newVoiceTestClient(t, srv, transport.RetryConfig{MaxAttempts: 1})
		_, err := client.Voice.ListVoices(context.Background(), &ListVoicesRequest{VoiceType: "invalid"})
		if err == nil {
			t.Fatal("ListVoices() error = nil, want non-nil")
		}

		var apiErr *protocol.APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("ListVoices() error type = %T, want *protocol.APIError", err)
		}

		if apiErr.StatusCode != 2013 || apiErr.StatusMsg != "invalid voice_type" {
			t.Fatalf("apiErr = %+v, want status_code=2013 status_msg=invalid voice_type", apiErr)
		}
	})

	t.Run("timeout canceled by context", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(120 * time.Millisecond)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"base_resp":{"status_code":0,"status_msg":"ok"},"voices":[]}`))
		}))
		defer srv.Close()

		client := newVoiceTestClient(t, srv, transport.RetryConfig{MaxAttempts: 1})
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
		defer cancel()

		_, err := client.Voice.ListVoices(ctx, &ListVoicesRequest{VoiceType: "all"})
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("ListVoices() error = %v, want context deadline exceeded", err)
		}
	})
}

func TestVoiceListVoicesValidation(t *testing.T) {
	t.Parallel()

	t.Run("negative page size is rejected", func(t *testing.T) {
		t.Parallel()

		client, err := NewClient(Config{BaseURL: "https://api.minimax.io"})
		if err != nil {
			t.Fatalf("NewClient() error = %v, want nil", err)
		}

		negative := -1
		_, err = client.Voice.ListVoices(context.Background(), &ListVoicesRequest{
			VoiceType: "all",
			PageSize:  &negative,
		})
		if err == nil {
			t.Fatal("ListVoices() error = nil, want non-nil")
		}

		if !strings.Contains(err.Error(), "page_size") {
			t.Fatalf("ListVoices() error = %v, want page_size validation", err)
		}
	})

	t.Run("nil service returns initialization error", func(t *testing.T) {
		t.Parallel()

		var service *VoiceService
		_, err := service.ListVoices(context.Background(), nil)
		if err == nil {
			t.Fatal("ListVoices() error = nil, want non-nil")
		}

		if !strings.Contains(err.Error(), "not initialized") {
			t.Fatalf("ListVoices() error = %v, want initialization error", err)
		}
	})
}

func newVoiceTestClient(t *testing.T, srv *httptest.Server, retry transport.RetryConfig) *Client {
	t.Helper()

	client, err := NewClient(Config{
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
		Retry:      retry,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v, want nil", err)
	}

	return client
}
