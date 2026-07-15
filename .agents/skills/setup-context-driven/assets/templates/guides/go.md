# Go

Keep command entry points thin and push behavior into cohesive internal
packages. Prefer stdlib APIs and justify every dependency.

Use context-first signatures for blocking and IO work. Wrap errors with the
failed operation and preserve matching through `errors.Is` or `errors.As` when
callers need it.

Use stdlib `testing`. CLI tests run commands through their public runner and
assert stdout, stderr, files, and exit codes.
