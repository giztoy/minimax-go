package codec

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEmptyHexAudio   = errors.New("empty hex audio data")
	ErrInvalidHexAudio = errors.New("invalid hex audio data")
)

// DecodeHexAudio decodes a hex-encoded audio string into bytes.
func DecodeHexAudio(input string) ([]byte, error) {
	normalized := strings.TrimSpace(input)
	normalized = strings.TrimPrefix(normalized, "0x")
	normalized = strings.TrimPrefix(normalized, "0X")

	if normalized == "" {
		return nil, ErrEmptyHexAudio
	}

	audio, err := hex.DecodeString(normalized)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidHexAudio, err)
	}

	return audio, nil
}
