---
task: task_02
spec: 0210-evidence-snapshots-that-stay-small
status: pending
type: docs
complexity: low
---

# Task 02: The guide describes one digest per declared input

## Overview

The "Evidence snapshots and row carry-forward" section of the Context-Driven
Development guide names the stale reason `input moved: <paths>` and says
nothing of the record's shape. After task_01, the Evidence Snapshot holds one
line per declared input and the reason names refs. This Task brings the guide
in line, and changes nothing else.

## Requirements

1. MUST add to the paragraph that defines the **Evidence Snapshot** in
   `docs/user-guide/context-driven-development.md` that the record holds one
   line per declared input: its ref, the number of files it matched and one
   SHA-256 digest over their sorted paths and content digests. It MUST add
   that its size therefore depends on rows and inputs, never on how many files
   a glob matches, and that a report recorded in the earlier per-file form is
   still read. The paragraph MUST contain the exact phrase
   `one line per declared input`.
2. MUST replace the reason `input moved: <paths>` in the same section with
   `input moved: <refs>`, and state that it names the declared inputs whose
   matched files changed.
3. MUST keep every other reason, phase and payload word of that section, so
   the Vocabulary Contract of Spec 0202 stays documented.
4. MUST NOT edit `CONTEXT.md`, any skill, the QA prompt contract, any test or
   any other section of the guide.

## Subtasks

- [ ] Describe the per-input record in the Evidence Snapshot paragraph.
- [ ] Replace the `input moved:` placeholder.

## Acceptance Criteria

- [ ] The guide contains `one line per declared input` and
      `input moved: <refs>`, and no longer contains `input moved: <paths>`.
- [ ] Every reason, phase and payload word Spec 0202 documented in that
      section is still present.

## Context

- interface: `docs/user-guide/context-driven-development.md`
- instruction: `docs/adr/0210-an-evidence-snapshot-records-one-digest-per-declared-input.md`

## Verification

- `for phrase in 'one line per declared input' 'input moved: <refs>' 'evidence_snapshots' 'Row carry-forward' 're-run: not pass' 'no evidence snapshot' 'evidence differs' 'prior_report' 'carried_rows' 'rerun_rows'; do tr -s '[:space:]' ' ' < docs/user-guide/context-driven-development.md | grep -qF -- "$phrase" || { printf 'missing phrase in guide: %s\n' "$phrase" >&2; exit 1; }; done && ! tr -s '[:space:]' ' ' < docs/user-guide/context-driven-development.md | grep -qF -- 'input moved: <paths>'` — expected: exit 0; before this Task the guide names neither new phrase and still holds `input moved: <paths>`, so the command fails.

## References

- [_prd.md](_prd.md) — Core Feature 5; Goal 1
- [_techspec.md](_techspec.md) — Vocabulary Contract; API Contract 1; API Contract 2; Build Order 2
- ADR-0210; ADR-0195
