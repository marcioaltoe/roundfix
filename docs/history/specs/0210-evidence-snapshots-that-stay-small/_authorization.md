---
status: approved
granted: 2026-09-30
action: record one digest per declared input in a QA Report's evidence snapshot, carry rows by comparing those digests, keep reading the per-file form, and describe the record in the Context-Driven Development guide
consuming: 0210-evidence-snapshots-that-stay-small
paths: []
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0210

On 2026-09-30 the maintainer asked for unattended work through every release
of the current program, v0.24.0 included, under the broad autonomy granted on
2026-09-29 for authoring, corrective work and delivery through merge. Spec
0202 ships in v0.24.0, and on 2026-10-01 its Evidence Snapshot grew Spec
0203's QA Report to 157,212 lines and broke the pre-PR review. This Spec
repairs that defect before the release.

## What is not governed

The set was measured with `GovernedPath` on the authoring branch at
`c3be3bc9`, through a `go test -overlay` probe that wrote nothing to the
repository. `internal/speccheck/report.go`,
`internal/speccheck/mechanical.go`, `internal/speccheck/evidence_record.go`,
`internal/speccheck/evidence_record_test.go`, the new
`internal/speccheck/evidence_digest_test.go`,
`internal/daemon/qa_evidence_snapshot_test.go`,
`internal/daemon/qa_two_pass_carry_test.go`,
`docs/user-guide/context-driven-development.md` and ADR-0210 are ordinary.
The governed `internal/speccheck/mechanical_test.go` is not touched; its carry
tests must pass unedited. No skill is edited.

## Limits

- No change to archived Specs, existing QA Reports, verdict rules, typed
  blocked counts, report naming, the `### QA settlement` section, archive
  eligibility or Task Carry-Forward.
- No edit to the Makefile, to lint, formatter or test-runner configuration,
  to a CI workflow or to `go.mod`.
- The live Run Database under `~/.roundfix` is never opened for writing by a
  Task or the gate, and no test reaches GitHub or a provider.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
