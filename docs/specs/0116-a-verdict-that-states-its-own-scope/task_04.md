---
status: pending
type: docs
---

# Task: The glossary names the Auditing Binary

A QA Report records `build:` and a reader cannot tell whether it identifies the
Roundfix that produced the verdict or the tree it audited. The distinction has
no name, which is why one field carried both readings.

## Work

- Add **Auditing Binary** to the domain context: the Roundfix that produced a
  verdict, as distinct from the tree it audited, carrying version and build
  identity, and reported with a staleness state that may be `unknown`.
- Record what to avoid, in the shape of its neighbours.
- Read the entry from the delivered code, not from the TechSpec draft. Where
  the two disagree, the code is the fact and the TechSpec is corrected.
- Check whether this Spec introduced, changed, or retired any other term the
  glossary should carry, and update the domain context when it did.

## References

- `_prd.md` → Core Feature 3
- `_techspec.md` → Build Order 4; Vocabulary Contract
- `docs/agents/domain.md` owns what the glossary check looks for

## Verification
- `grep -q "Auditing Binary" CONTEXT.md && grep -q "auditing_binary" CONTEXT.md && go test -count=1 -tags docscontract ./internal/docscontract`
