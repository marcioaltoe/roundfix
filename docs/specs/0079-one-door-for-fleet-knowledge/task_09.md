---
task: task_09
spec: 0079-one-door-for-fleet-knowledge
status: completed
type: chore
complexity: medium
---

# Task 09: Teach the guides what the operator was promised

## Overview

Corrective slice opened by the QA gate's finding F-002, the graph's last red
row. The PRD's User Experience section promises four operator behaviors that
never reached a durable instruction: Triage works the queue oldest first, an
empty inbox means rest, a findings directory holding only live work means
health, and a Rollup with no open members is a candidate for its own closure.
The contracts, checks, and pilot all shipped; the operator-facing reading of
them did not. This Task carries those four into the same modules the rest of
the flow lives in, under the same recorded authorization as task_02 and
task_06.

The measured trap this Task must not repeat: new clauses grow the Source
Baseline, and the maintained compatibility fixture's declared expectation
then names a stale count — the exact defect task_08 corrected. This Task
owns that correction for its own fallout, in its own commit, so it never
leaves the repository gate red for the gate that follows it.

## Requirements

1. MUST author the queue instructions in
   `internal/baseline/assets/modules/secondbrain.json`, using these literal
   tokens: Triage works pending entries `oldest first`; an
   `empty inbox` is rest, not a missing step.
2. MUST author the findings-side instructions in
   `internal/baseline/assets/modules/context-workflow.json`, using these
   literal tokens: a findings directory holding only live work is `health`,
   not loss, because rollups and the archive hold what was learned; a Rollup
   with `no open members` is a candidate for its own closure.
3. MUST write them as operator instructions a session can act on, not as
   restated contract — the reader must learn what to do with the queue and
   how to read an empty state.
4. MUST run the module chain per the Spec 0075 choreography: bootstrap the
   Source Baseline manifest rows for the new clauses, then run
   `make baseline-digests` twice, and adopt the regenerated managed-block
   postimages into `docs/agents/secondbrain.md` and
   `docs/agents/docs-layout.md`.
5. MUST move the maintained Source Baseline expectation in
   `internal/baseline/preservation_test.go` to the value the regenerated
   fixture carries whenever this Task's clauses grow the corpus, keeping the
   identity-agrees-with-entries and accounting assertions unweakened; the
   repository gate MUST be green when this Task settles.
6. MUST NOT restate or reword any clause task_02 or task_06 already
   authored; this Task adds the operator reading only.
7. MUST change only the authorization's bounded files, `internal/baseline/**`
   as sanctioned deterministic fallout plus the named expectation, and this
   task file.

## Subtasks

- [ ] Author the four operator instructions in the two modules.
- [ ] Run the chain; adopt both postimages.
- [ ] Re-record the maintained expectation if the corpus grew; prove green.

## Acceptance Criteria

- [ ] Both guides carry the four operator instructions with the literal
      tokens from Requirements 1–2.
- [ ] No clause authored by task_02 or task_06 changed.
- [ ] The digest chain is converged and `internal/baseline` is green.
- [ ] The diff stays inside the bounded files, the baseline package, and
      this task file.

## Context

- instruction: docs/workflow/authorizations/2026-08-06-fleet-knowledge-door.md
- interface: internal/baseline/assets/modules/secondbrain.json
- interface: internal/baseline/assets/modules/context-workflow.json
- interface: internal/baseline/preservation_test.go

## Verification

- `grep -q 'oldest first' docs/agents/secondbrain.md && grep -qi 'empty inbox' docs/agents/secondbrain.md`
  — expected: exit 0; the queue order and the empty-inbox reading are
  adopted in the secondbrain guide.
- `grep -q 'no open members' docs/agents/docs-layout.md && grep -qi 'health' docs/agents/docs-layout.md`
  — expected: exit 0; the rollup-closure candidacy and the healthy-directory
  reading are adopted in the layout guide.
- `go test -count=1 ./internal/baseline/...`
  — expected: exit 0; the corpus growth this Task caused did not leave the
  maintained expectation stale.
- `s1="$(git status --porcelain | sort)"; make baseline-digests >/dev/null 2>&1; s2="$(git status --porcelain | sort)"; [ "$s1" = "$s2" ]`
  — expected: exit 0; the digest chain is converged.
