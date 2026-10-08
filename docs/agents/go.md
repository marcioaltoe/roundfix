<!-- setup-context-driven:begin id=guide.go version=0.0.1 -->

# Go

These rules govern the repository's Go module: its commands, packages and
tests. Code in another language follows its own guide.

- **prohibited**: Do not hand-edit `go.mod` or `go.sum`; change them through the `go` command.

- **mandatory**: Keep each Go `main` package thin: it parses input, wires dependencies, and calls behavior that lives in cohesive packages.

- **mandatory**: Prefer the Go standard library. Add a third-party module only for a named job that the standard library and the modules already required cannot do, and record that reason in an ADR or another recorded repository decision in the change that adds it.

- **mandatory**: Give every goroutine an owner that waits for it and a cancellation path that stops it.

- **mandatory**: Pass a `context.Context` as the first parameter of blocking and IO work, and stop that work when the context is cancelled.

- **mandatory**: Wrap a returned error with the operation that failed using `%w`, and keep `errors.Is` and `errors.As` matching wherever a caller branches on the error.

- **mandatory**: When a change touches a file with a build constraint, build the affected non-test packages for every operating system its constraints name, for example with `GOOS=windows go build ./...`; a build or test run on the host compiles only the host's files.

- **mandatory**: Keep tests hermetic: set or clear with `t.Setenv` every environment variable the code under test reads, and never read the host's credentials, home directory, or tool state. Create a Unix socket under a short directory from `os.MkdirTemp("", ...)` rather than a deep `t.TempDir()`, because macOS refuses a socket path longer than 104 bytes, and never let a Verification depend on a test that can skip on the host.

- **mandatory**: Test observable package and command behavior through public entry points: stdout, stderr, files, exit codes, cancellation, and failure paths.

- **mandatory**: Run Go tests through `go test` with the standard `testing` package as the harness. An assertion or mocking library is a third-party module and needs its recorded reason.

<!-- setup-context-driven:end id=guide.go -->
