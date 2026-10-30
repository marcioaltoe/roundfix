---
task: task_01
spec: 0335-follow-up-takes-every-documented-type-and-an-absent-spec-root-is-empty
status: pending
type: backend
complexity: high
---

# Task 01: An absent Spec Root inside the repository reads as an empty one

## Overview

A repository whose last Spec was archived has no Spec Root in a fresh
checkout, because Git keeps no empty directory, and every command that only
enumerates Specs refuses there. This Task adds the absent value to the Spec
Root resolver, switches the enumerating readers to it, keeps every named-Spec
refusal byte-identical, and documents the rule in the configuration guide,
the four Roundfix skill references that describe those commands and the
glossary. It answers the second friction of the adopted Backlog Entry of
2026-10-30 (the follow-up command rejects documented message types and an
empty Spec Root).

## Requirements

1. MUST characterize first (lesson 9): write
   `TestAnAbsentSpecRootReadsLikeAnEmptyOne` in the new
   `internal/cli/absent_spec_root_test.go` and observe on the current tree
   that its `follow-up` subtest reports the `specs` source `unavailable` and
   that `spec check`, `deliver plan`, `deliver status` and `history
   sanitize` exit 2 with `which does not exist`; then implement.
2. MUST implement `_techspec.md` → API Contract 1 in
   `internal/config/config.go`: `SpecsRoot.Absent`,
   `ResolveSpecsRootOrAbsent` and `SpecsRoot.RequirePresent`, with
   `ResolveSpecsRoot` returning exactly the value or error bytes it returns
   today (Invariant 1). Absence needs `os.Lstat` to report not-exist at the
   resolved path and the nearest existing ancestor, after
   `filepath.EvalSymlinks`, to lie inside or at the evaluated repository
   root; a dangling link, a file, a stat error and a missing root outside the
   repository stay refused with today's text (Invariant 3).
3. MUST implement API Contract 2: `readFollowUpSpecs`, `runSpecCheckCommand`
   and `runDeliverPlan` without slugs, `loadDeliveryCommand`,
   `runHistoryCommand`, `printHistoryPlan` and `planBaselineUpdateHistory`
   call `ResolveSpecsRootOrAbsent`, and each produces for an absent root what
   it produces for an empty one (Invariant 2).
4. MUST implement API Contract 3: `runSpecCheckCommand` and `runDeliverPlan`
   with slugs, `runDeliverStart` and `runDeliverRetry` call `RequirePresent`
   before reading a Spec, so they refuse an absent root with today's text and
   exit code. Every other caller listed in `_techspec.md` → Current behavior
   2 keeps `ResolveSpecsRoot`. `archive` and `cancel` are not changed
   (Invariant 4).
5. MUST add the tests `_techspec.md` → Testing Approach names for task_01,
   each failing when its behavior is removed:
   - `internal/config/specs_root_absent_test.go`:
     `TestResolveSpecsRootOrAbsentReadsAnAbsentRootInsideTheRepository`,
     `TestResolveSpecsRootOrAbsentStillRefusesAMissingExternalRoot`,
     `TestResolveSpecsRootOrAbsentStillRefusesADanglingLink`,
     `TestResolveSpecsRootOrAbsentStillRefusesAFile`.
   - `internal/cli/absent_spec_root_test.go`:
     `TestAnAbsentSpecRootReadsLikeAnEmptyOne` (subtests `follow-up`,
     `follow-up --status`, `spec check`, `deliver plan`, `deliver status`,
     `history sanitize`; each runs in one temporary repository with no
     `docs/specs`, creates the directory empty, runs again, and compares
     stdout, stderr and exit code; the first run never prints `which does not
     exist`), `TestBaselineUpdateHistoryPlansAnAbsentSpecRootLikeAnEmptyOne`
     and `TestANamedSpecInAnAbsentSpecRootKeepsTheRefusal`.
   The follow-up subtests may start from `newFollowUpFixture` and remove its
   empty `docs/specs`; `deliver status` needs a Run Database, created with
   `store.Open` on the test home. Each top-level test calls `t.Parallel()`
   first or carries a `// Sequential: <reason>` comment; no test uses a
   real-time budget or the real `~/.roundfix`.
6. MUST add to `docs/user-guide/configuration.md`, section "Spec Root", and
   to the canonical `.agents/skills/roundfix/references/spec.md` (Spec
   Consistency Check), `deliver.md` (`deliver plan`), `archive.md` (History
   Sanitize Command) and `baseline.md` (the history plan) the sentence "An
   absent Spec Root inside the repository reads as an empty one", with the
   rule of ADR-0348: what counts as absent, which commands read it as empty,
   and that a named Spec, a dangling link or a missing external root is
   still refused. New text goes under existing headings, never inside a
   `### QA settlement` section.
7. MUST, through `domain-modeling`, revise **Spec Root** in `CONTEXT.md`:
   an absent Spec Root inside the repository holds no Spec, citing ADR-0348
   and never a Spec.
