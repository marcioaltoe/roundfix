# Static gates

Audited repository head: `2d9ab9b2d95d19ad06b2599be138b410ebeb7404`.

`rtk roundfix spec check 0131-a-failed-gate-accepts-its-repair --strict`
exited 0 before the matrix was populated:

```text
Spec 0131-a-failed-gate-accepts-its-repair
No findings. Authored Verification commands were not executed.
```

`rtk make verify` first exited 2 at `fmt-check`. The system toolchain was Go
1.27.1 and its `/opt/homebrew/bin/gofmt` proposed indentation changes in:

```text
internal/cli/baseline_skills_restore_test.go
internal/cli/baseline_assets_sync_test.go
```

Both paths are byte-unchanged across `62bb7874..2d9ab9b2`, and their last
change is commit `8d333a6c84080bc68d36ec17e5d60a2527a70c78` from 2026-08-03. The
repository's cached Go 1.26.7 toolchain produced no `gofmt -d` output for either
file. This classified the first failure as a formatter-path mismatch rather
than a change in this Spec.

The first compatible-toolchain `make verify` attempt reached the networked
skill check and the sandbox denied access to `api.github.com`. The unchanged
full-access rerun exited 0:

```text
go test -parallel 16 ./...
ok roundfix/internal/cli 66.542s
ok roundfix/internal/spec (cached)
ok roundfix/internal/daemon (cached)
ok roundfix/skills (cached)
Roundfix skill check passed: roundfix, write-idea, write-prd, write-techspec,
write-tasks, setup-context-driven, implement-task, implement-spec,
brainstorming, council, business-analyst, archive-spec, qa-gate, evidence-gate
go build ... -o bin/roundfix ./cmd/roundfix
```

Exact successful command:

```text
rtk make verify GO=/Users/marcio/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/go GOFMT=/Users/marcio/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/gofmt
```

The binary exercised by the public flows reported:

```text
roundfix 0.12.0 (2d9ab9b2-dirty, built 2026-09-10 09:14:09 -0300)
Go build toolchain: go1.26.7
SHA-256: 748c2a1ed248571d68eb5277ce30c3c23260ebbdc0ac722febde718d33b15aeb
```

The `-dirty` suffix records this in-progress QA report and evidence; the
application source is the audited head.

The Pull Request-boundary documentation control also exited 0 with the same Go
1.26.7 selection:

```text
rtk make verify-docs GO=/Users/marcio/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/go GOFMT=/Users/marcio/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/gofmt
```

It passed the documentation contracts, the three tagged repository
regeneration contracts, corpus budget, rebuilt CLI and every active Spec check.
