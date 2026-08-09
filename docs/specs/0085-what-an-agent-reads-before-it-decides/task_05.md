---
task: task_05
spec: 0085-what-an-agent-reads-before-it-decides
status: pending
type: docs
complexity: medium
---

# Task 05: Give every retired ADR a status and a forward pointer

## Overview

`docs/adr/` holds 105 records with no structural separation between the accepted,
the legacy-bodied, and the statusless. An Agent reading the directory for context
cannot tell which decisions are live. This Task standardizes the frontmatter
across the corpus and gives every retired record a pointer to what replaced it.

## Requirements

1. MUST give every ADR the same frontmatter fields, so a reader can decide
   liveness from the frontmatter alone.
2. MUST convert every legacy body-line status into that frontmatter without
   changing the decision text.
3. MUST give every retired record a forward pointer to the decision that
   replaced it, or an explicit statement that nothing did.
4. MUST NOT change the reasoning in any ADR; this Task edits frontmatter and
   adds pointers.
5. MUST leave a record whose status cannot be determined marked as such rather
   than guessed.

## Subtasks

- [ ] Standardize the frontmatter across the corpus.
- [ ] Convert legacy body-line statuses.
- [ ] Add the forward pointers.

## Acceptance Criteria

- [ ] Every ADR carries the standard frontmatter fields.
- [ ] No ADR states its status only in its body.
- [ ] Every retired ADR names its replacement or states there is none.
- [ ] No decision text changed.

## Bounded scope

This Task may create or modify only:

- `docs/adr/**`
- `docs/specs/0085-what-an-agent-reads-before-it-decides/task_05.md`

## Verification

- `test -z "$(grep -Lq '^status:' docs/adr/*.md 2>/dev/null; for f in docs/adr/*.md; do head -8 "$f" | grep -q '^status:' || echo "$f"; done)"` — expected: exits 0, proving every ADR carries a frontmatter status.
- `test -z "$(grep -lE '^\*?\*?Status: ' docs/adr/*.md)"` — expected: exits 0, proving no legacy body-line status remains.
- `test -z "$(for f in $(grep -lE '^status: (deprecated|superseded)' docs/adr/*.md); do head -8 "$f" | grep -q '^superseded_by:' || echo "$f"; done)"` — expected: exits 0, proving every retired record carries a forward pointer field.
- `GOCACHE="$PWD/.gocache" go test ./internal/speccheck -count=1 2>&1 | grep -q '^ok'` — expected: exits 0, proving the ADR corpus still parses for the consistency checks that read it.

## References

- `_prd.md` → the archive read path.
- `_techspec.md` → Build Order 5.