8. MUST copy the canonical skill to `skills/` with `make skills-sync`, then
   raise both version fields of `.agents/skills/roundfix/SKILL.md` and its
   mirror with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
   -record-skill-versions -count=1`, recording
   `skills/testdata/owned-skill-versions.json`.
9. MUST re-record the Behavior Surface fingerprint as the last step, after
   the last guide edit: `go test -tags docscontract ./internal/docscontract
   -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
   -count=1`, recording `docs/references/behavior-surfaces.json`.

## Subtasks

- [ ] Characterize the refusals in an absent Spec Root.
- [ ] Add the absent value, `ResolveSpecsRootOrAbsent` and `RequirePresent`.
- [ ] Switch the enumerating readers and guard the named-Spec readers.
- [ ] Add the named tests.
- [ ] Update the configuration guide, the four skill references and the glossary; sync, record the skill version, re-record the Behavior Surface fingerprint last.

## Acceptance Criteria

- [ ] Every enumerating command gives the same output and exit code with an absent and an empty Spec Root.
- [ ] A named Spec, a dangling link, a file and a missing external root are refused with today's text.
- [ ] The guide, the four references and **Spec Root** state the rule and cite ADR-0348.

## Context

- instruction: `docs/adr/0348-an-absent-spec-root-inside-the-repository-holds-no-spec.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- interface: `internal/config/config.go`
- creates: `internal/config/specs_root_absent_test.go`
- interface: `internal/cli/follow_up_inventory.go`
- interface: `internal/cli/spec_check.go`
- interface: `internal/cli/deliver.go`
- interface: `internal/cli/deliver_plan.go`
- interface: `internal/cli/history.go`
- interface: `internal/cli/baseline_update_history.go`
- creates: `internal/cli/absent_spec_root_test.go`
- interface: `docs/user-guide/configuration.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `CONTEXT.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- `config.ResolveSpecsRoot` and the new `ResolveSpecsRootOrAbsent` — consumers: `./internal/config`, `./internal/cli`
- The owned-skill version record and the Behavior Surface record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestResolveSpecsRoot.*|TestAnAbsentSpecRootReadsLikeAnEmptyOne|TestBaselineUpdateHistoryPlansAnAbsentSpecRootLikeAnEmptyOne|TestANamedSpecInAnAbsentSpecRootKeepsTheRefusal|TestFollowUp.*|TestDeliverPlan.*|TestHistorySanitize.*|TestSpecCheckWithoutStageIsUnchanged)$' ./internal/config ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestResolveSpecsRootOrAbsentReadsAnAbsentRootInsideTheRepository TestResolveSpecsRootOrAbsentStillRefusesAMissingExternalRoot TestResolveSpecsRootOrAbsentStillRefusesADanglingLink TestResolveSpecsRootOrAbsentStillRefusesAFile TestResolveSpecsRootRejectsInvalidRoots TestAnAbsentSpecRootReadsLikeAnEmptyOne TestBaselineUpdateHistoryPlansAnAbsentSpecRootLikeAnEmptyOne TestANamedSpecInAnAbsentSpecRootKeepsTheRefusal TestFollowUpListsEveryPendingItemOnce TestDeliverPlanReportsAnApprovedAndABlockedSpec TestHistorySanitizeLeavesACancelledRecord; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; GOOS=windows go build -buildvcs=false ./internal/config ./internal/cli` — expected: exit 0. Before this Task the absent-root tests do not exist, so the first pass check fails. Afterwards they pass with the existing resolver, follow-up, deliver plan, history and spec check tests that consume the resolver, and both packages build for Windows.
- `phrase='An absent Spec Root inside the repository reads as an empty one'; for file in docs/user-guide/configuration.md .agents/skills/roundfix/references/spec.md .agents/skills/roundfix/references/deliver.md .agents/skills/roundfix/references/archive.md .agents/skills/roundfix/references/baseline.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done; awk 'index($0, "**Spec Root**:") == 1 { f=1; next } f && index($0, "_Avoid_") == 1 { exit } f { print }' CONTEXT.md | tr -s '[:space:]' ' ' | grep -qF 'ADR-0348' || { echo 'CONTEXT.md **Spec Root** does not cite ADR-0348' >&2; exit 1; }; for ref in spec deliver archive baseline; do cmp ".agents/skills/roundfix/references/$ref.md" "skills/roundfix/references/$ref.md" || exit 1; done; cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no guide carries the phrase, so the first check fails. Afterwards the guide and the four references carry it, **Spec Root** cites ADR-0348, the mirrors equal their canonical copies, and the owned-skill version and documentation contracts, the Behavior Surface fingerprint among them, pass.

## References

- `_prd.md` → Goals; Core Features 2-3; User Experience; Non-Goals / Out of Scope; Success Metric 3; Success Metric 4; Glossary
- `_techspec.md` → Current behavior 1-2; Interfaces; API Contracts 1-3; Invariants 1-4; Testing Approach; Build Order 1
- ADR-0348
- [Adopted Backlog Entry](references/2026-10-30-the-follow-up-command-rejects-documented-message-types-and-an-empty-spec-root.md)
