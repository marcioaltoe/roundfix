---
type: fix
status: promoted
created: 2026-10-22
spec: 0303-the-linux-lane-names-the-step-that-failed
reason: null
---

# The Linux Lane reports a failed go list as a go test failure

## Symptom

When the `go list` that selects the candidate's packages inside the Linux Lane
container exits non-zero, the lane reports `Linux Lane change failure: go test
exited <n>` and parks the item `gate-failed`, although `go test` never ran.
The operator reading the blocker looks for a failing test that does not exist.
A `go list` ended by a signal is reported as `docker run exited <n>: the
container process ended by signal <s>`, which names neither step.

Quaestor reported it on 2026-10-21 at 19:02 BRT through an ai-memory message:
a Codex review on Quaestor pull request #7 of the vendored Roundfix reference
`.agents/skills/roundfix/references/deliver.md` (about lines 686–692) found
that "every started nonzero exit is reported as `go test exited <n>`, but
package selection runs `go list` first; a `go list` failure never reaches
`go test` and is misclassified".

## Where

`internal/cli/deliver_linux_lane.go`: `linuxLaneScript` runs
`pkgs="$(go list -e ...)" && ... && exec go test -count=1 "$@"` under
`sh -ceu`, so a failed `go list` ends the script with `go list`'s own exit
code; `classifyLinuxLaneOutcome` reads any exit other than `3`, `126`, `127`
and `128` or more as `go test exited <n>` (Park Class `finding`). The same
script and classifier serve the Delivery Queue gate and `roundfix deliver
lane`. The user guide `docs/user-guide/commands/deliver.md` and both copies of
the Roundfix `deliver` reference describe the same rule.

## Expected

The lane script marks a failed `go list` explicitly, with its exit code, and
the outcome names `go list` and that code; the Park Class of a `go list`
failure is decided and documented. Every `go test` outcome keeps its bytes.

## Evidence

The Quaestor report above, verified by the operator against
`linuxLaneScript` and `classifyLinuxLaneOutcome` on 2026-10-22. No Finding.
