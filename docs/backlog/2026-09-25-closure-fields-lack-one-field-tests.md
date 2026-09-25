---
type: fix
status: open
created: 2026-09-25
spec: null
reason: null
---

# No test pins the one-field and blank-field closure cases

## Symptom

The archived-Finding closure check is correct, but no test covers a Finding with only `closure_reason`, only `closure_evidence`, or a blank reason, so a regression would pass the suite.

## Where

`internal/speccheck/citations_test.go`.

## Expected

Table cases for reason-only, evidence-only and whitespace-only reason, each expecting `SC-ARCHIVE-LICENSE`.

## Evidence

Recorded limit in `docs/history/specs/0164-knowledge-lifecycle-and-capture/_prd.md`.
