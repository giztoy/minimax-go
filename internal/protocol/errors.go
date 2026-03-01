package protocol

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

const maxBodyPreviewSize = 512

// BaseResp represents business status metadata in Minimax responses.
type BaseResp struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

// APIError is the unified error model for HTTP and base_resp semantics.
type APIError struct {
	HTTPStatus int
	StatusCode int
	StatusMsg  string
	Body       string
	Cause      error
}

func (e *APIError) Error() string {
	if e == nil {
		return "<nil>"
	}

	if e.StatusCode != 0 {
		return fmt.Sprintf("minimax api error: status_code=%d status_msg=%q", e.StatusCode, e.StatusMsg)
	}

	if e.HTTPStatus != 0 {
		if e.Body != "" {
			return fmt.Sprintf("minimax http error: status=%d body=%q", e.HTTPStatus, e.Body)
		}
		return fmt.Sprintf("minimax http error: status=%d", e.HTTPStatus)
	}

	if e.Cause != nil {
		return fmt.Sprintf("minimax transport error: %v", e.Cause)
	}

	return "minimax api error"
}

func (e *APIError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func NewHTTPError(httpStatus int, body []byte) *APIError {
	return &APIError{
		HTTPStatus: httpStatus,
		StatusMsg:  http.StatusText(httpStatus),
		Body:       compactBody(body),
	}
}

func NewBaseRespError(httpStatus int, baseResp BaseResp, body []byte) *APIError {
	return &APIError{
		HTTPStatus: httpStatus,
		StatusCode: baseResp.StatusCode,
		StatusMsg:  baseResp.StatusMsg,
		Body:       compactBody(body),
	}
}

type responseEnvelope struct {
	BaseResp   *BaseResp `json:"base_resp"`
	StatusCode *int      `json:"status_code"`
	StatusMsg  *string   `json:"status_msg"`
}

// ParseBaseResp extracts base_resp fields from a response body.
func ParseBaseResp(body []byte) (BaseResp, bool) {
	if len(body) == 0 {
		return BaseResp{}, false
	}

	var envelope responseEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		return BaseResp{}, false
	}

	if envelope.BaseResp != nil {
		return *envelope.BaseResp, true
	}

	if envelope.StatusCode != nil || envelope.StatusMsg != nil {
		resp := BaseResp{}
		if envelope.StatusCode != nil {
			resp.StatusCode = *envelope.StatusCode
		}
		if envelope.StatusMsg != nil {
			resp.StatusMsg = *envelope.StatusMsg
		}
		return resp, true
	}

	return BaseResp{}, false
}

// CheckResponse normalizes HTTP and business status into a unified error.
func CheckResponse(httpStatus int, body []byte) error {
	if httpStatus < http.StatusOK || httpStatus >= http.StatusMultipleChoices {
		return NewHTTPError(httpStatus, body)
	}

	if baseResp, ok := ParseBaseResp(body); ok && baseResp.StatusCode != 0 {
		return NewBaseRespError(httpStatus, baseResp, body)
	}

	return nil
}

// IsRetryable reports whether an error is retryable.
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.HTTPStatus == http.StatusTooManyRequests || apiErr.HTTPStatus >= http.StatusInternalServerError
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return netErr.Timeout()
	}

	return false
}

func compactBody(body []byte) string {
	s := strings.TrimSpace(string(body))
	if len(s) <= maxBodyPreviewSize {
		return s
	}
	return s[:maxBodyPreviewSize] + "..."
}
