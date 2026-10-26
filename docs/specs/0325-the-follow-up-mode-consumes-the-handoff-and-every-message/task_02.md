---
task: task_02
spec: 0325-the-follow-up-mode-consumes-the-handoff-and-every-message
status: completed
type: docs
complexity: low
---

# Task 02: The glossary defines Session Handoff

## Overview

The `--follow-up` Skill Mode now claims the project's Session Handoff and
hands its items to the Follow-up Command, but `CONTEXT.md` has no term for
the handoff and its **Follow-up Item** entry does not name those items. This
Task adds **Session Handoff** and revises **Follow-up Item** through
`domain-modeling`.

## Requirements

1. MUST, through `domain-modeling`, add the **Session Handoff** entry of
   `_techspec.md` → Interfaces to `CONTEXT.md`, directly after the
   `_Avoid_: Inbox Entry, capture, note` line of **Cross-project Message**,
   separated by one blank line, with its link to ADR-0338.
2. MUST revise **Follow-up Item** as `_techspec.md` → Interfaces states,
   keeping the rest of that entry and every other entry unchanged.
3. MUST keep the glossary contracts passing.
4. MUST NOT change any file other than `CONTEXT.md`.

## Subtasks

- [ ] Add **Session Handoff**.
- [ ] Revise **Follow-up Item**.

## Acceptance Criteria

- [ ] `CONTEXT.md` defines **Session Handoff** with its ADR-0338 link and one `_Avoid_` line.
- [ ] **Follow-up Item** names a Session Handoff's next steps and open questions handed in as message entries.
- [ ] The glossary contracts pass.

## Consumers

None. `CONTEXT.md` is read by the glossary contracts and by people; no code
consumes the two entries.

## Context

- instruction: `docs/agents/domain.md`
- instruction: `docs/adr/0338-the-follow-up-mode-claims-the-handoff-and-leaves-no-handed-in-message-pending.md`
- interface: `CONTEXT.md`

## Verification

- `tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Session Handoff**:' && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Follow-up Item**: One pending unit the Follow-up Command lists: an ai-memory message the caller hands in, a next step or open question of a Session Handoff the caller hands in as a message entry' && go test -timeout 15m -count=1 -tags docscontract -run '^TestGlossary' ./internal/docscontract` — expected: exit 0. Before this Task the term is absent, so the first check fails; afterwards both entries are present and the glossary contracts pass.

## References

- `_prd.md` → Goals 5; Core Feature 5; Success Metric 5; Glossary
- `_techspec.md` → Interfaces; Glossary; Build Order 2

## Result

Implemented the two `CONTEXT.md` edits through `domain-modeling`, applied
verbatim from `_techspec.md` → Interfaces:

- Added the **Session Handoff** entry directly after the
  `_Avoid_: Inbox Entry, capture, note` line of **Cross-project Message**,
  separated by one blank line, with the ADR-0338 link and one `_Avoid_` line
  (`Handoff file, session note, resume brief`).
- Revised **Follow-up Item** to insert `a next step or open question of a
  Session Handoff the caller hands in as a message entry with id
  \`handoff:<handoff id>:<n>\`,` between `an ai-memory message the caller
  hands in,` and `an open Backlog Entry,`; every other byte of that entry and
  every other entry is unchanged.

Focused check (the Task's own Verification command, run while working):

```
tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Session Handoff**:' && \
  tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Follow-up Item**: One pending unit the Follow-up Command lists: an ai-memory message the caller hands in, a next step or open question of a Session Handoff the caller hands in as a message entry' && \
  go test -timeout 15m -count=1 -tags docscontract -run '^TestGlossary' ./internal/docscontract
```

Result: `ok roundfix/internal/docscontract 0.641s` (exit 0).

Acceptance criteria evidence:

- **`CONTEXT.md` defines Session Handoff with its ADR-0338 link and one
  `_Avoid_` line** — the entry is present with the
  `[ADR-0338](docs/adr/0338-...)` link and the single `_Avoid_: Handoff file,
  session note, resume brief` line.
- **Follow-up Item names a Session Handoff's next steps and open questions
  handed in as message entries** — the entry now reads "a next step or open
  question of a Session Handoff the caller hands in as a message entry with
  id `handoff:<handoff id>:<n>`" and keeps the rest of its bytes.
- **The glossary contracts pass** — `TestGlossary` passes under the
  `docscontract` build tag (exit 0).

Scope: only `CONTEXT.md` changed in the working tree (`git status --porcelain`
shows this task file's Result section as the sole other modification).

