package stream

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

// Event represents a single SSE event.
type Event struct {
	ID    string
	Event string
	Data  string
	Retry int
}

// Reader reads SSE events one by one.
type Reader struct {
	reader *bufio.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{reader: bufio.NewReader(r)}
}

// Next returns the next complete event, or io.EOF when exhausted.
func (r *Reader) Next() (Event, error) {
	var (
		event     Event
		dataLines []string
		hasField  bool
	)

	for {
		line, err := r.reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return Event{}, err
		}

		if errors.Is(err, io.EOF) && line == "" {
			if !hasField {
				return Event{}, io.EOF
			}
			event.Data = strings.Join(dataLines, "\n")
			return event, nil
		}

		line = strings.TrimRight(line, "\r\n")

		if line == "" {
			if !hasField {
				if errors.Is(err, io.EOF) {
					return Event{}, io.EOF
				}
				continue
			}

			event.Data = strings.Join(dataLines, "\n")
			return event, nil
		}

		if strings.HasPrefix(line, ":") {
			if errors.Is(err, io.EOF) && !hasField {
				return Event{}, io.EOF
			}
			continue
		}

		field, value := splitField(line)
		hasField = true

		switch field {
		case "id":
			event.ID = value
		case "event":
			event.Event = value
		case "data":
			dataLines = append(dataLines, value)
		case "retry":
			if retry, convErr := strconv.Atoi(value); convErr == nil && retry >= 0 {
				event.Retry = retry
			}
		}

		if errors.Is(err, io.EOF) {
			event.Data = strings.Join(dataLines, "\n")
			return event, nil
		}
	}
}

func splitField(line string) (string, string) {
	idx := strings.IndexByte(line, ':')
	if idx == -1 {
		return line, ""
	}

	field := line[:idx]
	value := line[idx+1:]
	if strings.HasPrefix(value, " ") {
		value = value[1:]
	}

	return field, value
}
