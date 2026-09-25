---
status: approved
granted: 2026-09-24
action: supersede ADR-0123, retire Review Artifacts on recorded evidence, close terminal records without a Spec, document every history family, and name the capture publication owner
consuming: 0164-knowledge-lifecycle-and-capture
paths:
  - internal/baseline/assets/modules/context-workflow.json
  - internal/baseline/assets/modules/secondbrain.json
  - internal/spec/archive.go
  - internal/spec/archive_test.go
  - internal/speccheck/backlog.go
  - internal/speccheck/backlog_test.go
  - internal/docscontract/publicdocs_test.go
  - docs/agents/docs-layout.md
  - docs/agents/secondbrain.md
  - docs/agents/setup-context.json
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0164

On 2026-09-24 the maintainer expressly authorized delivery 6 of the restructured
queue: the governed paths the portfolio Spec 0120 proposed in its
`_authorization.md`, the explicit revision of ADR-0123, and the sanctioned
regenerations below. This Spec carries Core Features 2, 3, 4 and 8 of Spec 0120.
The proposed paths `.agents/skills/archive-spec/SKILL.md` and
`skills/archive-spec/SKILL.md` are dropped because no Task changes them.

## Why each governed path is unavoidable

- `internal/baseline/assets/modules/context-workflow.json` — the canonical
  source of the history-family clause (Task 04).
- `internal/baseline/assets/modules/secondbrain.json` — the canonical source of
  the capture clause (Task 05).
- `internal/spec/archive.go`, `internal/spec/archive_test.go` — the History Root
  resolver gains the handoff family and the list of every family (Task 04).
- `internal/speccheck/backlog.go`, `internal/speccheck/backlog_test.go` — the
  Backlog check reports a terminal entry left in `docs/backlog/` (Task 03).
- `internal/docscontract/publicdocs_test.go` — the guide is tested against the
  resolver's families (Task 04).
- `docs/agents/docs-layout.md`, `docs/agents/secondbrain.md`,
  `docs/agents/setup-context.json` — the managed guides and manifest, rendered
  from the modules by the public Baseline update (Tasks 04 and 05).

## What is not governed

Measured with `GovernedPath`: `internal/spec/review_liveness.go`,
`internal/spec/retirement.go`, `internal/baseline/history_layout.go`,
`internal/speccheck/citations.go`, their tests, and the ADR files under
`docs/adr/` and `docs/history/adr/` are ordinary source.

## Explicit ADR revision

ADR-0123 is superseded, not edited in place: Task 01 adds ADR-0163, moves
ADR-0123 and the proposed ADR-0152 to `docs/history/adr/` with
`status: superseded` and `superseded_by: ADR-0163`, before Task 02 changes the
behaviour ADR-0123 governs.

## Sanctioned regeneration

The repository-owned commands resolve their generated outputs. These
declarations record the regeneration that follows the approved module edits;
they add no source paths. No skill is edited, so `make skills-sync` must leave
the tree unchanged.

```yaml
command: make baseline-digests
```

```yaml
command: make skills-sync
```

## Limits

- No action, operation or path beyond those above.
- No change to the companion Secondbrain repository, its jobs, scripts, mirrors
  or immutable sources.
- Original observations, legacy Spec bytes and existing absorption pointers stay
  intact; an invalid `absorbed_by` is never hidden by closure fields.
- No hand-edited digest pin.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
