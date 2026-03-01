package minimax

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/giztoy/minimax-go/internal/transport"
)

const (
	defaultVoiceListPath = "/v1/get_voice"
	defaultVoiceType     = "all"
)

type VoiceService struct {
	transport *transport.Client
	endpoint  string
}

type ListVoicesRequest struct {
	VoiceType string `json:"voice_type,omitempty"`
	PageSize  *int   `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
}

type Voice struct {
	VoiceID     string                     `json:"voice_id,omitempty"`
	VoiceName   string                     `json:"voice_name,omitempty"`
	Description []string                   `json:"description,omitempty"`
	CreatedTime string                     `json:"created_time,omitempty"`
	VoiceType   string                     `json:"voice_type,omitempty"`
	Raw         map[string]json.RawMessage `json:"-"`
}

type ListVoicesResponse struct {
	Voices        []Voice                    `json:"voices"`
	NextPageToken string                     `json:"next_page_token,omitempty"`
	HasMore       bool                       `json:"has_more,omitempty"`
	Raw           map[string]json.RawMessage `json:"-"`
}

type listVoicesWireRequest struct {
	VoiceType string `json:"voice_type"`
	PageSize  *int   `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
}

type listVoicesRawResponse struct {
	Voices          []Voice                    `json:"voices,omitempty"`
	SystemVoice     []Voice                    `json:"system_voice,omitempty"`
	VoiceCloning    []Voice                    `json:"voice_cloning,omitempty"`
	VoiceGeneration []Voice                    `json:"voice_generation,omitempty"`
	NextPageToken   string                     `json:"next_page_token,omitempty"`
	HasMore         bool                       `json:"has_more,omitempty"`
	Raw             map[string]json.RawMessage `json:"-"`
}

// ListVoices queries available voices with filter and pagination parameters.
func (s *VoiceService) ListVoices(ctx context.Context, request *ListVoicesRequest) (*ListVoicesResponse, error) {
	if s == nil || s.transport == nil {
		return nil, errors.New("voice service is not initialized")
	}

	payload, query, err := buildListVoicesPayload(request)
	if err != nil {
		return nil, err
	}

	var raw listVoicesRawResponse
	if err := s.transport.DoJSON(ctx, transport.JSONRequest{
		Method: http.MethodPost,
		Path:   s.endpoint,
		Query:  query,
		Body:   payload,
	}, &raw); err != nil {
		return nil, err
	}

	return &ListVoicesResponse{
		Voices:        collectVoices(raw),
		NextPageToken: strings.TrimSpace(raw.NextPageToken),
		HasMore:       raw.HasMore,
		Raw:           cloneRawMessages(raw.Raw),
	}, nil
}

func buildListVoicesPayload(request *ListVoicesRequest) (listVoicesWireRequest, url.Values, error) {
	payload := listVoicesWireRequest{VoiceType: defaultVoiceType}
	if request != nil {
		payload.VoiceType = strings.TrimSpace(request.VoiceType)
		payload.PageSize = request.PageSize
		payload.PageToken = strings.TrimSpace(request.PageToken)
	}

	if payload.VoiceType == "" {
		payload.VoiceType = defaultVoiceType
	}

	if payload.PageSize != nil && *payload.PageSize < 0 {
		return listVoicesWireRequest{}, nil, errors.New("list voices request page_size must be non-negative")
	}

	query := make(url.Values)
	query.Set("voice_type", payload.VoiceType)
	if payload.PageSize != nil {
		query.Set("page_size", strconv.Itoa(*payload.PageSize))
	}
	if payload.PageToken != "" {
		query.Set("page_token", payload.PageToken)
	}

	return payload, query, nil
}

func collectVoices(raw listVoicesRawResponse) []Voice {
	if raw.Voices != nil {
		return cloneVoices(raw.Voices, "")
	}

	voices := make([]Voice, 0, len(raw.SystemVoice)+len(raw.VoiceCloning)+len(raw.VoiceGeneration))
	voices = appendVoices(voices, raw.SystemVoice, "system")
	voices = appendVoices(voices, raw.VoiceCloning, "voice_cloning")
	voices = appendVoices(voices, raw.VoiceGeneration, "voice_generation")

	if voices == nil {
		return make([]Voice, 0)
	}

	return voices
}

func appendVoices(dst, src []Voice, voiceType string) []Voice {
	for _, item := range src {
		copied := item
		if copied.VoiceType == "" {
			copied.VoiceType = voiceType
		}
		copied.Raw = cloneRawMessages(copied.Raw)
		dst = append(dst, copied)
	}
	return dst
}

func cloneVoices(voices []Voice, defaultType string) []Voice {
	if voices == nil {
		return make([]Voice, 0)
	}

	cloned := make([]Voice, 0, len(voices))
	for _, item := range voices {
		copied := item
		if copied.VoiceType == "" && defaultType != "" {
			copied.VoiceType = defaultType
		}
		copied.Raw = cloneRawMessages(copied.Raw)
		cloned = append(cloned, copied)
	}

	return cloned
}

func cloneRawMessages(raw map[string]json.RawMessage) map[string]json.RawMessage {
	if len(raw) == 0 {
		return nil
	}

	cloned := make(map[string]json.RawMessage, len(raw))
	for key, value := range raw {
		clonedValue := make(json.RawMessage, len(value))
		copy(clonedValue, value)
		cloned[key] = clonedValue
	}

	return cloned
}

func (v *Voice) UnmarshalJSON(data []byte) error {
	type alias Voice

	var parsed alias
	if err := json.Unmarshal(data, &parsed); err != nil {
		return err
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	delete(raw, "voice_id")
	delete(raw, "voice_name")
	delete(raw, "description")
	delete(raw, "created_time")
	delete(raw, "voice_type")

	*v = Voice(parsed)
	if len(raw) > 0 {
		v.Raw = raw
	} else {
		v.Raw = nil
	}

	return nil
}

func (r *listVoicesRawResponse) UnmarshalJSON(data []byte) error {
	type alias listVoicesRawResponse

	var parsed alias
	if err := json.Unmarshal(data, &parsed); err != nil {
		return err
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	delete(raw, "voices")
	delete(raw, "system_voice")
	delete(raw, "voice_cloning")
	delete(raw, "voice_generation")
	delete(raw, "next_page_token")
	delete(raw, "has_more")
	delete(raw, "base_resp")
	delete(raw, "status_code")
	delete(raw, "status_msg")

	*r = listVoicesRawResponse(parsed)
	if len(raw) > 0 {
		r.Raw = raw
	} else {
		r.Raw = nil
	}

	return nil
}
