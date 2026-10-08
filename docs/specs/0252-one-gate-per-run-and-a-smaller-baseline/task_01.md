---
task: task_01
spec: 0252-one-gate-per-run-and-a-smaller-baseline
status: completed
type: backend
complexity: medium
---

# Task 01: A Daemon-assigned Task runs focused tests and leaves the full gate to the Daemon

## Overview

The core and Spec workflow clauses tell every Task agent to run the selected
incremental Verification before handoff, and this repository's incremental
command is the full suite. The Baseline audit of 2026-10-08 counted 74 such
agent runs (284 agent-minutes) in Specs 0225–0248, beside a Daemon gate that
already ran at every settlement (findings B01 and B15,
`docs/references/2026-10-08-baseline-audit.md`). This Task rewords the four
clauses in place with `_techspec.md` → Exact clause texts, aligns
`implement-task` §7 and revises **Incremental Verification** (ADR-0257). It is
verifiable on its own: this repository's guides scope both tiers, and an
adopter's update retains the four clauses.

This is an authorized tooling Task. It may change only the files in its
Context, the derived files the sanctioned regeneration rewrites, and this Task
file.

## Requirements

1. MUST answer findings B01 and B15 of the Baseline audit of 2026-10-08 by
   replacing the `guidance` of `clause.core.run-selected-verification` and
   `clause.core.verification-two-tiers` in
   `internal/baseline/assets/modules/core.json`, and of
   `clause.spec.verification-two-tiers` and
   `clause.spec.verification-fails-before-the-change` in
   `internal/baseline/assets/modules/spec-workflow.json`, with the exact texts
   of `_techspec.md` → Exact clause texts. It edits the strings in place, so
   every other byte stays as it is. Each clause MUST keep its `id` and its
   `mandatory` enforcement and MUST NOT gain a `replaces` list.
2. MUST raise by one, from the value on this Task's starting commit, the
   task_01 versions of `_techspec.md` → Version changes. It MUST then run
   `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`,
   then `make baseline-digests`, then
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`.
   A second refresh MUST report `File changes: 0`. It MUST NOT hand-edit a
   module version, pin, golden, snapshot or generated guide, and the Result
   MUST name every file the commands rewrote.
3. MUST replace step 3 of `.agents/skills/implement-task/SKILL.md` §7 with the
   text of `_techspec.md` → Skill and document texts. It MUST then run
   `make skills-sync` and
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   and `make skills-sync` again, so that both version fields rise and the
   mirror matches. No other line of the skill and no `### QA settlement`
   section of any skill changes.
4. MUST revise the **Incremental Verification** entry of `CONTEXT.md` through
   `domain-modeling` to the text of `_techspec.md` → Skill and document texts,
   keeping its `_Avoid_` line.
5. MUST create `internal/baseline/verification_tier_clauses_test.go` with
   `TestTheVerificationTierClausesCarryTheirText`,
   `TestTheVerificationTierClausesRenderInTheGuides` and
   `TestAnAdopterRetainsTheVerificationTierClauses`, following
   `internal/baseline/glossary_clauses_test.go`. The four texts are literals in
   the test. The rendering test MUST require each one exactly once as a
   `mandatory` bullet in the Standard TypeScript golden guide that renders it
   (`agent-instructions.md` or `spec-routing.md`). The retention test MUST
   require a ready Managed Refresh of a Source Baseline adopter that records
   all four `retained` and no clause `unaccounted`.
6. MUST change the sentence that
   `internal/baseline/promoted_spec_and_typescript_clauses_test.go` pins for
   `clause.spec.verification-fails-before-the-change` to the new text and
   nothing else in that file. That is a declared break.
7. MUST prove each new gate can fail. The Result MUST record one sabotage of a
   clause text and one of the retention, for example a changed enforcement,
   with the test that failed for each, and that the source was restored and
   regenerated.
8. MUST NOT change any other clause, the Source Baseline assets, the retention
   transition, the force record or any production Go file.

## Subtasks

- [ ] Reword the four clauses and raise the versions.
- [ ] Record, regenerate and refresh twice.
- [ ] Change `implement-task` §7, record its version and sync the mirror.
- [ ] Revise **Incremental Verification** in `CONTEXT.md`.
- [ ] Add the clause test file, update the pinned sentence and record each sabotage.

