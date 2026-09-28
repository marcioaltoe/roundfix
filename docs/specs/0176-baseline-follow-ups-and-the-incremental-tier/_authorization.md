---
status: approved
granted: 2026-09-28
action: publish the incremental verification tier as a required Baseline decision, route the generated QA Archive Override rule through its command, plan skill-lock changes from the lock present after the fetch, and give a blocked reconcile an empty plan and a documented exit
consuming: 0176-baseline-follow-ups-and-the-incremental-tier
paths:
  - internal/baseline/assets/decisions.json
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/modules/spec-workflow.json
  - internal/baseline/assets/profiles/go-cli-tui.json
  - internal/baseline/assets/profiles/rust-cli.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - internal/baseline/assets/templates/index.json
  - internal/baseline/assets/templates/guides/agent-instructions.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md
  - docs/agents/agent-instructions.md
  - docs/agents/docs-layout.md
  - docs/agents/spec-routing.md
  - docs/agents/setup-context.json
  - internal/baseline/plan_test.go
  - internal/cli/baseline_human_test.go
  - internal/cli/baseline_plan_test.go
  - internal/cli/baseline_release_gate_test.go
  - internal/docscontract/publicdocs_test.go
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - .agents/skills/setup-context-driven/SKILL.md
  - skills/setup-context-driven/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0176

The maintainer approved Onda 2 of the efficiency sequence in chat on 2026-09-28
("Siga para onda 2 como sugerido até o release"), and this Spec is part of it.
On 2026-09-25 the maintainer expressly authorized
`internal/cli/baseline_plan_test.go` and
`internal/cli/baseline_release_gate_test.go` for Spec 0121 Core Feature 6, which
this Spec takes over. The Baseline assets and the other Go test files ride the
standing grant of 2026-09-21 for governed source. The skill files ride the
standing grant of 2026-09-18 for keeping the shipped skills true to the CLI.
The set was measured with `GovernedPath` and by running the sanctioned
regeneration in a disposable copy of main `0160f70a` on 2026-09-28.

## Why each governed path is unavoidable

- `internal/baseline/assets/decisions.json` — declares the
  `verification.incremental` decision (task_04).
- `internal/baseline/assets/modules/core.json` — requires the decision and
  rewords the two-tier clause (task_04).
- `internal/baseline/assets/modules/spec-workflow.json` — carries the
  docs-layout override clause (task_03) and the Spec-routing two-tier clause
  (task_04).
- `internal/baseline/assets/profiles/go-cli-tui.json`,
  `internal/baseline/assets/profiles/rust-cli.json` — list the decision in
  `entryDecisions` (task_04).
- `internal/baseline/assets/profiles/standard-typescript-monorepo.json` — lists
  the decision (task_04). Its digest pin is also rewritten by
  `make baseline-digests` (task_03, task_04).
- `internal/baseline/assets/templates/index.json`,
  `internal/baseline/assets/templates/guides/agent-instructions.md` — the
  agent-instructions template renders the new token (task_04).
- The three formatter golden guides under
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/`
  — rewritten by `make baseline-digests`: `docs-layout.md` for the override
  clause (task_03), and `agent-instructions.md` and `spec-routing.md` for the
  incremental tier (task_04).
- `docs/agents/docs-layout.md`, `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`, `docs/agents/setup-context.json` — this
  repository's managed guides and Setup Manifest, rendered from the modules by
  the public Baseline update (task_03, task_04).
- `internal/cli/baseline_plan_test.go`,
  `internal/cli/baseline_release_gate_test.go` — authorized on 2026-09-25. Their
  complete decision lists and Makefile fixtures gain the new decision and
  target (task_04).
- `internal/baseline/plan_test.go` — measured on 2026-09-28: its shared request
  helpers pass complete decision lists. Without the new decision, 58 of its own
  top-level tests fail, and so do most other Baseline package tests that use
  those helpers (task_04).
- `internal/cli/baseline_human_test.go` — measured on 2026-09-28: its scripted
  prompt answers end before the new decision's prompt, and 17 of its
  top-level tests fail (task_04).
- `internal/docscontract/publicdocs_test.go` — measured on 2026-09-28:
  `TestBaselineDecisionExamples` and `TestProjectConstraintDocumentation` assert
  exactly 15 published decisions, and the published Decision Document must carry
  the new required decision (task_04).
- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` — the Roundfix
  skill documents the reconcile exit (task_01) and the lock refusal (task_02).
- `.agents/skills/setup-context-driven/SKILL.md`,
  `skills/setup-context-driven/SKILL.md` — the setup skill describes the new
  decision and its migration (task_04).

## What is not governed

Measured with `GovernedPath`, these files are ordinary:

- `internal/baseline/skills_reconcile.go`, `internal/baseline/skills_restore.go`,
  `internal/baseline/profile_alignment.go` and their tests;
- `internal/baseline/apply_test.go` and
  `internal/baseline/plan_characterization_test.go`;
- the new test files `internal/baseline/qa_override_clause_test.go`,
  `internal/baseline/incremental_verification_test.go`,
  `internal/baseline/skills_lock_read_test.go` and
  `internal/cli/baseline_incremental_verification_test.go`;
- `internal/baseline/testdata/catalog.*` and
  `internal/baseline/testdata/plan-characterization/*.golden.json`;
- `internal/cli/cli.go`, `internal/cli/baseline_apply_test.go`,
  `internal/cli/baseline_update_test.go`,
  `internal/cli/baseline_skills_reconcile_test.go` and
  `internal/cli/baseline_skills_restore_test.go`;
- `docs/user-guide/commands.md`, `docs/user-guide/context-driven-development.md`
  and `CONTEXT.md`.

## Sanctioned regeneration

The grant names each command, and the tree names its outputs (ADR-0149). No
digest pin is hand-edited. The formatter goldens, the Standard TypeScript
Monorepo digest pin and the catalog test data are fallout of the authorized
module edits (ADR-0081).

```yaml
command: make baseline-digests
```

```yaml
command: make skills-sync
```

This repository's managed guides and Setup Manifest are rendered by the public
Baseline update,
`go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`.
task_04 adds `--adopt-suggested`, and a second refresh must report
`File changes: 0`.

## Limits

- No change to the frozen parity corpus, to any Source Baseline, to
  `internal/baseline/assets/setups/` or to `docs/references/coverage-record.json`.
- No Task moves, edits, supersedes or archives a Spec 0121 file. The operator
  does that at this Spec's archive step.
- No Fiscus file or other adopter repository is changed. The Fiscus mirror is
  read only.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