- `output="$(git status --porcelain | awk '{print $NF}' | grep -vE '^(internal/baseline/|docs/agents/docs-layout\.md$|docs/agents/secondbrain\.md$|docs/agents/setup-context\.json$|docs/specs/0079-one-door-for-fleet-knowledge/task_09\.md$)')"; [ -z "$output" ]`
  — expected: exit 0; nothing outside the bounded files, the baseline
  package, and this task file changed.

## References

- `_prd.md` → User Experience; User Stories 2–3.
- `qa/qa-report-2026-08-06-01.md` → Finding F-002 and row D01.
- `_techspec.md` → Testing Approach (module choreography).
- ADR-0095, ADR-0081.

## Result

### Implementation

- Versioned the secondbrain module, its guide, and its read-only rule. Added
  separate operator clauses that direct Triage to work pending entries
  `oldest first` by `created_at` and treat an `empty inbox` as rest rather
  than an omitted step.
- Versioned the context-workflow module, its docs-layout guide, and its layout
  rule. Added separate operator clauses that read a findings directory holding
  only live work as `health` because Rollups and the archive preserve what was
  learned, and treat a Rollup with `no open members` as a candidate for its
  own closure through the existing Finding lifecycle.
- Bootstrapped four marker-delimited Source Baseline entries and manifest rows.
  The sanctioned regeneration replaced every temporary span and digest,
  regenerated the managed guide postimages, and moved the maintained Source
  Baseline identity from 114 to 118 entries while accounting stayed at 51.
- Adopted the regenerated `guide.secondbrain` and `guide.docs-layout` managed
  blocks into the two public guides. Moved only
  `maintainedSourceBaselineEntries` from 114 to 118; the independent identity,
  exact-entry, and accounting assertions remain intact.

### Focused checks

- The pre-change QA evidence and a fresh source inspection found none of the
  four promised operator readings in either canonical module. After the module
  edit, focused `jq -e` assertions over both modules exited 0 for all four
  literal tokens and the four distinct clause IDs.
- The first
  `GOCACHE=/private/tmp/roundfix-task09-gocache rtk make baseline-digests`
  pass exited 0 with `changed:true`. After adopting both postimages, the second
  pass exited 0 with `changed:false`.
- Before the expectation edit,
  `GOCACHE=/private/tmp/roundfix-task09-gocache rtk go test
  ./internal/baseline -count=1 -run
  TestReadoptionCompatibilityMaintainedFixture` reached the intended assertion
  and reported identity 118, entries 118, and accounting 51 against the stale
  expectation. After the edit, the focused compatibility, catalog, and
  formatter selection exited 0 with three passing tests.
- `GOCACHE=/private/tmp/roundfix-task09-gocache rtk go test
  ./internal/baseline -count=1` exited 0 with 560 passing tests.
- Exact object comparisons against `HEAD` exited 0 for every clause Task 02 or
  Task 06 authored in the two modules. Exact managed-block comparisons against
  the formatter goldens exited 0 for both adopted guides. A focused manifest
  assertion confirmed four new rows with valid nonzero spans and digests.
- `rtk git diff --check` exited 0.

### Acceptance-criterion evidence

1. The two public guide blocks are byte-identical to their regenerated module
   postimages. Focused searches locate `oldest first` and `empty inbox` in the
   Secondbrain guide, and `health` and `no open members` in the docs-layout
   guide.
2. Exact `HEAD` comparisons cover the Task 02 and Task 06 clause objects,
   including the Inbox Entry, Triage, Rollup/archive, fleet-flow, automatic
   capture, research-capture, and write-boundary clauses; all comparisons
   exited 0.
3. The required regeneration sequence ended with `changed:false`. Fresh
   focused runs passed the three named seams and all 560 tests in the
   `internal/baseline` package; the maintained count assertion still compares
   identity to entries, entries to the declared 118 value, and accounting to
   the separate declared 51 value.
4. The changed-path audit contains only the two authorized guides, this Task
   file, and `internal/baseline/**` module causes, sanctioned deterministic
   fallout, Source Baseline bootstrap artifacts, and the named maintained
   expectation. The Task-file frontmatter change was the pre-existing
   Daemon-owned status transition.

### Handoff boundary

- The Daemon-owned commands under `## Verification` were not run in this Agent
  turn. Task status remains Daemon-owned; no commit, push, or pull request was
  created.