## Acceptance Criteria

- [ ] `docs/agents/agent-instructions.md` and `docs/agents/spec-routing.md`
      scope both tiers to work outside a Run and name focused tests for a
      Daemon-assigned turn.
- [ ] `docs/agents/spec-routing.md` names `--strict --run-verification`.
- [ ] `implement-task` §7 forbids the selected commands and the full suite, and
      its mirror matches.
- [ ] `CONTEXT.md` carries the revised **Incremental Verification**.
- [ ] A Source Baseline adopter's refresh retains the four clauses, and a
      second refresh of this repository is a no-op.

## Context

- instruction: `docs/adr/0257-inside-a-run-the-daemon-is-the-only-full-gate-and-the-baseline-states-each-rule-once.md`
- instruction: `docs/adr/0250-a-module-version-is-chosen-when-recorded-and-the-coverage-record-lists-every-platform.md`
- instruction: `docs/references/2026-10-08-baseline-audit.md`
- instruction: `internal/baseline/glossary_clauses_test.go`
- interface: `internal/baseline/assets/modules/core.json`
- interface: `internal/baseline/assets/modules/spec-workflow.json`
- interface: `internal/baseline/module-versions.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `internal/baseline/promoted_spec_and_typescript_clauses_test.go`
- interface: `docs/agents/agent-instructions.md`
- interface: `docs/agents/spec-routing.md`
- interface: `docs/agents/setup-context.json`
- interface: `.agents/skills/implement-task/SKILL.md`
- interface: `skills/implement-task/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `CONTEXT.md`
- creates: `internal/baseline/verification_tier_clauses_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestTheVerificationTierClausesCarryTheirText|TestTheVerificationTierClausesRenderInTheGuides|TestAnAdopterRetainsTheVerificationTierClauses|TestTwoTierClausesNameTheSelectedCommands|TestTheSpecAndTypeScriptGuidesStateThePromotedRules|TestBaselineClauseForceIsCharacterized|TestNoTwoBaselineClausesShareText|TestEveryBaselineModuleVersionIsRecorded|TestCatalogCompatibility|TestFormatterComposition|TestBaselinePlanCharacterization)$" ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheVerificationTierClausesCarryTheirText TestTheVerificationTierClausesRenderInTheGuides TestAnAdopterRetainsTheVerificationTierClauses TestTwoTierClausesNameTheSelectedCommands TestTheSpecAndTypeScriptGuidesStateThePromotedRules TestBaselineClauseForceIsCharacterized TestNoTwoBaselineClausesShareText TestEveryBaselineModuleVersionIsRecorded TestCatalogCompatibility; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the three verification-tier tests do not exist, so the command fails.
- `for phrase in "Outside a Run, run the selected repository Verification before a completion claim" "In a Daemon-assigned turn, run focused tests of the changed packages and neither selected command"; do tr -s '[:space:]' ' ' < docs/agents/agent-instructions.md | grep -qF -- "$phrase" || { printf 'missing phrase in docs/agents/agent-instructions.md: %s\n' "$phrase" >&2; exit 1; }; done; for phrase in "a Daemon-assigned Task hands back after focused tests" "--strict --run-verification"; do tr -s '[:space:]' ' ' < docs/agents/spec-routing.md | grep -qF -- "$phrase" || { printf 'missing phrase in docs/agents/spec-routing.md: %s\n' "$phrase" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < .agents/skills/implement-task/SKILL.md | grep -qF -- "its incremental Verification, or the full test suite" || { printf 'implement-task lacks the handoff rule\n' >&2; exit 1; }; cmp -s .agents/skills/implement-task/SKILL.md skills/implement-task/SKILL.md || { printf 'the implement-task mirror differs\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "for validating the current change outside a Run while reusing safe local state" || { printf 'CONTEXT.md lacks the revised **Incremental Verification**\n' >&2; exit 1; }; go test -count=1 -run '^(TestEveryOwnedSkillVersionIsRecorded|TestAuthorialSkillSync)$' ./skills || exit 1; plan="$(go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json)" || { printf '%s\n' "$plan"; exit 1; }; printf '%s\n' "$plan" | grep -qF -- '"state":"current"' || { printf 'the guides are not refreshed: %s\n' "$plan" >&2; exit 1; }` — expected: exit 0; before this Task the agent-instructions guide asks every Task for the incremental tier, so the command fails at its first phrase.

## References

- [_prd.md](_prd.md) — Goal 1; User Story 1; Core Features 1 and 6; Success Metric 1
- [_techspec.md](_techspec.md) — API Contract 1; API Contract 5; Exact clause texts; Skill and document texts; Version changes; Retention; Derived files; Testing Approach; Build Order 1
- ADR-0257; ADR-0250; ADR-0244; ADR-0186


## Result

### Implementation handed back

Reworded the four Verification clauses in place with the exact TechSpec texts,
answering audit findings B01 and B15. Their identities, `mandatory`
enforcement, and absence of `replaces` are preserved. No other clause changed.
The rule and guide versions rose once from the starting commit:

- `rule.core.verification-selected`: 7 → 8;
  `guide.agent-instructions`: 13 → 14.
- `rule.spec.routing`: 7 → 8; `guide.spec-routing`: 11 → 12.
- The module recorder chose `core`: 18 → 19 and `spec-workflow`: 15 → 16.
- The skill recorder raised both `implement-task` version fields: 0.0.3 → 0.0.4.

Changed only §7 step 3 and the recorded version fields of `implement-task`.
Revised **Incremental Verification** through `domain-modeling` to the exact
TechSpec definition and preserved its `_Avoid_` line. Added the three required
clause tests using the existing glossary-test pattern, with independent
literal texts, forced golden bullets, and real temporary-adopter refreshes.
Updated only the declared pinned sentence in the promoted-clause test.

### Focused evidence by acceptance criterion

| Acceptance criterion | Evidence from this Agent turn |
| --- | --- |
| Both generated guides scope the tiers outside a Run and name focused tests inside one | An exact-byte Python inspection required each of the four TechSpec texts exactly once as a mandatory bullet in its repository guide. `TestTheVerificationTierClausesCarryTheirText` and `TestTheVerificationTierClausesRenderInTheGuides` passed after sabotage restoration and regeneration. |
| Spec routing names `--strict --run-verification` | Exact clause inspection passed; `TestTheSpecAndTypeScriptGuidesStateThePromotedRules` passed for both Standard TypeScript and Go profiles. |
| Skill §7 forbids the selected commands and full suite; mirror matches | Exact-byte inspection proved only step 3 plus the two recorded version fields changed, and canonical/mirror bytes match. `TestAuthorialSkillSync` passed. |
| Glossary carries the revised Incremental Verification | Exact-byte comparison to HEAD with only the authored definition substituted passed, including preservation of `_Avoid_`. |
| Source Baseline adopter retains all four clauses; repository refresh is a no-op | `TestAnAdopterRetainsTheVerificationTierClauses` passed, requiring `ready`, all four `retained` dispositions and retention evidence, and no `unaccounted` disposition. A second repository refresh, and a final refresh after restoration, both reported `File changes: 0` and verified idempotence. |

### Commands and outcomes

Go and Make commands used
`GOCACHE=/private/tmp/roundfix-0252-task01-go-cache`; command output was
preserved with `rtk proxy`.

- Starting focused check:
  `go test ./internal/baseline -run '^TestTheVerificationTierClauses(CarryTheirText|RenderInTheGuides)$' -count=1`
  exited 1 on all four original texts and their golden bullets.
- `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`
  exited 0 and recorded the two module versions.
- `make skills-sync`, then
  `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
  then `make skills-sync` each exited 0. The recorder had sandbox permission
  to write the protected canonical skill.
- `make baseline-digests` exited 0 after skill sync, and again after sabotage
  restoration. Its final strict catalog validation passed. No generated
  expectation or guide was edited by hand.
- `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`
  applied and verified three generated files. Its approved Plan Digest was
  `sha256:298dc0841edcc2c0a28b3fb36e98b997f0a1b18bd9c26765c3ed2f000e14d7a3`.
  The second and final invocations exited 0 with `File changes: 0`; their
  Plan Digest was
  `sha256:b50dbef7060a80c7eee0f69411d7b9f33195ce7107f07bb91cc629c27c1e8136`.
  These local refreshes reported zero semantic retention evidence; the
  temporary Source Baseline adopter test supplies that criterion's evidence.
- Final focused baseline check:
  `go test ./internal/baseline -run '^(TestTheVerificationTierClausesCarryTheirText|TestTheVerificationTierClausesRenderInTheGuides|TestAnAdopterRetainsTheVerificationTierClauses|TestTheSpecAndTypeScriptGuidesStateThePromotedRules)$' -count=1 -v`
  exited 0, with all four tests and their subtests passing.
- `go test ./skills -run '^TestAuthorialSkillSync$' -count=1` exited 0.
- Exact-byte source/guide/skill/glossary inspection passed. Changed-file
  postflight and whitespace inspection are recorded below.

### Sabotage evidence and restoration

1. Temporarily changed `clause.core.run-selected-verification` from
   “Outside a Run” to “Outside this Run” in `core.json`, then regenerated its
   derived wording. The focused three-test command exited 1:
   `TestTheVerificationTierClausesCarryTheirText` reported `force/text differs`,
   and `TestTheVerificationTierClausesRenderInTheGuides` reported
   `guide lacks exactly one forced clause`, both for that identity.
2. Temporarily injected `mutateCatalogClause` after the retention test's
   adopter creation, changing that clause's in-memory enforcement to
   `stop-and-ask`. `TestAnAdopterRetainsTheVerificationTierClauses` failed:
   the refresh was `action_required`, category `classification`, with
   `retention transition has 1 unaccounted clause(s): clause.core.run-selected-verification`.

Restored the module and new test file byte-for-byte from their pre-sabotage
copies. No sabotage was recorded as a module or skill version. Ran
`make baseline-digests` successfully after restoration, then reran the focused
checks successfully and confirmed the repository refresh still changes no
file. The test contains no sabotage injection.

### Files rewritten by sanctioned commands

The module-version recorder changed:

- `internal/baseline/assets/modules/core.json` (module version line)
- `internal/baseline/assets/modules/spec-workflow.json` (module version line)
- `internal/baseline/module-versions.json`

`make baseline-digests` changed these derived outputs across the successful
regeneration runs:

- `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`
- `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`
- `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- `internal/baseline/testdata/catalog.diagnostics.golden.json`
- `internal/baseline/testdata/catalog.digest`
- `internal/baseline/testdata/catalog.normalized.json`
- `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`

The repository refresh changed:

- `docs/agents/agent-instructions.md`
- `docs/agents/spec-routing.md`
- `docs/agents/setup-context.json`

The skill-version recorder changed:

- `.agents/skills/implement-task/SKILL.md` (both version fields)
- `skills/implement-task/SKILL.md` (both version fields)
- `skills/testdata/owned-skill-versions.json`

`make skills-sync` recopied the Makefile's owned-skill trees; only
`skills/implement-task/SKILL.md` has changed bytes. All other owned skill
mirrors remain identical to their starting bytes.

### Execution limits and recovered interruptions

- The first Go attempt could not read the shared Go build cache; the
  task-scoped cache resolved it.
- The first digest attempt refused the unsynced skill mirror; syncing and
  recording the skill resolved it before successful regeneration.
- The first Baseline apply could not create its Git-private transaction
  directory. The authorized retry with sandbox permission applied the same
  plan successfully.
- During sabotage setup, editing the new test while digest regeneration was
  still running triggered suiteguard's repository-boundary guard. That run
  was not counted as successful regeneration. After both sabotage edits were
  restored, sequential regeneration and focused tests passed.
- Refresh warnings identify the existing nested golden and Source Baseline
  instruction carriers, which the command leaves unchanged.
- The authored `## Verification` commands, selected repository Verification,
  selected incremental Verification, and full test suite were not run.
  Task status and settlement remain Daemon-owned; no commit, push, PR,
  Task Graph edit, or other Task edit was made.

Changed-file postflight: all 23 changed paths are declared in task_01 Context
or are the assigned Task file. `git diff --check` passed. No Source Baseline
asset, retention transition, force record, production Go file, other Task, or
Task Graph changed. The pre-existing Daemon-written `status: in_progress`
remains unchanged.

## Carry-forward provenance

- Source Run: `run_20261008T180729Z_1c5cbdf685b778d0`
- Source commit: `03c5992b0b6021250380fba480e9a44aa22c6421`
