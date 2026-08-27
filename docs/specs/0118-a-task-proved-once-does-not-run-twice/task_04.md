---
status: pending
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
