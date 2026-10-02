---
task: task_03
spec: 0216-baseline-wording-left-after-the-stack-wave
status: completed
type: backend
complexity: medium
---

# Task 03: The skill guide tells the Agent to ask for a skill only a person can start

## Overview

Four modules dispatch nine triggers to skills whose metadata turns off model invocation (`cut-release`, `handoff`, `grill-with-docs`, `implement-spec`, `setup-context-driven`, `architectural-analysis`, `refactoring-analysis`, `to-prompt`), so an Agent is told to activate a skill it cannot load. This Task adds one mandatory core clause to the skill guide of every profile: ask the person to run such a skill instead of activating it or carrying out its workflow. It gives the clause its Source Baseline row and refreshes this repository's own skill guide. It answers the Backlog Entry "A dispatch trigger names a skill the model cannot invoke" of 2026-09-30.

This is an authorized tooling Task. It may change only the files in its Context, the derived files the sanctioned regeneration rewrites, and this Task file.

## Requirements

1. MUST add `clause.core.ask-the-person-to-start-a-person-only-skill`, enforcement `mandatory`, with the guidance the TechSpec's "Exact texts" gives for task_03, as the last clause of `rule.core.skill-dispatch` in `internal/baseline/assets/modules/core.json`, editing in place.
2. MUST add the Source Baseline row the TechSpec's "Source Baseline rows" table gives for task 03.
3. MUST raise by one, from the value on the starting main, the versions the TechSpec's "Version changes" lists for task_03.
4. MUST add the clause to the force record and raise the maintained Source Baseline entry count by one.
5. MUST create `internal/baseline/person_only_skill_test.go` with the two tests the TechSpec's Testing Approach 3 names.
6. MUST run `make baseline-digests` twice (the second reports `"changed":false`), then the Managed Refresh twice; the second MUST report `File changes: 0`, and this repository's `docs/agents/skill-dispatch.md` MUST state the clause. MUST NOT hand-edit a snapshot, golden, digest or Source Baseline offset.
7. MUST NOT change, remove or mark any dispatch trigger or Skill Activation, MUST NOT edit any skill, and MUST NOT rename or remove a top-level test.

## Subtasks

- [ ] Add the clause and its Source Baseline row; raise versions.
- [ ] Update the force record and entry count; create the new test file.
- [ ] Regenerate and refresh twice.

## Acceptance Criteria

- [ ] The skill guides of the Rust CLI, Go CLI/TUI and Standard TypeScript Monorepo profiles state the clause with `mandatory`, and the Rust guide still dispatches `cut-release`.
- [ ] This repository's skill guide states the clause.
- [ ] A second regeneration and a second Managed Refresh change nothing.

## Context

- instruction: `docs/adr/0222-a-baseline-guide-says-only-what-holds-for-the-repository-that-reads-it.md`
- interface: `docs/agents/setup-context.json`
- interface: `docs/agents/skill-dispatch.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`
- interface: `internal/baseline/assets/modules/core.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/baseline.json`
- interface: `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/agent-instructions.md`
- interface: `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json`
- interface: `internal/baseline/assets/source-baselines/index.json`
- interface: `internal/baseline/clause_characterization_test.go`
- creates: `internal/baseline/person_only_skill_test.go`
- interface: `internal/baseline/preservation_test.go`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`

## Verification

- `out="$(go test -count=1 -v -run '^(TestEveryBuiltInProfileTellsTheAgentToAskForAPersonOnlySkill|TestTheSkillGuidesStateThePersonOnlySkillClause|TestBaselineClauseForceIsCharacterized|TestNoTwoBaselineClausesShareText|TestReadoptionCompatibilityMaintainedFixture|TestCatalogCompatibility|TestBaselinePlanCharacterization)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestEveryBuiltInProfileTellsTheAgentToAskForAPersonOnlySkill TestTheSkillGuidesStateThePersonOnlySkillClause TestReadoptionCompatibilityMaintainedFixture; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/agents/skill-dispatch.md | grep -qF -- "ask the person to run it instead of activating it or carrying out its workflow yourself." || { printf 'missing person-only clause in docs/agents/skill-dispatch.md\n' >&2; exit 1; }; go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task the two new tests do not exist and this repository's skill guide lacks the clause, so the command fails.

