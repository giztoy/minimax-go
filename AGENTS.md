# AGENTS.md
Guide for coding agents operating in this repository.
Scope: workflow, build/lint/test commands, and code style.

## 1) Repository model
- Local development + direct merge.
- No mandatory PR gate.
- No mandatory CI gate.
- Review basis: static code review + local test evidence.

## 2) Task source of truth
Read in this order:
1. `openteam/task.md` (goal + acceptance)
2. `openteam/plan.md` (plan + required fixes)
3. `openteam/test.md` (if present)
4. Repository code/tests
5. `openteam/worklog.md` (execution/test logs)
If `design_proposal.md` / `test.md` is missing, treat as non-blocking unless task text says otherwise.

## 3) Build / format / lint / test commands
Repository type: Go module (`go.mod` in root).

### Format all code
```bash
go fmt ./...
```

### Build all packages
```bash
go build ./...
```

### Lint baseline
```bash
go vet ./...
```
Optional if installed:
```bash
golangci-lint run
```

### Run all tests
```bash
go test ./...
```

### Run tests in one package
```bash
go test ./internal/transport
```

### Run one test function
```bash
go test ./internal/transport -run '^TestOpenStream$'
```

### Run one subtest
```bash
go test ./internal/transport -run 'TestOpenStream/http 200 with base_resp error should fail'
```

### Run tests by name pattern
```bash
go test ./... -run 'TestSpeechSynthesize'
```

### Run race detector
```bash
go test -race ./...
```

### Coverage snapshot
```bash
go test ./... -cover
```

### Run speech example
```bash
go run ./examples/speech -h
```

## 4) Coding style guidelines

### Language and comments
- Use English for code, comments, and CLI-facing text.
- Write comments for intent/constraints, not obvious mechanics.

### Imports
- Let `go fmt` manage import ordering.
- Avoid unused imports and dot imports.
- Alias imports only for collisions/clarity.

### Formatting
- Always run `go fmt ./...` before finishing.
- Keep structure and naming consistent across packages.
- Keep struct tags explicit and consistent.

### Types and API design
- Prefer concrete types over `any` when possible.
- For optional numeric JSON fields, use pointers (`*float64`, `*int`) so explicit zero values are representable.
- Keep exported API minimal and stable.
- Separate internal wire structs from public API structs when helpful.

### Naming
- Exported identifiers: `PascalCase`.
- Unexported identifiers: `camelCase`.
- Use Go-style acronyms (`APIError`, `HTTPStatus`, `ID`).
- Test functions: `TestXxx`; subtests should describe behavior.

### Error handling
- Never swallow errors.
- Add context with `%w` wrapping (`fmt.Errorf("...: %w", err)`).
- Use sentinel errors only when callers need `errors.Is`.
- Preserve root cause detail for debugging.

### Context usage
- Accept `context.Context` for network/IO boundaries.
- Propagate context downward.
- Do not create `context.Background()` in library flows.
- Honor cancellation/timeouts in retry and sleep logic.

### HTTP/transport behavior
- Normalize HTTP status errors and business errors (`base_resp`).
- Validate stream content type and parse non-stream responses as potential structured errors.
- Keep retry policy explicit (retryable kinds, max attempts, backoff).
- Never retry `context.Canceled` or `context.DeadlineExceeded`.

### Testing expectations
- Cover success and failure paths.
- Include timeout/cancel/retry edge cases for transport code.
- Assert meaningful outcomes (error type, status code, retry attempts, backoff calls).
- Keep unit tests deterministic and offline unless task explicitly requests online integration.

### Security and secrets
- Never commit secrets/tokens/credentials.
- Use environment variables for local secrets.
- Treat `openteam/` as local collaboration artifacts (gitignored in this repository).

## 5) Static review checklist
Before marking done, verify:
1. Functional correctness vs `openteam/task.md`
2. Error handling quality and diagnosability
3. Edge cases: empty/malformed input, timeout, cancel, retry limits
4. Test coverage for critical paths
5. Implementation completeness (no placeholder behavior)
6. Repository hygiene (no binaries/temp files/sensitive data)
Severity model:
- P0: must fix (correctness/stability/security)
- P1: should fix (maintainability/readability)
- P2: optional improvement

## 6) Required agent outputs
When finishing work, update:
1. `openteam/plan.md` (progress + fix checklist)
2. `openteam/worklog.md` (commands + summarized results)
3. `openteam/review.md` (when acting as reviewer)

## 7) Cursor / Copilot rule files
Checked in this repository:
- `.cursor/rules/` — not found
- `.cursorrules` — not found
- `.github/copilot-instructions.md` — not found
If these files are added later, update this AGENTS.md accordingly.
