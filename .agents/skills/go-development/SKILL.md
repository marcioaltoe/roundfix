---
name: go-development
description: Use when writing, modifying, or reviewing Go production code in this repository, including CLI commands, config loading, SQLite persistence, GitHub integration, ACP process execution, daemon loops, concurrency, error handling, and package boundaries.
---

# Go development

Build the smallest idiomatic Go change that satisfies the current Roundfix
contract.

## Workflow

1. Read `CONTEXT.md`, relevant ADRs, and the code that owns the behavior.
2. Search with `rtk rg` or `rtk rg --files` before changing package boundaries.
3. Keep `cmd/roundfix` as startup glue. Put behavior in `internal/...`.
4. Prefer the standard library. Add dependencies only when they remove real
   complexity or provide a mature domain primitive.
5. Add dependencies with `rtk go get`, not by editing `go.mod` manually.
6. Run `rtk gofmt -w <changed-go-files>` and `rtk go test ./...`.

## Design rules

- Keep functions small and named after the domain action they perform.
- Define narrow interfaces at the consuming package, not at the provider by
  default.
- Accept interfaces and return concrete structs when that reduces coupling.
- Pass `context.Context` first for blocking, IO, process, network, database,
  and daemon-boundary operations.
- Return errors with context using `%w`; match with `errors.Is` and
  `errors.As`, not string comparisons.
- Avoid package-level mutable state. If shared state is required, name the
  owner and shutdown path.
- Use `slog` for operational logs once logging exists. Use command output
  writers for CLI user-facing text.
- Do not use `panic` or `log.Fatal` in production paths outside unrecoverable
  startup in `main`.

## Concurrency rules

- Every goroutine needs an owner, cancellation path, and completion signal.
- Use `select` with `ctx.Done()` in long-running loops.
- Avoid unbounded channels and unbounded retries.
- Treat `max_run_duration`, token budget, and review rounds as different
  controls. Do not collapse them into one loop guard.
- Prefer explicit state transitions over implicit boolean flags.

## Roundfix boundaries

- `fetch` downloads review issues and never starts an ACP runtime.
- `resolve` works over downloaded unresolved issues and does not fetch.
- `watch` coordinates fetch and resolve over configured review rounds.
- Fail during preflight before remote fetches or agent startup when config,
  artifact directory, PR state, database path, or runtime availability is
  invalid.
- Keep markdown artifacts separate from the global `roundfix.db`.

## Verification

Use these checks before finishing Go work:

```bash
rtk gofmt -w <changed-go-files>
rtk go test ./...
```

If CLI behavior changed, also run:

```bash
rtk go run ./cmd/roundfix --help
```

If concurrency changed, also run:

```bash
rtk go test -race ./...
```
