---
status: accepted
created_at: 2026-10-18T17:30:00Z
updated_at: 2026-10-19T20:00:00Z
deprecated_at: null
superseded_by: null
---

# The Delivery Queue gate runs the changed Go packages in a Linux Lane

The Daemon's repository Verification runs on the maintainer's macOS machine,
so a defect that exists only on Linux reaches GitHub Actions first. On
2026-10-17 Spec 0276's fresh item binary copy failed only in CI on Linux with
`fork/exec <copy>: text file busy` (pull request #568). The test that proves
the fix, `internal/cli/deliver_item_binary_fork_linux_test.go`, carries
`//go:build linux`, so no Run on the maintainer's machine ever compiles it. The
operator proved the fix by hand in a `golang:1.26` container: a copy of the
tree, a fresh `git init`, and `go test` of the package, passing with the fix and
failing with it reverted.

## Decision

**The Delivery Queue's repository gate runs a Linux Lane when Project Config
opts in.** With `delivery.linux_lane: true`, after the repository Verification
passes, the gate runs `go test -count=1` over the candidate's changed Go
packages in a `golang:<go.mod go version>` container, once per gate, before the
candidate is pushed. The changed packages come from the candidate's diff
against its merge base with the refreshed default branch: each changed path
maps to the nearest directory at the candidate head, from its own directory
upward, that holds a `.go` file, and a change to `go.mod` or `go.sum` selects
`./...`.

The container is bounded:

- The image is used only when it is already present locally: the lane checks
  `docker image inspect` and runs with `--pull never`. It never pulls.
- `--network none`, `--rm`, and `--user 65534:65534`. Module sources come from
  the host module cache, mounted read-only, with `GOFLAGS=-mod=readonly`,
  `GOPROXY=off` and `GOTOOLCHAIN=local`.
- The tree arrives as `git archive` of the candidate head on standard input and
  is unpacked inside the container, so nothing is written into the checkout.
- The container receives only the variables `go test` needs: `HOME`, `GOCACHE`,
  `GOMODCACHE`, `GOFLAGS`, `GOPROXY`, `GOTOOLCHAIN` and `GOSUMDB`. The docker
  CLI itself keeps the Queue Owner's environment, so its context and socket
  resolve as the operator's do.
- A 20-minute cap and any cancellation remove the container with
  `docker rm -f`.

When the lane cannot run, it skips with a named reason on the Queue Owner's
console log and the gate continues: no `go.mod`, no `go` version in it, no
changed Go package, `go` or `docker` not on `PATH`, `docker info` failing, the
image absent, or `docker run` exiting `125`. Docker documents that exit code
`125` "indicates that the error is with Docker daemon itself"
(https://docs.docker.com/engine/containers/run/), so it is the machine's
failure, not the candidate's. Any other non-zero exit, and the cap, fail the
gate and park the item `gate-failed` with the log path.

## Considered Options

- Running the lane at every Task settlement was rejected by cost. Measured on
  2026-10-18 on the maintainer's machine, compiling the `internal/cli` and
  `internal/worktree` test binaries for Linux from an empty build cache took
  12.9 s, and `go test` of those two packages on the host took 156.9 s. Most
  Specs change `internal/cli`, so a Spec of four Tasks with retries would add
  ten minutes or more, against one run of about three minutes per candidate.
- A repository contract test that shells out to `docker` was rejected. A new
  Repository Contract Test needs a Makefile list entry, which the cycle's grant
  excludes, and a test that skips on a machine without Docker would make a
  Verification depend on the host.
- Pulling a missing image was rejected. A pull is a network fetch inside a gate
  the operator did not ask for; the operator pulls the image once.
- Mounting the checkout into the container was rejected. A Run Worktree's `.git`
  is a file pointing outside the mount, tests need a repository, and a writable
  mount would let the container write into the checkout.
- Turning the lane on for every Go repository was rejected. It would start
  containers on adopters' machines without their choice.

## Consequences

- A candidate whose changed packages fail on Linux parks before any Pull
  Request exists, and its log is under
  `<artifact dir>/delivery/<slug>/linux-lane.log`.
- On a machine without the image, every gate logs one skip line. The image
  `golang:1.26` was not present on the maintainer's machine on 2026-10-18.
- The lane tests only the changed packages, so a Linux-only defect in an
  unchanged package that a change breaks indirectly still reaches CI.
- The Delivery Queue reads `delivery.linux_lane` from the fetched default
  branch, so enabling it takes effect after the change that enables it merges.

**Amended by [ADR-0294](0294-the-linux-lane-runs-like-ci.md) (2026-10-19).** The container runs with `--init`; its image is `delivery.linux_lane_image` when set, read from the same fetched default-branch Project Config, and `golang:<go.mod go version>` otherwise; a test that needs a tool outside the Base Test Tools skips naming it; exit `3` from the lane's own setup, `126`, `127`, an exit of `128` or more and the 20-minute cap are a Lane Container Failure that parks `linux-lane-container` in the `environment` Park Class, while any other non-zero exit is a Lane Change Failure that parks `gate-failed` with its reason; and `roundfix deliver lane` runs the lane on demand.
