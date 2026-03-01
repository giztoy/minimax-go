# File Upload Example

`examples/file` demonstrates file uploads via `File.Upload`. It reads a local file, sends multipart upload request to Minimax, then prints uploaded metadata from the response.

## Quick start

```bash
export MINIMAX_API_KEY="your_api_key"

go run ./examples/file \
  -input /path/to/local/file.wav \
  -purpose voice_clone
```

If successful, the command prints upload status, file ID, file URL, and metadata fields.

## Show all CLI options

```bash
go run ./examples/file -h
```

## Common flags

- `-api-key`: Minimax API key (takes precedence over `MINIMAX_API_KEY`)
- `-base-url`: API endpoint (default: `https://api.minimax.io`)
- `-input`: local file path to upload
- `-file-name`: optional uploaded filename override (defaults to base name of `-input`)
- `-content-type`: optional MIME type override
- `-purpose`: purpose field (default: `voice_clone`)
- `-timeout`: request timeout (default: `30s`)

## Environment variables

- `MINIMAX_API_KEY`
- `MINIMAX_BASE_URL`
- `MINIMAX_FILE_INPUT`
- `MINIMAX_FILE_NAME`
- `MINIMAX_FILE_CONTENT_TYPE`
- `MINIMAX_FILE_PURPOSE`
- `MINIMAX_FILE_TIMEOUT`
