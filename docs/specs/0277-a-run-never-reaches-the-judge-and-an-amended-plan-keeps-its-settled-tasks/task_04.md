---
task: task_04
spec: 0277-a-run-never-reaches-the-judge-and-an-amended-plan-keeps-its-settled-tasks
status: pending
type: backend
complexity: medium
---

# Task 04: The history cleanup reads every ADR reference form adopters use

## Overview

The history cleanup resolves a superseded ADR's successor only from
`superseded_by: ADR-NNNN`, and counts a successor as naming its predecessor
only through `ADR-NNNN` or the predecessor's file name, so on 2026-10-17 the
fluxus adopter's ADR 0002 was refused although its successor 0064 already
named it (`_techspec.md` → Current behavior). This Task makes the plan that
`roundfix history sanitize` and `roundfix baseline update` share read
`superseded_by` and `supersedes` in every form adopters write. It is
verifiable on its own through planner tests on temporary repositories.

## Requirements

1. MUST answer the lead's measurement of 2026-10-17 on a scratch clone of the
   fluxus adopter, cited in `_prd.md` → Sources 5, by implementing API
   Contract 8 of `_techspec.md` in `internal/spec/history_entries.go`:
   `historyADRSuccessor` resolves `superseded_by` written as `ADR-<n>` or `<n>`
   to the accepted `docs/adr/<n>-*.md`, and as `<n>-<slug>`, `<n>-<slug>.md` or
   a path whose last element is `<n>-<slug>.md` to the accepted
   `docs/adr/<n>-<slug>.md` exactly; anything else, or a match that is not
   `status: accepted`, stays a refusal by path with today's reason.
2. MUST treat a successor as already naming the retired ADR, and append no
   lineage line for it, when its body names it as today or when its front
   matter `supersedes`, a string or a list of strings, holds `ADR-<n>`, `<n>`,
   `<stem>`, `<stem>.md` or a path ending in `/<stem>.md` for the retired
   file's number `<n>` and name `<stem>` without `.md`. The lineage line's
   bytes, the refusal text, the History Citations and every other unit MUST
   NOT change (ADR-0282).
3. MUST add the tests named in Verification, each calling `t.Parallel()`
   first, in the new `internal/spec/history_successor_forms_test.go`:
   `TestHistoryADRSuccessorResolvesEveryAcceptedForm` (each `superseded_by`
   form resolves its accepted successor, and a slug that names no file, a
   non-accepted match and an unknown form are each refused by path),
   `TestHistoryADRSuccessorThatSupersedesInFrontMatterGetsNoLineage` (each
   `supersedes` form, as a string and in a list, suppresses the line, and a
   `supersedes` naming another ADR does not) and
   `TestHistoryPlanRemovesTheFluxusShapeWithoutALineageLine` (a superseded ADR
   under a legacy ADR root with `superseded_by` written as number and slug, and
   an accepted successor whose `supersedes` names it by number and slug: the
   plan removes the ADR, refuses nothing and lists no lineage line). Every new
   test compiles for Windows: no Unix-only system call outside a
   `unix`-constrained file. Every existing history test keeps its name and
   expected text.
4. MUST document the accepted forms of `superseded_by` and `supersedes`,
   naming the `NNNN-slug.md` form, in `docs/user-guide/commands/history.md`,
   `docs/user-guide/commands/baseline.md`,
   `.agents/skills/roundfix/references/archive.md` and
   `.agents/skills/roundfix/references/baseline.md`, next to the refusal text
   they already quote. No text inside any `### QA settlement` section may
   change.
5. MUST run `make skills-sync`, raise and record the Roundfix skill version
   with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   and re-record the Behavior Surface Record with
   `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
   because this Task edits user guide pages. No record or mirror is
   hand-edited.

## Subtasks

- [ ] Resolve every `superseded_by` form.
- [ ] Read `supersedes` front matter when deciding whether a successor names its predecessor.
- [ ] Add the three planner tests.
- [ ] Document the forms in the guides and the Roundfix references.
- [ ] Sync the mirrors, record the version and the Behavior Surface Record.

## Acceptance Criteria

- [ ] Every accepted `superseded_by` form resolves its successor, and every other value is refused by path with today's text.
- [ ] A successor that names its predecessor in `supersedes` in any accepted form gets no lineage line.
- [ ] The fluxus shape leaves without a refusal or a lineage line.

## Consumers

- The history plan of `history sanitize` and `baseline update` — consumers: `./internal/spec`, `./internal/cli`

## Context

- instruction: `docs/adr/0282-the-history-cleanup-names-a-superseded-adr-in-its-successor-and-a-large-candidate-is-reviewed-in-segments.md`
- instruction: `docs/adr/0278-a-retirement-deletes-and-the-history-root-holds-only-archive-records.md`
- instruction: `docs/adr/0254-baseline-update-sanitizes-pending-history-in-the-change-it-plans.md`
- instruction: `internal/spec/history_lineage_test.go`
- interface: `internal/spec/history_entries.go`
- creates: `internal/spec/history_successor_forms_test.go`
- interface: `docs/user-guide/commands/history.md`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -count=1 -v -run '^(TestHistoryADRSuccessorResolvesEveryAcceptedForm|TestHistoryADRSuccessorThatSupersedesInFrontMatterGetsNoLineage|TestHistoryPlanRemovesTheFluxusShapeWithoutALineageLine)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestHistoryADRSuccessorResolvesEveryAcceptedForm TestHistoryADRSuccessorThatSupersedesInFrontMatterGetsNoLineage TestHistoryPlanRemovesTheFluxusShapeWithoutALineageLine; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for f in docs/user-guide/commands/history.md docs/user-guide/commands/baseline.md .agents/skills/roundfix/references/archive.md .agents/skills/roundfix/references/baseline.md; do grep -qF -- 'NNNN-slug.md' "$f" || { printf 'missing the NNNN-slug.md form in %s\n' "$f" >&2; exit 1; }; done; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -run 'History' ./internal/spec || exit 1; go test -count=1 -run 'History|BaselineUpdate' ./internal/cli || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the three planner tests do not exist, so the command fails at its first pass check; afterwards the planner, history command, documentation, skill version and Behavior Surface Record checks pass.

## References

- `_prd.md` → Sources; Goals; User Story 5; Core Feature 8; Core Feature 9; Success Metric 5
- `_techspec.md` → API Contract 8; API Contract 9; Invariant 7; Current behavior; Testing Approach; Build Order 4
- ADR-0282; ADR-0278; ADR-0254; ADR-0280; ADR-0189; ADR-0256; ADR-0259
