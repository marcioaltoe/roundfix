---
status: accepted
created_at: 2026-10-19T20:00:00Z
updated_at: 2026-10-20T01:00:00Z
deprecated_at: null
superseded_by: null
---

# The Linux Lane runs like CI

The first Linux Lane run that executed, on Spec 0280's candidate on
2026-10-18, parked the item `gate-failed` on two environment gaps while every
changed package passed. `TestGuidanceCompositionJourney/rust-cli` failed
because the `golang:1.26` image has no `rustfmt`, and
`TestCLITestHelpersEndWithTheirTestBinary` timed out after 594 s because a
helper orphaned inside the container was never reaped. The maintainer turned
the lane off the same day. On 2026-10-19 the maintainer chose "--init +
imagem configurável + skip com motivo (Recomendado)".

## Decision

**The lane runs under an init.** `docker run` gains `--init`, so
`docker-init`, which Docker documents as backed by tini, is PID 1 and reaps
every orphan. Without it, PID 1 is `go test`, which waits only for its own
children, and an orphan that has exited stays a zombie that `kill(pid, 0)`
still finds.

**The image is configurable.** `delivery.linux_lane_image` names the image.
It is read from the same fetched default-branch Project Config as
`delivery.linux_lane`, with User Config beneath it. Empty, the default, keeps
`golang:<go.mod go version>`. A value must be a container image reference
(`name[:tag][@digest]`, the grammar of the distribution reference library), so
it can never start with `-` or carry a space. The lane still never pulls: a
configured image that is not present locally is skipped like the default one.

**A test that needs an External Test Tool skips naming it.** A tool outside
the repository's Base Test Tools (`go`, `gofmt`, `git`, `make` and the POSIX
utilities the tests run) is resolved through `testfixture.RequireTool`, which
skips the test with `requires <tool>, which is not on PATH` when the tool is
absent, on any machine and in any container. A repository rule in
`internal/testfixture` refuses a test that resolves such a tool by a literal
name in any other way.

**A container failure is told apart from a change failure.** The lane script
checks for `tar`, `git` and `go`, unpacks the candidate and records it before
`go test` starts, and any failure there exits `3` with a
`roundfix-linux-lane: container: <cause>` line. Measured on 2026-10-19 in
`golang:1.26`, `go test` exited `1` for a failing test, a build failure and a
panic, and the go command keeps `2` for its own usage errors, so these
outcomes do not overlap:

- **Lane Container Failure**: exit `3`, `126` or `127`, an exit of `128` or
  more (the container process ended by a signal), or the 20-minute cap. The
  item parks `linux-lane-container: <cause>; see <lane log>` in the
  `environment` Park Class.
- **Lane Change Failure**: any other non-zero exit. The item parks
  `gate-failed: Linux Lane change failure: go test exited <n>; see <lane log>`
  in the `finding` Park Class.

The lane log ends with the outcome's line. The skips ADR-0288 defined, an
absent image and `docker run` exiting `125` among them, stay skips.

**The lane runs on demand.** `roundfix deliver lane [--all] [<revision>]` runs
the same lane against a revision of the current checkout, `HEAD` by default,
without a queue or a Spec and whatever `delivery.linux_lane` says. It reads
the image from the checkout's Project Config, selects the changed packages as
the gate does or every package with `--all`, and exits `0` when it passes, `1`
on a Lane Change Failure, and `3` when it skips or meets a Lane Container
Failure. The operator turns the lane back on in this repository only after an
on-demand run passes, in a change of its own.

## Considered Options

- Building and pinning a repository image with Rust and the other CI tools was
  rejected by the maintainer ("Sem build de imagem própria"). An operator who
  wants those tools names a richer image in `delivery.linux_lane_image`.
- Reaping inside the lane script (`sh` as PID 1 running `go test` as a child)
  was rejected. Docker ships an init for exactly this, and a shell that reaps
  only while it waits is not a contract.
- Treating every non-zero exit as the change's failure, as ADR-0288 did, was
  rejected. Exit `126` and `127` come from Docker, and a missing tool before
  `go test` is the container's, so parking them as findings sends the
  operator to fix a change that is correct.
- Turning an absent image or exit `125` into a container failure was rejected.
  ADR-0288 made them skips because the lane never ran; changing that would
  park every item on a machine without the image.
- Reading the image from the item branch was rejected for the same reason
  ADR-0288 reads the opt-in from the fetched default: an item branch must not
  choose the container its own gate runs in.
- A flag on `deliver start` instead of a subcommand was rejected: the operator
  needs the lane's outcome without recording a queue.

## Consequences

- With the default image, the Rust journey skips in the lane and still runs in
  CI, where the runner image carries `rustfmt`. A CI runner that lost
  `rustfmt` would skip it too.
- A change that breaks the module for an older configured image's toolchain
  fails as a Lane Change Failure, because `go test` reports it; the log names
  the toolchain.
- The Park Class table gains `linux-lane-container`, and `deliver --help`
  gains `deliver lane`.

**Amended by [ADR-0297](0297-the-linux-lane-reads-the-repository-history-read-only.md) (2026-10-19).** The container also mounts the repository object store that `git rev-parse --path-format=absolute --git-path objects` names, read-only at `/repo-objects`, and the lane script registers it as an alternate of the container repository right after `git init`, so a test that reads a historical object behaves as in CI. The working tree, refs and a writable store are still never mounted; a store whose path holds a colon, that is not a directory, or that borrows objects through `info/alternates` skips the lane with its reason.
