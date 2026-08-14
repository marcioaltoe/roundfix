# Static and race gates

Build: `ffd6852`.

The first managed-sandbox `make verify` attempt reported 2,740 passes and five
failures. All five failures were `/bin/ps` process-identity checks denied by
the sandbox. No product test failed.

The exact gate was rerun with local process-inspection access:

```text
GOCACHE=/private/tmp/roundfix-qa0042-rerun-gocache
GOFLAGS=-buildvcs=false
rtk make verify
```

Result:

```text
Go test: 2745 passed in 23 packages
Go test: 4 passed in 1 packages
Roundfix skill check passed: 14 shipped Skills
go build -buildvcs=false ... -o bin/roundfix
exit 0
```

The Spec-required race gate also passed with process-inspection access:

```text
rtk go test -race ./... -count=1
```

Every package passed, including `internal/agent`, `internal/baseline`,
`internal/cli`, `internal/daemon`, `internal/runevent`, `internal/store`,
`internal/tui`, `internal/worktree`, and `skills`.
