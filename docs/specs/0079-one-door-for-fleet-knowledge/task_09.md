---
task: task_09
spec: 0079-one-door-for-fleet-knowledge
status: pending
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
