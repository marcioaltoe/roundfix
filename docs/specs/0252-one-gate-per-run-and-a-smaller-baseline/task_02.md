---
task: task_02
spec: 0252-one-gate-per-run-and-a-smaller-baseline
status: pending
type: backend
complexity: medium
---

# Task 02: The recorded Verification and runtime values are the ones CI and the Project Config use

## Overview

This repository records `rtk make verify` and `rtk make verify-incremental`,
but `rtk` is a local output filter and CI runs `make verify` and
`make verify-changed`. The autonomous-work guide names `codex gpt-5.6-sol` and
`claude opus 5 xhigh`, which no Agent Selection Profile uses (findings B02 and
B03 of the Baseline audit of 2026-10-08,
`docs/references/2026-10-08-baseline-audit.md`). This Task changes the four
catalog values and the runtime header with `_techspec.md` → Data Models and
Template texts. It re-records this repository's four decisions, and it aligns
the `setup-context-driven` skill and the user guide (ADR-0257). Both runtime
decisions stay. It is verifiable on its own: the Setup Manifest, the guides and
the catalog name no `rtk` command and no stale model.

This is an authorized tooling Task. It may change only the files in its
Context, the derived files the sanctioned regeneration rewrites, and this Task
file.

## Requirements

1. MUST answer findings B02 and B03 of the Baseline audit of 2026-10-08 by
   changing, in `internal/baseline/assets/decisions.json`, the four values and
   the two runtime summaries of `_techspec.md` → Data Models. It MUST raise
   each of the four decisions' versions by one. Each decision keeps its `id`,
   `type` and effects. `runtime.backend` and `runtime.design` MUST stay in the
   `autonomous-work` module's `requiredDecisions`.
