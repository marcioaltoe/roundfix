---
task: task_02
spec: 0223-adjustments-the-adopters-asked-for
status: pending
type: backend
complexity: medium
---

# Task 02: The branch prefix becomes an optional Baseline decision

## Overview

Make `branch.prefix` an optional Baseline decision with the machinery Spec
0207 built for `frontend.layout`, as the TechSpec's "The optional branch
prefix" states: the core module stops requiring it, the agent-instructions
template's prefix sentence moves into the decision's renderer, a recorded
value renders today's bytes and an unrecorded one renders the commit-type rule.
This answers the Backlog Entry of 2026-10-03, "A repository cannot leave the
branch prefix to the commit-type rule", from the oraculum maintainer's request
of 2026-10-01.

## Requirements

1. MUST add `"optional": true` to `branch.prefix` in
   `internal/baseline/assets/decisions.json` and keep its `version`, `type`,
   `default`, `summary` and render binding byte for byte.
2. MUST remove `branch.prefix` from the core module's `requiredDecisions` in
   `internal/baseline/assets/modules/core.json` and keep it in the decisions of
   every built-in profile.
3. MUST replace, in `internal/baseline/assets/templates/guides/agent-instructions.md`,
   the five lines from "The branch-prefix pattern is" to "documented
   namespace." with the single line `{{branch.prefix}}`, without raising the
   template's version.
4. MUST add a `branch.prefix` case to `renderProjectDecision` returning Fixed
   text 1 for the recorded value, and to `renderUnrecordedProjectDecision`
   returning Fixed text 2, in `internal/baseline/project_decision_render.go`.
   Fixed text 1 MUST equal, for `<type>/` and `ma/`, the bytes the template
   rendered before this Task.
5. MUST regenerate the catalog snapshots and plan goldens with
   `make baseline-digests`, then refresh this repository's Setup Manifest with
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`;
   a second refresh MUST report no file change, and
   `docs/agents/agent-instructions.md` MUST keep its bytes. No derived file is
   hand-edited.
6. MUST change the case `decisions-absent-names-every-required-decision` of
   `TestBaselinePlanAdoptionAndDecisionCharacterizationCorpus` in
   `internal/cli/baseline_plan_test.go` so that a decision the embedded catalog
   declares optional is asserted not to be named, and every other decision of
   the profile is still asserted named. No other case of that file changes.
7. MUST add the tests of the TechSpec's Testing Approach 1 and 2 in the two new
   files, and MUST NOT edit `internal/cli/baseline_human_test.go`,
   `internal/baseline/plan_test.go` or any other existing test.
8. MUST NOT remove or rename any Baseline Normative Clause, change the human
   first-adoption questions, or edit `CONTEXT.md`.

## Subtasks

- [ ] Declare the decision optional and drop it from the core module's required list.
- [ ] Move the sentence into the renderer with both fixed texts.
- [ ] Regenerate the derived artifacts and refresh the Setup Manifest.
- [ ] Update the characterization case and add the five tests.

## Acceptance Criteria

- [ ] An unrecorded `branch.prefix` is never reported missing and renders
      Fixed text 2.
- [ ] `<type>/` and `ma/` render the bytes they rendered before.
- [ ] `baseline update` on a manifest without the decision asks nothing.
- [ ] This repository's agent-instructions guide is byte-identical.

## Context

- interface: `internal/baseline/assets/decisions.json`
- interface: `internal/baseline/assets/modules/core.json`
- interface: `internal/baseline/assets/templates/guides/agent-instructions.md`
- interface: `internal/baseline/project_decision_render.go`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/setup-context.json`
- interface: `internal/cli/baseline_plan_test.go`
- creates: `internal/baseline/branch_prefix_optional_test.go`
- creates: `internal/cli/baseline_branch_prefix_test.go`
- instruction: `docs/adr/0228-an-unrecorded-branch-prefix-follows-the-commit-type-rule.md`
- instruction: `docs/adr/0205-a-repository-records-its-frontend-layout-and-an-unrecorded-layout-follows-the-suggestion.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestAnUnrecordedBranchPrefixStatesTheCommitTypeRule|TestARecordedBranchPrefixKeepsItsSentence|TestBranchPrefixIsAnOptionalDecision|TestBaselineUpdateWithoutABranchPrefixAsksNothing|TestBaselineUpdateKeepsARecordedBranchPrefix|TestBaselinePlanAdoptionAndDecisionCharacterizationCorpus)$" ./internal/baseline ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestAnUnrecordedBranchPrefixStatesTheCommitTypeRule TestARecordedBranchPrefixKeepsItsSentence TestBranchPrefixIsAnOptionalDecision TestBaselineUpdateWithoutABranchPrefixAsksNothing TestBaselineUpdateKeepsARecordedBranchPrefix TestBaselinePlanAdoptionAndDecisionCharacterizationCorpus; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the five new tests do not exist, so the command fails; after it they pass and the characterization corpus passes with the changed case.

## References

- `_prd.md` → Core Features 1-3; User Stories 1-2; Success Metrics 1-2
- `_techspec.md` → The optional branch prefix; API Contract 4; Testing Approach 1; Testing Approach 2; Build Order 2
- ADR-0205; ADR-0228
