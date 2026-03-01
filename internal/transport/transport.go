package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
	"time"

	"github.com/giztoy/minimax-go/internal/protocol"
)

const (
	defaultMaxAttempts = 3
	defaultBaseDelay   = 100 * time.Millisecond
	defaultMaxDelay    = 2 * time.Second
)

type RetryConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
	ShouldRetry func(error) bool
	Sleep       func(context.Context, time.Duration) error
}

type Config struct {
	BaseURL        string
	APIKey         string
	HTTPClient     *http.Client
	DefaultHeaders http.Header
	Retry          RetryConfig
}

type Client struct {
	baseURL        string
	apiKey         string
	httpClient     *http.Client
	defaultHeaders http.Header
	retry          RetryConfig
}

type JSONRequest struct {
	Method  string
	Path    string
	Query   url.Values
	Headers http.Header
	Body    any
}

type StreamRequest struct {
	Method  string
	Path    string
	Query   url.Values
	Headers http.Header
	Body    any
}

type UploadRequest struct {
	Method          string
	Path            string
	Query           url.Values
	Headers         http.Header
	Fields          map[string]string
	FileField       string
	FileName        string
	FileContentType string
	FileData        []byte
}

func New(config Config) (*Client, error) {
	retry := withRetryDefaults(config.Retry)

	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}

	return &Client{
		baseURL:        strings.TrimSpace(config.BaseURL),
		apiKey:         strings.TrimSpace(config.APIKey),
		httpClient:     httpClient,
		defaultHeaders: config.DefaultHeaders.Clone(),
		retry:          retry,
	}, nil
}

