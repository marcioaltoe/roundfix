# Static gate and build evidence

- Audited head and Spec target head:
  `b2c1d03ecb7cf63e905f94bc850b615218347471`.
- Resolved toolchain: `GOTOOLCHAIN=go1.26.7`; `go version` returned
  `go version go1.26.7 darwin/arm64`.
- Formatter:
  `/Users/marcio/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/gofmt`.
- Installed `roundfix 0.12.0` ran
  `rtk roundfix spec check 0134-a-governed-deletion-the-gate-can-see --strict`
  and exited 0 with no findings. Authored Verification commands were not run.
- The sandboxed first Verification attempt was stopped when a test process
  reached `api.github.com`. The unchanged full-access run reached the suite but
  seven old Task-cycle tests timed out after 6-9 seconds under whole-suite
  scheduling pressure. Those seven tests then passed together on the unchanged
  tree in 1.30 seconds; their bodies predate this Spec.
- One final unchanged full-access run of `rtk env GOTOOLCHAIN=go1.26.7
  GOCACHE=/tmp/roundfix-go-build-0134-rerun make GO=go
  GOFMT=/Users/marcio/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/gofmt
  verify` exited 0. It passed every Go package, the focused Skills tests,
  `roundfix skills check`, and the build. `internal/daemon` passed in 6.456
  seconds.
- The rebuilt binary reported
  `roundfix 0.12.0 (b2c1d03e-dirty, built 2026-09-14 07:44:06 -0300)` and
  repeated the strict Spec check with no findings.
- `rtk git status --porcelain=v1 --untracked-files=no` produced no tracked
  changes. The binary's dirty marker records the untracked Daemon-seeded QA
  report and rerun evidence.

The clean final full gate is the Verification result. The earlier sandbox
denial and transient timeout run are retained as environment deviations rather
than omitted.
