# Voice Design Example

`examples/voice/design` demonstrates `Voice.DesignVoice`.

## Quick start

```bash
export MINIMAX_API_KEY="your_api_key"

go run ./examples/voice/design \
  -prompt "calm and warm podcast host" \
  -preview-text "hello, welcome to our podcast"
```

Use `-json` to print the typed response as JSON (`Raw` unknown fields are not included):

```bash
go run ./examples/voice/design \
  -prompt "clear male narrator" \
  -preview-text "this is a preview" \
  -json
```

## Show all CLI options

```bash
go run ./examples/voice/design -h
```

## Common flags

- `-api-key`: Minimax API key (takes precedence over `MINIMAX_API_KEY`)
- `-base-url`: API endpoint (default: `https://api.minimax.io`)
- `-prompt`: voice description prompt
- `-preview-text`: preview text
- `-voice-id`: optional custom voice ID
- `-timeout`: request timeout (default: `30s`)
- `-json`: print response as formatted JSON

## Environment variables

- `MINIMAX_API_KEY`
- `MINIMAX_BASE_URL`
- `MINIMAX_VOICE_DESIGN_PROMPT`
- `MINIMAX_VOICE_DESIGN_PREVIEW_TEXT`
- `MINIMAX_VOICE_DESIGN_VOICE_ID`
- `MINIMAX_VOICE_DESIGN_TIMEOUT`
