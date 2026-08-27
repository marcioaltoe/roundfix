---
status: pending
type: qa
---

# Task: QA gate

Verify every deliverable of this Spec against the running commands.

## Work

- An Unresolved Run's proved Tasks are handed back, and every proof still
  refuses what it refused before
- A Stopped Run behaves exactly as it did, proving the widening changed the
  selection and not the proofs
- A Run whose outcome is neither is refused with its outcome named
- An implement whose Spec has carriable stranded work is refused with no Run
  created, and the named command clears it
- An implement whose stranded work is not carriable proceeds with a report
- An implement whose carry-forward inspection fails proceeds, with the failure
  reported
- The refusal names the Run with the largest carriable set
- The glossary carries the term with its accepted outcomes, and the vocabulary
  detector runs
- The skill's changed paths stay inside the recorded authorization's bounded
  files, checked from Git evidence
- The glossary check: whether this Spec introduced, changed, or retired a term
  the domain context should carry

## Outside evidence

One acceptance row rests on evidence this Spec did not author. A `fiscus`
session on 2026-08-07/08 delivered one Spec through five terminal Runs, all
Unresolved, totalling 5h21m — a repository this Spec did not build, measured
before this Spec existed, recorded at
`references/2026-08-12-five-unresolved-runs-to-deliver-one-spec.md` and
provenanced in `references/_index.md`. The row records that this measurement is
what establishes the requirement, rather than a rehearsal of the Spec's own
premise.

## References

- All user stories and core features

## Verification
- `roundfix spec check 0118-a-task-proved-once-does-not-run-twice --strict && go test -count=1 ./internal/cli ./internal/spec ./internal/store 2>&1 | grep -q "^ok"`
