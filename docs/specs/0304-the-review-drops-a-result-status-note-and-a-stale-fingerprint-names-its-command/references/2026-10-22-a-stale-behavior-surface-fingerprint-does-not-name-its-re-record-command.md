---
type: fix
status: promoted
created: 2026-10-22
spec: 0304-the-review-drops-a-result-status-note-and-a-stale-fingerprint-names-its-command
reason: null
---

# A stale Behavior Surface fingerprint does not name its re-record command

## Symptom

When a guide, a command's help or the configuration keys change without the
Behavior Surface Record being re-recorded, `TestTheSkillCoverageMapIsCurrent`
reports only `<surface id>: stale fingerprint`. Neither the Agent that edited
the guide nor the operator reading the failure learns which command fixes it.
On 2026-10-22 Spec 0298's Delivery Queue item parked `gate-failed` and needed a
retry, and a Run of Spec 0299 ended Unresolved after its documentation Task
failed its Verification twice, both for a stale fingerprint after a guide edit.

## Where

`internal/docscontract/skill_coverage_test.go`, `skillCoverageProblems`, the
`stale fingerprint` problem. The command that re-records the record is
declared in `.roundfixrc.yml` under `delivery.derived_paths` for
`docs/references/behavior-surfaces.json`:
`go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.

## Expected

The stale fingerprint problem names that command and says to run it after the
last edit of the surface's sources. The test's other problems keep their
bytes.

## Evidence

none yet