## References

- `_prd.md` → Goal 3; Core Feature 3; Success Metric 3; Success Metric 4; Declared breaks; Acceptance evidence
- `_techspec.md` → Exact texts; Source Baseline rows; Version changes; Existing tests that change; API Contract 3; Testing Approach 3; Build Order 3
- ADR-0058, ADR-0060, ADR-0186, ADR-0202, ADR-0222


## Result

Added the mandatory person-only skill clause as the final clause of
`rule.core.skill-dispatch`, with the exact TechSpec text. Raised core 17 → 18,
the rule 5 → 6 and `guide.skill-dispatch` 4 → 5. Added the Source Baseline
corpus row, manifest row and index identity after the vendored-skills clause,
and updated the force characterization and maintained entry count 162 → 163.
The sanctioned generator filled all offsets, digests, identity records and
expected outputs; none were hand-edited. Managed Refresh updated this
repository's skill guide and Setup Manifest.

Acceptance evidence from focused implementation checks:

| Criterion | Evidence |
| --- | --- |
| Profile guides state the mandatory clause; Rust keeps `cut-release` | `GOCACHE=/private/tmp/roundfix-task03-go-cache rtk proxy go test ./internal/baseline -run 'TestEveryBuiltInProfileTellsTheAgentToAskForAPersonOnlySkill\|TestTheSkillGuidesStateThePersonOnlySkillClause' -count=1 -v` exited 0. Both top-level tests passed: all four built-in profiles select core and carry the mandatory clause; Rust CLI, Go CLI/TUI and Standard TypeScript Monorepo Plan postimages state the exact guidance. The Rust subtest also asserts the skill name and `trigger.rust.cut-release` remain rendered. Before editing assets, both tests failed on the missing clause. |
| This repository's guide states the clause | The first successful Managed Refresh reported `File changes: 2` for `docs/agents/skill-dispatch.md` and `docs/agents/setup-context.json`. Inspection of the generated guide confirms the exact text with `mandatory`. |
| Regeneration and refresh converge | Two `rtk proxy make baseline-digests` runs exited 0: first `"changed":true`, second `"changed":false`. Two successful `GOCACHE=/private/tmp/roundfix-task03-go-cache rtk proxy go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text` runs exited 0: the second reported `File changes: 0` and `Idempotence: verified`. |

The default-cache Managed Refresh attempt was blocked by sandbox cache access;
the task-local-cache attempt then reached the Git transaction lock and was
blocked by its read-only sandbox path. The two successful refreshes used
approved elevated execution for that lock. Skills were skipped as requested.
Direct source comparison confirms each version rose exactly once, the clause
is last, and the core's existing skill declarations are unchanged.

The Task's declared Verification was not run; it and Task settlement remain
Daemon-owned. Task status, the Task Graph, other Tasks, skills, dispatch
triggers and Skill Activations were left untouched.

Additional checks:

- Changed-path postflight found 20 changed paths, all within this Task's
  Context or its Task file; `rtk proxy git -c core.fsmonitor=false diff --check`
  exited 0.
- The first `GOCACHE=/private/tmp/roundfix-task03-go-cache rtk proxy make
  verify-incremental` exited 2. CLI force-stop tests were blocked by sandbox
  process-table access. The suite guard also detected this Result being edited
  while tests ran, so that invocation is not passing evidence. A fresh run
  uses approved elevated execution and leaves repository files unchanged for
  its entire duration.

- The fresh elevated `GOCACHE=/private/tmp/roundfix-task03-go-cache rtk proxy
  make verify-incremental` exited 0: formatting, vet, all Go test packages,
  skill synchronization/checks and the CLI build passed. The repository stayed
  unchanged throughout that check. Output is retained locally at
  `/private/tmp/roundfix-task03-incremental-elevated.log`.
- Final changed-path postflight still contains exactly the 20 authorized Task
  paths; no follow-up work was added to this diff. No implementation blocker
  remains for the Daemon's declared Verification.
