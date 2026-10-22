---
status: accepted
created_at: 2026-10-22T22:00:00Z
updated_at: 2026-10-22T22:00:00Z
deprecated_at: null
superseded_by: null
---

# The Linux Lane names a failed go list apart from go test

The lane script selects the candidate's packages with `go list` before it
runs `go test`, but a failed `go list` ends the script with its own exit code,
so the lane reports `go test exited <n>` for a test run that never started.
Quaestor reported it on 2026-10-21 from a review of the Roundfix `deliver`
reference. The lane script now marks a failed `go list` with its exit code and
a reserved exit, and the outcome names `go list`; its Park Class follows what
can make it fail.

## Decision

**A failed `go list` has its own exit and line.** In the lane script, a
non-zero exit of the package-selection `go list` prints
`roundfix-linux-lane: go list: exited <n>` and exits `4`. No other step of the
script exits `4`, and `go test`, which the script runs last with `exec`,
exits `1` for failing tests, build failures and panics and `2` for its usage
errors (ADR-0294's measurement), so the exit code decides that `go list`
failed and the line carries its code.

**Its Park Class follows the cause.** Measured on 2026-10-22 with the lane's
environment (`GOFLAGS=-mod=readonly`, `GOPROXY=off`, `GOTOOLCHAIN=local`,
`GOSUMDB=off`, an empty module cache), `go list -e` exited `1` for a `go.mod`
that does not parse, for a `go` directive above the running toolchain and for
a `go.work` that names a missing module, and exited `0` for a requirement
absent from the module cache, an import of an unknown module and a Go syntax
error, which `go test` then reports. The Go documentation of `go list` says
that with `-e` "the list command never prints errors to standard error and
instead processes the erroneous packages with the usual printing". A
non-zero `go list -e` exit is therefore a property of the candidate's module
files against the image's toolchain:

- **Lane Change Failure** for an exit `<n>` from `1` to `127`: the cause is
  `go list exited <n>`, the item parks `gate-failed: Linux Lane change
  failure: go list exited <n>; see <lane log>` in the `finding` Park Class,
  and `roundfix deliver lane` exits `1`.
- **Lane Container Failure** for `<n>` of `128` or more, a `go list` ended by
  a signal, as ADR-0294 classifies a container process ended by one: the cause
  is `go list exited <n>: the process ended by signal <n-128>`, the item parks
  `linux-lane-container` in the `environment` Park Class, and `roundfix
  deliver lane` exits `3`.
- An exit `4` without a well-formed line, which only a lane script that lost
  its line produces, is a Lane Change Failure with the cause `the lane script
  exited 4 without reporting the go list exit code`.

The last output line that starts with `roundfix-linux-lane: go list: exited `
and ends in decimal digits names `<n>`; a line that only contains the prefix is
ignored, as for the setup marker.

**Every other outcome keeps its bytes.** Exit `3`, `126`, `127`, `128` or
more, the cap, exit `125` and every other non-zero exit keep ADR-0294's
classification and text, `go test exited <n>` included.

## Considered Options

- Exiting with `go list`'s own code and only printing a line was rejected: the
  exit code would no longer tell which step failed, and a test's output after
  `go test` started could carry a line that looks like the marker.
- Reusing exit `3` with a second marker was rejected: ADR-0294 makes exit `3`
  the script's container setup, a Lane Container Failure, and the class must
  not depend on which line a test printed.
- Classifying every `go list` failure as a Lane Container Failure was
  rejected: every measured failure came from the candidate's `go.mod`,
  `go.work` or `go` directive, and an absent module does not fail `go list
  -e`, so parking in `environment` would send the operator to repair Docker
  for a change that is wrong.
- Telling causes apart by `go list`'s error text was rejected: the text is not
  a contract of the go command, and the classification stays structural.

## Consequences

- An exit `4` that the lane classified as `go test exited 4` before is now a
  `go list` failure; `go test` never exits `4`, so no `go test` outcome
  changes.
- A module the read-only cache lacks still fails as `go test exited 1`, as
  before, because `go list -e` does not fail for it.
- Lane logs written before keep their bytes; only new runs carry the line.

ADR-0294 is amended by this record.
