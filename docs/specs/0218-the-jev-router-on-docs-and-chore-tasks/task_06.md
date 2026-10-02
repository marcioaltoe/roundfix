---
task: task_06
spec: 0218-the-jev-router-on-docs-and-chore-tasks
status: pending
type: chore
complexity: medium
---

# Task 06: Run the Jev Router measurement in committed scratch clones

## Overview

Corrective Task for QA finding F1 of the 2026-10-02 QA Report (row 5). Task 04 produced no measurement: the replay protocol needs a Task-only graph in each scratch clone, `roundfix implement` reads the Spec graph from `HEAD`, and the Agent declined to commit inside the scratch clone because Task Agents never commit. The no-commit rule protects this repository and its Run Worktree, not a disposable clone. This Task states that a scratch clone may hold its own commits and runs the measurement Task 04 specified.

## Requirements

1. MUST follow Task 04's Requirements 1 to 7 and the TechSpec "The measurement protocol", with one clarification: inside each disposable scratch clone under the system temporary directory, the Agent MAY create the commits needed to build the replay pre-state, such as the Task-only graph and the scratch `.roundfixrc.yml`. Those commits are never pushed, never made in this repository or its Run Worktree, and are removed with the clone.
2. MUST replace the zero-row record at `docs/specs/0218-the-jev-router-on-docs-and-chore-tasks/measurement/jev-router.md` with the measured rows: at least one routed replay of a `docs` or `chore` Task with its router cost, or the refusal that stopped routed replays under the shared ceiling, recorded as Task 04 Requirement 4 requires.
3. MUST NOT print, echo or write `ROUNDFIX_OPENROUTER_API_KEY`, MUST NOT read `OPENROUTER_API_KEY`, and MUST NOT raise, bypass or work around the monthly ceiling.
4. MUST mint the Backlog Entry of Task 04 Requirement 6 only when its condition holds; otherwise `## Reading` says why not.

## Subtasks

- [ ] Build each replay pre-state with commits inside its scratch clone.
- [ ] Run the routed replays first, then the default replays, and record every row.

## Acceptance Criteria

- [ ] The record holds measured rows with categories and router costs, or the recorded ceiling refusal.
- [ ] No commit, configuration or profile of this repository changed outside the record and an optional Backlog Entry.

## Context

- interface: `docs/specs/0218-the-jev-router-on-docs-and-chore-tasks/measurement/jev-router.md`
- instruction: `docs/specs/0218-the-jev-router-on-docs-and-chore-tasks/task_04.md`

## Verification

- `record='docs/specs/0218-the-jev-router-on-docs-and-chore-tasks/measurement/jev-router.md'; test -f "$record" || { printf 'the measurement record is missing\n' >&2; exit 1; }; if grep -q 'No replay ran' "$record"; then printf 'the record still says no replay ran\n' >&2; exit 1; fi; grep -Eq '^[|] [^|]+ [|] (docs|chore) [|]' "$record" || grep -qi 'ceiling' "$record" || { printf 'no measured docs or chore row and no ceiling refusal\n' >&2; exit 1; }; grep -q '^## Reading' "$record" || { printf 'missing Reading section\n' >&2; exit 1; }` — expected: exit 0

## References

- QA Report 2026-10-02 → F1; row 5
- Task 04 → Requirements 1 to 7
- `_techspec.md` → The measurement protocol
