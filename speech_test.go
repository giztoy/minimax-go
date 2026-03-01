package minimax

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/giztoy/minimax-go/internal/codec"
	"github.com/giztoy/minimax-go/internal/protocol"
	"github.com/giztoy/minimax-go/internal/transport"
)

func TestSpeechSynthesize(t *testing.T) {
	t.Parallel()

	t.Run("success returns audio bytes", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != defaultSpeechSynthesizePath {
				t.Fatalf("path = %s, want %s", r.URL.Path, defaultSpeechSynthesizePath)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"base_resp":{"status_code":0,"status_msg":"ok"},"data":{"audio_hex":"48656c6c6f"}}`))
		}))
		defer srv.Close()

		client, err := NewClient(Config{
			BaseURL:    srv.URL,
			HTTPClient: srv.Client(),
			Retry: transport.RetryConfig{
				MaxAttempts: 1,
			},
		})
		if err != nil {
			t.Fatalf("NewClient() error = %v, want nil", err)
		}

		resp, err := client.Speech.Synthesize(context.Background(), SpeechRequest{Text: "hello"})
		if err != nil {
			t.Fatalf("Synthesize() error = %v, want nil", err)
		}

		if len(resp.Audio) == 0 {
			t.Fatal("len(resp.Audio) = 0, want > 0")
		}

		if string(resp.Audio) != "Hello" {
			t.Fatalf("resp.Audio = %q, want %q", string(resp.Audio), "Hello")
		}
	})

	t.Run("volume zero is serialized in request", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("ReadAll(r.Body) error = %v", err)
			}

			var payload map[string]any
			if err := json.Unmarshal(body, &payload); err != nil {
				t.Fatalf("Unmarshal(request body) error = %v", err)
			}

			voiceSetting, ok := payload["voice_setting"].(map[string]any)
			if !ok {
				t.Fatalf("voice_setting missing or invalid type: %T", payload["voice_setting"])
			}

			volValue, exists := voiceSetting["vol"]
			if !exists {
				t.Fatalf("voice_setting.vol missing in payload: %v", voiceSetting)
			}

			volNumber, ok := volValue.(float64)
			if !ok {
				t.Fatalf("voice_setting.vol type = %T, want float64", volValue)
			}

			if volNumber != 0 {
				t.Fatalf("voice_setting.vol = %v, want 0", volNumber)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"base_resp":{"status_code":0,"status_msg":"ok"},"data":{"audio_hex":"48656c6c6f"}}`))
		}))
		defer srv.Close()

		client, err := NewClient(Config{
			BaseURL:    srv.URL,
			HTTPClient: srv.Client(),
			Retry: transport.RetryConfig{
				MaxAttempts: 1,
			},
		})
		if err != nil {
			t.Fatalf("NewClient() error = %v, want nil", err)
		}

		zero := 0.0
		resp, err := client.Speech.Synthesize(context.Background(), SpeechRequest{
			Text:    "hello",
			VoiceID: "English_Graceful_Lady",
			Vol:     &zero,
		})
		if err != nil {
			t.Fatalf("Synthesize() error = %v, want nil", err)
		}

		if string(resp.Audio) != "Hello" {
			t.Fatalf("resp.Audio = %q, want %q", string(resp.Audio), "Hello")
		}
	})

	t.Run("api error returns unified error model", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"base_resp":{"status_code":2301,"status_msg":"invalid voice"}}`))
		}))
		defer srv.Close()

		client, err := NewClient(Config{
			BaseURL:    srv.URL,
			HTTPClient: srv.Client(),
			Retry: transport.RetryConfig{
				MaxAttempts: 1,
			},
		})
		if err != nil {
			t.Fatalf("NewClient() error = %v, want nil", err)
		}

		_, err = client.Speech.Synthesize(context.Background(), SpeechRequest{Text: "hello"})
		if err == nil {
			t.Fatal("Synthesize() error = nil, want non-nil")
		}

		var apiErr *protocol.APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("Synthesize() error type = %T, want *protocol.APIError", err)
		}

		if apiErr.StatusCode != 2301 || apiErr.StatusMsg != "invalid voice" {
			t.Fatalf("apiErr = %+v, want status_code=2301 status_msg=invalid voice", apiErr)
		}
	})

	t.Run("timeout canceled by context", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(120 * time.Millisecond)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"base_resp":{"status_code":0,"status_msg":"ok"},"data":{"audio_hex":"48656c6c6f"}}`))
		}))
		defer srv.Close()

		client, err := NewClient(Config{
			BaseURL:    srv.URL,
			HTTPClient: srv.Client(),
			Retry: transport.RetryConfig{
				MaxAttempts: 1,
			},
		})
		if err != nil {
			t.Fatalf("NewClient() error = %v, want nil", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
		defer cancel()

		_, err = client.Speech.Synthesize(ctx, SpeechRequest{Text: "hello"})
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("Synthesize() error = %v, want context deadline exceeded", err)
		}
	})

	t.Run("retry eventually success", func(t *testing.T) {
		t.Parallel()

		var attempts int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			current := atomic.AddInt32(&attempts, 1)
			if current < 3 {
				w.WriteHeader(http.StatusServiceUnavailable)
				_, _ = w.Write([]byte(`{"error":"temporary unavailable"}`))
				return
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"base_resp":{"status_code":0,"status_msg":"ok"},"data":{"audio_hex":"48656c6c6f"}}`))
		}))
		defer srv.Close()

		client, err := NewClient(Config{
			BaseURL:    srv.URL,
			HTTPClient: srv.Client(),
			Retry: transport.RetryConfig{
				MaxAttempts: 3,
				BaseDelay:   time.Millisecond,
				MaxDelay:    time.Millisecond,
				Sleep: func(context.Context, time.Duration) error {
					return nil
				},
			},
		})
		if err != nil {
			t.Fatalf("NewClient() error = %v, want nil", err)
		}

		resp, err := client.Speech.Synthesize(context.Background(), SpeechRequest{Text: "hello"})
		if err != nil {
			t.Fatalf("Synthesize() error = %v, want nil", err)
		}

		if string(resp.Audio) != "Hello" {
			t.Fatalf("resp.Audio = %q, want Hello", string(resp.Audio))
		}

		if got := atomic.LoadInt32(&attempts); got != 3 {
			t.Fatalf("attempts = %d, want 3", got)
		}
	})

	t.Run("retry exhausted returns last retryable error", func(t *testing.T) {
		t.Parallel()

		var attempts int32
		var sleepCalls int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&attempts, 1)
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"error":"temporary unavailable"}`))
		}))
		defer srv.Close()

		client, err := NewClient(Config{
			BaseURL:    srv.URL,
			HTTPClient: srv.Client(),
			Retry: transport.RetryConfig{
				MaxAttempts: 3,
				BaseDelay:   time.Millisecond,
				MaxDelay:    time.Millisecond,
				Sleep: func(context.Context, time.Duration) error {
					atomic.AddInt32(&sleepCalls, 1)
					return nil
				},
			},
		})
		if err != nil {
			t.Fatalf("NewClient() error = %v, want nil", err)
		}

		_, err = client.Speech.Synthesize(context.Background(), SpeechRequest{Text: "hello"})
		if err == nil {
			t.Fatal("Synthesize() error = nil, want non-nil")
		}

		var apiErr *protocol.APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("Synthesize() error type = %T, want *protocol.APIError", err)
		}

		if apiErr.HTTPStatus != http.StatusServiceUnavailable {
			t.Fatalf("apiErr.HTTPStatus = %d, want %d", apiErr.HTTPStatus, http.StatusServiceUnavailable)
		}

		if got := atomic.LoadInt32(&attempts); got != 3 {
			t.Fatalf("attempts = %d, want 3", got)
		}

		if got := atomic.LoadInt32(&sleepCalls); got != 2 {
			t.Fatalf("sleepCalls = %d, want 2", got)
		}
	})

	t.Run("invalid hex returns decode error", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"base_resp":{"status_code":0,"status_msg":"ok"},"data":{"audio_hex":"xyz"}}`))
		}))
		defer srv.Close()

		client, err := NewClient(Config{
			BaseURL:    srv.URL,
			HTTPClient: srv.Client(),
			Retry: transport.RetryConfig{
				MaxAttempts: 1,
			},
		})
		if err != nil {
			t.Fatalf("NewClient() error = %v, want nil", err)
		}

		_, err = client.Speech.Synthesize(context.Background(), SpeechRequest{Text: "hello"})
		if !errors.Is(err, codec.ErrInvalidHexAudio) {
			t.Fatalf("Synthesize() error = %v, want codec.ErrInvalidHexAudio", err)
		}
	})

	t.Run("explicit context cancel", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(120 * time.Millisecond)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"base_resp":{"status_code":0,"status_msg":"ok"},"data":{"audio_hex":"48656c6c6f"}}`))
		}))
		defer srv.Close()

		client, err := NewClient(Config{
			BaseURL:    srv.URL,
			HTTPClient: srv.Client(),
			Retry: transport.RetryConfig{
				MaxAttempts: 1,
			},
		})
		if err != nil {
			t.Fatalf("NewClient() error = %v, want nil", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err = client.Speech.Synthesize(ctx, SpeechRequest{Text: "hello"})
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Synthesize() error = %v, want context canceled", err)
		}
	})
}
