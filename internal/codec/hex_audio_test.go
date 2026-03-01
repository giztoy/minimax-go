package codec

import (
	"errors"
	"testing"
)

func TestDecodeHexAudio(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		audio, err := DecodeHexAudio("48656c6c6f")
		if err != nil {
			t.Fatalf("DecodeHexAudio() error = %v, want nil", err)
		}

		if string(audio) != "Hello" {
			t.Fatalf("DecodeHexAudio() = %q, want %q", string(audio), "Hello")
		}
	})

	t.Run("invalid hex", func(t *testing.T) {
		t.Parallel()

		_, err := DecodeHexAudio("xyz")
		if !errors.Is(err, ErrInvalidHexAudio) {
			t.Fatalf("DecodeHexAudio() error = %v, want ErrInvalidHexAudio", err)
		}
	})

	t.Run("empty hex", func(t *testing.T) {
		t.Parallel()

		_, err := DecodeHexAudio("   ")
		if !errors.Is(err, ErrEmptyHexAudio) {
			t.Fatalf("DecodeHexAudio() error = %v, want ErrEmptyHexAudio", err)
		}
	})
}
