---
status: approved
granted: 2026-09-30
action: scope the Spec citation projection to the Spec's own Task files and make the Relocation Citation scan read tracked files through an os.Root of the repository
consuming: 0186-citation-scans-that-read-only-what-they-mean-to
paths: []
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0186

On 2026-09-30 the maintainer instructed "Continue até o final para o release de
todas as implementações e ajustes". This Spec delivers two of the open
adjustments. The 2026-09-29 "Autonomia ampla" standing authority covers
authoring, corrective Tasks and delivery through merge.

## Governed paths

None. Measured against `GovernedPath` on 2026-09-30:
`internal/speccheck/citations.go`, `internal/baseline/history_citations.go`,
`internal/baseline/history_citations_open_unix.go`,
`internal/baseline/history_citations_open_windows.go`, the two new test files
and this Spec's own files are ordinary.

## What is not governed and must stay untouched

`internal/baseline/repository_test.go`, `internal/baseline/plan_test.go`,
`internal/cli/cli_test.go`, `docs/references/coverage-record.json`, `go.mod` and
the Makefile are governed or out of scope. No Task edits them. A Task that finds
it cannot avoid one of them stops and reports instead of editing it.

## Sanctioned regeneration

None.

## Limits

- No command output, flag, exit code or schema change.
- No new dependency and no `go.mod` edit.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
