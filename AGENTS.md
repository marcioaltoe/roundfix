# Agent instructions

This repository is Roundfix. It is a local-first Go CLI and future daemon for
fetching pull request review issues, resolving unresolved issues through the
user's selected ACP runtime, and pushing only after the pull request has no
remaining unresolved issues.

## High priority

- Use `rtk` before shell commands.
- Use `ma/` for every agent-created branch.
- Keep the project KISS. Prefer the smallest behavior that satisfies the
  documented product contract.
- Do not copy names, branding, package names, comments, examples, or generated
  artifacts from reference projects into this repository.
- Do not commit, push, rebase, reset, restore, clean, or remove tracked files
  unless the user explicitly asks for that git operation.
- If unexpected user changes exist, read them and work with them. Do not revert
  unrelated work.
- **ALWAYS USE** the `golang-pro` skill before writing any Go code

## Project map

- `CONTEXT.md` stores the project vocabulary and current decisions.
- `docs/product-brief.md` stores the product contract from the grill session.
- `docs/adr/` stores accepted architectural decisions.
- `cmd/roundfix/` is the standalone CLI entry point.
- `internal/app/` holds app metadata.
- `internal/cli/` owns CLI parsing, output, and exit behavior.

## Local skill router

Use the narrowest local skill that matches the task:

- `go-development` + ` before writing or changing Go production code.
- `go-cli` before changing command parsing, interactive prompts, command
  output, or exit codes.
- `go-testing` before writing, changing, or reviewing Go tests.
- `go-tui` + `tui-design` + `tui-glamorous` before building terminal UI, panes, keybindings, or streaming views.
- `systematic-debugging` and `no-workarounds` for bugs, regressions, and
  failing tests.
- `docs-writer` for markdown docs.
- `verification-before-completion` before claiming work is complete.

## Go workflow

- Search local code with `rtk rg` or `rtk rg --files`; do not use web search for
  local code.
- Keep `cmd/roundfix/main.go` thin. Push behavior into `internal/...`.
- Prefer the Go standard library until a dependency has a clear job.
- Add dependencies with `rtk go get`, not by editing `go.mod` manually.
- Pass `context.Context` as the first argument for blocking, IO, process,
  network, database, and daemon-boundary operations.
- Return explicit errors with context. Wrap underlying errors with `%w` and use
  `errors.Is` or `errors.As` for matching.
- Avoid `panic` and `log.Fatal` outside unrecoverable startup in `main`.
- Own every goroutine with cancellation and a clear shutdown path.
- Keep exported comments short and useful. Comment invariants and protocol edge
  cases, not obvious assignments.

## Verification

Run the smallest relevant gate, then the broad gate before finishing:

```bash
rtk gofmt -w <changed-go-files>
rtk go test ./...
rtk go run ./cmd/roundfix --help
```

If concurrency changed, also run:

```bash
rtk go test -race ./...
```

If any required gate fails, report the failing command and do not claim the task
is complete.
