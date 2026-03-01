# Voice Example

`examples/voice` demonstrates voice list queries via `Voice.ListVoices`, including voice type filter and pagination parameters.

## Quick start

```bash
export MINIMAX_API_KEY="your_api_key"

go run ./examples/voice \
  -voice-type all \
  -page-size 20
```

Use `-json` to print the formatted typed response as JSON (`Raw` unknown fields are not included):

```bash
go run ./examples/voice \
  -voice-type system \
  -json
```

## Show all CLI options

```bash
go run ./examples/voice -h
```

## Common flags

- `-api-key`: Minimax API key (takes precedence over `MINIMAX_API_KEY`)
- `-base-url`: API endpoint (default: `https://api.minimax.io`)
- `-voice-type`: voice type filter (`system`, `voice_cloning`, `voice_generation`, `all`)
- `-page-size`: optional page size
- `-page-token`: optional next-page token
- `-timeout`: request timeout (default: `30s`)
- `-json`: print response as formatted JSON

## Environment variables

- `MINIMAX_API_KEY`
- `MINIMAX_BASE_URL`
- `MINIMAX_VOICE_TYPE`
- `MINIMAX_VOICE_PAGE_SIZE`
- `MINIMAX_VOICE_PAGE_TOKEN`
- `MINIMAX_VOICE_TIMEOUT`
