---
status: pending
type: docs
---

# Task: Document Hook Strictness Invariant

Write invariant into Baseline module.

## Work
- Add section to docs/agents/autonomous-work.md
- Text: "A commit hook must not be stricter"
- Cite ADR-0098
- Use managed marker boundaries

## Verification
- `grep -q "commit hook must not be stricter" docs/agents/autonomous-work.md`


## References

- User Story 3: Invariant documented
- Core Feature 4: Invariant Documentation

## Result
Invariant documented in rendered guidance.
