---
type: fix
status: open
created: 2026-09-25
spec: null
reason: null
---

# The blocked required-removed reconcile result breaks the output contract

## Symptom

When a required skill was removed upstream, `baseline skills reconcile --format
json` prints `"plannedChanges": null` instead of `[]`, and exits 3, which the
help documents only for plan confirmation, so an agent looks for a plan digest
that does not exist.

## Where

`internal/baseline/skills_reconcile.go` (required-removed branch) and the exit
mapping in `internal/cli/cli.go`.

## Expected

`plannedChanges: []`, and either a documented exit for the required-removed
block or its own category; a CLI test for the JSON shape and exit code.

## Evidence

Recorded limit of Spec 0163 (pre-PR review of 2026-09-25).