2. MUST replace `internal/baseline/assets/templates/guides/autonomous-work.md`
   with the whole text of `_techspec.md` → Template texts, keeping its token
   list. It MUST raise the task_02 versions of `_techspec.md` → Version
   changes. It MUST then run
   `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`,
   `make baseline-digests` and
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`.
3. MUST re-record this repository's decisions as `_authorization.md` →
   Recorded decision change says. It MUST run
   `go run -buildvcs=false ./cmd/roundfix baseline plan --profile go-cli-tui --decision preservation.mode=preservation`
   with one `--decision` for every decision `docs/agents/setup-context.json`
   records. It keeps each recorded value except the four new ones of
   `_techspec.md` → Data Models, and passes `--format json` into a temporary
   file outside the repository. It MUST then run
   `go run -buildvcs=false ./cmd/roundfix baseline apply --plan <file> --confirm-plan <digest>`.
   The Result MUST name every file the apply wrote. A following refresh MUST
   report `current`, and MUST NOT hand-edit the manifest or a guide.
4. MUST change the two `setup-context-driven` sentences and the user-guide
   rows and sentence named in `_techspec.md` → Skill and document texts, and
   nothing else in either file. It MUST run `make skills-sync`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   then `make skills-sync` again.
5. MUST create `internal/baseline/runtime_and_verification_decisions_test.go`
   with `TestTheDecisionCatalogProposesNoRtkCommand`,
   `TestTheRuntimeDecisionsPreferTheProjectConfigTuples` and
   `TestTheAutonomousWorkGuideDefersToAgentSelectionProfiles`, as
   `_techspec.md` → Testing Approach describes. The values are literals in the
   test.
6. MUST change the expected values that
   `internal/baseline/incremental_verification_test.go`,
   `internal/cli/baseline_human_test.go` and
   `internal/cli/baseline_incremental_verification_test.go` pin, and nothing
   else in them, to the new suggestion and defaults. Those are declared
   breaks. The cases that test parsing a `rtk`-prefixed command elsewhere stay
   unchanged.
7. MUST prove each new gate can fail. The Result MUST record one sabotage of a
   catalog value, for example restoring `rtk make verify`, and one of the
   template, each with the test that failed, and that the source was restored
   and regenerated.
8. MUST NOT change `.roundfixrc.yml`, the Makefile, any production Go file, the
   user guide's decision-document example, or any module other than
   `autonomous-work`.

## Subtasks

- [ ] Change the four catalog values, the summaries and the template, and raise the versions.
- [ ] Record, regenerate and refresh.
- [ ] Re-record this repository's four decisions through plan and apply.
- [ ] Align `setup-context-driven` and the user guide, and record the skill version.
- [ ] Add the decision tests, update the pinned values and record each sabotage.

## Acceptance Criteria

- [ ] `docs/agents/setup-context.json` records `make verify`,
      `make verify-changed`, `codex gpt-6.1-sol high` and `claude opus high`.
- [ ] `docs/agents/agent-instructions.md` names `make verify` and
      `make verify-changed`.
- [ ] `docs/agents/autonomous-work.md` defers to Agent Selection Profiles and
      names the Light Tier.
- [ ] No catalog default or suggestion, guide, skill or user-guide row names an
      `rtk`-prefixed command.
- [ ] A second refresh of this repository is a no-op.

## Context

- instruction: `docs/adr/0257-inside-a-run-the-daemon-is-the-only-full-gate-and-the-baseline-states-each-rule-once.md`
- instruction: `docs/adr/0238-a-light-tier-runs-low-complexity-tasks-on-an-open-model.md`
- instruction: `docs/references/2026-10-08-baseline-audit.md`
- instruction: `.roundfixrc.yml`
- instruction: `.github/workflows/ci-verify.yml`
- interface: `internal/baseline/assets/decisions.json`
- interface: `internal/baseline/assets/templates/guides/autonomous-work.md`
- interface: `internal/baseline/assets/templates/index.json`
- interface: `internal/baseline/assets/modules/autonomous-work.json`
- interface: `internal/baseline/module-versions.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `internal/baseline/incremental_verification_test.go`
- interface: `internal/cli/baseline_human_test.go`
- interface: `internal/cli/baseline_incremental_verification_test.go`
- interface: `docs/agents/agent-instructions.md`
- interface: `docs/agents/autonomous-work.md`
- interface: `docs/agents/setup-context.json`
- interface: `.agents/skills/setup-context-driven/SKILL.md`
- interface: `skills/setup-context-driven/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/context-driven-development.md`
- creates: `internal/baseline/runtime_and_verification_decisions_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestTheDecisionCatalogProposesNoRtkCommand|TestTheRuntimeDecisionsPreferTheProjectConfigTuples|TestTheAutonomousWorkGuideDefersToAgentSelectionProfiles|TestIncrementalVerificationDecisionIsDeclaredWithoutADefault|TestEveryBaselineModuleVersionIsRecorded|TestCatalogCompatibility|TestFormatterComposition|TestBaselinePlanCharacterization|TestHumanBaselineDecisionDefaults|TestBaselineUpdateNamesTheMissingIncrementalDecision|TestBaselineUpdateAdoptsADeclaredIncrementalSuggestion)$" ./internal/baseline ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheDecisionCatalogProposesNoRtkCommand TestTheRuntimeDecisionsPreferTheProjectConfigTuples TestTheAutonomousWorkGuideDefersToAgentSelectionProfiles TestIncrementalVerificationDecisionIsDeclaredWithoutADefault TestEveryBaselineModuleVersionIsRecorded TestHumanBaselineDecisionDefaults TestBaselineUpdateNamesTheMissingIncrementalDecision TestBaselineUpdateAdoptsADeclaredIncrementalSuggestion; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the three decision tests do not exist, so the command fails.
- `manifest="$(tr -s '[:space:]' ' ' < docs/agents/setup-context.json)"; for phrase in '"verification.gate": { "value": "make verify" }' '"verification.incremental": { "value": "make verify-changed" }' '"runtime.backend": { "value": "codex gpt-6.1-sol high" }' '"runtime.design": { "value": "claude opus high" }'; do printf '%s\n' "$manifest" | grep -qF -- "$phrase" || { printf 'missing decision in docs/agents/setup-context.json: %s\n' "$phrase" >&2; exit 1; }; done; for phrase in "Agent Selection Profiles in Project Config choose each Agent Session's ACP Runtime, model, and reasoning effort" "runs on the Light Tier when it is available"; do tr -s '[:space:]' ' ' < docs/agents/autonomous-work.md | grep -qF -- "$phrase" || { printf 'missing phrase in docs/agents/autonomous-work.md: %s\n' "$phrase" >&2; exit 1; }; done; if grep -qF -e 'gpt-5.6-sol' -e 'opus 5 xhigh' docs/agents/autonomous-work.md; then printf 'the autonomous-work guide still names a stale model\n' >&2; exit 1; fi; if grep -qF 'rtk make' internal/baseline/assets/decisions.json docs/agents/agent-instructions.md docs/user-guide/context-driven-development.md .agents/skills/setup-context-driven/SKILL.md; then printf 'an rtk-prefixed command remains\n' >&2; exit 1; fi; cmp -s .agents/skills/setup-context-driven/SKILL.md skills/setup-context-driven/SKILL.md || { printf 'the setup-context-driven mirror differs\n' >&2; exit 1; }; go test -count=1 -run '^(TestEveryOwnedSkillVersionIsRecorded|TestAuthorialSkillSync)$' ./skills || exit 1; plan="$(go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json)" || { printf '%s\n' "$plan"; exit 1; }; printf '%s\n' "$plan" | grep -qF -- '"state":"current"' || { printf 'the guides are not refreshed: %s\n' "$plan" >&2; exit 1; }` — expected: exit 0; before this Task the Setup Manifest records the rtk-prefixed gate, so the command fails at its first phrase.

## References

- [_prd.md](_prd.md) — Goals 2 and 3; User Stories 2 and 3; Core Features 2 and 3; Success Metrics 2 and 3
- [_techspec.md](_techspec.md) — API Contract 2; API Contract 3; Data Models; Template texts; Skill and document texts; Version changes; Derived files; Testing Approach; Build Order 2
- ADR-0257; ADR-0238; ADR-0250
