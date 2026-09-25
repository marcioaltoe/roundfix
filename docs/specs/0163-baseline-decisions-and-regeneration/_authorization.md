---
status: approved
granted: 2026-09-24
action: keep HTTP decision fields on a mode change, refuse unreachable Greenfield adoption, declare skill regeneration outputs, reconcile obsolete skill-lock entries on proven absence, and update the shipped skills
consuming: 0163-baseline-decisions-and-regeneration
paths:
  - internal/cli/baseline_human_test.go
  - internal/baseline/plan_test.go
  - internal/baseline/derived_ownership_test.go
  - internal/speccheck/mechanical_test.go
  - internal/cli/baseline_documentation_contract_test.go
  - skills/_ownership.yml
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

# Approved authority for Spec 0163

On 2026-09-24 the maintainer expressly authorized the governed paths listed in
Spec 0121's `_authorization.md`, with sanctioned regeneration by
`make baseline-digests` and `make skills-sync`. This record keeps the paths
that this Spec's Tasks change. It drops the paths that belonged only to Spec
0121 Core Feature 2 (cut) or Core Feature 6 (not carried here). Narrowing
removes authority and needs no further approval.

## Why each governed path is unavoidable

- `internal/cli/baseline_human_test.go` — holds the interactive Baseline
  command's tests. The HTTP mode change and the Greenfield refusal are proved
  there. Governed by history (ADR-0130).
- `internal/baseline/plan_test.go` — holds the `BuildPlan` tests, where the
  JSON planning outcome of the Greenfield refusal is proved. Governed by
  history.
- `internal/baseline/derived_ownership_test.go` — holds the `OutputsFor`
  tests, where the `make skills-sync` output set is proved. Governed by
  history.
- `skills/_ownership.yml` — the new code-generator declaration that assigns
  the skill mirrors to `make skills-sync`. `GovernedPath` reports it as
  ordinary source; it is kept because the maintainer's list names it.
- `.agents/skills/roundfix/SKILL.md` — the canonical Roundfix skill. It
  documents `baseline skills reconcile` and the changed Baseline behaviour.
  Required by the skill-sync hard rule.
- `skills/roundfix/SKILL.md` — the shipped mirror, written only by
  `make skills-sync`.
- `.agents/skills/setup-context-driven/SKILL.md` — the canonical adoption
  skill. It describes the HTTP mode change, the Greenfield refusal and lock
  reconciliation.
- `skills/setup-context-driven/SKILL.md` — the shipped mirror, written only by
  `make skills-sync`.

## What is not governed

Measured with `GovernedPath`: `internal/cli/baseline_human.go`,
`internal/baseline/preservation.go`, `internal/baseline/derived_ownership.go`,
`internal/speccheck/mechanical.go`, the new
`internal/baseline/skills_reconcile.go` and
`internal/cli/baseline_skills_reconcile.go` with their tests,
`internal/cli/baseline_profile.go`, `internal/cli/cli.go`,
`internal/baseline/preservation_test.go` and `docs/user-guide/commands.md`
are ordinary source.

## Sanctioned regeneration

```yaml
command: make baseline-digests
```

```yaml
command: make skills-sync
```

## Added after the first Run

`internal/speccheck/mechanical_test.go` rides the standing grant of 2026-09-21
for governed paths a slice genuinely needs: the mechanical audit's tests call
`baseline.OutputsFor`, so declaring the skill regeneration outputs changes their
expectations. The first Run proved the need: Task 03 edited it and the QA
authorization audit refused the unbounded path.

`internal/cli/baseline_documentation_contract_test.go` also rides the standing
grant of 2026-09-21: it parses every documented `roundfix baseline` example, and
the new `baseline skills reconcile` example needs its case. The second Run's QA
precondition proved the need.

## Limits

- No action, operation or path beyond those above.
- No Baseline asset under `internal/baseline/assets/`, no `Makefile`, and no
  file under `internal/suiteguardcontract/` changes.
- No real repository's `skills-lock.json` and no installed skill tree is
  changed; tests use disposable repositories and local Git sources.
- No paid API use, network call in tests, release, tag, deployment or
  branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
