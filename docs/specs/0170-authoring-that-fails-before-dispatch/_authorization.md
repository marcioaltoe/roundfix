---
status: approved
granted: 2026-09-25
action: refuse undeclared Governed Paths and CLI changes without their guide at authoring, close three checker honesty gaps, and align the authoring templates and skills with the checker
consuming: 0170-authoring-that-fails-before-dispatch
paths:
  - internal/speccheck/constraints.go
  - internal/speccheck/coherence.go
  - internal/docscontract/testdata/corpus-golden.json
  - internal/spec/archive_layout_characterization_test.go
  - .agents/skills/qa-gate/SKILL.md
  - skills/qa-gate/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - .agents/skills/write-prd/references/prd-template.md
  - skills/write-prd/references/prd-template.md
  - .agents/skills/write-techspec/references/techspec-template.md
  - skills/write-techspec/references/techspec-template.md
  - .agents/skills/write-tasks/references/task-template.md
  - .agents/skills/write-tasks/SKILL.md
  - skills/write-tasks/SKILL.md
  - skills/baseline_skill_contract_test.go
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0170

The maintainer approved the 2026-09-25 efficiency sequence (Ondas 0 and 1) in
chat on 2026-09-25; this Spec is part of Onda 1. The Go sources and test files
ride the standing grant of 2026-09-21 for governed source a slice genuinely
needs; the skill and template files ride the standing grant of 2026-09-18 for
keeping the shipped skills true to the CLI. The set was measured with
`GovernedPath`.

## Why each governed path is unavoidable

- `internal/speccheck/constraints.go` — `Check` and the Tooling authority row
  parser live there; the new detectors are called from `Check`, the row must
  carry its bounded paths, and `SC-TOOLING-UNDECLARED` is declared beside the
  other tooling codes.
- `internal/speccheck/coherence.go` — `stagedDetectors` registers each new code
  as a Task-stage detector, and `SC-CLI-UNDOCUMENTED` is declared beside the
  other Task Graph codes.
- `internal/docscontract/testdata/corpus-golden.json`,
  `internal/spec/archive_layout_characterization_test.go` — the corpus golden
  counts every characterized code and its pin must match it; both new codes
  join with count `0`.
- `.agents/skills/qa-gate/SKILL.md`, `skills/qa-gate/SKILL.md` — the table of
  authoring rules removed from the QA matrix names each `spec check` code.
- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` — the Spec
  Consistency Check section lists the stable codes the CLI emits.
- The PRD and TechSpec templates and their mirrors — the Tooling authority row
  label is the defect.
- `.agents/skills/write-tasks/references/task-template.md`,
  `.agents/skills/write-tasks/SKILL.md`, `skills/write-tasks/SKILL.md` — the
  Verification and Context guidance authors follow.
- `skills/baseline_skill_contract_test.go` — the contract test that pins the
  template label.

## What is not governed

`internal/speccheck/undeclared.go`, `internal/speccheck/surface.go`,
`internal/speccheck/citations.go`, `internal/speccheck/verification.go`, their
tests, `internal/spec/collision.go` and its test,
`internal/docscontract/corpus_test.go`, `CONTEXT.md` and the mirror
`skills/write-tasks/references/task-template.md` are ordinary.

## Sanctioned regeneration

```yaml
command: make skills-sync
```

## Limits

- No change to `GovernedPath`, the governed set, or the changed-path audit.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
