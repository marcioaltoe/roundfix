<!-- setup-context-driven:begin id=guide.go version=0.0.1 -->

# Go

- Keep Go command entry points thin and behavior in cohesive packages. Prefer the standard library; add a dependency only for a named job it cannot perform, and change `go.mod` and `go.sum` through Go tooling rather than hand edits.

- Use context-first signatures for blocking and IO work. Give every goroutine an owner and cancellation path. Wrap errors with the failed operation using `%w`, and preserve `errors.Is` or `errors.As` matching where callers need it.

- Use stdlib `testing`. Test observable package and CLI behavior through public runners, including stdout, stderr, files, exit codes, cancellation, and failure paths; do not add production-only test hooks.

<!-- setup-context-driven:end id=guide.go -->

<!-- roundfix:repository-rule:begin id=rule.b8db22d5fd25040757476f158bf97daf69dd6bbe5773390c73bb4128c4386fbb -->
- **NEVER** hand-edit `go.mod`/`go.sum`. Use `rtk go get` / `rtk go mod tidy`.

<!-- roundfix:repository-rule:end id=rule.b8db22d5fd25040757476f158bf97daf69dd6bbe5773390c73bb4128c4386fbb -->

<!-- roundfix:repository-rule:begin id=rule.e8d9aa165083fcc61cc6a852fa5cdca5c52bc4f12055cbc01931292897763f7b -->
- **Zero test dependencies**: stdlib `testing` only — table tests, hand-rolled
  fakes, buffer-captured CLI runs (`Run(args, &stdout, &stderr) int`). **Do
  NOT introduce** testify, mockery, or TUI test harnesses.

<!-- roundfix:repository-rule:end id=rule.e8d9aa165083fcc61cc6a852fa5cdca5c52bc4f12055cbc01931292897763f7b -->

<!-- roundfix:repository-rule:begin id=rule.596377c10f505c35df02e372aea92bcccb98991ac1ea055d03301bb1b9b231a9 -->
- Keep `cmd/roundfix/main.go` thin. Push behavior into `internal/...`.

<!-- roundfix:repository-rule:end id=rule.596377c10f505c35df02e372aea92bcccb98991ac1ea055d03301bb1b9b231a9 -->
