# Static gate and build evidence

- Audited head: `6ebb40e4f7eedb53c6c153a02461be1302578f67`.
- Resolved toolchain: `GOTOOLCHAIN=go1.26.7`; `go version` returned `go version go1.26.7 darwin/arm64`.
- Formatter: `/Users/marcio/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/gofmt`.
- `rtk roundfix spec check 0134-a-governed-deletion-the-gate-can-see --strict` exited 0 with `No findings. Authored Verification commands were not executed.`
- The rebuilt `bin/roundfix spec check 0134-a-governed-deletion-the-gate-can-see --strict` also exited 0 with the same result.
- The first full-gate attempt was stopped when the sandbox blocked a test process from reaching `api.github.com`. One unchanged full-access retry of `rtk env GOTOOLCHAIN=go1.26.7 GOCACHE=/tmp/roundfix-go-build-0134-qa make GO=go GOFMT=/Users/marcio/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/gofmt verify` exited 0. The run completed `go test -parallel 16 ./...`, the focused Skills tests, `roundfix skills check`, and the build.
- `rtk env GOTOOLCHAIN=go1.26.7 GOCACHE=/tmp/roundfix-go-build-0134-qa make GO=go GOFMT=/Users/marcio/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/gofmt fmt-check` exited 0.
- The rebuilt binary reported `roundfix 0.12.0 (6ebb40e4-dirty, built 2026-09-14 07:06:38 -0300)`. The dirty marker is the untracked Daemon-seeded QA directory; `git status --porcelain=v1 --untracked-files=no` was empty.

The strict checker retained skips for absent Finding, rollup, backlog, and references artifacts. It also skipped the vocabulary detector, so the QA matrix independently checks the no-new-token contract.
