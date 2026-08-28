---
status: completed
type: docs
---

# Task: The glossary names the accepted outcomes

The term **Task Carry-Forward** and the corrected Reconcile Command entry
landed with this Spec's authoring, because the Vocabulary Contract detector
refuses a Spec whose coined term is undocumented and because the claim that
`--apply` is the command's only mutation switch was already false. What could
not be written then is which Run outcomes the act accepts: the sentence only
became true when task_01 widened the guard.

## Work

- Record in the Task Carry-Forward entry that the act accepts a Run whose
  outcome is Stopped or Unresolved, and that every other terminal outcome is
  refused.
- Read the claim from the delivered code, not from the TechSpec draft. Where
  the two disagree, the code is the fact and the TechSpec is corrected.
- Check whether this Spec introduced, changed, or retired any other term the
  glossary should carry, and update the domain context when it did.
- Keep the entry in the shape of its neighbours: definition, then what to
  avoid.

## References

- `_prd.md` → Core Feature 6
- `_techspec.md` → Build Order 4; Vocabulary Contract
- `docs/agents/domain.md` owns what the glossary check looks for

## Verification
- `grep -q "Stopped or Unresolved" CONTEXT.md && grep -q "Task Carry-Forward" CONTEXT.md && go test -count=1 -tags docscontract ./internal/docscontract 2>&1 | grep -q "^ok"`

## Result

Implementation:

- Updated the Task Carry-Forward glossary definition to state that it accepts a
  Run whose outcome is Stopped or Unresolved and refuses every other terminal
  outcome.
- Corrected the TechSpec's stale System Architecture wording to match the
  delivered `carryForwardAcceptedStates` set and both membership guards.
- Audited the Spec's vocabulary against the full glossary and the delivered
  implementation. No additional Spec-introduced, changed, or retired term
  needs a glossary entry; the existing Reconcile Command entry already names
  all three mutation switches.

Focused checks:

- Pre-change signal: `rtk rg -n -F 'Stopped or Unresolved' CONTEXT.md` exited
  `1`; `rtk rg -n -F 'not Stopped' docs/specs/0118-a-task-proved-once-does-not-run-twice/_techspec.md`
  found the stale architecture sentence.
- `rtk rg -n -F 'Stopped or Unresolved' CONTEXT.md` exited `0` and found the
  accepted-outcome sentence at line 442.
- `rtk rg -n -C 2 'carryForwardAcceptedStates|slices.Contains\\(carryForwardAcceptedStates'
  internal/cli/carryforward.go internal/cli/reconcile.go` exited `0` and
  found the delivered accepted set plus both guards.
- The stale TechSpec searches for `not Stopped` and `non-Stopped` each exited
  `1`; the corrected `Stopped nor Unresolved` search exited `0`.
- `rtk git diff --check` exited `0`.

Acceptance evidence:

1. The glossary entry names `Stopped or Unresolved` as the accepted outcomes
   and states that every other terminal outcome is refused.
2. The implementation search confirms the wording follows the delivered
   accepted-state set and guard membership checks, while the TechSpec now
   describes those facts rather than its pre-task draft.
3. The glossary audit found no other term requiring an entry, and the existing
   Reconcile Command entry retains the `--apply`, `--discard-superseded`, and
   `--carry-forward` mutation switches.

The Daemon-owned Verification command was not run during this Agent turn.
