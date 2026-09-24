---
status: approved
granted: 2026-09-24
action: add the authorized QA Archive Override to the archive command, refuse claimed ADR ordinals in Spec check, and state settlement and authoring rules in the skills
consuming: 0159-archive-override-and-authoring-rules
paths:
  - internal/spec/archive.go
  - internal/spec/archive_test.go
  - internal/cli/cli_test.go
  - internal/speccheck/coherence.go
  - internal/docscontract/testdata/corpus-golden.json
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/write-tasks/references/task-template.md
  - skills/write-tasks/SKILL.md
  - .agents/skills/qa-gate/SKILL.md
  - skills/qa-gate/SKILL.md
  - .agents/skills/archive-spec/SKILL.md
  - skills/archive-spec/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0159

Two grants are consumed here.

The skill files ride the maintainer's express authorization of 2026-09-24,
given for exactly this content: the two declarations of Spec 0158, the real
`archive --qa-override` flag with its provenance, one settlement table across
the QA gate, archive and Roundfix skills, and guidance on temporal
prerequisites, property-shaped acceptance and CI feasibility.

The Go sources and goldens ride the standing grant of 2026-09-21, which covers
any governed path a slice genuinely needs, on the condition that the bounded
set is measured rather than predicted and recorded here with its reason. The
set was measured with `GovernedPath` itself.

## Why each governed path is unavoidable

- `internal/spec/archive.go`, `internal/spec/archive_test.go` — the override is
  an archive eligibility rule, and archive eligibility lives there.
- `internal/cli/cli_test.go` — the archive command's help changes, and the
  command surface's tests live there.
- `internal/speccheck/coherence.go` — the detector registry; a new code is
  registered there.
- `internal/docscontract/testdata/corpus-golden.json` — pins per-code finding
  counts over the active corpus; a new code may need its entry.
- The eight skill files — the authorized guidance; `.agents/skills/` is
  canonical and `skills/` is its distributed mirror.

## What is not governed

`internal/cli/archive.go`, `internal/cli/archive_test.go`,
`internal/speccheck/ordinal.go` and its test,
`skills/write-tasks/references/task-template.md`,
`skills/settlement_guidance_repocontract_test.go`,
`docs/user-guide/commands.md` and `docs/adr/` are ordinary.

## Sanctioned regeneration

```yaml
command: make skills-sync
```

## Limits

- No action, operation or path beyond those above.
- No Baseline module, linter or Verification configuration edit.
- An override never changes a QA Task's status or a QA Report's verdict.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
