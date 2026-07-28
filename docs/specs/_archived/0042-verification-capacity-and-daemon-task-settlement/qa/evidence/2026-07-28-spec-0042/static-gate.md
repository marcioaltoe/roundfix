# Static gate

Build: `859300203565dc17bfbf01ae4e7a2512e573c17c`

Exact repository gate:

```text
GOCACHE=/private/tmp/roundfix-qa0042-gocache
GOFLAGS=-buildvcs=false
rtk make verify
```

The managed sandbox run reported 2,722 passing Go tests and five failures.
Every failure was a `/bin/ps` `operation not permitted` denial in the
owner-process identity checks. The exact same command with process-inspection
access passed:

```text
Go test: 2727 passed in 23 packages
Go test: 4 passed in 1 packages
Roundfix skill check passed: 14 shipped Skills
go build -buildvcs=false ... -o bin/roundfix
exit 0
```

The full race command was also run:

```text
rtk go test -race ./... -count=1
```

Its sandboxed attempt isolated the same five `/bin/ps` denials. The exact
command with process-inspection access passed every package, including
`internal/agent`, `internal/baseline`, `internal/cli`, `internal/daemon`,
`internal/store`, `internal/tui`, `internal/worktree`, and `skills`.
