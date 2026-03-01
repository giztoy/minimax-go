# Speech Stream Example

`examples/speech_stream` demonstrates streaming text-to-speech via `Speech.OpenStream`. It reads SSE chunks incrementally, decodes audio bytes chunk by chunk, and writes merged bytes to a local file.

## Quick start

```bash
export MINIMAX_API_KEY="your_api_key"

go run ./examples/speech_stream \
  -text "hello from minimax-go stream example" \
  -voice-id "male-qn-qingse" \
  -output /tmp/speech_stream_output.audio
```

If successful, the command prints the output file path, total bytes, and chunk count.

## Show all CLI options

```bash
go run ./examples/speech_stream -h
```

## Common flags

- `-api-key`: Minimax API key (takes precedence over `MINIMAX_API_KEY`)
- `-base-url`: API endpoint (default: `https://api.minimax.io`)
- `-text`: text to synthesize
- `-model`: model name (default: `speech-2.6-hd`)
- `-voice-id`: optional voice ID
- `-speed`: optional speech speed
- `-volume`: optional speech volume
- `-timeout`: request timeout (default: `30s`)
- `-output`: output file path (default: `speech_stream_output.audio`)

## Environment variables

Primary stream variables:

- `MINIMAX_API_KEY`
- `MINIMAX_BASE_URL`
- `MINIMAX_SPEECH_STREAM_TEXT`
- `MINIMAX_SPEECH_STREAM_MODEL`
- `MINIMAX_SPEECH_STREAM_VOICE_ID`
- `MINIMAX_SPEECH_STREAM_SPEED`
- `MINIMAX_SPEECH_STREAM_VOLUME`
- `MINIMAX_SPEECH_STREAM_TIMEOUT`
- `MINIMAX_SPEECH_STREAM_OUTPUT`

Backward-compatible fallback variables are also supported for convenience:

- `MINIMAX_SPEECH_TEXT`
- `MINIMAX_SPEECH_MODEL`
- `MINIMAX_SPEECH_VOICE_ID`
- `MINIMAX_SPEECH_SPEED`
- `MINIMAX_SPEECH_VOLUME`
- `MINIMAX_SPEECH_TIMEOUT`
- `MINIMAX_SPEECH_OUTPUT`
