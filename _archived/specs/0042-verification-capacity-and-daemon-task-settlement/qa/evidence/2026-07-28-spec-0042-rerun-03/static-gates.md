# Static and race gates

Build: `1b1bfc345af11138c8240d1ce62bd5ddd0065d32`.

The first `rtk make verify` attempt proved the managed sandbox cannot write
the default macOS Go build cache. With
`GOCACHE=/private/tmp/roundfix-qa-0042-rerun03-gocache`, the sandboxed gate
then reached 2,755 passing tests and failed only the five owner-process tests
whose `/bin/ps` calls were denied.

The exact gate reran with the same writable cache and local process-inspection
access:

```text
rtk env GOCACHE=/private/tmp/roundfix-qa-0042-rerun03-gocache \
  GOFLAGS=-buildvcs=false make verify
```

Result: exit 0; 2,760 Go tests passed in 24 packages, four Skill policy tests
passed, all fourteen shipped Skills passed `roundfix skills check`, and the
CLI built at `bin/roundfix`.

The Spec-required race gate also exited 0 with process-inspection access:

```text
rtk env GOCACHE=/private/tmp/roundfix-qa-0042-rerun03-gocache \
  GOFLAGS=-buildvcs=false go test -race ./... -count=1
```

Every package passed, including `internal/cli`, `internal/daemon`,
`internal/runevent`, `internal/store`, `internal/tui`, `internal/worktree`,
and `skills`. No race, blocked worker, or permit leak appeared.
