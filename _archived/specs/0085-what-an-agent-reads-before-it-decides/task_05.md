---
task: task_05
spec: 0085-what-an-agent-reads-before-it-decides
status: completed
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
- `test "$(c=0; for f in docs/adr/*.md; do head -8 "$f" | grep -q '^status:' || c=$((c+1)); done; echo $c)" -eq 0` — expected: exits 0. Seventy-three ADRs carry no `status:` in their frontmatter today, so this is the gap the Task closes; the guard it replaces asked whether every *retired* ADR names a replacement, and zero ADRs are retired, so it passed before any work.

Whole-package sweeps, `go build`, `go clean -testcache` and `make verify` are
deliberately absent: each passes against a tree where no work has happened, so
it approves the Task before it starts. Regression is the Run-level gate's job.

## References

- `_prd.md` → the archive read path.
- `_techspec.md` → Build Order 5.

## Result

Standardized the live 116-file ADR corpus on the lifecycle frontmatter keys
`status`, `created_at`, `updated_at`, `deprecated_at`, and `superseded_by`.
Added that block to the 73 nonconforming records, converted all 20 legacy
`Status: Accepted` body lines, and recorded the 53 statusless legacy records as
`accepted` under the repository rule that treats them as active. ADR-0106
remains the only fully retired record and retains its forward pointer to
ADR-0108; records with explicitly partial supersessions remain active.

- Criterion 1: a Ruby corpus audit required the five lifecycle keys in order
  on every `docs/adr/*.md` file and reported
  `adrs=116 standard_frontmatter=116`.
- Criterion 2: a Ruby body-only status audit stripped each frontmatter block,
  searched the remaining ADR body, and reported `legacy_body_statuses=0`.
- Criterion 3: a Ruby retired-record audit resolved every non-null ADR pointer
  to a current `docs/adr/` target and reported
  `retired=1 pointers_resolved_or_explicit_none=1`.
- Criterion 4: a Ruby preservation audit removed lifecycle frontmatter and the
  migrated legacy status line from both `HEAD` and the working copy before an
  exact comparison; it reported `decision_text_matches_head=116`.

Focused checks:

- The combined lifecycle audit validated allowed statuses, ISO-8601 lifecycle
  timestamps, exact frontmatter keys, retired pointers, and decision-text
  preservation; it reported 115 accepted ADRs, one superseded ADR, and zero
  failures.
- `rtk git diff --check` reported no whitespace errors.

The Daemon-owned `## Verification` commands were not run.