// DoJSON sends a JSON request and unmarshals the response into out.
func (c *Client) DoJSON(ctx context.Context, request JSONRequest, out any) error {
	method := request.Method
	if method == "" {
		method = http.MethodPost
	}

	var payload []byte
	var err error
	if request.Body != nil {
		payload, err = json.Marshal(request.Body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
	}

	err = c.withRetry(ctx, func() error {
		req, reqErr := c.buildRequest(ctx, method, request.Path, request.Query, bytes.NewReader(payload))
		if reqErr != nil {
			return reqErr
		}

		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
		mergeHeaders(req.Header, request.Headers)

		resp, doErr := c.httpClient.Do(req)
		if doErr != nil {
			return doErr
		}
		defer resp.Body.Close()

		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return fmt.Errorf("read response body: %w", readErr)
		}

		if checkErr := protocol.CheckResponse(resp.StatusCode, body); checkErr != nil {
			return checkErr
		}

		if out == nil || len(body) == 0 {
			return nil
		}

		if unmarshalErr := json.Unmarshal(body, out); unmarshalErr != nil {
			return fmt.Errorf("decode response body: %w", unmarshalErr)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// OpenStream opens a streaming connection; caller must close the returned body.
func (c *Client) OpenStream(ctx context.Context, request StreamRequest) (io.ReadCloser, error) {
	method := request.Method
	if method == "" {
		method = http.MethodGet
	}

	var payload []byte
	var err error
	if request.Body != nil {
		payload, err = json.Marshal(request.Body)
		if err != nil {
			return nil, fmt.Errorf("marshal stream body: %w", err)
		}
	}

	var lastErr error
	for attempt := 1; attempt <= c.retry.MaxAttempts; attempt++ {
		req, reqErr := c.buildRequest(ctx, method, request.Path, request.Query, bytes.NewReader(payload))
		if reqErr != nil {
			return nil, reqErr
		}

		req.Header.Set("Accept", "text/event-stream")
		if len(payload) > 0 {
			req.Header.Set("Content-Type", "application/json")
		}
		mergeHeaders(req.Header, request.Headers)

		resp, doErr := c.httpClient.Do(req)
		if doErr != nil {
			lastErr = doErr
			if !c.shouldRetry(doErr) || attempt == c.retry.MaxAttempts {
				return nil, doErr
			}

			if sleepErr := c.retry.Sleep(ctx, c.retryDelay(attempt)); sleepErr != nil {
				return nil, sleepErr
			}
			continue
		}

		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
			body, readErr := io.ReadAll(resp.Body)
			resp.Body.Close()
			if readErr != nil {
				lastErr = fmt.Errorf("read stream error response: %w", readErr)
			} else {
				lastErr = protocol.CheckResponse(resp.StatusCode, body)
			}

			if !c.shouldRetry(lastErr) || attempt == c.retry.MaxAttempts {
				return nil, lastErr
			}

			if sleepErr := c.retry.Sleep(ctx, c.retryDelay(attempt)); sleepErr != nil {
				return nil, sleepErr
			}
			continue
		}

		contentType := resp.Header.Get("Content-Type")
		if !isEventStreamContentType(contentType) {
			body, readErr := io.ReadAll(resp.Body)
			resp.Body.Close()
			if readErr != nil {
				lastErr = fmt.Errorf("read non-stream response body: %w", readErr)
			} else if checkErr := protocol.CheckResponse(resp.StatusCode, body); checkErr != nil {
				lastErr = checkErr
			} else {
				lastErr = fmt.Errorf("unexpected stream content type: %q", contentType)
			}

			if !c.shouldRetry(lastErr) || attempt == c.retry.MaxAttempts {
				return nil, lastErr
			}

			if sleepErr := c.retry.Sleep(ctx, c.retryDelay(attempt)); sleepErr != nil {
				return nil, sleepErr
			}
			continue
		}

		return resp.Body, nil
	}

	if lastErr == nil {
		lastErr = errors.New("open stream failed")
	}

	return nil, lastErr
}

// Upload sends a multipart/form-data request.
func (c *Client) Upload(ctx context.Context, request UploadRequest, out any) error {
	method := request.Method
	if method == "" {
		method = http.MethodPost
	}

	if request.FileField == "" {
		return errors.New("upload request requires FileField")
	}

	if request.FileName == "" {
		return errors.New("upload request requires FileName")
	}

	return c.withRetry(ctx, func() error {
		var payload bytes.Buffer
		writer := multipart.NewWriter(&payload)

		for key, value := range request.Fields {
			if err := writer.WriteField(key, value); err != nil {
				return fmt.Errorf("write field %s: %w", key, err)
			}
		}

		header := textproto.MIMEHeader{}
		contentDisposition := fmt.Sprintf(`form-data; name=%q; filename=%q`, request.FileField, request.FileName)
		header.Set("Content-Disposition", contentDisposition)
		if request.FileContentType != "" {
			header.Set("Content-Type", request.FileContentType)
		}

		part, err := writer.CreatePart(header)
		if err != nil {
			return fmt.Errorf("create file part: %w", err)
		}

		if _, err := part.Write(request.FileData); err != nil {
			return fmt.Errorf("write file data: %w", err)
		}

		if err := writer.Close(); err != nil {
			return fmt.Errorf("close multipart writer: %w", err)
		}

		req, reqErr := c.buildRequest(ctx, method, request.Path, request.Query, bytes.NewReader(payload.Bytes()))
		if reqErr != nil {
			return reqErr
		}

		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", writer.FormDataContentType())
		mergeHeaders(req.Header, request.Headers)

		resp, doErr := c.httpClient.Do(req)
		if doErr != nil {
			return doErr
		}
		defer resp.Body.Close()

		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return fmt.Errorf("read upload response: %w", readErr)
		}

		if checkErr := protocol.CheckResponse(resp.StatusCode, body); checkErr != nil {
			return checkErr
		}

		if out == nil || len(body) == 0 {
			return nil
		}

		if unmarshalErr := json.Unmarshal(body, out); unmarshalErr != nil {
			return fmt.Errorf("decode upload response: %w", unmarshalErr)
		}

		return nil
	})
}

func (c *Client) withRetry(ctx context.Context, op func() error) error {
	var lastErr error

	for attempt := 1; attempt <= c.retry.MaxAttempts; attempt++ {
		err := op()
		if err == nil {
			return nil
		}

		lastErr = err
		if !c.shouldRetry(err) || attempt == c.retry.MaxAttempts {
			return err
		}

		if sleepErr := c.retry.Sleep(ctx, c.retryDelay(attempt)); sleepErr != nil {
			return sleepErr
		}
	}

	if lastErr == nil {
		return errors.New("request failed")
	}

	return lastErr
}

func (c *Client) shouldRetry(err error) bool {
	if err == nil {
		return false
	}

	if c.retry.ShouldRetry != nil {
		return c.retry.ShouldRetry(err)
	}

	return protocol.IsRetryable(err)
}

func (c *Client) retryDelay(attempt int) time.Duration {
	delay := c.retry.BaseDelay
	for i := 1; i < attempt; i++ {
		delay *= 2
		if delay >= c.retry.MaxDelay {
			return c.retry.MaxDelay
		}
	}

	if delay > c.retry.MaxDelay {
		return c.retry.MaxDelay
	}

	return delay
}

func withRetryDefaults(retry RetryConfig) RetryConfig {
	if retry.MaxAttempts <= 0 {
		retry.MaxAttempts = defaultMaxAttempts
	}

	if retry.BaseDelay <= 0 {
		retry.BaseDelay = defaultBaseDelay
	}

	if retry.MaxDelay <= 0 {
		retry.MaxDelay = defaultMaxDelay
	}

	if retry.Sleep == nil {
		retry.Sleep = sleepWithContext
	}

	return retry
}

func (c *Client) buildRequest(ctx context.Context, method, path string, query url.Values, body io.Reader) (*http.Request, error) {
	resolvedURL, err := c.resolveURL(path)
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(resolvedURL)
	if err != nil {
		return nil, fmt.Errorf("parse request url: %w", err)
	}

	if query != nil {
		q := parsedURL.Query()
		for key, values := range query {
			for _, value := range values {
				q.Add(key, value)
			}
		}
		parsedURL.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, parsedURL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	mergeHeaders(req.Header, c.defaultHeaders)
	if c.apiKey != "" && req.Header.Get("Authorization") == "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	return req, nil
}

func (c *Client) resolveURL(path string) (string, error) {
	trimmedPath := strings.TrimSpace(path)
	if trimmedPath == "" {
		return "", errors.New("request path is empty")
	}

	if strings.HasPrefix(trimmedPath, "http://") || strings.HasPrefix(trimmedPath, "https://") {
		return trimmedPath, nil
	}

	if c.baseURL == "" {
		return "", errors.New("baseURL is empty for relative path request")
	}

	return strings.TrimRight(c.baseURL, "/") + "/" + strings.TrimLeft(trimmedPath, "/"), nil
}

func mergeHeaders(dst, src http.Header) {
	if dst == nil || src == nil {
		return
	}

	for key, values := range src {
		dst.Del(key)
		for idx, value := range values {
			if idx == 0 {
				dst.Set(key, value)
				continue
			}
			dst.Add(key, value)
		}
	}
}

func isEventStreamContentType(contentType string) bool {
	trimmed := strings.TrimSpace(contentType)
	if trimmed == "" {
		return false
	}

	mediaType, _, err := mime.ParseMediaType(trimmed)
	if err == nil {
		return strings.EqualFold(mediaType, "text/event-stream")
	}

	return strings.HasPrefix(strings.ToLower(trimmed), "text/event-stream")
}

func sleepWithContext(ctx context.Context, delay time.Duration) error {
	if delay <= 0 {
		return nil
	}

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
