# Static gate evidence

Build: `c6a5ea156c70d33d59e9e5e1d98fcc2fe24cb318`.

- `rtk roundfix spec check 0130-documentation-cleanup-compatibility --strict`
  exited 0 with `No findings. Authored Verification commands were not executed.`
  It reported the missing `docs/backlog`, Vocabulary Contract and references index
  inputs as skipped detectors; QA retained the applicable manual checks.
- `rtk make verify` exited 2 at `fmt-check`. It named
  `internal/cli/baseline_skills_restore_test.go` and
  `internal/cli/baseline_assets_sync_test.go` as needing formatting.
- `rtk gofmt -d internal/cli/baseline_skills_restore_test.go
  internal/cli/baseline_assets_sync_test.go` exited 1 and showed indentation-only
  differences in the two files.
- `rtk git diff --name-status 54614682..c6a5ea15 --
  internal/cli/baseline_skills_restore_test.go
  internal/cli/baseline_assets_sync_test.go` returned no paths. The formatter
  failure therefore predates the authorized repair and is outside its scope.
- `rtk make verify-docs` exited 0. It built the CLI, passed `docscontract`, passed
  the tagged Baseline and Skill regeneration contracts, passed the corpus budget,
  and reported no findings for active Spec consistency checks.


## 2026-09-09 — same-commit toolchain control

The Supervisor independently ran the full pipeline in a clean detached clone
of `c6a5ea156c70d33d59e9e5e1d98fcc2fe24cb318`:

```sh
rtk proxy env GOTOOLCHAIN=go1.26.7 GOCACHE=/Users/marcio/Library/Caches/go-build make GOFMT=/Users/marcio/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/gofmt verify verify-docs
```

Result: exit 0. Formatting, the complete Go suite, Skill checks, build,
documentation contracts, regeneration contracts, corpus budget and Spec checks
passed. The retained [transcript](05-toolchain-control.txt) contains the command
outputs. An independent `git rev-parse HEAD` still returned the audited commit;
`git status --porcelain` exited 0 with no output.

The control used Go/gofmt 1.26.7. The original system formatter listed two CLI
test files, while the explicit 1.26.7 formatter listed none on identical bytes.
The original `make verify` failure remains an observed result, but the inference
that these source files require formatting is withdrawn. F-001 is resolved as an
execution-environment mismatch. No source/configuration edit or waiver was used.

Transcript SHA-256: `dae0c1b09cca878b1975b7f4332232b1553a86f5b35545af016ebe9d30423363`.
