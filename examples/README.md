# Examples

This directory contains runnable examples for `github.com/giztoy/minimax-go`.

## `speech` example

`examples/speech` demonstrates synchronous text-to-speech synthesis via `Speech.Synthesize`, then writes the returned audio bytes to a local file.

### Prerequisites

- Go 1.26+
- A valid Minimax API key

### Quick start

```bash
export MINIMAX_API_KEY="your_api_key"

go run ./examples/speech \
  -text "hello from minimax-go" \
  -voice-id "male-qn-qingse" \
  -output /tmp/speech_output.audio
```

If successful, the command prints the output file path and byte size.

### Show all CLI options

```bash
go run ./examples/speech -h
```

### Common flags

- `-api-key`: Minimax API key (takes precedence over `MINIMAX_API_KEY`)
- `-base-url`: API endpoint (default: `https://api.minimax.io`)
- `-text`: text to synthesize
- `-model`: model name (default: `speech-2.6-hd`)
- `-voice-id`: optional voice ID
- `-speed`: optional speech speed
- `-volume`: optional speech volume
- `-timeout`: request timeout (default: `30s`)
- `-output`: output file path (default: `speech_output.audio`)

### Environment variables

You can configure the same options via environment variables:

- `MINIMAX_API_KEY`
- `MINIMAX_BASE_URL`
- `MINIMAX_SPEECH_TEXT`
- `MINIMAX_SPEECH_MODEL`
- `MINIMAX_SPEECH_VOICE_ID`
- `MINIMAX_SPEECH_SPEED`
- `MINIMAX_SPEECH_VOLUME`
- `MINIMAX_SPEECH_TIMEOUT`
- `MINIMAX_SPEECH_OUTPUT`
