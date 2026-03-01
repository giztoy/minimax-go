package stream

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func TestReaderNext(t *testing.T) {
	t.Parallel()

	t.Run("single event", func(t *testing.T) {
		t.Parallel()

		r := NewReader(strings.NewReader("id: 1\nevent: message\ndata: hello\n\n"))

		event, err := r.Next()
		if err != nil {
			t.Fatalf("Next() error = %v, want nil", err)
		}

		if event.ID != "1" || event.Event != "message" || event.Data != "hello" {
			t.Fatalf("event = %+v, want id=1 event=message data=hello", event)
		}
	})

	t.Run("multiline data", func(t *testing.T) {
		t.Parallel()

		r := NewReader(strings.NewReader("data: line1\ndata: line2\n\n"))

		event, err := r.Next()
		if err != nil {
			t.Fatalf("Next() error = %v, want nil", err)
		}

		if event.Data != "line1\nline2" {
			t.Fatalf("event.Data = %q, want %q", event.Data, "line1\nline2")
		}
	})

	t.Run("eof", func(t *testing.T) {
		t.Parallel()

		r := NewReader(strings.NewReader(""))

		_, err := r.Next()
		if !errors.Is(err, io.EOF) {
			t.Fatalf("Next() error = %v, want io.EOF", err)
		}
	})
}
