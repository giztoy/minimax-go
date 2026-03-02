# AGENTS.md
Guide for coding agents in this repository.
Scope: workflow, build/lint/test commands, style rules, and review checklist.

## 1) Repository model
- Go module repo (`go.mod` at root).
- Local development + direct merge to `main`.
- No mandatory PR/CI gate.
- Quality gate = local review + local command evidence.

## 2) Task source of truth
Read in this order before coding:
1. `openteam/task.md`
2. `openteam/plan.md`
3. `openteam/test.md` (if present)
4. `openteam/design_proposal/*`
5. Existing code/tests/examples
6. `openteam/worklog.md`

If a non-critical doc is missing, treat as non-blocking unless task text says it is required.

## 3) Build / lint / test commands

### Baseline (always run before finishing)
```bash
go fmt ./...
go build ./...
go vet ./...
go test ./...
```

### Run tests by scope
Run one package:
```bash
go test ./internal/transport
go test ./examples/speech
```

Run a single test function (important):
```bash
go test ./ -run '^TestSpeechAsync$'
go test ./examples/speech -run '^TestWaitTask$'
```

Run one subtest:
```bash
go test ./ -run 'TestGetAsyncTask/succeeded with url and no audio bytes is valid'
```

Run by regex pattern across packages:
```bash
go test ./... -run 'TestSpeechAsync|TestGetAsyncTask'
```

Race / coverage / no-cache:
```bash
go test -race ./...
go test ./... -cover
go test ./... -count=1
```

### Example command checks
```bash
go run ./examples/speech -h
go run ./examples/speech async -h
go run ./examples/speech stream -h
go run ./examples/speech http -h
go run ./examples/voice/list -h
go run ./examples/file -h
```

## 4) Code style and design guidelines

### 4.1 Language and comments
- Use English for identifiers, comments, and CLI/help text.
- Comment intent/constraints; avoid narrating obvious mechanics.

### 4.2 Imports and formatting
- Use `go fmt` to manage imports and formatting.
- No unused imports, no dot imports.
- Alias imports only for collision/clarity.

### 4.3 Layering
- Public API at repo root.
- Reusable internals in `internal/transport`, `internal/protocol`, `internal/stream`, `internal/codec`.
- Do not bypass transport/protocol abstractions in new features.

### 4.4 Types and JSON modeling
- Prefer explicit structs over `map[string]any` for stable fields.
- Keep raw payload maps when forward-compatibility is needed.
- Optional numeric JSON fields should use pointers (`*int`, `*float64`) so explicit zero remains representable.
- Separate public structs from wire/raw structs when shapes differ.
- Use typed enums for normalized states (e.g., async task state).

### 4.5 Naming
- Exported: `PascalCase`; unexported: `camelCase`.
- Acronyms in Go style: `APIError`, `HTTPStatus`, `TaskID`, `URL`.
- Tests: `TestXxx`; subtests describe behavior.

### 4.6 Error handling
- Never swallow errors.
- Wrap with context via `%w`:
```go
return fmt.Errorf("query async task: %w", err)
```
- Use `errors.Is` / `errors.As` for classification.
- Fail fast on local validation errors before network calls.

### 4.7 Context handling
- Accept `context.Context` at network/IO boundaries.
- Propagate context downward unchanged.
- No `context.Background()` inside library flows.
- Respect cancellation/deadline in polling/retry loops.

### 4.8 HTTP/protocol behavior
- Normalize HTTP errors and `base_resp` business errors consistently.
- For stream APIs, validate stream semantics and parse structured errors when possible.
- Retry logic must not hide `context.Canceled` / `context.DeadlineExceeded`.

## 5) Testing policy
- Cover success, failure, and boundary cases for each feature.
- Mandatory for network flows: timeout, cancel, retry boundaries, protocol errors.
- Assert meaningful outcomes (state mapping, error type, key fields), not only `err == nil`.
- Unit tests should be offline via `httptest`; online tests must be opt-in by env switch.

## 6) Security and hygiene
- Never commit real secrets/tokens/credentials.
- Placeholder examples like `your_api_key` are acceptable.
- No binaries/temp artifacts/log dumps in commits.
- Treat `openteam/` as local collaboration assets in this repository context.

## 7) Review checklist before merge
1. Meets `openteam/task.md` acceptance criteria.
2. No placeholder logic (`TODO` stubs/fake returns) on required paths.
3. Error handling is contextual and diagnosable.
4. Edge cases covered (invalid input, timeout, cancel, protocol errors).
5. Test evidence is reproducible via commands.
6. No sensitive data or unrelated files; examples remain runnable.

## 8) Required agent outputs
When finishing work, update as applicable:
1. `openteam/plan.md`
2. `openteam/worklog.md`
3. `openteam/review.md` (if acting as reviewer)
4. design proposal status (`draft/wip/done/freeze`) when milestone changes

## 9) Cursor / Copilot rule files
Checked at update time:
- `.cursor/rules/` — not found
- `.cursorrules` — not found
- `.github/copilot-instructions.md` — not found

If these files are added later, merge their instructions into this AGENTS.md and keep this section updated.
