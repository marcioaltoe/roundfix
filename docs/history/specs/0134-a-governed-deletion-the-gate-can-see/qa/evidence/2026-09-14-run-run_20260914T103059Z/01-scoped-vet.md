# Scoped analyzer evidence

- Audited head: `b2c1d03ecb7cf63e905f94bc850b615218347471`.
- `rtk env GOTOOLCHAIN=go1.26.7
  GOCACHE=/tmp/roundfix-go-build-0134-rerun go version` reported
  `go version go1.26.7 darwin/arm64`.
- `rtk env GOTOOLCHAIN=go1.26.7
  GOCACHE=/tmp/roundfix-go-build-0134-rerun go vet ./internal/daemon
  ./internal/cli ./internal/speccheck ./internal/suiteguardcontract` exited 0
  with no diagnostics.

This is the analyzer scope declared by amended Task 04. The earlier
repository-wide 31 `copylocks` diagnostics are all in `internal/agent`; Spec
0123 owns their correction and Spec 0124 owns later continuous analyzer
coverage.
