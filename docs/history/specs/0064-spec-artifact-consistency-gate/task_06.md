---
task: task_06
spec: 0064-spec-artifact-consistency-gate
status: completed
type: docs
complexity: low
---

# Task 06: Name the check in the glossary

## Overview

The check mints user-facing vocabulary — the command, its two severities, and
the `SC-*` diagnostic codes — and the repository's rule is that user-facing
vocabulary lives in the glossary before it is used in help text, findings, and
reports. This slice adds those entries so every later reader of a finding can
resolve its terms in one place.

## Requirements

1. MUST add a Spec Consistency Check glossary entry defining the check as
   read-only, pre-Run, citation-based, and never a QA verdict, with its
   `_Avoid_` line.
2. MUST define the two finding severities: an `error` names a contradiction
   whose two sides the check located; a `gap` names a candidate the check
   surfaced but cannot settle.
3. MUST state that the `SC-*` diagnostic codes are stable and never renumbered
   once shipped.
4. MUST use the existing glossary entry shape — bold term, definition,
   `_Avoid_` line — and place the entries beside the other support-command
   terms.
5. MUST NOT change any other glossary entry.

## Subtasks

- [x] Add the Spec Consistency Check entry.
- [x] Add the severity entries.
- [x] Record the diagnostic-code stability rule.

## Acceptance Criteria

- [x] The glossary defines Spec Consistency Check with an `_Avoid_` line.
- [x] The glossary defines both severities with their distinct meanings.
- [x] The glossary states the `SC-*` codes are stable once shipped.
- [x] No pre-existing glossary entry is modified.

## Verification

- `grep -q "^\*\*Spec Consistency Check\*\*:" CONTEXT.md` — expected: exit 0.
- `grep -A 3 "^\*\*Spec Consistency Check\*\*:" CONTEXT.md | grep -q "^_Avoid_:"`
  — expected: exit 0; the entry carries its `_Avoid_` line.
- `grep -q "^\*\*Consistency Finding Severity\*\*:" CONTEXT.md` — expected:
  exit 0; the two severities are defined.
- `grep -A 4 "^\*\*Consistency Finding Severity\*\*:" CONTEXT.md | grep -q "gap"`
  — expected: exit 0; the gap severity is distinguished from the error.
- `grep -q "SC-" CONTEXT.md` — expected: exit 0; the diagnostic-code family is
  named.
- `git diff --name-only HEAD | grep -v "^CONTEXT.md$" | grep -v "^docs/specs/0064-spec-artifact-consistency-gate/" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; this Task changed only the glossary and its own Task file.

## References

- `_prd.md` → Project Constraints (identifier strategy).
- `_techspec.md` → Project Constraints; Integration Points; Build Order 8.
- `docs/agents/domain.md`.

## Result

Implemented the two glossary entries beside the existing support-command
terms. Spec Consistency Check now names the read-only, pre-Run,
citation-based boundary and excludes artifact mutation and QA verdicts;
Consistency Finding Severity distinguishes a located contradiction (`error`)
from an unsettled candidate (`gap`) and records the stability rule for shipped
`SC-*` diagnostic codes.

Focused checks:

- Pre-change `rtk grep -n "Spec Consistency Check|Consistency Finding Severity|SC-" CONTEXT.md` returned no matches, establishing that the vocabulary was absent before this slice.
- `rtk proxy sed -n '372,398p' CONTEXT.md` showed both new bold-term entries, each with its definition and `_Avoid_` line, between Doctor Command and Archive Command.
- `rtk git diff -- CONTEXT.md` reported eight insertions and no deletions, so no pre-existing glossary entry changed.
- The focused banned-pattern scan over `CONTEXT.md` returned no matches for the tech-writer skill's prohibited filler vocabulary.

Acceptance evidence:

- Spec Consistency Check is defined as a read-only, pre-Run support command that compares written citations, declarations, and cross-references, never edits artifacts, and never emits a QA verdict; the entry includes `_Avoid_`.
- Consistency Finding Severity defines `error` as a contradiction with both sides located and `gap` as a surfaced candidate the check cannot settle; the entry includes `_Avoid_`.
- The severity entry states that each `SC-*` diagnostic code is stable and never renumbered once shipped.
- The focused diff contains only the two new entries, with no deletion or replacement of existing glossary text.

Not run: the commands under `## Verification`; the Roundfix Daemon owns that
verification and Task settlement.
