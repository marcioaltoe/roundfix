---
task: task_02
spec: 0252-one-gate-per-run-and-a-smaller-baseline
status: completed
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


## Result

Implemented the Task 02 slice for Daemon Verification. Task status remains
Daemon-owned; no declared Verification command, selected repository or
incremental Verification, full suite, commit, push, or PR was run.

### Implementation

- The gate default is `make verify`; the incremental suggestion is
  `make verify-incremental`. Runtime defaults are `codex gpt-6.1-sol high`
  and `claude opus high`, with the two authored preference summaries.
  Each decision rose from version 1 to 2; IDs, types and effects are unchanged.
  Both runtime decisions remain required by `autonomous-work`.
- The autonomous-work template now contains the whole authored text, with
  its original token list. Its template version rose from 1 to 2 and its
  supporting-guide version from 12 to 13. The Module Version Record step
  chose module version 13 and recorded it.
- Updated only the two specified skill sentences and five specified user-guide
  substitutions. The skill recorder raised both version fields from 0.0.3
  to 0.0.4. The explicit decision-document example is byte-identical to HEAD.
- Added the three authored decision/guide tests with literal expectations.
  Updated only pinned suggestions, defaults and the two runtime prompt-summary
  expectations in the three existing test files. RTK parser cases elsewhere
  were not changed.

### Regeneration and recorded decisions

Commands below used `GOCACHE=/private/tmp/roundfix-task02-gocache` for Go
builds and raw output through `rtk proxy`.

1. `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`
   exited 0.
2. `make skills-sync`,
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   then `make skills-sync` exited 0. The recorder required sandbox escalation
   to write the protected canonical `.agents` skill.
3. `make baseline-digests` exited 0 and regenerated the declared catalog,
   profile, formatter and plan-characterization artifacts. An earlier attempt
   stopped on skill-mirror drift before the synchronization above; no assertion
   was weakened. Regeneration succeeded after synchronization.
4. `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`
   exited 0 and verified the refresh. It changed
   `docs/agents/autonomous-work.md` and `docs/agents/setup-context.json`.
   Sandbox escalation was required for the Git-private transaction lock.
5. Planned with `go run -buildvcs=false ./cmd/roundfix baseline plan --profile go-cli-tui --decision preservation.mode=preservation`,
   one `--decision` for all twelve recorded decisions, and `--format json`.
   Every value was preserved except the four authored changes. The plan was
   saved outside the repository at `/private/tmp/roundfix-task02-plan.json`.
   Reviewed digest:
   `sha256:fce89c15830bf52dfc3a70c280c8983f97a2615aec71a494c3f339ecfc667725`.
6. `go run -buildvcs=false ./cmd/roundfix baseline apply --plan /private/tmp/roundfix-task02-plan.json --confirm-plan sha256:fce89c15830bf52dfc3a70c280c8983f97a2615aec71a494c3f339ecfc667725`
   exited 0, verified fourteen postimages and wrote exactly these three changed
   files: `docs/agents/agent-instructions.md`,
   `docs/agents/autonomous-work.md`, `docs/agents/setup-context.json`.
   The other eleven postimages already matched and were not replaced.
   No manifest or generated guide was hand-edited.

### Focused checks and acceptance evidence

- `go test ./internal/baseline -run '^TestThe(DecisionCatalog|RuntimeDecisions|AutonomousWorkGuide)' -count=1`
  initially exited 1 with all three tests failing on the original commands,
  runtimes and guide header.
- After restoration and final regeneration,
  `go test ./internal/baseline -run '^(TestTheDecisionCatalogProposesNoRtkCommand|TestTheRuntimeDecisionsPreferTheProjectConfigTuples|TestTheAutonomousWorkGuideDefersToAgentSelectionProfiles|TestIncrementalVerificationDecisionIsDeclaredWithoutADefault)$' -count=1 -v`
  exited 0 and named all four passing tests.
- `go test ./internal/cli -run '^(TestHumanBaselineDecisionDefaults|TestHumanBaselineFirstAdoptionPromptSequenceCharacterization|TestBaselineUpdateNamesTheMissingIncrementalDecision|TestBaselineUpdateAdoptsADeclaredIncrementalSuggestion)$' -count=1`
  exited 0.
- Acceptance criterion 1: parsed the generated manifest and confirmed exactly
  `make verify`, `make verify-changed`, `codex gpt-6.1-sol high`, and
  `claude opus high`. Compared all eight other decision values with HEAD;
  every value is unchanged.
- Acceptance criterion 2: inspected the generated agent-instructions header;
  it names `make verify` and `make verify-changed`.
