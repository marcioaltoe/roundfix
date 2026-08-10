<!-- setup-context-driven:begin id=guide.go version=0.0.1 -->

# Go

- Keep Go command entry points thin and behavior in cohesive packages. Prefer the standard library; add a dependency only for a named job it cannot perform, and change `go.mod` and `go.sum` through Go tooling rather than hand edits.

- Use context-first signatures for blocking and IO work. Give every goroutine an owner and cancellation path. Wrap errors with the failed operation using `%w`, and preserve `errors.Is` or `errors.As` matching where callers need it.

- Use stdlib `testing`. Test observable package and CLI behavior through public runners, including stdout, stderr, files, exit codes, cancellation, and failure paths; do not add production-only test hooks.

<!-- setup-context-driven:end id=guide.go -->