- Acceptance criterion 3: normalized and inspected the generated autonomous-work
  header; it carries both authored profile-selection and Light Tier phrases.
  The rendered golden test also passes.
- Acceptance criterion 4: inspected catalog values, every repository guide,
  both setup-context-driven skill copies, and the user guide; none contains
  `rtk make`. The catalog test rejects any string default or suggestion
  beginning with `rtk `. The skill copies are byte-identical.
- Acceptance criterion 5: ran
  `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format text`
  twice after restoration. Both exited 0, reported `current` and zero file
  changes, and preserved every tracked file's bytes. Logs are
  `/private/tmp/roundfix-task02-refresh-1.log` and
  `/private/tmp/roundfix-task02-refresh-2.log`.
- Final scope inspection found 26 changed paths, all declared by Task 02,
  including the pre-existing Daemon status change in this Task file.
  `git diff --check` exited 0. No other module, task file, Task Graph,
  production Go file, Project Config, Makefile or CI workflow changed.

### Sabotage evidence

Each test below exited 1 and named the expected failed test. Each source was
restored; `make baseline-digests` then regenerated the final artifacts and
exited 0. All three new tests subsequently passed together.

- Catalog: temporarily restored the gate default to `rtk make verify`.
  `go test ./internal/baseline -run '^TestTheDecisionCatalogProposesNoRtkCommand$' -count=1`
  failed on the local-filter prefix and wrong literal default. Evidence:
  `/private/tmp/roundfix-task02-sabotage-catalog.log`.
- Runtime: temporarily restored the backend default to `codex gpt-5.6-sol`.
  `go test ./internal/baseline -run '^TestTheRuntimeDecisionsPreferTheProjectConfigTuples$' -count=1`
  failed on the backend tuple. Evidence:
  `/private/tmp/roundfix-task02-sabotage-runtime.log`.
- Template: temporarily replaced `runs on the Light Tier` with
  `runs on the Standard Tier` in the template and ran `make baseline-digests`
  to render that regression.
  `go test ./internal/baseline -run '^TestTheAutonomousWorkGuideDefersToAgentSelectionProfiles$' -count=1`
  failed on the missing Light Tier phrase. Evidence:
  `/private/tmp/roundfix-task02-sabotage-template.log`.

Baseline retained its existing two nested-carrier warnings for the formatter
fixture and Source Baseline corpus. No blocking implementation issue remains;
terminal Verification and settlement belong to the Daemon.


### Verification Feedback repair — attempt 1

Inspected both Daemon diagnostic artifacts. The changed skill made the
TechSpec's before/after description look like an unproven receipt of its
current contents (`SC-RECEIPT-UNPROVEN`). The changed user guide also invalidated
its generated Behavior Surface fingerprint. These caused the corpus and
skill-coverage checks to fail in the Daemon's attempt.

- Reworded only the Task 02 skill-change paragraph in `_techspec.md` as a
  change description, preserving both old and new values. It no longer
  presents the old sentence as a receipt of the current skill.
- Regenerated `docs/references/behavior-surfaces.json` with Project Config's
  declared command:
  `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
  It exited 0 and changed only the context-driven-development guide's
  fingerprint. No generated digest was hand-edited, and no corpus golden,
  test assertion, production code or configuration changed.
- Focused repair check:
  `go test -count=1 -tags docscontract ./internal/docscontract -run '^(TestCheckCorpusGolden|TestCheckActiveCorpusHasNoErrors|TestTheSkillCoverageMapIsCurrent)$'`
  exited 0. The corpus checks exercise the real active Spec through
  `speccheck.Check`, including receipt consistency, and the coverage check
  validates the regenerated record without its record flag.
- These checks used `GOCACHE=/private/tmp/roundfix-task02-gocache` and raw
  output through `rtk proxy`. `git diff --check` exited 0 after the repair.
  The repair adds only the Task 02 paragraph correction in the bundled
  TechSpec and the sanctioned derived Behavior Surface record to the prior
  diff. Task status, the Task's authored requirements and Verification,
  other task files, and the Task Graph remain unchanged.

The full `make verify-changed` and settlement checks were not rerun by the
Agent. Their next attempt and terminal verdict remain Daemon-owned.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `docs/references/behavior-surfaces.json`
- `docs/specs/0252-one-gate-per-run-and-a-smaller-baseline/_techspec.md`

## Carry-forward provenance

- Source Run: `run_20261008T180729Z_1c5cbdf685b778d0`
- Source commit: `7521a8403846ffc489d02fc8fc0c3e999b1f99fb`
